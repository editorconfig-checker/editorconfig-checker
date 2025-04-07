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

const (
	// HitMissRatioForUTF1632Checks defines the ratio of hits to misses that must
	// be met for the IsUTF* functions to return true. For example, the number 1.1
	// requires the function find 10% more hits than misses. This number was chosen
	// empiricallly, but may need to be adjusted as we gather more test results.
	HitMissRatioForUTF1632Checks = 1.01

	// HitMissRatioIfBOMFound is returned if a file begins with a BOM
	// (Byte Order Mark).
	HitMissRatioIfBOMFound = 1023.0

	// MinConfidenceForUTF1632Checks defines the minimum confidence factor to use
	// our own check for if a file is UTF16/UTF32 encoded. 100 represents
	// a 100% confidence level, which is the highest confidence returned by our
	// chardet library.
	// testdata/wpt/resources/utf-32-big-endian-nobom.html is an 85.0.
	MinConfidenceForUTF1632Checks = 90.0

	// BinaryData contains the string to return if the data contains binary
	// data (data that is not decodable, by any of the decoders that golang provides).
	BinaryData = "binary"

	// UnknownEncoding is returned if the encoding could not be determined.
	UnknownEncoding = "unknown"

	// See https://spec.editorconfig.org/#supported-pairs
	// CharsetUnset defines the value allowing for file encoding.
	CharsetUnset = "unset"
	// CharsetLatin1 defines the value for the ISO-8859-1 file encoding.
	CharsetLatin1 = "latin1"
	// CharsetUTF8 defines the value for the UTF-8 encoding, without a
	// Byte Order Mark (BOM).
	CharsetUTF8 = "utf-8"
	// CharsetUTF8BOM defines the value for the UTF-8 encoding, with a
	// Byte Order Mark (BOM).
	CharsetUTF8BOM = "utf-8-bom"
	// CharsetUTF16BE defines the value for the UTF-16BE (Big Endian) encoding
	// with, or without, a Byte Order Mark (BOM).
	CharsetUTF16BE = "utf-16be"
	// CharsetUTF16LE defines the value for the UTF-16LE (Lil Endian) encoding
	// with, or without, a Byte Order Mark (BOM).
	CharsetUTF16LE = "utf-16le"

	// CharsetUTF32BE defines the value for the UTF-32BE (Big Endian) encoding
	// with, or without, a Byte Order Mark (BOM).
	CharsetUTF32BE = "utf-32be"
	// CharsetUTF32LE defines the value for the UTF-32LE (Lil Endian) encoding
	// with, or without, a Byte Order Mark (BOM).
	CharsetUTF32LE = "utf-32le"
)

var (
	// ValidCharsets contains an array of the allowable values for the charset key
	// in an .editorconfig file.
	ValidCharsets = []string{
		CharsetLatin1,
		CharsetUTF8,
		CharsetUTF8BOM,
		CharsetUTF16BE,
		CharsetUTF16LE,
	}

	// decoders contains a map of all the character encoder/decoders known to go.
	decoders = map[string]encoding.Encoding{
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
		"utf8bom": unicode.UTF8BOM,
		"utf32be": utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM),
		"utf32le": utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM),
	}

	// In https://github.com/golang/text/blob/HEAD/encoding/ianaindex/ianaindex.go#L156
	// but not included above:
	// 	 enc3:    asciiEnc,
	//   enc1015: unicode.UTF16(unicode.BigEndian, unicode.UseBOM)

	// See https://en.wikipedia.org/wiki/C0_and_C1_control_codes
	c0Chars = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		// Allow TAB (ASCII 9), LF (10), FF (12), and CR (13)
		0x08 /*TAB  LF*/, 0x0b /*FF   CR*/, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}

	c1Chars = []byte{
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97,
		0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
	}

	hiChars = []byte{
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97,
		0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
		0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7,
		0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
		0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,
		0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf,
		0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7,
		0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf,
		0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7,
		0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf,
		0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7,
		0xe8, 0xe9, 0xea, 0xeb, 0xec, 0xed, 0xee, 0xef,
		0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7,
		0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff,
	}

	supportedUTFEncodings = []string{
		consts.UTF8,    // "UTF-8"
		consts.UTF8SIG, // "UTF-8-SIG"
		consts.UTF16,   // "UTF-16"
		consts.UTF16Le, // "UTF-16LE"
		consts.UTF16Be, // "UTF-16BE"
		"utf8bom",
	}

	utfEncodingMap map[string]string

	// encodingToDecoderMap maps wlynxg/chardet's encoding name to a name used by go.
	encodingToDecoderMap map[string]string

	supportedUTFEncodingMap map[string]string

	latin1Encodings = []string{
		CharsetLatin1,
		"ascii",
		"iso88591",
	}
)

