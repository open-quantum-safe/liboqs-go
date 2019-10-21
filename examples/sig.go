// signature Go example
package main

import (
    "fmt"
    "github.com/open-quantum-safe/liboqs-go/oqs"
)

func main() {
    fmt.Println("Supported signatures:")
    fmt.Println(oqs.GetSupportedSIGs())

    fmt.Println("\nEnabled signatures:")
    fmt.Println(oqs.GetEnabledSIGs())

    sigName := "DEFAULT"
    signer := oqs.Signature{}
    defer signer.Clean() // clean up even in case of panic

    signer.Init(sigName, []byte{})
    fmt.Printf("\nSignature details:\n%#v\n", signer.GetDetails())

    msg := []byte("This is the message to sign")
    pubKey := signer.GenerateKeypair()
    fmt.Printf("\nSigner public key:\n% X ... % X\n", pubKey[0:8],
        pubKey[len(pubKey)-8:])

    signature := signer.Sign(msg)
    fmt.Printf("\nSignature:\n% X ... % X\n", signature[0:8],
        signature[len(signature)-8:])

    verifier := oqs.Signature{}
    defer verifier.Clean() // clean up even in case of panic

    verifier.Init(sigName, []byte{})
    isValid := verifier.Verify(msg, signature, pubKey)

    fmt.Println("\nValid signature? ", isValid)
}
