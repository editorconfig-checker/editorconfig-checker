package encoding

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"
	"text/tabwriter"
	"unicode/utf8"

	"github.com/wlynxg/chardet"
	"github.com/wlynxg/chardet/consts"
)

const minConfidenceToFailTests = 1

const minConfidenceToSetEncoding = 10

const defaultConfidence = 1

const testResultsJson = "test-results.json"

const addNewFilesEnvVar = "EDITORCONFIG_ADD_NEW_FILES"

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

var addNewFiles bool

var normalizedEncodings map[string]string

var encodingToCharsetMap map[string]string

func init() {
	addNewFiles = os.Getenv(addNewFilesEnvVar) != ""

	normalizedEncodings = make(map[string]string, len(encodings))
	for _, encoding := range encodings {
		normalized := normalizeName(encoding)
		normalizedEncodings[normalized] = encoding
	}

	encodingToCharsetMap = make(map[string]string, len(encodings))
	for _, encoding := range encodings {
		normalized := normalizeName(encoding)
		encodingToCharsetMap[normalized] = CharsetLatin1
	}

	encodingToCharsetMap[normalizeName(consts.UTF8)] = CharsetUTF8
	encodingToCharsetMap[normalizeName(consts.UTF8SIG)] = CharsetUTF8BOM
	encodingToCharsetMap[normalizeName(consts.UTF16)] = CharsetUTF16LE
	encodingToCharsetMap[normalizeName(consts.UTF16Be)] = CharsetUTF16BE
	encodingToCharsetMap[normalizeName(consts.UTF16Le)] = CharsetUTF16LE

	// We don't want these to default to latin1.
	encodingToCharsetMap[normalizeName(consts.UTF32)] = CharsetUTF32LE
	encodingToCharsetMap[normalizeName(consts.UTF32Be)] = CharsetUTF32BE
	encodingToCharsetMap[normalizeName(consts.UTF32Le)] = CharsetUTF32LE

	encodingToCharsetMap["utf8bom"] = CharsetUTF8BOM
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

// Do not reorder these tests, as they are in this order when
// EDITORCONFIG_ADD_NEW_FILES is set to generate new entries in
// test-results.json.

func TestCharsetsMatch(t *testing.T) {
	type test struct {
		charset1 string
		charset2 string
		want     bool
	}
	tests := []test{
		{CharsetLatin1, "latin-1", true},
		{CharsetLatin1, "iso88591", false},
		{CharsetLatin1, "iso8859-1", false},
		{CharsetLatin1, "iso-8859-1", false},
		{CharsetUTF8, "utf8", true},
		{CharsetUTF8BOM, "utf8bom", true},
		{CharsetUTF8BOM, "utf-8bom", true},
		{CharsetUTF8BOM, "utf8-bom", true},

		{CharsetUTF16BE, "utf16be", true},
		{CharsetUTF16BE, "utf16-be", true},
		{CharsetUTF16BE, "utf-16-be", true},
		{CharsetUTF16LE, "utf16le", true},
		{CharsetUTF16LE, "utf16-le", true},
		{CharsetUTF16LE, "utf-16-le", true},
		{CharsetUTF32LE, "utf32le", true},
		{CharsetUTF32BE, "utf32be", true},

		{CharsetUTF8, "utf8bom", false},
		{CharsetUTF8BOM, "utf8", false},

		{CharsetUTF16BE, "utf16be", true},
		{CharsetUTF16LE, "utf16le", true},

		{CharsetUTF16BE, "utf16le", false},
		{CharsetUTF16LE, "utf16be", false},

		{CharsetUTF16LE, "utf32le", false},
		{CharsetUTF16BE, "utf32be", false},
	}

	for _, charset1 := range ValidCharsets {
		for _, charset2 := range ValidCharsets {
			if charset1 != charset2 {
				if charset1 == "latin1" && charset2 == "utf-8" {
					continue
				}
				tests = append(tests, test{charset1, charset2, false})
			}
		}
	}

	charset2s := []string{"", BinaryData, UnknownEncoding, "latin", "utf", "utf32", "unset"}
	for _, charset2 := range charset2s {
		for _, charset1 := range ValidCharsets {
			tests = append(tests, test{charset1, charset2, false})
		}
	}
	for _, tt := range tests {
		got := CharsetsMatch(tt.charset1, tt.charset2)
		if got != tt.want {
			t.Errorf("CharsetsMatch(%q, %q): got %v, want %v", tt.charset1, tt.charset2, got, tt.want)
		}
	}
}

func TestDetectByBOM(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= minConfidenceToFailTests

		if !strings.Contains(strings.ToLower(tt.Filename), "bom") {
			continue
		}
		if strings.Contains(strings.ToLower(tt.Filename), "nobom") {
			continue
		}

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		encoding := DetectByBOM(fileContent)

		if !equal(encoding, tt.Encoding) {
			encoding2, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("DetectByBOM (%v): got %v, want %v", tt.Filename, encoding, tt.Encoding)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, encoding2 %v)", confidence, language, tt.Comment, encoding2)
			if failTest {
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					if encoding != "" {
						charset := encodingToCharset(encoding)
						if charset == CharsetUTF8 {
							charset = CharsetUTF8BOM
						}
						tests[i].Encoding = encoding
						tests[i].Charset = charset
						tests[i].Binary = IsBinary(fileContent)
						tests[i].Confidence = confidence
						t.Logf("DetectByBOM: setting test=%+v", tests[i])
						continue
					} else {
						tests[i].Comment += encoding2 + " (bom),"
						// There's no BOM, so keep testing.
					}
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}
	}
}

