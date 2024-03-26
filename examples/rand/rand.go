// Various RNGs Go example
package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/open-quantum-safe/liboqs-go/oqs"
	oqsrand "github.com/open-quantum-safe/liboqs-go/oqs/rand" // RNG support
)

// CustomRNG provides a (trivial) custom random number generator; the memory is
// provided by the caller, i.e. oqsrand.RandomBytes or
// oqsrand.RandomBytesInPlace
func CustomRNG(randomArray []byte, bytesToRead int) {
	for i := 0; i < bytesToRead; i++ {
		randomArray[i] = byte(i % 256)
	}
}

func main() {
	fmt.Println("liboqs version: " + oqs.LiboqsVersion())

	if err := oqsrand.RandomBytesSwitchAlgorithm("system"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%18s% X\n", "System (default): ", oqsrand.RandomBytes(32))
	if err := oqsrand.RandomBytesCustomAlgorithm(CustomRNG); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-18s% X\n", "Custom RNG: ", oqsrand.RandomBytes(32))

	// We do not yet support OpenSSL under Windows
	if runtime.GOOS != "windows" {
		if err := oqsrand.RandomBytesSwitchAlgorithm("OpenSSL"); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-18s% X\n", "OpenSSL: ", oqsrand.RandomBytes(32))
	}
}
