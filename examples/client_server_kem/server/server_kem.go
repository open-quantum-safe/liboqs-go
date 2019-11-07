// key encapsulation TCP server Go example
// run with "go run server_kem.go <port number> <KEM name>"
package main

import (
	"errors"
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: server_kem <port number> [KEM name (optional)]")
		os.Exit(-1)
	}
	port := os.Args[1]
	kemName := "DEFAULT"
	if len(os.Args) > 2 {
		kemName = os.Args[2]
	}

	fmt.Println("Launching KEM", kemName, "server on port", port)
	{
		kem := oqs.KeyEncapsulation{}
		kem.Init(kemName, nil)
		fmt.Println(kem.Details())
		kem.Clean()
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	// listen indefinitely (until explicitly stopped, e.g. with CTRL+C in UNIX)
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		// handle connections concurrently
		go handleConnection(conn, kemName)
	}
}

func handleConnection(conn net.Conn, kemName string) {
	defer conn.Close() // clean up even in case of panic

	// send KEM name to client first
	_, err := fmt.Fprintln(conn, kemName)
	if err != nil {
		panic(errors.New("server cannot send the KEM name to the client"))
	}

	// construct and initialize the KEM server
	server := oqs.KeyEncapsulation{}
	defer server.Clean() // clean up even in case of panic
	server.Init(kemName, nil)

	// read the public key sent by the client
	clientPublicKey := make([]byte, server.Details().LengthPublicKey)
	n, err := io.ReadFull(conn, clientPublicKey)
	if err != nil {
		panic(err)
	} else if n != server.Details().LengthPublicKey {
		panic(errors.New("server expected to read " + string(server.Details().
			LengthPublicKey) + " bytes, but instead read " + string(n)))
	}

	// encapsulate the secret
	ciphertext, sharedSecretServer := server.EncapSecret(clientPublicKey)

	// then send ciphertext to client and close the connection
	n, err = conn.Write(ciphertext)
	if err != nil {
		panic(err)
	} else if n != server.Details().LengthCiphertext {
		panic(errors.New("server expected to write " + string(server.
			Details().LengthCiphertext) + " bytes, but instead wrote " + string(n)))
	}

	fmt.Printf("\nServer shared secret:\n% X ... % X\n",
		sharedSecretServer[0:8], sharedSecretServer[len(sharedSecretServer)-8:])
}
