package custom_comparator

//#include "comparator.h"
import "C"

import (
	"unsafe"
)

var CompareBytesReversedPointer unsafe.Pointer = unsafe.Pointer(C.cmp_variable)
