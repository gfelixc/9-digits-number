package logger

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Logger struct {
	file *os.File
}

var (
	exactNineDecimalDigitsNumber = regexp.MustCompile(`^[0-9]{9}$`)

	ErrNonExactDecimalDigitsNumber = errors.New("only decimal digits are allowed")
)

const (
	newlineSequence = "\n"
	logFilename     = "numbers.log"
)

func New() Logger {
	file, _ := os.Create(logFilename)

	return Logger{
		file: file,
	}
}

func (l Logger) OnlyNumbers(s string) error {
	if !exactNineDecimalDigitsNumber.MatchString(s) {
		return ErrNonExactDecimalDigitsNumber
	}

	leftZerosStripped := strings.TrimLeft(s, "0")

	_, err := l.file.Write([]byte(leftZerosStripped + newlineSequence))

	if err != nil {
		return fmt.Errorf("error writing the log file: %w", err)
	}

	return nil
}
