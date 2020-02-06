// Package oqstests provides unit testing for the oqs Go package.
package oqstests

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// wgKEM groups goroutines and blocks the caller until all goroutines finish.
var wgKEM sync.WaitGroup

// testKEM tests a specific KEM.
func testKEM(kemName string, t *testing.T) {
	defer wgKEM.Done()
	var client, server oqs.KeyEncapsulation
	defer client.Clean()
	defer server.Clean()
	// ignore potential errors everywhere
	_ = client.Init(kemName, nil)
	_ = server.Init(kemName, nil)
	clientPublicKey, _ := client.GenerateKeyPair()
	ciphertext, sharedSecretServer, _ := server.EncapSecret(clientPublicKey)
	sharedSecretClient, _ := client.DecapSecret(ciphertext)
	if !bytes.Equal(sharedSecretClient, sharedSecretServer) {
		t.Fatal(kemName + ": shared secrets do not coincide")
	}
}

// TestKeyEncapsulation tests all enabled KEMs.
func TestKeyEncapsulation(t *testing.T) {
	wgKEM.Add(len(oqs.EnabledKEMs()))
	// test LEDAcryptKEM-LT52 first in the main goroutine
	// (stack size too small otherwise)
	testKEM("LEDAcryptKEM-LT52", t)
	// test the rest of the KEMs
	for _, kemName := range oqs.EnabledKEMs() {
		fmt.Println(kemName)
		// issues with stack size being too small in macOS
		if kemName != "LEDAcryptKEM-LT52" {
			go testKEM(kemName, t)
		}
	}
	wgKEM.Wait()
}

// TestUnsupportedKeyEncapsulation tests that an unsupported KEM emits an error.
func TestUnsupportedKeyEncapsulation(t *testing.T) {
	client := oqs.KeyEncapsulation{}
	defer client.Clean()
	if err := client.Init("unsupported_kem", nil); err == nil {
		t.Fatal("Unsupported KEM should have emitted an error")
	}
}
