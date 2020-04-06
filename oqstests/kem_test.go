// Package oqstests provides unit testing for the oqs Go package.
package oqstests

import (
	"bytes"
	"log"
	"runtime"
	"sync"
	"testing"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// disabledKEMPatterns lists KEMs for which unit testing is disabled
var disabledKEMPatterns []string

// noThreadKEMPatterns lists KEMs that have issues running in a separate thread
var noThreadKEMPatterns = []string{"LEDAcryptKEM-LT52"}

// wgKEM groups goroutines and blocks the caller until all goroutines finish.
var wgKEM sync.WaitGroup

// testKEM tests a specific KEM.
func testKEM(kemName string, threading bool, t *testing.T) {
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
	if !bytes.Equal(sharedSecretClient, sharedSecretServer) {
		// t.Errorf is thread-safe
		t.Errorf(kemName + ": shared secrets do not coincide")
	}
}

// TestKeyEncapsulation tests all enabled KEMs.
func TestKeyEncapsulation(t *testing.T) {
	// disable some KEMs in macOS/OSX
	if runtime.GOOS == "darwin" {
		disabledKEMPatterns = []string{"Classic-McEliece"}
	}
	// first test KEMs that belong to noThreadKEMPatterns[] in the main
	// goroutine (stack size is 8Mb on macOS), due to issues with stack size
	// being too small in macOS (512Kb for threads)
	cnt := 0
	for _, kemName := range oqs.EnabledKEMs() {
		if stringMatchSlice(kemName, disabledKEMPatterns) {
			cnt++
			continue
		}
		// issues with stack size being too small in macOS
		if stringMatchSlice(kemName, noThreadKEMPatterns) {
			cnt++
			testKEM(kemName, false, t)
		}
	}
	// test the remaining KEMs in separate goroutines
	wgKEM.Add(len(oqs.EnabledKEMs()) - cnt)
	for _, kemName := range oqs.EnabledKEMs() {
		if stringMatchSlice(kemName, disabledKEMPatterns) {
			continue
		}
		if !stringMatchSlice(kemName, noThreadKEMPatterns) {
			go testKEM(kemName, true, t)
		}
	}
	wgKEM.Wait()
}

// TestUnsupportedKeyEncapsulation tests that an unsupported KEM emits an error.
func TestUnsupportedKeyEncapsulation(t *testing.T) {
	client := oqs.KeyEncapsulation{}
	defer client.Clean()
	if err := client.Init("unsupported_kem", nil); err == nil {
		t.Errorf("Unsupported KEM should have emitted an error")
	}
}
