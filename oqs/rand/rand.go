// Package rand provides support for various RNG
package rand // import "github.com/open-quantum-safe/liboqs-go/oqs/rand"

/*
#cgo pkg-config: liboqs
#include <oqs/oqs.h>
*/
import "C"
import (
	"errors"
)

/**************** Randomness ****************/

// RandomBytes generates bytesToRead random bytes. This implementation uses
// whichever algorithm has been selected by RandomBytesSwitchAlgorithm, the
// default being "system".
func RandomBytes(bytesToRead int) []byte {
	result := make([]byte, bytesToRead)
	C.OQS_randombytes((*C.uint8_t)(&result[0]), C.size_t(bytesToRead))
	return result
}

// RandomBytesSwitchAlgorithm switches the core OQS_randombytes to use the
// specified algorithm. Possible values are "system", "NIST-KAT", "OpenSSL".
// See <oqs/rand.h> C header for more details.
func RandomBytesSwitchAlgorithm(algName string) {
	if C.OQS_randombytes_switch_algorithm(C.CString(algName)) != C.OQS_SUCCESS {
		panic(errors.New("can not switch algorithm"))
	}
}

// RandomBytesNistKatInit initializes the NIST DRBG with the entropyInput seed,
// which must by 48 exactly bytes long. The personalizationString is an optional
// personalization string.
func RandomBytesNistKatInit(entropyInput [48]byte,
	personalizationString []byte) {
	if len(personalizationString) > 0 {
		C.OQS_randombytes_nist_kat_init((*C.uint8_t)(&entropyInput[0]),
			(*C.uint8_t)(&personalizationString[0]), 256)
		return
	}
	C.OQS_randombytes_nist_kat_init((*C.uint8_t)(&entropyInput[0]),
		(*C.uint8_t)(nil), 256)
}

// RandomBytesCustomAlgorithm switches RandomBytes to use the given function.
// This allows additional custom RNGs besides the provided ones. The provided
// RNG function must have the same signature as RandomBytes.
//func RandomBytesCustomAlgorithm(rngFunc func(int) []byte) {
//	C.OQS_randombytes_custom_algorithm(rngFunc)
//}

/**************** END Randomness ****************/
