// Key encapsulation TCP client Go example
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: client_kem <address> <port number>")
		os.Exit(1)
	}
	address := os.Args[1]
	port := os.Args[2]

	fmt.Println("Launching KEM client on", address+":"+port)
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		log.Fatal(errors.New("client cannot connect " +
			"to " + address + ":" + port))
	}
	defer conn.Close() // clean up even in case of panic

	// Construct the KEM client
	client := oqs.KeyEncapsulation{}
	defer client.Clean() // clean up even in case of panic

	// Receive the KEM name from the server
	kemName, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(errors.New("client cannot receive the " +
			"KEM name from the server"))
	}
	kemName = kemName[:len(kemName)-1] // remove the '\n'

	// Initialize the KEM client and generate the key pairs
	if err := client.Init(kemName, nil); err != nil {
		log.Fatal(err)
	}
	clientPublicKey, err := client.GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	// Send the client public key to the server
	_, err = conn.Write(clientPublicKey)
	if err != nil {
		log.Fatal(errors.New("client cannot send the public key to the " +
			"server"))
	}

	// Listen for reply from the server, e.g. for the encapsulated secret
	ciphertext := make([]byte, client.Details().LengthCiphertext)
	n, err := io.ReadFull(conn, ciphertext)
	if err != nil {
		log.Fatal(err)
	} else if n != client.Details().LengthCiphertext {
		log.Fatal(errors.New("client expected to read " +
			fmt.Sprintf("%v", client.Details().LengthCiphertext) + " bytes, but instead " +
			"read " + fmt.Sprintf("%v", n)))
	}

	// Decapsulate the secret and extract the shared secret
	sharedSecretClient, err := client.DecapSecret(ciphertext)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(client.Details())
	fmt.Printf("\nClient shared secret:\n% X ... % X\n",
		sharedSecretClient[0:8], sharedSecretClient[len(sharedSecretClient)-8:])
}
