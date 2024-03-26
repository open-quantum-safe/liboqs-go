// Signature Go example
package main

import (
	"fmt"
	"log"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func main() {
	fmt.Println("liboqs version: " + oqs.LiboqsVersion())
	fmt.Println("Enabled signatures:")
	fmt.Println(oqs.EnabledSigs())

	sigName := "Dilithium2"
	signer := oqs.Signature{}
	defer signer.Clean() // clean up even in case of panic

	if err := signer.Init(sigName, nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSignature details:")
	fmt.Println(signer.Details())

	msg := []byte("This is the message to sign")
	pubKey, err := signer.GenerateKeyPair()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nSigner public key:\n% X ... % X\n", pubKey[0:8],
		pubKey[len(pubKey)-8:])

	signature, _ := signer.Sign(msg)
	fmt.Printf("\nSignature:\n% X ... % X\n", signature[0:8],
		signature[len(signature)-8:])

	verifier := oqs.Signature{}
	defer verifier.Clean() // clean up even in case of panic

	if err := verifier.Init(sigName, nil); err != nil {
		log.Fatal(err)
	}

	isValid, err := verifier.Verify(msg, signature, pubKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nValid signature?", isValid)
}
