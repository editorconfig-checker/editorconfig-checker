// package encoding contains all the encoding functions
package encoding

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/baulk/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

const BinaryData = "binary"

var encodings = map[string]encoding.Encoding{
	// In https://github.com/golang/text/blob/HEAD/encoding/htmlindex/map.go#L64 and
	//    https://github.com/golang/text/blob/HEAD/encoding/ianaindex/ianaindex.go#L156 :
	"utf8":        unicode.UTF8,
	"ibm866":      charmap.CodePage866,
	"iso88592":    charmap.ISO8859_2,
	"iso88593":    charmap.ISO8859_3,
	"iso88594":    charmap.ISO8859_4,
	"iso88595":    charmap.ISO8859_5,
	"iso88596":    charmap.ISO8859_6,
	"iso88597":    charmap.ISO8859_7,
	"iso88598":    charmap.ISO8859_8,
	"iso885910":   charmap.ISO8859_10,
	"iso885913":   charmap.ISO8859_13,
	"iso885914":   charmap.ISO8859_14,
	"iso885915":   charmap.ISO8859_15,
	"iso885916":   charmap.ISO8859_16,
	"koi8r":       charmap.KOI8R,
	"koi8u":       charmap.KOI8U,
	"macintosh":   charmap.Macintosh,
	"windows874":  charmap.Windows874,
	"windows1250": charmap.Windows1250,
	"windows1251": charmap.Windows1251,
	"windows1252": charmap.Windows1252,
	"windows1253": charmap.Windows1253,
	"windows1254": charmap.Windows1254,
	"windows1255": charmap.Windows1255,
	"windows1256": charmap.Windows1256,
	"windows1257": charmap.Windows1257,
	"windows1258": charmap.Windows1258,
	"gbk":         simplifiedchinese.GBK,
	"gb18030":     simplifiedchinese.GB18030,
	"big5":        traditionalchinese.Big5,
	"eucjp":       japanese.EUCJP,
	"iso2022jp":   japanese.ISO2022JP,
	"shiftjis":    japanese.ShiftJIS,
	"euckr":       korean.EUCKR,
	"utf16be":     unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	"utf16le":     unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
	// Not in https://github.com/golang/text/blob/HEAD/encoding/htmlindex/map.go#L64 :
	"iso88591":  charmap.ISO8859_1,
	"ibm037":    charmap.CodePage037,
	"ibm437":    charmap.CodePage437,
	"ibm850":    charmap.CodePage850,
	"ibm852":    charmap.CodePage852,
	"ibm855":    charmap.CodePage855,
	"ibm858":    charmap.CodePage858,
	"ibm860":    charmap.CodePage860,
	"ibm862":    charmap.CodePage862,
	"ibm863":    charmap.CodePage863,
	"ibm865":    charmap.CodePage865,
	"ibm1047":   charmap.CodePage1047,
	"ibm1140":   charmap.CodePage1140,
	"iso88596e": charmap.ISO8859_6E,
	"iso88596i": charmap.ISO8859_6I,
	"iso88598e": charmap.ISO8859_8E,
	"iso88598i": charmap.ISO8859_8I,
	"iso88599":  charmap.ISO8859_9,
	"hzgb2312":  simplifiedchinese.HZGB2312,
	// Not https://github.com/golang/text/blob/HEAD/encoding/ianaindex/ianaindex.go#L156 :
	"macintoshcyrillic": charmap.MacintoshCyrillic,
	// Not in https://github.com/golang/text/blob/HEAD/encoding/htmlindex/map.go#L64 or
	//        https://github.com/golang/text/blob/HEAD/encoding/ianaindex/ianaindex.go#L156 :
	"utf8bom": unicode.UTF8,
	"utf32be": utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
	"utf32le": utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM),
}

// In https://github.com/golang/text/blob/HEAD/encoding/ianaindex/ianaindex.go#L156
// but not included above:
// 	 enc3:    asciiEnc,
//   enc1015: unicode.UTF16(unicode.BigEndian, unicode.UseBOM)

// DecodeBytes converts a byte array to a string
func DecodeBytes(contentBytes []byte) (string, string, error) {
	contentString := string(contentBytes)

	charset, err := detectText(contentBytes)
	if err != nil {
		if IsBinaryFile(contentBytes) {
			return contentString, BinaryData, nil
		}
		return contentString, charset, err
	}
	decodedContentString, err := decodeText(contentBytes, charset)
	if err != nil {
		if IsBinaryFile(contentBytes) {
			return contentString, BinaryData, nil
		}
		return contentString, charset, err
	}
	return decodedContentString, charset, nil
}

func detectText(contentBytes []byte) (string, error) {
	detector := chardet.NewTextDetector()
	results, err := detector.DetectAll(contentBytes)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("Failed to determine charset")
	}
	confidence := -1
	keys := make([]string, 0, len(results))
	for _, result := range results {
		_, ok := getEncoding(result.Charset)
		if !ok {
			continue
		}
		if result.Confidence < confidence {
			break
		}
		confidence = result.Confidence
		keys = append(keys, result.Charset)
	}
	sort.Strings(keys)
	return keys[0], nil
}

func decodeText(contentBytes []byte, charset string) (string, error) {
	enc, ok := getEncoding(charset)
	if !ok {
		return "", fmt.Errorf("unrecognized charset %s", charset)
	}
	var err error
	contentBytes, err = enc.NewDecoder().Bytes(contentBytes)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(contentBytes) {
		return "", fmt.Errorf("the file is not a valid UTF-8 encoded file")
	}
	return string(contentBytes), nil
}

func getEncoding(charset string) (encoding.Encoding, bool) {
	r := strings.NewReplacer("-", "", "_", "")
	key := strings.ToLower(r.Replace(charset))
	enc, ok := encodings[key]
	return enc, ok
}

var binaryChars = [256]bool{}

func init() {
	// Allow tab (9), lf (10), ff (12), and cr (13)
	trues := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 11, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	for _, i := range trues {
		binaryChars[i] = true
	}
}

// IsBinaryFile returns true if the bytes contain \x00-\x08,\x0b,\x0e-\x1f
func IsBinaryFile(rawFileContent []byte) bool {
	for _, b := range rawFileContent {
		if binaryChars[b] {
			return true
		}
	}
	return false
}