func TestIsUTF32BE(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= minConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		want := equal(tt.Encoding, consts.UTF32Be)
		got := IsUTF32BE(fileContent) >= HitMissRatioForUTF1632Checks

		if got != want {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsUTF32BE (%v): got %v, want %v", tt.Filename, got, want)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					if got {
						tests[i].Encoding = consts.UTF32Be
						tests[i].Charset = encodingToCharset(tests[i].Encoding)
						tests[i].Binary = true
						tests[i].Confidence = confidence
						t.Logf("IsUTF32BE: setting test=%+v", tests[i])
						continue
					} else {
						tests[i].Comment += "utf32be,"
					}
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}
	}
}

func TestIsUTF32LE(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= minConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		want := equal(tt.Encoding, consts.UTF32Le)
		got := IsUTF32LE(fileContent) >= HitMissRatioForUTF1632Checks

		if got != want {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsUTF32LE (%v): got %v, want %v", tt.Filename, got, want)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					if got {
						tests[i].Encoding = consts.UTF32Le
						tests[i].Charset = encodingToCharset(tests[i].Encoding)
						tests[i].Binary = true
						tests[i].Confidence = confidence
						t.Logf("IsUTF32LE: setting test=%+v", tests[i])
						continue
					} else {
						tests[i].Comment += "utf32le,"
					}
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}
	}
}

func TestIsUTF16BE(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= minConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		want := equal(tt.Encoding, consts.UTF16Be)
		got := IsUTF16BE(fileContent) >= HitMissRatioForUTF1632Checks

		if got != want {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsUTF16BE (%v): got %v, want %v", tt.Filename, got, want)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					if got {
						tests[i].Encoding = consts.UTF16Be
						tests[i].Charset = encodingToCharset(tests[i].Encoding)
						tests[i].Binary = true
						tests[i].Confidence = confidence
						t.Logf("IsUTF16BE: setting test=%+v", tests[i])
						continue
					} else {
						tests[i].Comment += "utf16be,"
					}
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}
	}
}

func TestIsUTF16LE(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= minConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		want := equal(tt.Encoding, consts.UTF16Le)
		got := IsUTF16LE(fileContent) >= HitMissRatioForUTF1632Checks

		if got != want {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsUTF16LE (%v): got %v, want %v", tt.Filename, got, want)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					if got {
						tests[i].Encoding = consts.UTF16Le
						tests[i].Charset = encodingToCharset(tests[i].Encoding)
						tests[i].Binary = true
						tests[i].Confidence = confidence
						t.Logf("IsUTF16LE: setting test=%+v", tests[i])
						continue
					} else {
						tests[i].Comment += "utf16le,"
					}
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}
	}
}

func TestIsBinary(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= minConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		binary := IsBinary(fileContent)

		if binary != tt.Binary {
			_, encoding, _ := Decode(fileContent)
			_, confidence, language := Detect(fileContent)

			msg := fmt.Sprintf("IsBinary (%v): got %v, want %v", tt.Filename, binary, tt.Binary)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					tests[i].Binary = true
					if confidence > minConfidenceToSetEncoding {
						tests[i].Encoding = encoding
						tests[i].Charset = encodingToCharset(tests[i].Encoding)
						tests[i].Confidence = confidence
						t.Logf("IsBinary: setting test=%+v", tests[i])
						continue
					} else {
						tests[i].Comment += "binary,"
					}
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}
	}
}

