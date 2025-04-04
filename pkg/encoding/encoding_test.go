package encoding

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/wlynxg/chardet"
	"github.com/wlynxg/chardet/consts"
)

const MinConfidenceToFailTests = 0.1

const DefaultConfidence = 50.0

const testResultsJson = "test-results.json"

var encodings = []string{
	// Do not reorder this list, so we can easily identify changes at
	// https://github.com/wlynxg/chardet/blob/204b5da7/consts/const.go#L16
	consts.Ascii,   // "Ascii"
	consts.UTF8,    // "UTF-8"
	consts.UTF8SIG, // "UTF-8-SIG"
	consts.UTF16,   // "UTF-16"
	consts.UTF16Le, // "UTF-16LE"
	consts.UTF16Be, // "UTF-16BE"
	consts.UTF32,   // "UTF-32"
	consts.UTF32Be, // "UTF-32BE"
	consts.UTF32Le, // "UTF-32LE"

	consts.GB2312,   // "GB2312"
	consts.HzGB2312, // "HZ-GB-2312"
	consts.ShiftJis, // "SHIFT_JIS"
	consts.Big5,     // "Big5"
	consts.Johab,    // "Johab"
	consts.Koi8R,    // "KOI8-R"
	consts.TIS620,   // "TIS-620"

	consts.MacCyrillic, // "MacCyrillic"
	consts.MacRoman,    // "MacRoman"

	consts.EucTw, // "EUC-TW"
	consts.EucKr, // "EUC-KR"
	consts.EucJp, // "EUC-JP"

	consts.CP932, // "CP932"
	consts.CP949, // "CP949"

	consts.Windows1250, // "Windows-1250"
	consts.Windows1251, // "Windows-1251"
	consts.Windows1252, // "Windows-1252"
	consts.Windows1253, // "Windows-1253"
	consts.Windows1254, // "Windows-1254"
	consts.Windows1255, // "Windows-1255"
	consts.Windows1256, // "Windows-1256"
	consts.Windows1257, // "Windows-1257"

	consts.ISO88591,  // "ISO-8859-1"
	consts.ISO88592,  // "ISO-8859-2"
	consts.ISO88595,  // "ISO-8859-5"
	consts.ISO88596,  // "ISO-8859-6"
	consts.ISO88597,  // "ISO-8859-7"
	consts.ISO88598,  // "ISO-8859-8"
	consts.ISO88599,  // "ISO-8859-9"
	consts.ISO885913, // "ISO-8859-13"
	consts.ISO2022CN, // "ISO-2022-CN"
	consts.ISO2022JP, // "ISO-2022-JP"
	consts.ISO2022KR, // "ISO-2022-KR"
	consts.UCS43412,  // "X-ISO-10646-UCS-4-3412"
	consts.UCS42143,  // "X-ISO-10646-UCS-4-2143"

	consts.IBM855, // "IBM855"
	consts.IBM866, // "IBM866"
}

var encodingToCharsetMap = map[string]string{
	consts.Ascii:    CharsetLatin1,
	consts.UTF8:     CharsetUTF8,
	consts.UTF8SIG:  CharsetUTF8BOM,
	consts.UTF16:    CharsetUTF16LE,
	consts.UTF16Le:  CharsetUTF16LE,
	consts.UTF16Be:  CharsetUTF16BE,
	consts.ISO88591: CharsetLatin1,
}

type test struct {
	Filename   string
	Encoding   string
	Charset    string
	Errored    bool
	Binary     bool
	Confidence float64
	Comment    string
}

var tests = []test{}

var normalizedNames map[string]string

func init() {
	normalizedNames = make(map[string]string, len(encodings))
	for _, encoding := range encodings {
		normalized := normalizeName(encoding)
		normalizedNames[normalized] = encoding
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestCharsetsEqual(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		_, encoding, _ := Decode(fileContent)
		charset := normalizeCharsetName(encoding)

		if tt.Charset != "" {
			bomEncoding := DetectByBOM(fileContent)

			if !CharsetsEqual(charset, tt.Charset, bomEncoding) {
				_, confidence, language := Detect(fileContent)

				msg := fmt.Sprintf("CharsetsEqual (%v): got %v, want %v", tt.Filename, charset, tt.Charset)
				msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)
				if failTest {
					t.Error(msg)
					// dumpInfo(fileContent, t)
				} else {
					t.Log(msg)
				}

				if tt.Comment == "" {
					tests[i].Comment = "unequal"
				}
			}
			continue
		}

		// tt.Charset == "", but charset is valid
		if IsSupportedCharset(charset) {
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("CharsetsEqual (%v): got %v, want %v", tt.Filename, charset, tt.Charset)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				t.Error(msg)
				// dumpInfo(fileContent, t)
			} else {
				t.Log(msg)
			}

			if tt.Comment == "" {
				tests[i].Comment = "unequal"
			}
		}
	}
}

