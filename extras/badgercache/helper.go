package badgercache

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/dgraph-io/badger"
	"github.com/golang/snappy"
	"github.com/hashicorp/go-msgpack/codec"
	"gopkg.in/kothar/brotli-go.v0/enc"
	// "github.com/rohanthewiz/roencoding"
	// log "github.com/sirupsen/logrus"
)

const (
	// GzipMinSize gzip min size
	GzipMinSize = 1024
	// CacheFormatRaw raw
	CacheFormatRaw = 0
	// CacheFormatRawGzip raw gzip
	CacheFormatRawGzip = 1
	// CacheFormatJSON json
	CacheFormatJSON = 10
	// CacheFormatJSONGzip json gzip
	CacheFormatJSONGzip = 11
)

//GetChecksum gets the checksum.
func getChecksum(data string) string {
	checksum := sha1.Sum([]byte(data))
	return fmt.Sprintf("%x", checksum)
}

func Compress(data []byte) ([]byte, error) {
	return snappy.Encode([]byte{}, data), nil
}

func Decompress(data []byte) ([]byte, error) {
	return snappy.Decode([]byte{}, data)
}

func ungzipData(data []byte) ([]byte, error) {
	raw := bytes.NewBuffer(data)
	r, err := gzip.NewReader(raw)
	if err != nil && err != io.EOF {
		// log.Fatalln("badgercache.Get().ungzipData().gzip.NewReader(), ERROR: ", err)
		return []byte{}, err
	}
	defer r.Close()

	resp, err := ioutil.ReadAll(r)
	if err != nil {
		// log.Fatalln("badgercache.Get().ungzipData().ioutil.ReadAll(), ERROR: ", err)
		return nil, err
	}
	return resp, nil
}

/*
func ungzipData2(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Fatalln("badgercache.Get().ungzipData().gzip.NewReader(), ERROR: ", err)
		return nil, err
	}
	defer r.Close()
	data, err = ioutil.ReadAll(r)
	if err != nil {
		log.Fatalln("badgercache.Get().ungzipData().ioutil.ReadAll(), ERROR: ", err)
		return nil, err
	}
	return data, nil
}
*/

func gzipData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		// log.Fatalln("badgercache.Get().ungzipData().gzip.Write(), ERROR: ", err)
		return nil, err
	}
	err = w.Close()
	if err != nil {
		// log.Fatalln("badgercache.Get().ungzipData().w.Close(), ERROR: ", err)
		return nil, err
	}
	return b.Bytes(), nil
}

func fmtToJsonArr(s []byte) []byte {
	s = bytes.Replace(s, []byte("{"), []byte("[{"), 1)
	s = bytes.Replace(s, []byte("}"), []byte("},"), -1)
	s = bytes.TrimSuffix(s, []byte(","))
	s = append(s, []byte("]")...)
	return s
}

type compressed struct {
	bufBr   *bytes.Buffer
	bufRaw  *bytes.Buffer
	bufGzip *bytes.Buffer
}

func compressWithGzip(b []byte) *bytes.Buffer {
	buf := &bytes.Buffer{}
	zw := gzip.NewWriter(buf)
	_, err := zw.Write(b)

	if err != nil {
		// log.Println("Gzip compression error: ", err)
		return buf
	}

	err = zw.Close()

	if err != nil {
		// log.Println("Gzip compression error: ", err)
		return buf
	}

	return buf
}

func compressWithBrotli(input []byte) *bytes.Buffer {
	params := enc.NewBrotliParams()
	// brotli supports quality values from 0 to 11 included
	// 0 is the fastest, 11 is the most compressed but slowest
	params.SetQuality(11)
	compressed, _ := enc.CompressBuffer(params, input, make([]byte, 0))
	buf := bytes.NewBuffer(compressed)
	return buf
}

// Decode reverses the encode operation on a byte slice input
func decodeMsgPack(buf []byte, out interface{}) error {
	r := bytes.NewBuffer(buf)
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(r, &hd)
	return dec.Decode(out)
}

// Encode writes an encoded object to a new bytes buffer
func encodeMsgPack(in interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	hd := codec.MsgpackHandle{}
	enc := codec.NewEncoder(buf, &hd)
	err := enc.Encode(in)
	return buf, err
}

// Converts bytes to an integer
func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Converts a uint to a byte slice
func uint64ToBytes(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

func safeKey(item *badger.Item) []byte {
	key := item.Key()
	dst := make([]byte, len(key))
	copy(dst, key)
	return dst
}

/*
	Refs:
	- https://github.com/deepakjois/badgerpp/blob/master/cmd/root.go
*/
/*
func (c *Cache) SetString(key, val string) error {
	if c.debug {
		log.WithFields(log.Fields{
			"key": key,
		}).Debug("badgercache.SetString()")
	}
	return c.db.Set([]byte(key), []byte(val))
}

func (c *Cache) GetString(key string) (out string, err error) {
	var item badger.KVItem
	err = c.Get([]byte(key), &item)
	if err != nil {
		return
	}
	out = string(item.Value())
	return
}

func (c *Cache) SetBytes(k, v []byte) error {
	if c.debug {
		log.WithFields(log.Fields{
			"key": k,
		}).Debug("badgercache.SetString()")
	}
	return c.db.Set(k, v)
}

func (c *Cache) GetBytes(k []byte) (out []byte, err error) {
	var item badger.KVItem
	err = c.Get(k, &item)
	if err != nil {
		return
	}
	out = item.Value()
	return
}
*/

/*
// Add a hashed key to the store if it doesn't already exist
func (c *Cache) TouchHashed(in string) (err error) {
	return c.db.Touch([]byte(roencoding.XXHash(in)))
}

// Does hash of key exist in the store?
func (c *Cache) ExistsHashed(in string) (exists bool, err error) {
	return c.db.Exists([]byte(roencoding.XXHash(in)))
}
*/