func TestDecode(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= minConfidenceToFailTests

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}
		_, encoding, err := Decode(fileContent)

		_, confidence, language := Detect(fileContent)

		errored := err != nil
		if errored != tt.Errored {
			errmsg := "nil"
			if err != nil {
				errmsg = err.Error()
			}
			msg := fmt.Sprintf("Decode (%v): got %v, want %v (%v)", tt.Filename, errored, tt.Errored, errmsg)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, encoding %v, tt.Confidene %f)", confidence, language, tt.Comment, encoding, tests[i].Confidence)

			if failTest {
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					tests[i].Encoding = encoding
					tests[i].Charset = encodingToCharset(tests[i].Encoding)
					tests[i].Errored = errored
					if tests[i].Encoding != tt.Encoding {
						tests[i].Comment = addEncodingToComment(tests[i].Comment, tt.Encoding, "!")
					}
					tests[i].Confidence = 0
					t.Logf("Decode: setting test=%+v", tests[i])
					continue
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}

		want := normalizeName(tt.Encoding)
		got := normalizeName(encoding)

		if want != got {
			msg := fmt.Sprintf("Decode2 (%v): got %v, want %v", tt.Filename, encoding, tt.Encoding)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, json.Confidence %f)", confidence, language, tt.Comment, tt.Confidence)

			if failTest {
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					tests[i].Encoding = encoding
					tests[i].Charset = encodingToCharset(tests[i].Encoding)
					if tests[i].Encoding != tt.Encoding {
						tests[i].Comment = addEncodingToComment(tests[i].Comment, tt.Encoding, "")
					}
					if confidence > minConfidenceToSetEncoding {
						tests[i].Confidence = confidence
						t.Logf("Decode2: setting test=%+v", tests[i])
						continue
					}
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}

		if want == got && confidence != tt.Confidence {
			// The test succeeded, but the confidence level changed.
			tests[i].Encoding = encoding
			tests[i].Charset = encodingToCharset(tests[i].Encoding)
			msg := fmt.Sprintf("Decode3 (%v): Confidence changed from %f to %f", tt.Filename, tt.Confidence, confidence)
			msg += fmt.Sprintf(" (encoding %v charset %v language %q, comment %q)", tests[i].Encoding, tests[i].Charset, language, tt.Comment)
			t.Log(msg)
			tests[i].Confidence = confidence
			t.Logf("Decode3: setting test=%+v", tests[i])
		}
	}
}

func TestDetect(t *testing.T) {
	for i, tt := range tests {
		failTest := tt.Confidence >= minConfidenceToFailTests

		if tt.Binary {
			continue
		}

		fileContent, err := readFile(tt.Filename)
		if err != nil {
			t.Fatalf("%s: %s", tt.Filename, err.Error())
		}

		encoding, confidence, language := Detect(fileContent)

		want := normalizeName(tt.Encoding)
		got := normalizeName(encoding)

		if addNewFiles && got == want {
			msg := fmt.Sprintf("GOOD: Detect: %-15s %-10s %6.2f %-15s %-20s %6.2f: %v",
				tt.Encoding, tt.Charset, confidence, language, tt.Comment, tt.Confidence, tt.Filename)
			t.Log(msg)
			if confidence > minConfidenceToSetEncoding {
				tests[i].Charset = encodingToCharset(encoding)
				tests[i].Confidence = confidence
				t.Logf("Detect: setting test=%+v", tests[i])
				continue
			}
		}

		if got != want {
			msg := fmt.Sprintf("Detect (%v): got %v, want %v", tt.Filename, encoding, tt.Encoding)
			msg += fmt.Sprintf(" (confidence %f, language %q, comment %q, encoding %v)", confidence, language, tt.Comment, encoding)
			if failTest {
				t.Log(msg)
				if addNewFiles && tests[i].Confidence == defaultConfidence {
					tests[i].Encoding = encoding
					tests[i].Charset = encodingToCharset(tests[i].Encoding)
					if tests[i].Encoding != tt.Encoding {
						tests[i].Comment = addEncodingToComment(tests[i].Comment, tt.Encoding, "@")
					}
					t.Logf("Detect2: confidence=%f", confidence)
					if confidence > minConfidenceToSetEncoding {
						tests[i].Confidence = confidence
						t.Logf("Detect: setting test=%+v", tests[i])
						continue
					} else {
						t.Logf("Detect2: confidence=%f", confidence)
						tests[i].Confidence = 0
						t.Logf("Detect2: setting test=%+v", tests[i])
					}
				}
				t.Error("FAIL: " + msg)
			} else {
				t.Log(msg)
			}
		}
	}
}

