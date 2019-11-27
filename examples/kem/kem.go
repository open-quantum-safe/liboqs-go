// key encapsulation Go example
package main

import (
	"bytes"
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func main() {
	fmt.Println("Supported KEMs:")
	fmt.Println(oqs.SupportedKEMs())

	fmt.Println("\nEnabled KEMs:")
	fmt.Println(oqs.EnabledKEMs())

	kemName := "DEFAULT"
	client := oqs.KeyEncapsulation{}
	defer client.Clean() // clean up even in case of panic

	client.Init(kemName, nil)
	clientPublicKey := client.GenerateKeyPair()
	fmt.Println("\nKEM details:")
	fmt.Println(client.Details())

	server := oqs.KeyEncapsulation{}
	defer server.Clean() // clean up even in case of panic

	server.Init(kemName, nil)
	ciphertext, sharedSecretServer := server.EncapSecret(clientPublicKey)
	sharedSecretClient := client.DecapSecret(ciphertext)

	fmt.Printf("\nClient shared secret:\n% X ... % X\n",
		sharedSecretClient[0:8], sharedSecretClient[len(sharedSecretClient)-8:])
	fmt.Printf("\nServer shared secret:\n% X ... % X\n",
		sharedSecretServer[0:8], sharedSecretServer[len(sharedSecretServer)-8:])

	isValid := bytes.Equal(sharedSecretClient, sharedSecretServer)
	fmt.Println("\nShared secrets coincide? ", isValid)
}
