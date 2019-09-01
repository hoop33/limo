package badgercache

/*
import (
	"bufio"
	"net/textproto"
	"strings"

	"github.com/json-iterator/go"
	// "github.com/opennota/substring"
	"github.com/rai-project/linguist"
	"github.com/wttw/orderedheaders"
	"gopkg.in/src-d/enry.v1"
	yaml "gopkg.in/yaml.v2"

	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"

	// "github.com/alvaroloes/enumer"
	// "github.com/mlposey/dictionary"
	"github.com/tidwall/match"
	// "github.com/foszor/stringutils"
)
*/

/*
	Refs:
	- https://github.com/martinlindhe/formats/blob/master/parse/archive/arc_lzma.go
	- https://raw.githubusercontent.com/mohae/magicnum/master/compress/compress_test.go
	- https://github.com/kyurdakok/go-detect/blob/master/main.go
*/

//func SubString(pattern string, in []byte) bool {
//	m := substring.NewMatcher(pattern)
//	return m.Match(string(in))
//}
/*
func Match(pattern string, in []byte) bool {
	return match.Match(string(in), pattern)
}

func reader(s string) *textproto.Reader {
	return textproto.NewReader(bufio.NewReader(strings.NewReader(s)))
}

func ReadHeader(content string) {
	r := reader(content)
	hdrs, err := orderedheaders.ReadHeader(r)
	if err != nil {
		log.WithFields(log.Fields{
			"content": content,
		}).Errorln("badgercache.ReadHeader(), ERROR: ", err)
	}
	pp.Println(hdrs)
	return
}

func DetectLang(content string) string {
	lang := linguist.Detect(content)
	log.WithFields(log.Fields{
		"lang": lang,
	}).Info("badgercache.DetectLang()")
	return lang
}

func DetectType(name string, content string) (string, bool) {
	log.WithFields(log.Fields{
		"name": name,
	}).Info("badgercache.DetectType()")
	lang, safe := enry.GetLanguageByContent(name, []byte(content))
	return lang, safe
}

func IsJSON(b []byte) bool {
	var js jsoniter.RawMessage
	return jsoniter.Unmarshal(b, &js) == nil
}

func IsYAML(b []byte) bool {
	var y yaml.MapSlice
	return yaml.Unmarshal(b, &y) == nil
}

func IsGZIP(b []byte) bool {
	if b[0] != 0x1f || b[1] != 0x8b {
		return false
	}
	return true
}

func IsZIP(b []byte) bool {
	if b[0] != 'P' || b[1] != 'K' || b[2] != 3 || b[3] != 4 {
		return false
	}
	return true
}

func IsISO(b []byte) bool {
	pos := 0x8000
	if b[pos] != 1 || b[pos+1] != 'C' || b[pos+2] != 'D' {
		return false
	}
	return true
}

func IsDEB(b []byte) bool {
	s := string(b[0:21])
	return s == "!<arch>\n"+"debian-binary"
}

func IsBZIP2(b []byte) bool {
	if b[0] != 'B' || b[1] != 'Z' {
		return false
	}
	if b[2] != 'h' {
		// only huffman encoding is used in the format
		return false
	}
	return true
}

func IsXZ(b []byte) bool {
	if b[0] != 0xfd || b[1] != '7' || b[2] != 'z' || b[3] != 'X' ||
		b[4] != 'Z' || b[5] != 0x00 {
		return false
	}
	return true
}

func IsZLib(b []byte) bool {
	// XXX only matches zlib streams without dictionary.. this dont always work
	if b[0] != 0x78 {
		return false
	}
	if b[1] != 0x01 && b[1] != 0x9c && b[1] != 0xda {
		// compression level
		return false
	}
	return true
}

func IsXAR(b []byte) bool {
	if b[0] != 'x' || b[1] != 'a' || b[2] != 'r' || b[3] != '!' {
		return false
	}
	return true
}

func IsRAR(b []byte) bool {
	if b[0] != 'R' || b[1] != 'a' || b[2] != 'r' || b[3] != '!' {
		return false
	}
	// RAR 4.x signature
	//if (ReadByte() != 0x1A || ReadByte() != 0x07 || ReadByte() != 0x00)
	//    return false;
	// RAR 5.0 signature
	//if (ReadByte() != 0x1A || ReadByte() != 0x07 || ReadByte() != 0x01 || ReadByte() != 0x00)
	//    return false;
	return true
}

func IsLZMA(b []byte) bool {
	// XXX not proper magic , need other check
	if b[0] != 0x5d || b[1] != 0x00 {
		return false
	}
	return true
}

func IsLZH(b []byte) bool {
	if b[2] != '-' || b[3] != 'l' {
		return false
	}
	if b[4] == 'h' || b[4] == 'z' {
		return true
	}
	return false
}
*/
