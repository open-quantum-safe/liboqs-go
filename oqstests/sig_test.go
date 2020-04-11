package oqstests

import (
	"log"
	"runtime"
	"sync"
	"testing"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// disabledSigPatterns lists sigs for which unit testing is disabled
var disabledSigPatterns []string

// noThreadSigPatterns lists sigs that have issues running in a separate thread
var noThreadSigPatterns []string

// wgSig groups goroutines and blocks the caller until all goroutines finish.
var wgSig sync.WaitGroup

// testSig tests a specific signature.
func testSig(sigName string, msg []byte, threading bool, t *testing.T) {
	log.Println(sigName) // thread-safe
	if threading == true {
		defer wgSig.Done()
	}
	var signer, verifier oqs.Signature
	defer signer.Clean()
	defer verifier.Clean()
	// ignore potential errors everywhere
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

// TestSignature tests all enabled signatures.
func TestSignature(t *testing.T) {
	// disable some sigs in macOS/OSX
	if runtime.GOOS == "darwin" {
		disabledSigPatterns = []string{"Rainbow-IIIc", "Rainbow-Vc"}
	}
	// disable some sigs in Windows
	if runtime.GOOS == "windows" {
		disabledSigPatterns = []string{"Rainbow-IIIc", "Rainbow-Vc"}
	}
	msg := []byte("This is our favourite message to sign")
	// first test sigs that belong to noThreadSigPatterns[] in the main
	// goroutine, due to issues with stack size being too small in macOS or
	// Windows
	cnt := 0
	for _, sigName := range oqs.EnabledSigs() {
		if stringMatchSlice(sigName, disabledSigPatterns) {
			cnt++
			continue
		}
		// issues with stack size being too small
		if stringMatchSlice(sigName, noThreadSigPatterns) {
			cnt++
			testSig(sigName, msg, false, t)
		}
	}
	// test the remaining sigs in separate goroutines
	wgSig.Add(len(oqs.EnabledSigs()) - cnt)
	for _, sigName := range oqs.EnabledSigs() {
		if stringMatchSlice(sigName, disabledSigPatterns) {
			continue
		}
		if !stringMatchSlice(sigName, noThreadSigPatterns) {
			go testSig(sigName, msg, true, t)
		}
	}
	wgSig.Wait()
}

// TestUnsupportedSignature tests that an unsupported signature emits an error.
func TestUnsupportedSignature(t *testing.T) {
	signer := oqs.Signature{}
	defer signer.Clean()
	if err := signer.Init("unsupported_sig", nil); err == nil {
		t.Fatal("Unsupported signature should have emitted an error")
	}
}
