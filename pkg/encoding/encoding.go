// package encoding contains all the encoding functions
package encoding

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/wlynxg/chardet"
	"github.com/wlynxg/chardet/consts"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

// HitMissRatioForUTF1632Checks defines the ratio of hits to misses that must
// be met for the IsUTF* functions to return true. For example, the number 1.1
// requires the function find 10% more hits than misses. This number was chosen
// empiricallly, but may need to be adjusted as we gather more test results.
const HitMissRatioForUTF1632Checks = 1.01

// MinConfidenceForUTF1632Checks defines the minimum confidence factor to use
// our own check for if a file is UTF16/UTF32 encoded. The number 1.0 represents
// a 100% confidence level, and may be too high a setting, as it may not ever be
// returned by our chardet library.
const MinConfidenceForUTF1632Checks = 0.5

// BinaryData contains the string to return if the data contains binary
// data (data that is not decodable, by any of the decoders that golang provides).
const BinaryData = "binary"

const (
	// See https://spec.editorconfig.org/#supported-pairs
	// CharsetLatin1 defines the value for the ISO-8859-1 character set.
	CharsetLatin1 = "latin1"
	// CharsetUTF8 defines the value for the UTF-8 character set, without a
	// Byte Order Mark (BOM).
	CharsetUTF8 = "utf-8"
	// CharsetUTF8BOM defines the value for the UTF-8 character set, with a
	// Byte Order Mark (BOM).
	CharsetUTF8BOM = "utf-8-bom"
	// CharsetUTF16BE defines the value for the UTF-16BE (Big Endian) character
	// set, with, or without, a Byte Order Mark (BOM).
	CharsetUTF16BE = "utf-16be"
	// CharsetUTF16LE defines the value for the UTF-16LE (Lil Endian) character
	// set, with, or without, a Byte Order Mark (BOM).
	CharsetUTF16LE = "utf-16le"
)

// ValidCharsets contains an array of the allowable values for the charset key
// in an .editorconfig file.
var ValidCharsets = []string{
	CharsetLatin1,
	CharsetUTF8,
	CharsetUTF8BOM,
	CharsetUTF16BE,
	CharsetUTF16LE,
}

