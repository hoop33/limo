package custom_comparator

//#include <stdlib.h>
import "C"

import (
	"bytes"
	"unsafe"
)

//export CompareBytesReversed
func CompareBytesReversed(k1 unsafe.Pointer, k1len C.size_t, k2 unsafe.Pointer, k2len C.size_t) int {
	key1 := C.GoBytes(k1, C.int(k1len))
	key2 := C.GoBytes(k2, C.int(k2len))
	return -bytes.Compare(key1, key2)
}
