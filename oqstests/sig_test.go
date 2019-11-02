package oqstests

import (
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"log"
	"testing"
)

// testSig tests a specific signature.
func testSig(sigName string, msg []byte, t *testing.T) {
	log.Println(sigName)
	var signer, verifier oqs.Signature
	defer signer.Clean()
	defer verifier.Clean()
	signer.Init(sigName, []byte{})
	verifier.Init(sigName, []byte{})
	pubKey := signer.GenerateKeypair()
	signature := signer.Sign(msg)
	isValid := verifier.Verify(msg, signature, pubKey)
	if !isValid {
		t.Fatal("Signature verification failed")
	}
}

// TestSignature tests all enabled signatures.
func TestSignature(t *testing.T) {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	msg := []byte("This is our favourite message to sign")
	for _, sigName := range oqs.GetEnabledSIGs() {
		testSig(sigName, msg, t)
	}
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
	signer.Init("unsupported_sig", []byte{})
}
