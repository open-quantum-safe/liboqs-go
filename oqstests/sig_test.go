package oqstests

import (
	"log"
	"runtime"
	"sync"
	"testing"

	"github.com/open-quantum-safe/liboqs-go/oqs"
	"github.com/open-quantum-safe/liboqs-go/oqs/rand"
)

// disabledSigPatterns lists sigs for which unit testing is disabled
var disabledSigPatterns []string

// noThreadSigPatterns lists sigs that have issues running in a separate thread
var noThreadSigPatterns []string

// wgSigCorrectness groups goroutines and blocks the caller until all goroutines finish.
var wgSigCorrectness sync.WaitGroup

// wgSigWrongSignature groups goroutines and blocks the caller until all goroutines finish.
var wgSigWrongSignature sync.WaitGroup

// wgSigWrongPublicKey groups goroutines and blocks the caller until all goroutines finish.
var wgSigWrongPublicKey sync.WaitGroup

// testSigCorrectness tests a specific signature.
func testSigCorrectness(sigName string, msg []byte, threading bool, t *testing.T) {
	log.Println("Correctness - ", sigName) // thread-safe
	if threading == true {
		defer wgSigCorrectness.Done()
	}
	var signer, verifier oqs.Signature
	defer signer.Clean()
	defer verifier.Clean()
	// Ignore potential errors everywhere
	_ = signer.Init(sigName, nil)
	_ = verifier.Init(sigName, nil)
	pubKey, _ := signer.GenerateKeyPair()
	signature, _ := signer.Sign(msg)
	isValid, _ := verifier.Verify(msg, signature, pubKey)
	if !isValid {
		// t.Errorf is thread-safe
		t.Errorf(sigName + ": signature verification failed")
	}
}

// testSigWrongSignature tests the wrong signature regime of a specific signature.
func testSigWrongSignature(sigName string, msg []byte, threading bool, t *testing.T) {
	log.Println("Wrong signature - ", sigName) // thread-safe
	if threading == true {
		defer wgSigWrongSignature.Done()
	}
	var signer, verifier oqs.Signature
	defer signer.Clean()
	defer verifier.Clean()
	// Ignore potential errors everywhere
	_ = signer.Init(sigName, nil)
	_ = verifier.Init(sigName, nil)
	pubKey, _ := signer.GenerateKeyPair()
	signature, _ := signer.Sign(msg)
	wrongSignature := rand.RandomBytes(len(signature))
	isValid, _ := verifier.Verify(msg, wrongSignature, pubKey)
	if isValid {
		// t.Errorf is thread-safe
		t.Errorf(sigName + ": signature verification should have failed")
	}
}

// testSigWrongPublicKey tests the wrong public key regime of a specific signature.
func testSigWrongPublicKey(sigName string, msg []byte, threading bool, t *testing.T) {
	log.Println("Wrong public key - ", sigName) // thread-safe
	if threading == true {
		defer wgSigWrongPublicKey.Done()
	}
	var signer, verifier oqs.Signature
	defer signer.Clean()
	defer verifier.Clean()
	// Ignore potential errors everywhere
	_ = signer.Init(sigName, nil)
	_ = verifier.Init(sigName, nil)
	pubKey, _ := signer.GenerateKeyPair()
	wrongPubKey := rand.RandomBytes(len(pubKey))
	signature, _ := signer.Sign(msg)
	isValid, _ := verifier.Verify(msg, signature, wrongPubKey)
	if isValid {
		// t.Errorf is thread-safe
		t.Errorf(sigName + ": signature verification should have failed")
	}
}