func init() {
	encodingToDecoderMap = make(map[string]string, 3)

	encodingToDecoderMap[normalizeName(consts.Ascii)] = "iso88591"  // "Ascii"
	encodingToDecoderMap[normalizeName(consts.UTF8SIG)] = "utf8bom" // "UTF-8-SIG"
	encodingToDecoderMap[normalizeName(consts.MacRoman)] = "macintosh"

	supportedUTFEncodingMap = make(map[string]string, len(supportedUTFEncodings))
	for _, k := range supportedUTFEncodings {
		supportedUTFEncodingMap[normalizeName(k)] = k
	}

	// encodings with no matching decoder:
	// CP932     = "CP932"
	// CP949     = "CP949"
	// EucTw     = "EUC-TW"
	// GB2312    = "GB2312"
	// IBM866    = "IBM866"
	// ISO2022CN = "ISO-2022-CN"
	// ISO2022KR = "ISO-2022-KR"
	// Johab     = "Johab"
	// TIS620    = "TIS-620"
	// UCS42143  = "X-ISO-10646-UCS-4-2143"
	// UCS43412  = "X-ISO-10646-UCS-4-3412"
}

// CharsetsMatch checks the charset found is effectively a match for the
// desired charset. For now, all non-UTF8/16 encodings will match `latin1`
// (at least until we can improve our success rate on encoding detection).
// Per https://github.com/editorconfig-checker/editorconfig-checker/pull/457#issuecomment-2779587476
func CharsetsMatch(charsetFound, charsetWanted string) bool {
	charsetFound = normalizeCharsetName(charsetFound)
	charsetWanted = normalizeCharsetName(charsetWanted)

	if charsetFound == "utf8sig" {
		charsetFound = "utf8bom"
	}

	if !Supported(charsetFound) {
		if charsetWanted == CharsetLatin1 {
			return true
		}
	}

	// latin1 (iso88591) files are utf8 files, too.
	if charsetWanted == "utf8" && slices.Contains(latin1Encodings, charsetFound) {
		return true
	}

	return charsetFound == charsetWanted
}

// Decode attempts to determine a file's character encoding.
// If successful, it returns the content as a UTF-8 encoded string, and the
// name of the encoding.
// If not, it returns the content as a string, the encoding name, and an error.
func Decode(contentBytes []byte) (string, string, error) {
	encoding, _, _ := Detect(contentBytes)

	decodedContentString, err := decodeText(contentBytes, encoding)
	if err != nil {
		if IsBinary(contentBytes) {
			return string(contentBytes), BinaryData, nil
		}

		return string(contentBytes), encoding, err
	}

	return decodedContentString, encoding, nil
}

// DecodeBytes is deprecated and may be removed in the future.
// Use Decode instead.
func DecodeBytes(contentBytes []byte) (string, string, error) {
	return Decode(contentBytes)
}

// Detect returns the character encoding, a confidence level, and the
// language.
func Detect(contentBytes []byte) (string, float64, string) {
	result := chardet.Detect(contentBytes)
	encoding := result.Encoding
	confidence := 100.0 * result.Confidence

	if encoding == "" {
		encoding = UnknownEncoding
	}

	for {
		detected := DetectByBOM(contentBytes)
		if detected != "" {
			encoding = detected
			break
		}

		if encoding == consts.Ascii {
			break
		}

		// We need to check for UTF16/32 encodings first, as
		// UTF16/32 encoded first can be valid UTF8 files (surprisingly).
		// For example, without the logic below,
		// testdata/wpt/resources/utf-32-big-endian-nobom.html
		// is reported to be UTF-8 with a .85 (85%) confidence level.
		if confidence <= MinConfidenceForUTF1632Checks {
			if IsUTF32LE(contentBytes) >= HitMissRatioForUTF1632Checks {
				encoding = consts.UTF32Le
				break
			}
			if IsUTF32BE(contentBytes) >= HitMissRatioForUTF1632Checks {
				encoding = consts.UTF32Be
				break
			}
			if IsUTF16LE(contentBytes) >= HitMissRatioForUTF1632Checks {
				encoding = consts.UTF16Le
				break
			}
			if IsUTF16BE(contentBytes) >= HitMissRatioForUTF1632Checks {
				encoding = consts.UTF16Be
				break
			}
		}

		hasC0 := containsAnyByte(contentBytes, c0Chars)
		hasC1 := containsAnyByte(contentBytes, c1Chars)
		hasHI := containsAnyByte(contentBytes, hiChars)

		if utf8.Valid(contentBytes) {
			if !hasC0 && !hasC1 && !hasHI {
				encoding = consts.Ascii
				break
			}

			if !hasHI {
				encoding = consts.ISO88591
				break
			}

			// We prioritize identifying UTF-8 over non-UTF8 encodings,
			// so let's skip these checks.
			//
			// if containsAnyByte(contentBytes, c0Chars) {
			// 	// eg: ISO-2022-JP
			// 	break
			// }

			// if !containsAnyByte(contentBytes, hiChars) {
			// 	// eg: HZ-GB-2312
			// 	break
			// }

			encoding = consts.UTF8
			break
		}

		if strings.HasPrefix(strings.ToLower(encoding), "windows") && !hasC1 {
			// This is a false positive, as all Windows-* encodings include
			// C1 chars.
			encoding = consts.ISO88591
			break
		}

		break
	}

	return encoding, confidence, result.Language
}

