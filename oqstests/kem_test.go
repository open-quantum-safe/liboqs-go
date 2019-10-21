// Package oqstests provides unit testing for the oqs Go package
package oqstests // import "github.com/open-quantum-safe/liboqs-go/oqstests"

import (
    "bytes"
    "fmt"
    "github.com/open-quantum-safe/liboqs-go/oqs"
    "testing"
)

// Test all enabled KEMs
func TestKeyEncapsulation(t *testing.T) {
    var client, server oqs.KeyEncapsulation

    for _, kemName := range oqs.GetEnabledKEMs() {
        fmt.Println(kemName)
        client.Init(kemName, []byte{})
        server.Init(kemName, []byte{})
        clientPublicKey := client.GenerateKeypair()
        ciphertext, sharedSecretServer := server.EncapSecret(clientPublicKey)
        sharedSecretClient := client.DecapSecret(ciphertext)
        isValid := bytes.Compare(sharedSecretClient, sharedSecretServer) == 0
        if !isValid {
            t.Fatal("Shared secrets do not coincide")
        }
        client.Clean()
        server.Clean()
    }
}

// Test that an unsupported KEM emits a panic
func TestUnsupportedKeyEncapsulation(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("Unsupported KEM should have generated a panic")
        }
    }()
    client := oqs.KeyEncapsulation{}
    client.Init("unsupported_kem", []byte{})
}
