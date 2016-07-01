// Package forestdb uses the cgo compilation facilities to build the
// ForestDB library.
package forestdb

import (
	// explicit because these Go libraries do not export any Go symbols.
	_ "github.com/couchbaselabs/c-snappy"
)

// #cgo CPPFLAGS: -Iinternal/include -Iinternal/src -Iinternal/option -Iinternal/utils
// #cgo CPPFLAGS: -I../../couchbaselabs/c-snappy/internal
// #cgo CPPFLAGS: -DSNAPPY
// #cgo CXXFLAGS: -std=c++11 -DHAVE_GCC_ATOMICS -fno-omit-frame-pointer -momit-leaf-frame-pointer
// #cgo darwin LDFLAGS: -Wl,-undefined -Wl,dynamic_lookup
// #cgo !darwin LDFLAGS: -Wl,-unresolved-symbols=ignore-all
import "C"