func TestPrintSupportedEncodings(t *testing.T) {
	if !addNewFiles {
		return
	}

	entries := encodings

	i := slices.Index(entries, consts.UTF8SIG)
	if i != -1 {
		entries = slices.Delete(entries, i, i+1)
	}

	for i := range entries {
		entries[i] = strings.ToLower(entries[i])
	}
	entries = append(entries, "latin1")
	entries = append(entries, "utf-8-bom")
	sort.Strings(entries)

	width := 0
	for _, entry := range entries {
		width = max(width, len(entry))
	}
	w := tabwriter.NewWriter(os.Stdout, width, width, 1, ' ', 0) // tabwriter.Debug)
	columns := 75 / width
	for i := 0; i < len(entries); i += columns {
		s := entries[i]
		for j := 1; j < columns; j++ {
			if i+j > len(entries)-1 {
				break
			}
			s += "\t" + entries[i+j]
		}
		fmt.Fprintln(w, s)
	}
	w.Flush()
}

func addEncodingToComment(comment, encoding, suffix string) string {
	if slices.Contains([]string{consts.Ascii, UnknownEncoding}, encoding) {
		return comment
	}

	if strings.Contains(comment, encoding) {
		if suffix != "" {
			return comment + suffix + ","
		}

		return comment
	}

	return comment + "Should be " + encoding + "?" + suffix + ","
}

func readFile(filename string) ([]byte, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return fileContent, nil
}

func setup() {
	if exists(testResultsJson) {
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
	}

	if !addNewFiles {
		if len(tests) == 0 {
			fmt.Printf("File not found: %q\n", testResultsJson)
		}
		return
	}

	for i := range tests {
		tests[i].Confidence = defaultConfidence
		tests[i].Comment = ""
	}

	_ = filepath.Walk("testdata", func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		filename = filepath.ToSlash(filename)

		dir, _ := path.Split(filename)

		if dir == "testdata/" {
			// There are no test files in the root of testdata/ any more.
			return nil
		}
		for _, test := range tests {
			if strings.HasSuffix(test.Filename, filename) {
				return nil
			}
		}

		normalized := normalizeName(filename)

		encoding, ok := normalizedEncodings[normalized]
		if encoding == "" {
			for k, v := range normalizedEncodings {
				if ok, _ := regexp.MatchString(k, normalized); ok {
					encoding = v
					break
				}
			}
		}

		if encoding == "" {
			for k := range decoders {
				if ok, _ := regexp.MatchString(k, normalized); ok {
					encoding = k
					break
				}
			}
		}

		if encoding == "" {
			encoding = UnknownEncoding // Will be updated in tests.
		}

		charset, ok := encodingToCharsetMap[encoding]
		if !ok {
			charset = UnknownEncoding // Will be updated in tests.
		}

		tests = append(tests, test{filename, encoding, charset, false, false, defaultConfidence, ""})

		return nil
	})

	for _, tt := range tests {
		if tt.Encoding == UnknownEncoding {
			fmt.Printf("%-12s %-10s: %s\n", tt.Encoding, tt.Charset, tt.Filename)
		}
	}
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
	if !addNewFiles {
		return
	}

	for i := range tests {
		tests[i].Comment = strings.TrimSuffix(tests[i].Comment, ",")
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

func dump(contentBytes []byte, t *testing.T) { //nolint:unused
	t.Helper()

	result := chardet.Detect(contentBytes)
	isValidUTF8 := utf8.Valid(contentBytes)
	t.Logf("chardet.Detect()=%+v utf8.Valid()=%v", result, isValidUTF8)

	c0Index := containsAnyByteIndex(contentBytes, c0Chars)
	c1Index := containsAnyByteIndex(contentBytes, c1Chars)
	hiIndex := containsAnyByteIndex(contentBytes, hiChars)
	t.Logf("length =0x%04x", len(contentBytes))
	t.Logf("c0Index=0x%04x", c0Index)
	t.Logf("c1Index=0x%04x", c1Index)
	t.Logf("hiIndex=0x%04x", hiIndex)
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
	if hiIndex >= 0 {
		last := min(hiIndex+dump, len(contentBytes))
		fmt.Printf("hiIndex: 0x%04x:\n%v\n", hiIndex, hex.Dump(contentBytes[hiIndex:last]))
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

func encodingToCharset(encoding string) string {
	encoding = normalizeName(encoding)

	charset, ok := encodingToCharsetMap[encoding]
	if ok {
		return charset
	}

	return CharsetLatin1
}

func equal(name1, name2 string) bool {
	return normalizeName(name1) == normalizeName(name2)
}
