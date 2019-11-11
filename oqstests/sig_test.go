package oqstests

import (
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"sync"
	"testing"
)

// wgSig groups goroutines and blocks the caller until all goroutines finish.
var wgSig sync.WaitGroup

// testSig tests a specific signature.
func testSig(sigName string, msg []byte, t *testing.T) {
	defer wgSig.Done()
	var signer, verifier oqs.Signature
	defer signer.Clean()
	defer verifier.Clean()
	signer.Init(sigName, nil)
	verifier.Init(sigName, nil)
	pubKey := signer.GenerateKeypair()
	signature := signer.Sign(msg)
	isValid := verifier.Verify(msg, signature, pubKey)
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

// TestUnsupportedSignature tests that an unsupported signature emits a panic.
func TestUnsupportedSignature(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Unsupported signature should have emitted a panic")
		}
	}()
	signer := oqs.Signature{}
	defer signer.Clean()
	signer.Init("unsupported_sig", nil)
}