func TestDecode(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}
		_, encoding, err := Decode(fileContent)

		want := normalizeName(tt.Encoding)
		got := normalizeName(encoding)

		_, confidence, language := Detect(fileContent)

		errored := err != nil
		if want != got {
			msg := fmt.Sprintf("Decode (%v): got %v, want %v", tt.Filename, encoding, tt.Encoding)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, json.Confidence %v)", confidence, language, tt.Comment, tt.Confidence)

			if failTest {
				t.Error(msg)
				// dumpInfo(fileContent, t)
			} else {
				t.Log(msg)
			}
			if tt.Comment == "" {
				tests[i].Comment = encoding
			}
		}

		if want == got && tt.Confidence == 0 {
			// The test succeeded, but the json file says the test should fail (Confidence == 0).
			msg := fmt.Sprintf("Decode (%v): got %v, want %v (Confidence == 0)", tt.Filename, encoding, tt.Encoding)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s)", confidence, language, tt.Comment)
			if failTest {
				t.Error(msg)
			} else {
				t.Log(msg)
			}
		}

		if errored != tt.Errored {
			errmsg := "nil"
			if err != nil {
				errmsg = err.Error()
			}
			msg := fmt.Sprintf("Decode (%v): got %v, want %v (%v)", tt.Filename, errored, tt.Errored, errmsg)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)

			if failTest {
				t.Error(msg)
				// dumpInfo(fileContent, t)
			} else {
				t.Log(msg)
			}

			if tt.Comment == "" {
				tests[i].Comment = "errored"
			}
		}
	}
}

func TestDetect(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		if tt.Binary {
			continue
		}

		encoding, confidence, language := Detect(fileContent)

		want := normalizeName(tt.Encoding)
		got := normalizeName(encoding)

		if got != want {
			msg := fmt.Sprintf("Detect (%v): got %v, want %v", tt.Filename, encoding, tt.Encoding)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				t.Error(msg)
			} else {
				t.Log(msg)
			}

			if tt.Comment == "" {
				tests[i].Comment = encoding
			}
		}
	}
}

func TestDetectByBOM(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		if !strings.Contains(strings.ToLower(tt.Filename), "bom") {
			continue
		}

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		encoding := DetectByBOM(fileContent)

		want := normalizeName(tt.Encoding)
		got := normalizeName(encoding)

		if got != want {
			encoding2, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("DetectByBOM (%v): got %v, want %v", tt.Filename, encoding, tt.Encoding)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding2 %v)", confidence, language, tt.Comment, encoding2)
			if failTest {
				t.Error(msg)
			} else {
				t.Log(msg)
			}

			if tt.Comment == "" {
				tests[i].Comment = encoding + " (bom)"
			}
		}
	}
}

func TestIsBinary(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		binary := IsBinary(fileContent)

		if binary != tt.Binary {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsBinary (%v): got %v, want %v", tt.Filename, binary, tt.Binary)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				t.Error(msg)
			} else {
				t.Log(msg)
			}
		}

		if tt.Comment == "" {
			tests[i].Comment = "binary"
		}

	}
}

func TestIsUTF16BE(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		charset := normalizeCharsetName(tt.Charset)

		want := charset == "utf16be"
		got := IsUTF16BE(fileContent)

		if got != want {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsUTF16BE (%v): got %v, want %v", tt.Filename, got, want)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				t.Error(msg)
			} else {
				t.Log(msg)
			}
		}

		if tt.Comment == "" {
			tests[i].Comment = "utf16be"
		}
	}
}

func TestIsUTF16LE(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		charset := normalizeCharsetName(tt.Charset)

		want := charset == "utf16le"
		got := IsUTF16LE(fileContent)

		if got != want {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsUTF16LE (%v): got %v, want %v", tt.Filename, got, want)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				t.Error(msg)
			} else {
				t.Log(msg)
			}
		}

		if tt.Comment == "" {
			tests[i].Comment = "utf16le"
		}
	}
}

func TestIsUTF32BE(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		filename := normalizeName(tt.Filename)
		want := strings.Contains(filename, "utf") && strings.Contains(filename, "32be")
		got := IsUTF32BE(fileContent)

		if got != want {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsUTF32BE (%v): got %v, want %v", tt.Filename, got, want)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				t.Error(msg)
			} else {
				t.Log(msg)
			}
		}

		if tt.Comment == "" {
			tests[i].Comment = "utf32be"
		}
	}
}

func TestIsUTF32LE(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= MinConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		filename := normalizeName(tt.Filename)
		want := strings.Contains(filename, "utf") && strings.Contains(filename, "32le")
		got := IsUTF32LE(fileContent)

		if got != want {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsUTF32LE (%v): got %v, want %v", tt.Filename, got, want)
			msg += fmt.Sprintf(" (confidence %v, language %v, comment %s, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				t.Error(msg)
			} else {
				t.Log(msg)
			}
		}

		if tt.Comment == "" {
			tests[i].Comment = "utf32le"
		}
	}
}

