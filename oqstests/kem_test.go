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
var noThreadKEMPatterns = []string{}

// wgKEMCorrectness groups goroutines and blocks the caller until all goroutines finish.
var wgKEMCorrectness sync.WaitGroup

// wgKEMWrongCiphertext groups goroutines and blocks the caller until all goroutines finish.
var wgKEMWrongCiphertext sync.WaitGroup

// testKEMCorrectness tests the correctness of a specific KEM.
func testKEMCorrectness(kemName string, threading bool, t *testing.T) {
	log.Println("Correctness - ", kemName) // thread-safe
	if threading == true {
		defer wgKEMCorrectness.Done()
	}
	var client, server oqs.KeyEncapsulation
	defer client.Clean()
	defer server.Clean()
	// Ignore potential errors everywhere
	_ = client.Init(kemName, nil)
	_ = server.Init(kemName, nil)
	clientPublicKey, _ := client.GenerateKeyPair()
	ciphertext, sharedSecretServer, _ := server.EncapSecret(clientPublicKey)
	sharedSecretClient, _ := client.DecapSecret(ciphertext)
	if !bytes.Equal(sharedSecretClient, sharedSecretServer) {
		// t.Errorf is thread-safe
		t.Errorf("%s: shared secrets do not coincide", kemName)
	}
}

// testKEMWrongCiphertext tests the wrong ciphertext regime of a specific KEM.
func testKEMWrongCiphertext(kemName string, threading bool, t *testing.T) {
	log.Println("Wrong ciphertext - ", kemName) // thread-safe
	if threading == true {
		defer wgKEMWrongCiphertext.Done()
	}
	var client, server oqs.KeyEncapsulation
	defer client.Clean()
	defer server.Clean()
	// Ignore potential errors everywhere
	_ = client.Init(kemName, nil)
	_ = server.Init(kemName, nil)
	clientPublicKey, _ := client.GenerateKeyPair()
	ciphertext, sharedSecretServer, _ := server.EncapSecret(clientPublicKey)
	wrongCiphertext := oqs.RandomBytes(len(ciphertext))
	sharedSecretClient, _ := client.DecapSecret(wrongCiphertext)
	if bytes.Equal(sharedSecretClient, sharedSecretServer) {
		// t.Errorf is thread-safe
		t.Errorf("%s: shared secrets should not coincide", kemName)
	}
}

// TestKeyEncapsulationCorrectness tests the correctness of all enabled KEMs.
func TestKeyEncapsulationCorrectness(t *testing.T) {
	// Disable some KEMs in macOS/OSX
	if runtime.GOOS == "darwin" {
		disabledKEMPatterns = []string{}
	}
	// Disable some KEMs in Windows
	if runtime.GOOS == "windows" {
		disabledKEMPatterns = []string{}
	}
	// First test KEMs that belong to noThreadKEMPatterns[] in the main
	// goroutine, due to issues with stack size being too small in macOS or
	// Windows
	cnt := 0
	for _, kemName := range oqs.EnabledKEMs() {
		if stringMatchSlice(kemName, disabledKEMPatterns) {
			cnt++
			continue
		}
		// issues with stack size being too small
		if stringMatchSlice(kemName, noThreadKEMPatterns) {
			cnt++
			testKEMCorrectness(kemName, false, t)
		}
	}
	// Test the remaining KEMs in separate goroutines
	wgKEMCorrectness.Add(len(oqs.EnabledKEMs()) - cnt)
	for _, kemName := range oqs.EnabledKEMs() {
		if stringMatchSlice(kemName, disabledKEMPatterns) {
			continue
		}
		if !stringMatchSlice(kemName, noThreadKEMPatterns) {
			go testKEMCorrectness(kemName, true, t)
		}
	}
	wgKEMCorrectness.Wait()
}

// TestKeyEncapsulationWrongCiphertext tests the wrong ciphertext regime of all enabled KEMs.
func TestKeyEncapsulationWrongCiphertext(t *testing.T) {
	// disable some KEMs in macOS/OSX
	if runtime.GOOS == "darwin" {
		disabledKEMPatterns = []string{}
	}
	// Disable some KEMs in Windows
	if runtime.GOOS == "windows" {
		disabledKEMPatterns = []string{}
	}
	// First test KEMs that belong to noThreadKEMPatterns[] in the main
	// goroutine, due to issues with stack size being too small in macOS or
	// Windows
	cnt := 0
	for _, kemName := range oqs.EnabledKEMs() {
		if stringMatchSlice(kemName, disabledKEMPatterns) {
			cnt++
			continue
		}
		// Issues with stack size being too small
		if stringMatchSlice(kemName, noThreadKEMPatterns) {
			cnt++
			testKEMWrongCiphertext(kemName, false, t)
		}
	}
	// Test the remaining KEMs in separate goroutines
	wgKEMWrongCiphertext.Add(len(oqs.EnabledKEMs()) - cnt)
	for _, kemName := range oqs.EnabledKEMs() {
		if stringMatchSlice(kemName, disabledKEMPatterns) {
			continue
		}
		if !stringMatchSlice(kemName, noThreadKEMPatterns) {
			go testKEMWrongCiphertext(kemName, true, t)
		}
	}
	wgKEMWrongCiphertext.Wait()
}

// TestUnsupportedKeyEncapsulation tests that an unsupported KEM emits an error.
func TestUnsupportedKeyEncapsulation(t *testing.T) {
	client := oqs.KeyEncapsulation{}
	defer client.Clean()
	if err := client.Init("unsupported_kem", nil); err == nil {
		t.Errorf("Unsupported KEM should have emitted an error")
	}
}
