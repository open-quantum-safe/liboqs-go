package rand

// C callbacks, DO NOT CHANGE

/*
#include <stdint.h>
#include <stddef.h>
void algorithmPtr_cgo(uint8_t* random_array, size_t bytes_to_read) {
	void algorithmPtr(uint8_t*, size_t);
	algorithmPtr(random_array, bytes_to_read);
}
*/
import "C"
