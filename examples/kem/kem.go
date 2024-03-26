// Key encapsulation Go example
package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func main() {
	fmt.Println("liboqs version: " + oqs.LiboqsVersion())
	fmt.Println("Enabled KEMs:")
	fmt.Println(oqs.EnabledKEMs())

	kemName := "Kyber512"
	client := oqs.KeyEncapsulation{}
	defer client.Clean() // clean up even in case of panic

	if err := client.Init(kemName, nil); err != nil {
		log.Fatal(err)
	}

	clientPublicKey, err := client.GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nKEM details:")
	fmt.Println(client.Details())

	server := oqs.KeyEncapsulation{}
	defer server.Clean() // clean up even in case of panic

	if err := server.Init(kemName, nil); err != nil {
		log.Fatal(err)
	}

	ciphertext, sharedSecretServer, err := server.EncapSecret(clientPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	sharedSecretClient, err := client.DecapSecret(ciphertext)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nClient shared secret:\n% X ... % X\n",
		sharedSecretClient[0:8], sharedSecretClient[len(sharedSecretClient)-8:])
	fmt.Printf("\nServer shared secret:\n% X ... % X\n",
		sharedSecretServer[0:8], sharedSecretServer[len(sharedSecretServer)-8:])

	isValid := bytes.Equal(sharedSecretClient, sharedSecretServer)
	fmt.Println("\nShared secrets coincide?", isValid)
}
