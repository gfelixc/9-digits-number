package logger

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

type Logger struct {
	file     *os.File
	index    map[string]struct{}
	indexMux *sync.Mutex
}

var (
	exactNineDecimalDigitsNumber = regexp.MustCompile(`^[0-9]{9}$`)

	ErrDuplicatedNumber            = errors.New("duplicated number")
	ErrNonExactDecimalDigitsNumber = errors.New("only decimal digits are allowed")
)

const (
	newlineSequence = "\n"
	logFilename     = "numbers.log"
)

func New() Logger {
	file, _ := os.Create(logFilename)

	return Logger{
		index:    make(map[string]struct{}),
		indexMux: &sync.Mutex{},
		file:     file,
	}
}

func (l Logger) OnlyNumbers(s string) error {
	if !exactNineDecimalDigitsNumber.MatchString(s) {
		return ErrNonExactDecimalDigitsNumber
	}

	leftZerosStripped := strings.TrimLeft(s, "0")

	l.indexMux.Lock()
	defer l.indexMux.Unlock()

	if _, ok := l.index[leftZerosStripped]; ok {
		return ErrDuplicatedNumber
	}

	_, err := l.file.Write([]byte(leftZerosStripped + newlineSequence))

	if err != nil {
		return fmt.Errorf("error writing the log file: %w", err)
	}

	l.index[leftZerosStripped] = struct{}{}

	return nil
}
