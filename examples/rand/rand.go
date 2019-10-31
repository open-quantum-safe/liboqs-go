// various RNGs Go examples
package main

import (
	"fmt"
	oqsrand "github.com/open-quantum-safe/liboqs-go/oqs/rand"
)

// CustomRNG provides a (trivial) custom random number generator; the memory is
// provided by the caller, i.e. oqsrand.RandomBytes or
// oqsrand.RandomBytesInPlace
func CustomRNG(randomArray []byte, bytesToRead int) {
	for i := 0; i < bytesToRead; i++ {
		randomArray[i] = byte(i)
	}
}

func main() {
	oqsrand.RandomBytesSwitchAlgorithm("NIST-KAT")
	entropySeed := [48]byte{0: 100, 20: 200, 47: 150}
	oqsrand.RandomBytesNistKatInit(entropySeed, nil)
	fmt.Printf("%-18s% X\n", "NIST-KAT: ", oqsrand.RandomBytes(32))

	oqsrand.RandomBytesCustomAlgorithm(CustomRNG)
	fmt.Printf("%-18s% X\n", "Custom RNG: ", oqsrand.RandomBytes(32))

	oqsrand.RandomBytesSwitchAlgorithm("OpenSSL")
	fmt.Printf("%-18s% X\n", "OpenSSL: ", oqsrand.RandomBytes(32))

	oqsrand.RandomBytesSwitchAlgorithm("system")
	fmt.Printf("%18s% X\n", "System (default): ", oqsrand.RandomBytes(32))
}
