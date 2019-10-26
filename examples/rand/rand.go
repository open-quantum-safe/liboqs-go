package main

import (
	"fmt"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	oqsrand "github.com/open-quantum-safe/liboqs-go/oqs/rand"
)

func main() {

	rnd := [48]byte{0: 100, 20: 200, 47: 150}
	oqsrand.RandomBytesNistKatInit(rnd, []byte("initMsg"))
	r := oqsrand.RandomBytesSwitchAlgorithm("NIST-KAT")
	if r == oqs.ERROR {
		fmt.Println("Error")
	} else {
		fmt.Println(r)
	}
	fmt.Println(oqsrand.RandomBytes(10))
}
