package oqs

// C callbacks, DO NOT CHANGE

/*
#include <stdint.h>
#include <stddef.h>
void randAlgorithmPtr_cgo(uint8_t* random_array, size_t bytes_to_read) {
	void randAlgorithmPtr(uint8_t*, size_t);
	randAlgorithmPtr(random_array, bytes_to_read);
}
*/
import "C"
