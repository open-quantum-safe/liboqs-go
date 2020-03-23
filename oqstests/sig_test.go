package oqstests

import (
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// noThreadSigPatterns lists sigs that have issues running in a separate thread
var noThreadSigPatterns = []string{"Rainbow-IIIc", "Rainbow-Vc"}

// wgSig groups goroutines and blocks the caller until all goroutines finish.
var wgSig sync.WaitGroup

// testSig tests a specific signature.
func testSig(sigName string, msg []byte, threading bool) bool {
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
	return isValid
}

// TestSignature tests all enabled signatures.
func TestSignature(t *testing.T) {
	msg := []byte("This is our favourite message to sign")
	// first test sigs that belong to noThreadSigPatterns[] in the main
	// goroutine (stack size is 8Mb on macOS), due to issues with stack size
	// being too small in macOS (512Kb for threads)
	cnt := 0
	for _, sigName := range oqs.EnabledSigs() {
		for _, noThreadSig := range noThreadSigPatterns {
			if strings.Contains(sigName, noThreadSig) {
				cnt++
				if testSig(sigName, msg, false) == false {
					t.Fatal(sigName + ": signature verification failed")
				}
				break
			}
		}
	}
	// test the remaining sigs in separate goroutines
	wgSig.Add(len(oqs.EnabledSigs()) - cnt)
	for _, sigName := range oqs.EnabledSigs() {
		// issues with stack size being too small in macOS
		supportsThreads := true
		for _, noThreadSig := range noThreadSigPatterns {
			if strings.Contains(sigName, noThreadSig) {
				supportsThreads = false
				break
			}
		}
		if supportsThreads == true {
			log.Println(sigName)
			go testSig(sigName, msg, true)
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
