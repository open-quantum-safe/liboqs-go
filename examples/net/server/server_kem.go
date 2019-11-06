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
	if len(os.Args) != 3 {
		fmt.Println("specify both port number and KEM name")
		return
	}
	port := os.Args[1]
	kemName := os.Args[2]

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
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn, kemName)
	}
}

func handleConnection(conn net.Conn, kemName string) {
	// send KEM name to client first
	fmt.Fprintln(conn, kemName)

	server := oqs.KeyEncapsulation{}
	defer server.Clean() // clean up even in case of panic
	server.Init(kemName, nil)

	clientPublicKey := make([]byte, server.Details().LengthPublicKey)
	n, err := io.ReadFull(conn, clientPublicKey)

	if err != nil {
		panic(err)
	} else if n != server.Details().LengthPublicKey {
		panic(errors.New("server expected to read " + string(server.Details().
			LengthPublicKey) + " bytes, but instead read " + string(n)))
	}

	ciphertext, sharedSecretServer := server.EncapSecret(clientPublicKey)

	fmt.Printf("\nServer shared secret:\n% X ... % X\n",
		sharedSecretServer[0:8], sharedSecretServer[len(sharedSecretServer)-8:])

	// then send ciphertext to client and close the connection
	n, err = conn.Write(ciphertext)
	{
		if err != nil {
			panic(err)
		} else if n != server.Details().LengthCiphertext {
			panic(errors.New("server expected to write " + string(server.
				Details().LengthCiphertext) + " bytes, but instead wrote " + string(n)))
		}
	}
	conn.Close()
}
