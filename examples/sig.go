// signature Go example
package main

import (
    "fmt"
    "oqs"
)

func main() {
    fmt.Println("Supported signatures:")
    fmt.Println(oqs.GetSupportedSIGs())

    fmt.Println("\nEnabled signatures:")
    fmt.Println(oqs.GetEnabledSIGs())

    sigName := "DEFAULT"
    signer := oqs.Signature{}
    signer.Init(sigName, oqs.Bytes{})
    fmt.Printf("\nSignature details:\n%#v\n", signer.GetDetails())

    msg := oqs.Bytes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    pubKey := signer.GenerateKeypair()
    fmt.Printf("\nSigner public key:\n% X ... % X\n", pubKey[0:8],
        pubKey[len(pubKey)-8:])

    signature := signer.Sign(msg)
    fmt.Printf("\nSignature:\n% X ... % X\n", signature[0:8],
        signature[len(signature)-8:])

    verifier := oqs.Signature{}
    verifier.Init(sigName, oqs.Bytes{})
    isValid := verifier.Verify(msg, signature, pubKey)

    fmt.Println("\nValid signature? ", isValid)
}
