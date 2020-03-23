// Package oqstests provides unit testing for the oqs Go package.
package oqstests

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// noThreadKEMPatterns lists KEMs that have issues running in a separate thread
var noThreadKEMPatterns = []string{"Classic-McEliece", "LEDAcryptKEM-LT52"}

// wgKEM groups goroutines and blocks the caller until all goroutines finish.
var wgKEM sync.WaitGroup

// testKEM tests a specific KEM.
func testKEM(kemName string, threading bool) bool {
	log.Println(kemName) // thread-safe
	if threading == true {
		defer wgKEM.Done()
	}
	var client, server oqs.KeyEncapsulation
	defer client.Clean()
	defer server.Clean()
	// ignore potential errors everywhere
	_ = client.Init(kemName, nil)
	_ = server.Init(kemName, nil)
	clientPublicKey, _ := client.GenerateKeyPair()
	ciphertext, sharedSecretServer, _ := server.EncapSecret(clientPublicKey)
	sharedSecretClient, _ := client.DecapSecret(ciphertext)
	return bytes.Equal(sharedSecretClient, sharedSecretServer)
}

// TestKeyEncapsulation tests all enabled KEMs.
func TestKeyEncapsulation(t *testing.T) {
	// first test KEMs that belong to noThreadKEMPatterns[] in the main
	// goroutine (stack size is 8Mb on macOS), due to issues with stack size
	// being too small in macOS (512Kb for threads)
	cnt := 0
	for _, kemName := range oqs.EnabledKEMs() {
		for _, noThreadKem := range noThreadKEMPatterns {
			if strings.Contains(kemName, noThreadKem) {
				cnt++
				testKEM(kemName, false)
				break
			}
		}
	}
	// test the remaining KEMs in separate goroutines
	wgKEM.Add(len(oqs.EnabledKEMs()) - cnt)
	for _, kemName := range oqs.EnabledKEMs() {
		// issues with stack size being too small in macOS
		supportsThreads := true
		for _, noThreadKem := range noThreadKEMPatterns {
			if strings.Contains(kemName, noThreadKem) {
				supportsThreads = false
				break
			}
		}
		if supportsThreads == true {
			go testKEM(kemName, true)
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
