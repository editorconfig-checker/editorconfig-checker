package encoding

import (
	"os"
	"reflect"
	"testing"

	"github.com/baulk/chardet"
)

type decodeTest struct {
	filename string
	charset  string
	changed  bool
	errored  bool
}

var decodeTextTests = []decodeTest{
	// Encodings
	{"8859_1_da.html", "ISO-8859-1", false, false},
	{"8859_1_de.html", "ISO-8859-1", false, false},
	{"8859_1_en.html", "ISO-8859-1", false, false},
	{"8859_1_es.html", "ISO-8859-1", false, false},
	{"8859_1_fr.html", "ISO-8859-1", false, false},
	{"8859_1_pt.html", "ISO-8859-1", false, false},
	{"ascii.txt", "ISO-8859-1", false, false},
	{"big5.html", "Big5", false, false},
	{"candide-gb18030.txt", "windows-1252", false, false}, // should be GB18030
	{"candide-utf-16le.txt", "ISO-8859-1", false, false},  // should be UTF-16LE
	{"candide-utf-32be.txt", "UTF-32BE", false, false},
	{"candide-utf-8.txt", "UTF-8", false, false},
	{"candide-windows-1252.txt", "ISO-8859-1", false, false}, // should be windows-1252
	{"cp865.txt", "windows-1252", false, false},
	{"euc_jp.html", "EUC-JP", false, false},
	{"euc_kr.html", "EUC-KR", false, false},
	{"gb18030.html", "GB18030", false, false},
	{"html.html", "ISO-8859-1", false, false},
	{"html.iso88591.html", "ISO-8859-1", false, false},
	{"html.svg.html", "ISO-8859-1", false, false},
	{"html.usascii.html", "ISO-8859-1", false, false},
	{"html.utf8bomdetect.html", "UTF-8", false, false},
	{"html.utf8bom.html", "UTF-8", false, false},
	{"html.utf8bomws.html", "UTF-8", false, false},
	{"html.utf8.html", "ISO-8859-1", false, false}, // should be UTF-8 ?
	{"html.withbr.html", "ISO-8859-1", false, false},
	{"iso88591.txt", "ISO-8859-1", false, false},
	{"koi8_r.txt", "ISO-8859-1", false, false}, // should be KOI8-R
	{"latin1.txt", "ISO-8859-1", false, false},
	{"rashomon-euc-jp.txt", "EUC-JP", false, false},
	{"rashomon-iso-2022-jp.txt", "ISO-2022-JP", false, false},
	{"rashomon-shift-jis.txt", "Shift_JIS", false, false},
	{"rashomon-utf-8.txt", "UTF-8", false, false},
	{"shift_jis.html", "Shift_JIS", false, false},
	{"sunzi-bingfa-gb-levels-1-and-2-hz-gb2312.txt", "Big5", false, false}, // should be GB18030
	{"sunzi-bingfa-gb-levels-1-and-2-utf-8.txt", "UTF-8", false, false},
	{"sunzi-bingfa-simplified-gbk.txt", "GB18030", false, false},
	{"sunzi-bingfa-simplified-utf-8.txt", "UTF-8", false, false},
	{"sunzi-bingfa-traditional-big5.txt", "Big5", false, false},
	{"sunzi-bingfa-traditional-utf-8.txt", "UTF-8", false, false},
	{"unsu-joh-eun-nal-euc-kr.txt", "EUC-KR", false, false},
	{"unsu-joh-eun-nal-utf-8.txt", "UTF-8", false, false},
	{"utf16bebom.txt", "UTF-16BE", false, false},
	{"utf16lebom.txt", "UTF-16LE", false, false},
	{"utf16.txt", "UTF-16LE", false, false},
	{"utf32bebom.txt", "UTF-32BE", false, false},
	{"utf32lebom.txt", "UTF-32LE", false, false},
	{"utf8_bom.html", "UTF-8", false, false},
	{"utf8.html", "UTF-8", false, false},
	{"utf8-sdl.txt", "windows-1252", false, false}, // should be UTF-8
	{"utf8.txt", "UTF-8", false, false},
	{"utf8.txt-encoding-test-files.txt", "UTF-8", false, false},
	// Issues
	{"issue252.txt", "UTF-8", false, false},
}