// TestSignatureCorrectness tests all enabled signatures.
func TestSignatureCorrectness(t *testing.T) {
	// Disable some sigs in macOS/OSX
	if runtime.GOOS == "darwin" {
		disabledSigPatterns = []string{"Rainbow-III", "Rainbow-V"}
	}
	// Disable some sigs in Windows
	if runtime.GOOS == "windows" {
		disabledSigPatterns = []string{"Rainbow-V"}
	}
	msg := []byte("This is our favourite message to sign")
	// First test sigs that belong to noThreadSigPatterns[] in the main
	// goroutine, due to issues with stack size being too small in macOS or
	// Windows
	cnt := 0
	for _, sigName := range oqs.EnabledSigs() {
		if stringMatchSlice(sigName, disabledSigPatterns) {
			cnt++
			continue
		}
		// Issues with stack size being too small
		if stringMatchSlice(sigName, noThreadSigPatterns) {
			cnt++
			testSigCorrectness(sigName, msg, false, t)
		}
	}
	// Test the remaining sigs in separate goroutines
	wgSigCorrectness.Add(len(oqs.EnabledSigs()) - cnt)
	for _, sigName := range oqs.EnabledSigs() {
		if stringMatchSlice(sigName, disabledSigPatterns) {
			continue
		}
		if !stringMatchSlice(sigName, noThreadSigPatterns) {
			go testSigCorrectness(sigName, msg, true, t)
		}
	}
	wgSigCorrectness.Wait()
}

// TestSignatureWrongSignature tests the wrong signature regime of all enabled
// signatures.
func TestSignatureWrongSignature(t *testing.T) {
	// Disable some sigs in macOS/OSX
	if runtime.GOOS == "darwin" {
		disabledSigPatterns = []string{"Rainbow-III", "Rainbow-V"}
	}
	// Disable some sigs in Windows
	if runtime.GOOS == "windows" {
		disabledSigPatterns = []string{"Rainbow-V"}
	}
	msg := []byte("This is our favourite message to sign")
	// First test sigs that belong to noThreadSigPatterns[] in the main
	// goroutine, due to issues with stack size being too small in macOS or
	// Windows
	cnt := 0
	for _, sigName := range oqs.EnabledSigs() {
		if stringMatchSlice(sigName, disabledSigPatterns) {
			cnt++
			continue
		}
		// Issues with stack size being too small
		if stringMatchSlice(sigName, noThreadSigPatterns) {
			cnt++
			testSigWrongSignature(sigName, msg, false, t)
		}
	}
	// Test the remaining sigs in separate goroutines
	wgSigWrongSignature.Add(len(oqs.EnabledSigs()) - cnt)
	for _, sigName := range oqs.EnabledSigs() {
		if stringMatchSlice(sigName, disabledSigPatterns) {
			wgSigWrongSignature.Done()
			continue
		}
		if !stringMatchSlice(sigName, noThreadSigPatterns) {
			go testSigWrongSignature(sigName, msg, true, t)
		}
	}
	wgSigWrongSignature.Wait()
}

// TestSignatureWrongPublicKey tests the wrong public key regime of all
// enabled signatures.
func TestSignatureWrongPublicKey(t *testing.T) {
	// Disable some sigs in macOS/OSX
	if runtime.GOOS == "darwin" {
		disabledSigPatterns = []string{"Rainbow-III", "Rainbow-V"}
	}
	// Disable some sigs in Windows
	if runtime.GOOS == "windows" {
		disabledSigPatterns = []string{"Rainbow-V"}
	}
	msg := []byte("This is our favourite message to sign")
	// First test sigs that belong to noThreadSigPatterns[] in the main
	// goroutine, due to issues with stack size being too small in macOS or
	// Windows
	cnt := 0
	for _, sigName := range oqs.EnabledSigs() {
		if stringMatchSlice(sigName, disabledSigPatterns) {
			cnt++
			continue
		}
		// Issues with stack size being too small
		if stringMatchSlice(sigName, noThreadSigPatterns) {
			cnt++
			testSigWrongPublicKey(sigName, msg, false, t)
		}
	}
	// Test the remaining sigs in separate goroutines
	wgSigWrongPublicKey.Add(len(oqs.EnabledSigs()) - cnt)
	for _, sigName := range oqs.EnabledSigs() {
		if stringMatchSlice(sigName, disabledSigPatterns) {
			continue
		}
		if !stringMatchSlice(sigName, noThreadSigPatterns) {
			go testSigWrongPublicKey(sigName, msg, true, t)
		}
	}
	wgSigWrongPublicKey.Wait()
}

// TestUnsupportedSignature tests that an unsupported signature emits an error.
func TestUnsupportedSignature(t *testing.T) {
	signer := oqs.Signature{}
	defer signer.Clean()
	if err := signer.Init("unsupported_sig", nil); err == nil {
		t.Fatal("Unsupported signature should have emitted an error")
	}
}
