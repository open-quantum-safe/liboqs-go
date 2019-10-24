// key encapsulation Go example
package main

import (
	"bytes"
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func main() {
	var rnd [48]byte = [48]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
		15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
		33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48}
	oqs.RandomBytesNistKatInit(rnd, []byte("avlada"), 256)
	r := oqs.RandomBytesSwitchAlgorithm("NIST-KAT")
	if r == oqs.ERROR {
		fmt.Println("Error")
	} else {
		fmt.Println(r)
	}

	kemName := "DEFAULT"
	client := oqs.KeyEncapsulation{}
	defer client.Clean() // clean up even in case of panic

	client.Init(kemName, []byte{})
	clientPublicKey := client.GenerateKeypair()
	fmt.Println("\nKEM details:")
	fmt.Println(client.GetDetails())

	server := oqs.KeyEncapsulation{}
	defer server.Clean() // clean up even in case of panic

	server.Init(kemName, []byte{})
	ciphertext, sharedSecretServer := server.EncapSecret(clientPublicKey)
	sharedSecretClient := client.DecapSecret(ciphertext)

	fmt.Printf("\nClient shared secret:\n% X ... % X\n",
		sharedSecretClient[0:8], sharedSecretClient[len(sharedSecretClient)-8:])
	fmt.Printf("\nServer shared secret:\n% X ... % X\n",
		sharedSecretServer[0:8], sharedSecretServer[len(sharedSecretServer)-8:])

	isValid := bytes.Equal(sharedSecretClient, sharedSecretServer)
	fmt.Println("\nShared secrets coincide? ", isValid)

	fmt.Println(oqs.RandomBytes(10))
}