func TestDecodeBytesText(t *testing.T) {
	for _, tt := range decodeTextTests {
		filename := "testdata/" + tt.filename
		rawFileContent, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		// decodedFileContent, charset, err := DecodeBytes(rawFileContent)
		decodedFileContent, charset, err := DecodeBytes(rawFileContent)
		changed := reflect.DeepEqual(rawFileContent, decodedFileContent)
		errored := err != nil
		if charset != tt.charset {
			t.Errorf("DecodeBytes(%s): charset: expected: %v, got: %v (changed=%v, errored=%v)", tt.filename, tt.charset, charset, changed, errored)
			listAllCharsets(t, rawFileContent)
		}
		if changed != tt.changed {
			t.Errorf("DecodeBytes(%s): content changed: expected: %v, got: %v (charset=%v, errored=%v)", tt.filename, tt.changed, changed, charset, errored)
		}
		if errored != tt.errored {
			t.Errorf("DecodeBytes(%s): errored: expected: %v, got: %v (err=%v, charset=%v, changed=%v)", tt.filename, tt.errored, errored, err.Error(), charset, changed)
		}
	}
}

func listAllCharsets(t *testing.T, contentBytes []byte) {
	detector := chardet.NewTextDetector()
	results, err := detector.DetectAll(contentBytes)
	if err != nil {
		t.Logf("  err=%s", err)
		return
	}
	for i, result := range results {
		t.Logf("  result[%d]=%+v", i, result)
	}
}

type isBinaryFileTest struct {
	filename string
	binary   bool
}

var isBinaryFileTests = []isBinaryFileTest{
	{"8859_1_da.html", false},
	{"8859_1_de.html", false},
	{"8859_1_en.html", false},
	{"8859_1_es.html", false},
	{"8859_1_fr.html", false},
	{"8859_1_pt.html", false},
	{"ascii.txt", false},
	{"big5.html", false},
	{"candide-gb18030.txt", false},
	{"candide-utf-16le.txt", true},
	{"candide-utf-32be.txt", true},
	{"candide-utf-8.txt", false},
	{"candide-windows-1252.txt", false},
	{"cp865.txt", false},
	{"euc_jp.html", false},
	{"euc_kr.html", false},
	{"gb18030.html", false},
	{"html.html", false},
	{"html.iso88591.html", false},
	{"html.svg.html", false},
	{"html.usascii.html", false},
	{"html.utf8bomdetect.html", false},
	{"html.utf8bom.html", false},
	{"html.utf8bomws.html", false},
	{"html.utf8.html", false},
	{"html.withbr.html", false},
	{"iso88591.txt", false},
	{"koi8_r.txt", false},
	{"latin1.txt", false},
	{"rashomon-euc-jp.txt", false},
	{"rashomon-iso-2022-jp.txt", true}, // byte 89 is an Esc (ASCII 27)
	{"rashomon-shift-jis.txt", false},
	{"rashomon-utf-8.txt", false},
	{"shift_jis.html", false},
	{"sunzi-bingfa-gb-levels-1-and-2-hz-gb2312.txt", false},
	{"sunzi-bingfa-gb-levels-1-and-2-utf-8.txt", false},
	{"sunzi-bingfa-simplified-gbk.txt", false},
	{"sunzi-bingfa-simplified-utf-8.txt", false},
	{"sunzi-bingfa-traditional-big5.txt", false},
	{"sunzi-bingfa-traditional-utf-8.txt", false},
	{"unsu-joh-eun-nal-euc-kr.txt", false},
	{"unsu-joh-eun-nal-utf-8.txt", false},
	{"utf16bebom.txt", true},
	{"utf16lebom.txt", true},
	{"utf16.txt", true},
	{"utf32bebom.txt", true},
	{"utf32lebom.txt", true},
	{"utf8_bom.html", false},
	{"utf8.html", false},
	{"utf8-sdl.txt", false},
	{"utf8.txt", false},
	{"utf8.txt-encoding-test-files.txt", false},
}

func TestIsBinaryFile(t *testing.T) {
	for _, tt := range isBinaryFileTests {
		filename := "testdata/" + tt.filename
		rawFileContent, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		// decodedFileContent, charset, err := DecodeBytes(rawFileContent)
		binary := IsBinaryFile(rawFileContent)
		if binary != tt.binary {
			t.Errorf("IsBinaryFile(%s): expected: %v, got: %v", tt.filename, tt.binary, binary)
		}
	}
}
