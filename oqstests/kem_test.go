// Package oqstests provides unit testing for the oqs Go package.
package oqstests

import (
	"bytes"
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"sync"
	"testing"
)

// wgKEM groups goroutines and blocks the caller until all goroutines finish.
var wgKEM sync.WaitGroup

// testKEM tests a specific KEM.
func testKEM(kemName string, t *testing.T) {
	defer wgKEM.Done()
	var client, server oqs.KeyEncapsulation
	defer client.Clean()
	defer server.Clean()
	client.Init(kemName, nil)
	server.Init(kemName, nil)
	clientPublicKey := client.GenerateKeypair()
	ciphertext, sharedSecretServer := server.EncapSecret(clientPublicKey)
	sharedSecretClient := client.DecapSecret(ciphertext)
	if !bytes.Equal(sharedSecretClient, sharedSecretServer) {
		t.Fatal(kemName + ": shared secrets do not coincide")
	}
}

// TestKeyEncapsulation tests all enabled KEMs.
func TestKeyEncapsulation(t *testing.T) {
	wgKEM.Add(len(oqs.EnabledKEMs()))
	for _, kemName := range oqs.EnabledKEMs() {
		fmt.Println(kemName)
		go testKEM(kemName, t)
	}
	wgKEM.Wait()
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
	client.Init("unsupported_kem", nil)
}
