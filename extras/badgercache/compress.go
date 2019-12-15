package badgercache

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	// "github.com/TerraTech/gzipInfo/pkg/gzipInfo"
	// log "github.com/sirupsen/logrus"
)

const (
	Unknown    Format = iota // unknown format
	GZip                     // Gzip compression format
	BZip2                    // Bzip2 compression
	LZ4                      // LZ4 compression
	Tar                      // Tar format; normally used
	Tar1                     // Tar1 magicnum format; normalizes to Tar
	Tar2                     // Tar1 magicnum format; normalizes to Tar
	Zip                      // Zip archive
	ZipEmpty                 // Empty Zip Archive
	ZipSpanned               // Spanned Zip Archive
)

// Magic numbers for compression and archive formats
var (
	magicnumGZip       = []byte{0x1f, 0x8b}
	magicnumBZip2      = []byte{0x42, 0x5a, 0x68}
	magicnumLZ4        = []byte{0x18, 0x4d, 0x22, 0x04}
	magicnumTar1       = []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x30, 0x30} // offset: 257
	magicnumTar2       = []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x20, 0x00} // offset: 257
	magicnumZip        = []byte{0x50, 0x4b, 0x03, 0x04}
	magicnumZipEmpty   = []byte{0x50, 0x4b, 0x05, 0x06}
	magicnumZipSpanned = []byte{0x50, 0x4b, 0x07, 0x08}
	//magicnumLZW        = []byte{0x1F, 0x9d}
)

var (
	ErrUnknown = errors.New("unknown compression format")
	ErrEmpty   = errors.New("no data to read")
)

type Format int

const formatName = "UnknownGZipBZip2LZ4TarTar1Tar2ZipEmpty ZipSpanned Zip"

var formatIndex = [...]uint8{0, 7, 11, 16, 19, 22, 26, 30, 33, 42, 53}

func (i Format) String() string {
	if i < 0 || i >= Format(len(formatIndex)-1) {
		return fmt.Sprintf("Format(%d)", i)
	}
	return formatName[formatIndex[i]:formatIndex[i+1]]
}

// Ext returns the extension for the format. Formats may have more than one
// accepted extension; alternate extensiona are not supported.
func (f Format) Ext() string {
	switch f {
	case GZip:
		return ".gz"
	case BZip2:
		return ".bz2"
	case LZ4:
		return ".lz4"
	case Tar, Tar1, Tar2:
		return ".tar"
	case Zip, ZipEmpty, ZipSpanned:
		return ".zip"
		//case LZW:
		//	return ".Z"
	}
	return "unknown"
}

// ParseFormat takes a string and returns the format or unknown. Any compressed
// tar extensions are returned as the compression format and not tar.
//
// If the passed string starts with a '.', it is removed.
// All strings are lowercased
func ParseFormat(s string) Format {
	if len(s) == 0 {
		return Unknown
	}
	if s[0] == '.' {
		s = s[1:]
	}
	s = strings.ToLower(s)
	switch s {
	case "gzip", "tar.gz", "tgz":
		return GZip
	case "tar":
		return Tar
	case "bz2", "tbz", "tb2", "tbz2", "tar.bz2":
		return BZip2
	case "lz4", "tar.lz4", "tz4":
		return LZ4
	case "zip":
		return Zip
	}
	return Unknown
}

// GetFormat tries to match up the data in the Reader to a supported
// magic number, if a match isn't found, UnsupportedFmt is returned
//
// For zips, this will also match on files with empty zip or spanned zip magic
// numbers.  If you need to distinguich between the various zip formats, use
// something else.
func GetFormat(r io.ReaderAt) (Format, error) {
	// see if the reader contains anything
	b := make([]byte, 1)
	if _, err := r.ReadAt(b, 0); err == io.EOF {
		return Unknown, ErrEmpty
	}

	ok, err := IsLZ4(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return LZ4, nil
	}
	ok, err = IsGZip(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return GZip, nil
	}
	ok, err = IsZip(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return Zip, nil
	}
	ok, err = IsTar(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return Tar, nil
	}
	ok, err = IsBZip2(r)
	if err != nil {
		return Unknown, err
	}
	if ok {
		return BZip2, nil
	}
	//ok, err = IsLZW(r)
	//if err != nil {
	//	return Unknown, err
	//}
	//if ok {
	//	return LZW, nil
	//}
	return Unknown, ErrUnknown
}

