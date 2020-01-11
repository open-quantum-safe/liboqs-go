package oqstests

import (
	"fmt"
	"sync"
	"testing"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// wgSig groups goroutines and blocks the caller until all goroutines finish.
var wgSig sync.WaitGroup

// testSig tests a specific signature.
func testSig(sigName string, msg []byte, t *testing.T) {
	defer wgSig.Done()
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
		t.Fatal(sigName + ": signature verification failed")
	}
}

// TestSignature tests all enabled signatures.
func TestSignature(t *testing.T) {
	wgSig.Add(len(oqs.EnabledKEMs()))
	msg := []byte("This is our favourite message to sign")
	for _, sigName := range oqs.EnabledSigs() {
		fmt.Println(sigName)
		go testSig(sigName, msg, t)
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
