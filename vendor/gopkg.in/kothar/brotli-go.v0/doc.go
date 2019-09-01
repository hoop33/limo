// Package brotli contains bindings for the Brotli compression library
//
// This is a very basic Cgo wrapper for the enc and dec directories from the Brotli sources. I've made a few minor changes to get
// things working with Go.
//
// 1. The default dictionary has been extracted to a separate 'shared' package to allow linking the enc and dec cgo modules if you use both. Otherwise there are duplicate symbols, as described in the dictionary.h header files.
//
// 2. The dictionary variable name for the dec package has been modified for the same reason, to avoid linker collisions.
package brotli // import "gopkg.in/kothar/brotli-go.v0"

import (
	"gopkg.in/kothar/brotli-go.v0/dec"
	"gopkg.in/kothar/brotli-go.v0/enc"
	"gopkg.in/kothar/brotli-go.v0/shared"
)

var (
	_ = enc.CompressBuffer
	_ = dec.DecompressBuffer
	_ = shared.GetDictionary
)