// decoders contains a map of all the character encoder/decoders known to go.
var decoders = map[string]encoding.Encoding{
	// Leave these in the order found in the referenced repositories, so if a
	// new encoding is added, it will be easier to spot.

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
	"maccyrillic": charmap.MacintoshCyrillic,
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

// encodingMap maps wlynxg/chardet's encoding name to a name used by go.
var encodingMap = map[string]string{
	consts.Ascii:   consts.ISO88591, // "Ascii"
	consts.UTF8SIG: consts.UTF8,     // "UTF-8-SIG"
	// consts.MacRoman -> "macintosh"

	// no matching decoder:
	// GB2312   = "GB2312"
	// Johab    = "Johab"
	// TIS620   = "TIS-620"
	// EucTw = "EUC-TW"
	// CP932 = "CP932"
	// CP949 = "CP949"
	// ISO2022CN = "ISO-2022-CN"
	// ISO2022KR = "ISO-2022-KR"
	// UCS43412  = "X-ISO-10646-UCS-4-3412"
	// UCS42143  = "X-ISO-10646-UCS-4-2143"
	// IBM866 = "IBM866"
}

// Allow tab (9), lf (10), ff (12), and cr (13)
var c0Chars = []byte{
	'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07',
	'\x08' /*tab*/ /*lf */, '\x0b' /*ff */ /*cr */, '\x0e', '\x0f',
	'\x10', '\x11', '\x12', '\x13', '\x14', '\x15', '\x16', '\x17',
	'\x18', '\x19', '\x1a', '\x1b', '\x1c', '\x1d', '\x1e', '\x1f',
}

var c1Chars = []byte{
	'\x80', '\x81', '\x82', '\x83', '\x84', '\x85', '\x86', '\x87',
	'\x88', '\x89', '\x8a', '\x8b', '\x8c', '\x8d', '\x8e', '\x8f',
	'\x90', '\x91', '\x92', '\x93', '\x94', '\x95', '\x96', '\x97',
	'\x98', '\x99', '\x9a', '\x9b', '\x9c', '\x9d', '\x9e', '\x9f',
}

var hibitChars = []byte{
	'\x80', '\x81', '\x82', '\x83', '\x84', '\x85', '\x86', '\x87',
	'\x88', '\x89', '\x8a', '\x8b', '\x8c', '\x8d', '\x8e', '\x8f',
	'\x90', '\x91', '\x92', '\x93', '\x94', '\x95', '\x96', '\x97',
	'\x98', '\x99', '\x9a', '\x9b', '\x9c', '\x9d', '\x9e', '\x9f',
	'\xa0', '\xa1', '\xa2', '\xa3', '\xa4', '\xa5', '\xa6', '\xa7',
	'\xa8', '\xa9', '\xaa', '\xab', '\xac', '\xad', '\xae', '\xaf',
	'\xb0', '\xb1', '\xb2', '\xb3', '\xb4', '\xb5', '\xb6', '\xb7',
	'\xb8', '\xb9', '\xba', '\xbb', '\xbc', '\xbd', '\xbe', '\xbf',
	'\xc0', '\xc1', '\xc2', '\xc3', '\xc4', '\xc5', '\xc6', '\xc7',
	'\xc8', '\xc9', '\xca', '\xcb', '\xcc', '\xcd', '\xce', '\xcf',
	'\xd0', '\xd1', '\xd2', '\xd3', '\xd4', '\xd5', '\xd6', '\xd7',
	'\xd8', '\xd9', '\xda', '\xdb', '\xdc', '\xdd', '\xde', '\xdf',
	'\xe0', '\xe1', '\xe2', '\xe3', '\xe4', '\xe5', '\xe6', '\xe7',
	'\xe8', '\xe9', '\xea', '\xeb', '\xec', '\xed', '\xee', '\xef',
	'\xf0', '\xf1', '\xf2', '\xf3', '\xf4', '\xf5', '\xf6', '\xf7',
	'\xf8', '\xf9', '\xfa', '\xfb', '\xfc', '\xfd', '\xfe', '\xff',
}

// CharsetsEqual checks the two charset names for equality.
func CharsetsEqual(charsetFound, charsetWanted string, bomEncoding string) bool {
	charsetFound = normalizeCharsetName(charsetFound)
	charsetWanted = normalizeCharsetName(charsetWanted)
	if charsetWanted == "utf8bom" && bomEncoding != consts.UTF8 {
		return false
	}
	if charsetFound == "utf8" && bomEncoding == consts.UTF8 {
		charsetFound += "bom"
	}

	// latin1 (iso88591) files are utf8 files, too.
	if charsetWanted == "utf8" && charsetFound == CharsetLatin1 {
		return true
	}

	return charsetFound == charsetWanted
}

// Decode attempts to determine a file's character encoding. If successful,
// it returns the content as a UTF-8 encoded string, and the name of the encoding.
// If not, it returns an error.
func Decode(contentBytes []byte) (string, string, error) {
	contentString := string(contentBytes)

	encoding, _, _ := Detect(contentBytes)

	decodedContentString, err := decodeText(contentBytes, encoding)
	if err != nil {
		if IsBinary(contentBytes) {
			return contentString, BinaryData, nil
		}

		return contentString, encoding, err
	}

	return decodedContentString, encoding, nil
}

// DecodeBytes is deprecated. Use Decode instead.
func DecodeBytes(contentBytes []byte) (string, string, error) {
	return Decode(contentBytes)
}

// Detect returns the character encoding, a confidence level, and the
// language.
func Detect(contentBytes []byte) (string, float64, string) {
	result := chardet.Detect(contentBytes)
	encoding := result.Encoding

	mapped, ok := encodingMap[encoding]
	if ok {
		encoding = mapped
	}

	if !IsSupportedCharset(encoding) {
		detected := DetectByBOM(contentBytes)
		if detected != "" {
			return detected, result.Confidence, result.Language
		}
	}

	if utf8.Valid(contentBytes) {
		if containsAnyByte(contentBytes, c0Chars) {
			// eg: ISO-2022-JP
			return encoding, result.Confidence, result.Language
		}

		if !containsAnyByte(contentBytes, hibitChars) {
			// eg: HZ-GB-2312
			return encoding, result.Confidence, result.Language
		}

		return consts.UTF8, result.Confidence, result.Language
	}

	if result.Confidence <= MinConfidenceForUTF1632Checks {
		if IsUTF32LE(contentBytes) {
			return consts.UTF32Le, result.Confidence, result.Language
		}
		if IsUTF32BE(contentBytes) {
			return consts.UTF32Be, result.Confidence, result.Language
		}
		if IsUTF16LE(contentBytes) {
			return consts.UTF16Le, result.Confidence, result.Language
		}
		if IsUTF16BE(contentBytes) {
			return consts.UTF16Be, result.Confidence, result.Language
		}
	}

	if strings.HasPrefix(strings.ToLower(encoding), "windows") {
		// Windows-* encodings should include some C1 chars.
		if !containsAnyByte(contentBytes, c1Chars) {
			return consts.ISO88591, result.Confidence, result.Language
		}
	}

	return encoding, result.Confidence, result.Language
}

type bomEntry struct {
	bom     []byte
	charset string
}

var (
	// https://en.wikipedia.org/wiki/Byte_order_mark
	UTF8BOM     = []byte{'\xEF', '\xBB', '\xBF'}
	UTF32LEBOM  = []byte{'\xFF', '\xFE', '\x00', '\x00'}
	UTF32BEBOM  = []byte{'\x00', '\x00', '\xFE', '\xFF'}
	UCS43412BOM = []byte{'\xFE', '\xFF', '\x00', '\x00'}
	UCS42143BOM = []byte{'\x00', '\x00', '\xFF', '\xFE'}
	UTF16LEBOM  = []byte{'\xFF', '\xFE'}
	UTF16BEBOM  = []byte{'\xFE', '\xFF'}

	// List in descending length order so the longest matches are attempted first.
	bomEntries = []bomEntry{
		{UTF8BOM, consts.UTF8},
		{UTF32LEBOM, consts.UTF32Le},
		{UTF32BEBOM, consts.UTF32Be},
		{UCS43412BOM, consts.UCS43412},
		{UCS42143BOM, consts.UCS42143},
		{UTF16LEBOM, consts.UTF16Le},
		{UTF16BEBOM, consts.UTF16Be},
	}
)

// DetectByBOM detects the file's encoding solely by BOM (byte order mark).
func DetectByBOM(contentBytes []byte) string {
	for _, entry := range bomEntries {
		if bytes.HasPrefix(contentBytes, entry.bom) {
			return entry.charset
		}
	}

	return ""
}

// IsBinary returns true if the bytes contain \x00-\x08,\x0b,\x0e-\x1f .
func IsBinary(rawFileContent []byte) bool {
	return containsAnyByte(rawFileContent, c0Chars)
}

// IsBinaryFile is deprecated. Use IsBinary instead.
func IsBinaryFile(rawFileContent []byte) bool {
	return IsBinary(rawFileContent)
}

var normalizedCharsets []string

func init() {
	normalizedCharsets = make([]string, len(ValidCharsets))

	for i, charset := range ValidCharsets {
		normalizedCharsets[i] = normalizeCharsetName(charset)
	}
}

// IsSupportedCharset checks charset is supported by the editorconfig spec.
func IsSupportedCharset(charset string) bool {
	charset = normalizeCharsetName(charset)
	return slices.Contains(normalizedCharsets, charset)
}

// IsUTF16BE returns true if the file is UTF16BE encoded.
func IsUTF16BE(b []byte) bool {
	if DetectByBOM(b) == consts.UTF16Be {
		return true
	}

	if len(b) < 2 {
		return false
	}

	if len(b)%2 != 0 {
		return false
	}

	hit := 0
	miss := 0
	for i := 0; i < len(b)-1; i += 2 {
		if b[i] == 0x00 && b[i+1] >= 0x20 && b[i+1] <= 0x7E {
			hit++
			continue
		}
		miss++
	}

	return float64(hit)/float64(miss) >= HitMissRatioForUTF1632Checks
}

// IsUTF16LE returns true if the file is UTF16LE encoded.
func IsUTF16LE(b []byte) bool {
	if DetectByBOM(b) == consts.UTF16Le {
		return true
	}

	if len(b) < 2 {
		return false
	}

	if len(b)%2 != 0 {
		return false
	}

	hit := 0
	miss := 0
	for i := 0; i < len(b)-1; i += 2 {
		if b[i+1] == 0x00 && b[i] >= 0x20 && b[i] <= 0x7E {
			hit++
			continue
		}
		miss++
	}

	return float64(hit)/float64(miss) >= HitMissRatioForUTF1632Checks
}

// IsUTF32BE returns true if the file is UTF32BE encoded.
func IsUTF32BE(b []byte) bool {
	if DetectByBOM(b) == consts.UTF32Be {
		return true
	}

	if len(b) < 4 {
		return false
	}

	if len(b)%4 != 0 {
		return false
	}

	hit := 0
	miss := 0
	for i := 0; i+3 < len(b); i += 4 {
		if b[i] == 0x00 && b[i+1] == 0x00 && b[i+2] == 0x00 &&
			b[i+3] >= 0x20 && b[i+3] <= 0x7E {
			hit++
			continue
		}
		miss++
	}

	return float64(hit)/float64(miss) >= HitMissRatioForUTF1632Checks
}

// IsUTF32LE returns true if the file is UTF32LE encoded.
func IsUTF32LE(b []byte) bool {
	if DetectByBOM(b) == consts.UTF32Le {
		return true
	}

	if len(b) < 4 {
		return false
	}

	if len(b)%4 != 0 {
		return false
	}

	hit := 0
	miss := 0
	for i := 0; i+3 < len(b); i += 4 {
		if b[i] >= 0x20 && b[i] <= 0x7E &&
			b[i+1] == 0x00 && b[i+2] == 0x00 && b[i+3] == 0x00 {
			hit++
			continue
		}
		miss++
	}

	return float64(hit)/float64(miss) >= HitMissRatioForUTF1632Checks
}

func decodeText(contentBytes []byte, encoding string) (string, error) {
	enc, ok := getDecoder(encoding)
	if !ok {
		return "", fmt.Errorf("unrecognized character encoding %q", encoding)
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

func getDecoder(encoding string) (encoding.Encoding, bool) {
	normalized := normalizeName(encoding)
	if normalized == "macroman" {
		normalized = "macintosh"
	}
	decoder, ok := decoders[normalized]

	return decoder, ok
}

func normalizeName(charset string) string {
	r := strings.NewReplacer("-", "", "_", "", ".", "")
	return strings.ToLower(r.Replace(charset))
}

var normalizeCharsetMap = map[string]string{
	"iso88591": CharsetLatin1,
}

func normalizeCharsetName(charset string) string {
	normalized := normalizeName(charset)

	mapped, ok := normalizeCharsetMap[normalized]
	if ok {
		return mapped
	}

	return normalized
}

func containsAnyByte(a, b []byte) bool {
	lookup := [256]bool{}

	for _, c := range b {
		lookup[c] = true
	}

	for _, c := range a {
		if lookup[c] {
			return true
		}
	}

	return false
}