func TestIsSupportedCharset(t *testing.T) {
	type test struct {
		charset string
		want    bool
	}

	tests := []test{
		{CharsetLatin1, true},
		{CharsetUTF8, true},
		{CharsetUTF8BOM, true},
		{CharsetUTF16BE, true},
		{CharsetUTF16LE, true},
		{"latin-1", true},
		{"iso88591", true},
		{"iso8859-1", true},
		{"iso-8859-1", true},
		{"utf8", true},
		{"utf8bom", true},
		{"utf8-bom", true},
		{"utf-8bom", true},
		{"utf16be", true},
		{"utf16le", true},
		{"utf16-be", true},
		{"utf16-le", true},
		{"utf-16-be", true},
		{"utf-16-le", true},
		{"utf16", false},
		{"utf32", false},
		{"utf", false},
		{"unset", false},
		{BinaryData, false},
	}

	for _, tt := range tests {
		got := IsSupportedCharset(tt.charset)
		if got != tt.want {
			t.Errorf("IsSupportedCharset (%v): got %v, want %v", tt.charset, got, tt.want)
		}

		charset := strings.ToUpper(tt.charset)
		got = IsSupportedCharset(charset)
		if got != tt.want {
			t.Errorf("IsSupportedCharset (%v): got %v, want %v", tt.charset, got, tt.want)
		}
	}
}

func readFile(filename string) ([]byte, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return fileContent, nil
}

func setup() {
	f, err := os.Open(testResultsJson)
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(f).Decode(&tests)
	if err != nil {
		panic(err)
	}
	f.Close()
	sort.Slice(tests, func(i, j int) bool {
		return tests[i].Filename < tests[j].Filename
	})

	if os.Getenv("EDITORCONFIG_ADD_NEW_FILES") == "" {
		return
	}

	_ = filepath.Walk("testdata", func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		filename = filepath.ToSlash(filename)
		dir, base := path.Split(filename)

		if dir == "testdata/" {
			return nil
		}
		for _, test := range tests {
			if strings.HasSuffix(test.Filename, filename) {
				return nil
			}
		}

		base = base[:len(base)-len(filepath.Ext(base))]
		normalized := normalizeName(base)

		encoding, ok := normalizedNames[normalized]
		if !ok {
			encoding = consts.ISO88591
		}
		charset, ok := encodingToCharsetMap[encoding]
		if !ok {
			charset = CharsetLatin1
		}
		tests = append(tests, test{filename, encoding, charset, false, false, DefaultConfidence, ""})

		return nil
	})
}

func find(filename string) string { //nolint:unused
	if exists(filename) {
		return filepath.ToSlash(filename)
	}

	var found string

	_ = filepath.Walk("testdata", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == filename {
			found = path

			return filepath.SkipDir // stop walking
		}

		return nil
	})

	return filepath.ToSlash(found)
}

func exists(path string) bool { //nolint:unused
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func teardown() {
	if os.Getenv("EDITORCONFIG_ADD_NEW_FILES") == "" {
		return
	}

	sort.Slice(tests, func(i, j int) bool {
		return tests[i].Filename < tests[j].Filename
	})

	f, err := os.Create(testResultsJson)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	// Need 4 spaces to pass tests.
	enc.SetIndent("", "    ")
	_ = enc.Encode(tests)
}

func dumpInfo(contentBytes []byte, t *testing.T) { //nolint:unused
	t.Helper()

	result := chardet.Detect(contentBytes)
	isValidUTF8 := utf8.Valid(contentBytes)
	t.Logf("chardet.Detect()=%+v utf8.Valid()=%v", result, isValidUTF8)

	c0Index := containsAnyByteIndex(contentBytes, c0Chars)
	c1Index := containsAnyByteIndex(contentBytes, c1Chars)
	hibitIndex := containsAnyByteIndex(contentBytes, hibitChars)
	t.Logf("length    =0x%04x", len(contentBytes))
	t.Logf("c0Index   =0x%04x", c0Index)
	t.Logf("c1Index   =0x%04x", c1Index)
	t.Logf("hibitIndex=0x%04x", hibitIndex)
	dump := 256
	last := min(dump, len(contentBytes))
	fmt.Printf("Start:      0x%04x:\n%s\n", 0, hex.Dump(contentBytes[:last]))
	if c0Index >= 0 {
		last := min(c0Index+dump, len(contentBytes))
		fmt.Printf("c0Index:    0x%04x:\n%v\n", c0Index, hex.Dump(contentBytes[c0Index:last]))
	}
	if c1Index >= 0 {
		last := min(c1Index+dump, len(contentBytes))
		fmt.Printf("c1Index:    0x%04x:\n%v\n", c1Index, hex.Dump(contentBytes[c1Index:last]))
	}
	if hibitIndex >= 0 {
		last := min(hibitIndex+dump, len(contentBytes))
		fmt.Printf("hibitIndex: 0x%04x:\n%v\n", hibitIndex, hex.Dump(contentBytes[hibitIndex:last]))
	}
}

func containsAnyByteIndex(a, b []byte) int { //nolint:unused
	lookup := [256]bool{}
	for _, c := range b {
		lookup[c] = true
	}
	for i, c := range a {
		if lookup[c] {
			return i
		}
	}
	return -1
}
