package badgercache

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"os"
	"testing"

	"github.com/pierrec/lz4"
)

var testVal = []byte(`
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
 incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
 nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
 Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore
 eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt
 in culpa qui officia deserunt mollit anim id est laborum.
`)

// The tests to check format involve creating a compressed version using the
// desired algorithm and then checking its header.
//
// All algorithm specific tests will also call GetFormat() to validate its
// behavior for that algorithm.

// this test uses a tarball compressed with bzip2 because compress/bzip2
// doesn't have a compressor.
func TestIsBZip2(t *testing.T) {
	f, err := os.Open("../test_files/test.bz2")
	if err != nil {
		t.Errorf("open test.bz2: expected no error, got %s", err)
		return
	}
	defer f.Close()
	ok, err := IsBZip2(f)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("expected ok to be true for bzip2, got false")
	}
	format, err := GetFormat(f)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != BZip2 {
		t.Errorf("expected format to be bzip2, got %s", format)
	}
}

func TestIsGZip(t *testing.T) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	n, err := w.Write(testVal)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if n != 452 {
		t.Errorf("Expected 452 bytes to be written; %d were", n)
	}
	w.Close()
	r := bytes.NewReader(buf.Bytes())
	ok, err := IsGZip(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	format, err := GetFormat(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != GZip {
		t.Errorf("Expected format to be gzip got %s", format)
	}
}

func TestIsLZ4(t *testing.T) {
	var buf bytes.Buffer
	lw := lz4.NewWriter(&buf)
	n, err := lw.Write(testVal)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if n != 452 {
		t.Errorf("Expected 452 bytes to be written; %d were", n)
	}
	lw.Close()
	r := bytes.NewReader(buf.Bytes())
	ok, err := IsLZ4(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	format, err := GetFormat(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != LZ4 {
		t.Errorf("Expected format to be LZ4 got %s", format)
	}
}

func TestIsTar(t *testing.T) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	hdr := &tar.Header{
		Name: "lorem.txt",
		Mode: 644,
		Size: int64(len(testVal)),
	}
	err := tw.WriteHeader(hdr)
	if err != nil {
		t.Errorf("unexpected error while writing tar header: %s", err)
	}
	_, err = tw.Write(testVal)
	if err != nil {
		t.Errorf("unexpected error while writing file to tar: %s", err)
	}
	err = tw.Close()
	if err != nil {
		t.Errorf("unexpected error while closing tar: %s", err)
	}
	r := bytes.NewReader(buf.Bytes())
	ok, err := IsTar(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	format, err := GetFormat(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != Tar {
		t.Errorf("Expected format to be tar got %s", format)
	}
}

// a file is used for the non-empty test because creating one using the test
// data in this func resulted in the zip empty header...probably an error on
// my part.
func TestIsZip(t *testing.T) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	err := w.Close()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	r := bytes.NewReader(buf.Bytes())
	ok, err := IsZip(r)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	f, err := os.Open("../test_files/test.zip")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	defer f.Close()
	ok, err = IsZip(f)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	format, err := GetFormat(f)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != Zip {
		t.Errorf("Expected format to be gzip got %s", format)
	}
}

// TODO: commented out because LZW output doesn't have the magic number.
// Figure out why and resolve.
// Another example http://play.golang.org/p/zGLAj1ruoh
/*
func TestIsLZW(t *testing.T) {
	var buf bytes.Buffer
	lw := lzw.NewWriter(&buf, lzw.LSB, 8)
	n, err := lw.Write(testVal)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if n != 452 {
		t.Errorf("Expected 452 bytes to be copied; %d were", n)
	}
	lw.Close()
	t.Errorf("%x", buf.Bytes())
	rr := bytes.NewReader(buf.Bytes())
	ok, err := IsLZW(rr)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if !ok {
		t.Error("Expected ok to be true, got false")
	}
	format, err := GetFormat(rr)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if format != LZW {
		t.Errorf("Expected format to be LZW got %s", format)
	}
}
*/

func TestString(t *testing.T) {
	tests := []struct {
		f        Format
		expected string
	}{
		{Format(-1), "Format(-1)"},
		{Format(99), "Format(99)"},
		{Unknown, "Unknown"},
		{GZip, "GZip"},
		{BZip2, "BZip2"},
		{LZ4, "LZ4"},
		{Tar, "Tar"},
		{Tar1, "Tar1"},
		{Tar2, "Tar2"},
		{Zip, "Zip"},
		{ZipEmpty, "Empty Zip"},
		{ZipSpanned, "Spanned Zip"},
	}
	for _, test := range tests {
		s := test.f.String()
		if s != test.expected {
			t.Errorf("got %q; want %q", s, test.expected)
		}
	}
}

func TestExt(t *testing.T) {
	tests := []struct {
		f        Format
		expected string
	}{
		{Unknown, "unknown"},
		{Format(-1), "unknown"},
		{Format(99), "unknown"},
		{GZip, ".gz"},
		{BZip2, ".bz2"},
		{LZ4, ".lz4"},
		{Tar, ".tar"},
		{Tar1, ".tar"},
		{Tar2, ".tar"},
		{Zip, ".zip"},
		{ZipEmpty, ".zip"},
		{ZipSpanned, ".zip"},
	}
	for _, test := range tests {
		s := test.f.Ext()
		if s != test.expected {
			t.Errorf("got %s; want %s", s, test.expected)
		}
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		v string
		f Format
	}{
		{"", Unknown},
		{"z", Unknown},
		{"gzip", GZip},
		{"tar.gz", GZip},
		{"tgz", GZip},
		{"tar", Tar},
		{"bz2", BZip2},
		{"tbz", BZip2},
		{"tb2", BZip2},
		{"tbz2", BZip2},
		{"tar.bz2", BZip2},
		{"lz4", LZ4},
		{"tar.lz4", LZ4},
		{"tz4", LZ4},
		{"zip", Zip},
	}
	for _, test := range tests {
		f := ParseFormat(test.v)
		if f != test.f {
			t.Errorf("got %s; want %s", f, test.f)
		}
	}
}

func TestGetFormat(t *testing.T) {
	tests := []struct {
		name   Format
		format Format
		bytes  []byte
		offset int
		err    error
	}{
		{Unknown, Unknown, []byte{}, 0, ErrEmpty},
		{Unknown, Unknown, []byte{}, 0, ErrUnknown},
		{Unknown, Unknown, []byte{0x10, 0x11}, 0, ErrUnknown},
		{GZip, GZip, magicnumGZip, 0, nil},
		{BZip2, BZip2, magicnumBZip2, 0, nil},
		{LZ4, LZ4, []byte{0x04, 0x22, 0x4d, 0x18}, 0, nil},
		{Tar1, Tar, magicnumTar1, 257, nil},
		{Tar2, Tar, magicnumTar2, 257, nil},
		{Zip, Zip, magicnumZip, 0, nil},
		{ZipEmpty, Zip, magicnumZipEmpty, 0, nil},
		{ZipSpanned, Zip, magicnumZipSpanned, 0, nil},
	}

	for i, test := range tests {
		var b []byte
		if i != 0 {
			b = make([]byte, 512)
			var j int
			for i := test.offset; i < test.offset+len(test.bytes); i++ {
				b[i] = test.bytes[j]
				j++
			}
		}
		r := bytes.NewReader(b)
		f, err := GetFormat(r)
		if err != nil {
			if err != test.err {
				t.Errorf("%s: got %q; want %q", test.name, err, test.err)
			}
			continue
		}
		if test.err != nil {
			t.Errorf("%s: no error; expected %q", test.name, test.err)
			continue
		}
		if f != test.format {
			t.Errorf("%s: got %s; want %s", test.name, f, test.format)
		}
	}
}
