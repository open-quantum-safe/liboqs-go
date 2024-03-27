// Key encapsulation TCP server Go example
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// Counter is a thread-safe counter.
type Counter struct {
	mu  sync.Mutex
	cnt uint64
}

// Add increments the counter.
func (c *Counter) Add() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cnt++
}

// Val retrieves the counter's value.
func (c *Counter) Val() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	cnt := c.cnt
	return cnt
}

// counter is a thread-safe connection counter.
var counter Counter

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: server_kem <port number> [KEM name (optional)]")
		os.Exit(1)
	}
	port := os.Args[1]
	kemName := "Kyber512"
	if len(os.Args) > 2 {
		kemName = os.Args[2]
	}

	log.SetOutput(os.Stdout) // log to stdout instead the default stderr
	fmt.Println("Launching KEM", kemName, "server on port", port)
	{
		kem := oqs.KeyEncapsulation{}
		if err := kem.Init(kemName, nil); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n\n", kem.Details())
		kem.Clean()
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	// Listen indefinitely (until explicitly stopped, e.g. with CTRL+C in UNIX)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle connections concurrently
		go handleConnection(conn, kemName)
	}
}

func handleConnection(conn net.Conn, kemName string) {
	defer conn.Close() // clean up even in case of panic

	// Send KEM name to client first
	_, err := fmt.Fprintln(conn, kemName)
	if err != nil {
		log.Fatal(errors.New("server cannot send the KEM name to the client"))
	}

	// Construct and initialize the KEM server
	server := oqs.KeyEncapsulation{}
	defer server.Clean() // clean up even in case of panic
	if err := server.Init(kemName, nil); err != nil {
		log.Fatal(err)
	}

	// Read the public key sent by the client
	clientPublicKey := make([]byte, server.Details().LengthPublicKey)
	n, err := io.ReadFull(conn, clientPublicKey)
	if err != nil {
		log.Fatal(err)
	} else if n != server.Details().LengthPublicKey {
		log.Fatal(errors.New("server expected to read " +
			fmt.Sprintf("%v", server.Details().LengthPublicKey) + " bytes, but instead " +
			"read " + fmt.Sprintf("%v", n)))
	}

	// Encapsulate the secret
	ciphertext, sharedSecretServer, err := server.EncapSecret(clientPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	// Then send ciphertext to client and close the connection
	n, err = conn.Write(ciphertext)
	if err != nil {
		log.Fatal(err)
	} else if n != server.Details().LengthCiphertext {
		log.Fatal(errors.New("server expected to write " + fmt.Sprintf("%v", server.
			Details().LengthCiphertext) + " bytes, but instead wrote " + fmt.Sprintf("%v", n)))
	}

	log.Printf("\nConnection #%d - server shared secret:\n% X ... % X\n\n",
		counter.Val(), sharedSecretServer[0:8],
		sharedSecretServer[len(sharedSecretServer)-8:])

	// Increment the connection number
	counter.Add()
}
