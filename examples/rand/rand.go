package main

import (
	"fmt"
	oqsrand "github.com/open-quantum-safe/liboqs-go/oqs/rand"
)

func main() {
	rnd := [48]byte{0: 100, 20: 200, 47: 150}
	oqsrand.RandomBytesNistKatInit(rnd, []byte("initMsg"))
	oqsrand.RandomBytesSwitchAlgorithm("NIST-KAT")
	fmt.Println(oqsrand.RandomBytes(10))
}