// IsBZip2 checks to see if the received reader's contents are in bzip2 format
// by checking the magic numbers.
func IsBZip2(r io.ReaderAt) (bool, error) {
	h := make([]byte, 3)
	// Read the first 3 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var hb [3]byte
	// check for bzip2
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &hb)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched bzip2's magic number: %s", err)
	}
	var cb [3]byte
	cbuf := bytes.NewBuffer(magicnumBZip2)
	err = binary.Read(cbuf, binary.BigEndian, &cb)
	if err != nil {
		return false, fmt.Errorf("error while converting bzip2 magic number for comparison: %s", err)
	}
	if hb == cb {
		return true, nil
	}
	return false, nil
}

// IsGZip checks to see if the received reader's contents are in gzip format
// by checking the magic numbers.
func IsGZip(r io.ReaderAt) (bool, error) {
	h := make([]byte, 2)
	// Read the first 2 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h16 uint16
	// check for gzip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h16)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched bzip2's magic number: %s", err)
	}
	var c16 uint16
	cbuf := bytes.NewBuffer(magicnumGZip)
	err = binary.Read(cbuf, binary.BigEndian, &c16)
	if err != nil {
		return false, fmt.Errorf("error while converting bzip2 magic number for comparison: %s", err)
	}
	if h16 == c16 {
		return true, nil
	}
	return false, nil
}

// IsLZ4 checks to see if the received reader's contents are in LZ4 foramt by
// checking the magic numbers.
func IsLZ4(r io.ReaderAt) (bool, error) {
	h := make([]byte, 4)
	// Read the first 4 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h32 uint32
	// check for lz4
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &h32)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched LZ4's magic number: %s", err)
	}
	var c32 uint32
	cbuf := bytes.NewBuffer(magicnumLZ4)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting LZ4 magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	return false, nil
}

// IsLZW checks to see if the received reader's contents are in LZ4 format by
// checking the magic numbers.
//
// TODO: unsupported until I have a better understanding of how to handle LZW
/*
func IsLZW(r io.ReaderAt) (bool, error) {
	h := make([]byte, 2)
	// Reat the first 8 bytes since that's where most magic numbers are
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h16 uint16
	// check for lzw
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.LittleEndian, &h16)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched LZW's magic number: %s", err)
	}
	var c16 uint16
	cbuf := bytes.NewBuffer(magicnumLZW)
	err = binary.Read(cbuf, binary.BigEndian, &c16)
	if err != nil {
		return false, fmt.Errorf("error while converting LZW magic number for comparison: %s", err)
	}
	if h16 == c16 {
		return true, nil
	}
	return false, nil
}
*/

// IsTar checks to see if the received reader's contents are in the tar format
// by checking the magic numbers. This evaluates using both tar1 and tar2 magic
// numbers.
func IsTar(r io.ReaderAt) (bool, error) {
	h := make([]byte, 8)
	// Read the first 8 bytes at offset 257
	_, err := r.ReadAt(h, 257)
	if err != nil {
		return false, err
	}
	var h64 uint64
	// check for Zip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h64)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched tar's magic number: %s", err)
	}
	var c64 uint64
	cbuf := bytes.NewBuffer(magicnumTar1)
	err = binary.Read(cbuf, binary.BigEndian, &c64)
	if err != nil {
		return false, fmt.Errorf("error while converting the tar magic number for comparison: %s", err)
	}
	if h64 == c64 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(magicnumTar2)
	err = binary.Read(cbuf, binary.BigEndian, &c64)
	if err != nil {
		return false, fmt.Errorf("error while converting the empty tar magic number for comparison: %s", err)
	}
	if h64 == c64 {
		return true, nil
	}
	return false, nil
}

// IsZip checks to see if the received reader's contents are in the zip format
// by checking the magic numbers. This will match on zip, empty zip and spanned
// zip magic numbers. If you need to distinguish between those, use something
// else.
func IsZip(r io.ReaderAt) (bool, error) {
	h := make([]byte, 4)
	// Read the first 4 bytes
	_, err := r.ReadAt(h, 0)
	if err != nil {
		return false, err
	}
	var h32 uint32
	// check for Zip
	hbuf := bytes.NewReader(h)
	err = binary.Read(hbuf, binary.BigEndian, &h32)
	if err != nil {
		return false, fmt.Errorf("error while checking if input matched zip's magic number: %s", err)
	}
	var c32 uint32
	cbuf := bytes.NewBuffer(magicnumZip)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(magicnumZipEmpty)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the empty zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	cbuf = bytes.NewBuffer(magicnumZipSpanned)
	err = binary.Read(cbuf, binary.BigEndian, &c32)
	if err != nil {
		return false, fmt.Errorf("error while converting the spanned zip magic number for comparison: %s", err)
	}
	if h32 == c32 {
		return true, nil
	}
	return false, nil
}
