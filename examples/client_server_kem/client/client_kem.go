// key encapsulation TCP client Go example
// run with "go run client_kem.go <host address> <port number>"
package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"io"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: client_kem <address> <port number>")
		os.Exit(-1)
	}
	address := os.Args[1]
	port := os.Args[2]

	fmt.Println("Launching KEM client on", address+":"+port)
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		panic(errors.New("client cannot connect to " + address + ":" + port))
	}
	defer conn.Close() // clean up even in case of panic

	// construct the KEM client
	client := oqs.KeyEncapsulation{}
	defer client.Clean() // clean up even in case of panic

	// receive the KEM name from the server
	kemName, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(errors.New("client cannot receive the KEM name from the server"))
	}
	kemName = kemName[:len(kemName)-1] // remove the '\n'

	// initialize the KEM client and generate the key pairs
	client.Init(kemName, nil)
	clientPublicKey := client.GenerateKeypair()

	// send the client public key to the server
	_, err = conn.Write(clientPublicKey)
	if err != nil {
		panic(errors.New("client cannot send the public key to the server"))
	}

	// listen for reply from the server, e.g. for the encapsulated secret
	ciphertext := make([]byte, client.Details().LengthCiphertext)
	n, err := io.ReadFull(conn, ciphertext)
	if err != nil {
		panic(err)
	} else if n != client.Details().LengthCiphertext {
		panic(errors.New("client expected to read " + string(client.Details().
			LengthCiphertext) + " bytes, but instead read " + string(n)))
	}

	// decapsulate the secret and extract the shared secret
	sharedSecretClient := client.DecapSecret(ciphertext)

	fmt.Println(client.Details())
	fmt.Printf("\nClient shared secret:\n% X ... % X\n",
		sharedSecretClient[0:8], sharedSecretClient[len(sharedSecretClient)-8:])
}