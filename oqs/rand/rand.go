// Package rand provides support for various RNG-related functions.
package rand // import "github.com/open-quantum-safe/liboqs-go/oqs/rand"

/**************** Callbacks ****************/

/*
#cgo pkg-config: liboqs-go
#include <oqs/oqs.h>
typedef void (*algorithm_ptr)(uint8_t*, size_t);
void algorithmPtr_cgo(uint8_t*, size_t);
*/
import "C"

import (
	"errors"
	"unsafe"
)

// algorithmPtrCallback is a global RNG algorithm callback set by
// RandomBytesCustomAlgorithm.
var algorithmPtrCallback func([]byte, int)

// algorithmPtr is automatically invoked by RandomBytesCustomAlgorithm. When
// invoked, the memory is provided by the caller, i.e. RandomBytes or
// RandomBytesInPlace.
//
//export algorithmPtr
func algorithmPtr(randomArray *C.uint8_t, bytesToRead C.size_t) {
	// TODO optimize the copying if possible!
	result := make([]byte, int(bytesToRead))
	algorithmPtrCallback(result, int(bytesToRead))
	p := unsafe.Pointer(randomArray)
	for _, v := range result {
		*(*C.uint8_t)(p) = C.uint8_t(v)
		p = unsafe.Pointer(uintptr(p) + 1)
	}
}

/**************** END Callbacks ****************/

/**************** Randomness ****************/

// RandomBytes generates bytesToRead random bytes. This implementation uses
// either the default RNG algorithm ("system"), or whichever algorithm has been
// selected by RandomBytesSwitchAlgorithm.
func RandomBytes(bytesToRead int) []byte {
	result := make([]byte, bytesToRead)
	C.OQS_randombytes((*C.uint8_t)(unsafe.Pointer(&result[0])),
		C.size_t(bytesToRead))
	return result
}

// RandomBytesInPlace generates bytesToRead random bytes. This implementation
// uses either the default RNG algorithm ("system"), or whichever algorithm has
// been selected by RandomBytesSwitchAlgorithm. If bytesToRead exceeds the size
// of randomArray, only len(randomArray) bytes are read.
func RandomBytesInPlace(randomArray []byte, bytesToRead int) {
	if bytesToRead > len(randomArray) {
		bytesToRead = len(randomArray)
	}
	C.OQS_randombytes((*C.uint8_t)(unsafe.Pointer(&randomArray[0])),
		C.size_t(bytesToRead))
}

// RandomBytesSwitchAlgorithm switches the core OQS_randombytes to use the
// specified algorithm. Possible values are "system" and "OpenSSL".
// See <oqs/rand.h> liboqs header for more details.
func RandomBytesSwitchAlgorithm(algName string) error {
	if C.OQS_randombytes_switch_algorithm(C.CString(algName)) != C.OQS_SUCCESS {
		return errors.New("can not switch to \"" + algName + "\" algorithm")
	}
	return nil
}

// RandomBytesCustomAlgorithm switches RandomBytes to use the given function.
// This allows additional custom RNGs besides the provided ones. The provided
// RNG function must have the same signature as RandomBytesInPlace,
// i.e. func([]byte, int).
func RandomBytesCustomAlgorithm(fun func([]byte, int)) error {
	if fun == nil {
		return errors.New("the RNG algorithm callback can not be nil")
	}
	algorithmPtrCallback = fun
	C.OQS_randombytes_custom_algorithm(
		(C.algorithm_ptr)(unsafe.Pointer(C.algorithmPtr_cgo)))
	return nil
}

/**************** END Randomness ****************/
