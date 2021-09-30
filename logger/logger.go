package logger

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	file *os.File

	writerBatch    []string
	writerBatchMux *sync.Mutex

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

	logger := Logger{
		writerBatchMux: &sync.Mutex{},
		index:          make(map[string]struct{}),
		indexMux:       &sync.Mutex{},
		file:           file,
	}

	go logger.writeBatchInFile()

	return logger
}

func (l *Logger) Flush() error {
	l.writeBatch()

	if err := l.file.Sync(); err != nil {
		return err
	}

	return l.file.Close()
}

func (l *Logger) writeBatch() {
	l.writerBatchMux.Lock()
	defer l.writerBatchMux.Unlock()

	content := strings.Join(l.writerBatch, newlineSequence)

	_, err := l.file.Write([]byte(content))
	if err != nil {
		println(err.Error())
		return
	}

	l.writerBatch = nil

	return
}

func (l *Logger) addToBatch(s string) error {
	l.writerBatchMux.Lock()
	l.writerBatch = append(l.writerBatch, s)
	l.writerBatchMux.Unlock()

	return nil
}

func (l *Logger) OnlyNumbers(s string) error {
	if !exactNineDecimalDigitsNumber.MatchString(s) {
		return ErrNonExactDecimalDigitsNumber
	}

	leftZerosStripped := strings.TrimLeft(s, "0")

	l.indexMux.Lock()
	defer l.indexMux.Unlock()

	if _, ok := l.index[leftZerosStripped]; ok {
		return ErrDuplicatedNumber
	}

	err := l.addToBatch(s)
	if err != nil {
		return err
	}
	// _, err := l.file.Write([]byte(leftZerosStripped + newlineSequence))

	// if err != nil {
	// 	return fmt.Errorf("error writing the log file: %w", err)
	// }

	l.index[leftZerosStripped] = struct{}{}

	return nil
}

func (l *Logger) writeBatchInFile() {
	go func() {
		timer := time.Tick(1 * time.Second)
		for {
			<- timer
			l.writeBatch()
		}
	}()
}
