// Various RNGs Go example
package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/open-quantum-safe/liboqs-go/oqs"
	"github.com/open-quantum-safe/liboqs-go/oqs/rand" // RNG support
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

	if err := rand.RandomBytesSwitchAlgorithm("system"); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%18s% X\n", "System (default): ", rand.RandomBytes(32))
	if err := rand.RandomBytesCustomAlgorithm(CustomRNG); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-18s% X\n", "Custom RNG: ", rand.RandomBytes(32))

	// We do not yet support OpenSSL under Windows
	if runtime.GOOS != "windows" {
		if err := rand.RandomBytesSwitchAlgorithm("OpenSSL"); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-18s% X\n", "OpenSSL: ", rand.RandomBytes(32))
	}
}
