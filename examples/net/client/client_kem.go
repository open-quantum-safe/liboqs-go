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
		fmt.Println("specify both address and port number")
		return
	}
	address := os.Args[1]
	port := os.Args[2]

	fmt.Println("Launching KEM client on", address+":"+port)

	// connect to this socket
	conn, err := net.Dial("tcp", address+":"+port)
	if err != nil {
		panic(errors.New("client cannot connect to " + address + ":" + port))
	}

	client := oqs.KeyEncapsulation{}
	defer client.Clean() // clean up even in case of panic

	// receive the KEM name from the server
	kemName, _ := bufio.NewReader(conn).ReadString('\n')
	kemName = kemName[:len(kemName)-1]
	client.Init(kemName, nil)
	clientPublicKey := client.GenerateKeypair()

	// send to socket
	_, err = conn.Write([]byte(clientPublicKey))
	if err != nil {
		panic(errors.New("client cannot send the public key to the server"))
	}

	// listen for reply
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

	conn.Close()
}
