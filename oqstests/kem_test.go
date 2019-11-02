// Package oqstests provides unit testing for the oqs Go package.
package oqstests

import (
	"bytes"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"log"
	"testing"
)

// testKem tests a specific KEM.
func testKem(kemName string, t *testing.T) {
	log.Println(kemName)
	var client, server oqs.KeyEncapsulation
	defer client.Clean()
	defer server.Clean()
	client.Init(kemName, []byte{})
	server.Init(kemName, []byte{})
	clientPublicKey := client.GenerateKeypair()
	ciphertext, sharedSecretServer := server.EncapSecret(clientPublicKey)
	sharedSecretClient := client.DecapSecret(ciphertext)
	if !bytes.Equal(sharedSecretClient, sharedSecretServer) {
		t.Fatal("Shared secrets do not coincide")
	}
}

// TestKeyEncapsulation tests all enabled KEMs.
func TestKeyEncapsulation(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	for _, kemName := range oqs.GetEnabledKEMs() {
		testKem(kemName, t)
	}
}

// TestUnsupportedKeyEncapsulation tests that an unsupported KEM emits a panic.
func TestUnsupportedKeyEncapsulation(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Unsupported KEM should have emitted a panic")
		}
	}()
	client := oqs.KeyEncapsulation{}
	defer client.Clean()
	client.Init("unsupported_kem", []byte{})
}
