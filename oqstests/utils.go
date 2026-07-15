package oqstests

import (
	"strings"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

// extMuLen is the fixed external-mu (mu) input length, in bytes, required by the
// ML-DSA *-extmu variants.
const extMuLen = 64

// sigMessage returns the message to sign/verify for sigName. The ML-DSA *-extmu
// variants require a fixed 64-byte external-mu input; all other sigs use msg.
func sigMessage(sigName string, msg []byte) []byte {
	if strings.HasSuffix(sigName, "-extmu") {
		return oqs.RandomBytes(extMuLen)
	}
	return msg
}

// stringMatchSlice returns true if str contains as a substring some element of
// the string slice s, and false otherwise. For example, the function returns
// true for the case str = "test" and s = ["testing", "element"], and false for
// the case str = "test" and s = ["happy", "dog"].
func stringMatchSlice(str string, s []string) bool {
	for _, pattern := range s {
		if strings.Contains(str, pattern) {
			return true
		}
	}
	return false
}