type bomEntry struct {
	bom     []byte
	charset string
}

var (
	// https://en.wikipedia.org/wiki/Byte_order_mark
	UTF32LEBOM   = []byte{0xFF, 0xFE, 0x00, 0x00}
	UTF32BEBOM   = []byte{0x00, 0x00, 0xFE, 0xFF}
	UCS43412BOM  = []byte{0xFE, 0xFF, 0x00, 0x00}
	UCS42143BOM  = []byte{0x00, 0x00, 0xFF, 0xFE}
	UTFEBCDICBOM = []byte{0xDD, 0x73, 0x66, 0x73}
	GB18030BOM   = []byte{0x84, 0x31, 0x95, 0x33}
	UTF8BOM      = []byte{0xEF, 0xBB, 0xBF}
	UTF7BOM      = []byte{0x74, 0x64, 0x4C}
	UTF1BOM      = []byte{0xF7, 0x64, 0x4C}
	SCSUBOM      = []byte{0x0E, 0xFE, 0xFF}
	BOCU1BOM     = []byte{0xFB, 0xEE, 0x28}
	UTF16LEBOM   = []byte{0xFF, 0xFE}
	UTF16BEBOM   = []byte{0xFE, 0xFF}

	// List in descending length order so the longest matches are attempted first.
	bomEntries = []bomEntry{
		{UTF8BOM, consts.UTF8SIG},
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

// IsBinaryFile is deprecated and may be removed in the future.
// Use IsBinary instead.
func IsBinaryFile(rawFileContent []byte) bool {
	return IsBinary(rawFileContent)
}

// IsUTF16BE returns a hit/miss ratio to identify if the file is UTF16BE encoded.
func IsUTF16BE(b []byte) float64 {
	if DetectByBOM(b) == consts.UTF16Be {
		return HitMissRatioIfBOMFound
	}

	if len(b) < 2 {
		return 0.0
	}

	if len(b)%2 != 0 {
		return 0.0
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

	return float64(hit) / float64(miss)
}

// IsUTF16LE returns a hit/miss ratio to identify if the file is UTF16LE encoded.
func IsUTF16LE(b []byte) float64 {
	if DetectByBOM(b) == consts.UTF16Le {
		return HitMissRatioIfBOMFound
	}

	if len(b) < 2 {
		return 0.0
	}

	if len(b)%2 != 0 {
		return 0.0
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

	return float64(hit) / float64(miss)
}

// IsUTF32BE returns a hit/miss ratio to identify if the file is UTF32BE encoded.
func IsUTF32BE(b []byte) float64 {
	if DetectByBOM(b) == consts.UTF32Be {
		return HitMissRatioIfBOMFound
	}

	if len(b) < 4 {
		return 0.0
	}

	if len(b)%4 != 0 {
		return 0.0
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

	return float64(hit) / float64(miss)
}

// IsUTF32LE returns hit/miss ratio to identify if the file is UTF32LE encoded.
func IsUTF32LE(b []byte) float64 {
	if DetectByBOM(b) == consts.UTF32Le {
		return HitMissRatioIfBOMFound
	}

	if len(b) < 4 {
		return 0.0
	}

	if len(b)%4 != 0 {
		return 0.0
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

	return float64(hit) / float64(miss)
}

func Supported(encoding string) bool {
	normalized := normalizeName(encoding)
	_, ok := supportedUTFEncodingMap[normalized]

	return ok
}

// UnrecogizedEncodingError is returned if the encountered a character set
// we don't have a decoder for.
type UnrecogizedEncodingError struct {
	encoding string
}

func (e *UnrecogizedEncodingError) Error() string {
	return fmt.Sprintf("unrecognized character encoding %q", e.encoding)
}

func decodeText(contentBytes []byte, encoding string) (string, error) {
	enc, ok := getDecoder(encoding)
	if !ok {
		return string(contentBytes), &UnrecogizedEncodingError{encoding}
	}

	validUTF8, err := enc.NewDecoder().Bytes(contentBytes)
	if err != nil {
		return "", err
	}

	return string(validUTF8), nil
}

func getDecoder(encoding string) (encoding.Encoding, bool) {
	normalized := normalizeName(encoding)

	decoder, ok := decoders[normalized]

	if !ok {
		normalized, _ = encodingToDecoderMap[normalized]
	}

	decoder, ok = decoders[normalized]

	return decoder, ok
}

func normalizeName(name string) string {
	r := strings.NewReplacer("-", "", "_", "", ".", "")
	return strings.ToLower(r.Replace(name))
}

func normalizeCharsetName(charset string) string {
	normalized := normalizeName(charset)

	if normalized == "utf8sig" {
		normalized = "utf8bom"
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
