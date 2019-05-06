// Package fix provides ...
package fix

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/editorconfig-checker/editorconfig-checker/pkg/utils"
)

func FinalNewline(filePath string, insert string, endOfLine string) error {
	if insert == "true" {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("Cant read file %s\n\n%s", filePath, err)
		}

		eolChar := utils.GetEolChar(endOfLine)
		if _, err := f.Write([]byte(eolChar)); err != nil {
			return fmt.Errorf("Cant write final newline to file %s\n\n%s", filePath, err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("Cant close file %s\n\n%s", filePath, err)
		}
	} else if insert == "false" {
		fmt.Println("Should remove final newline")

		input, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("Cant read file %s\n\n%s\n", filePath, err)
		}

		eolChar := utils.GetEolChar(endOfLine)
		output := strings.TrimRight(string(input), eolChar)

		err = ioutil.WriteFile(filePath, []byte(output), 0644)
		if err != nil {
			return fmt.Errorf("Cant write file %s\n\n%s\n", filePath, err)
		}

	}

	return nil
}

func LineEnding(filePath string, eolChar string) error {
	return nil
}

func TrailingWhitespace(filePath string, lineNumber int, eolChar string) error {
	// input, err := ioutil.ReadFile(filePath)
	// if err != nil {
	// 	return fmt.Errorf("Cant read file %s\n\n%s\n", filePath, err)
	// }

	// lines := strings.Split(string(input), eolChar)
	// lines[lineNumber] = strings.TrimRight(lines[lineNumber], " \t")
	// output := strings.Join(lines, "\n")

	// err = ioutil.WriteFile(filePath, []byte(output), 0644)
	// if err != nil {
	// 	return fmt.Errorf("Cant write file %s\n\n%s\n", filePath, err)
	// }

	return nil
}

func Indentation(filePath string, targetIndentation string, size int) error {
	return nil
}
