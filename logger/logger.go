package logger

import (
	"errors"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	logWriter io.Writer

	ticker            *time.Ticker
	pendingToWrite    []string
	pendingToWriteMux *sync.Mutex

	index    map[string]struct{}
	indexMux *sync.Mutex
}

var (
	exactNineDecimalDigitsNumber = regexp.MustCompile(`^[0-9]{9}$`)

	ErrDuplicatedNumber            = errors.New("duplicated number")
	ErrNonExactDecimalDigitsNumber = errors.New("only decimal digits are allowed")
	ErrTerminateSequenceDetected   = errors.New("terminate sequence detected")
)

const (
	terminateSequence   = "terminate"
	frequencyOfWritings = 1 * time.Second
	newlineSequence     = "\n"
)

// New returns a Logger which triggers a go routine to write the data in memory
// to the logWriter every (logger.frequencyOfWritings duration)
// In order to prevent data loss, Logger.Shutdown() MUST be called during graceful shutdown
func New(logWriter io.Writer) *Logger {
	logger := &Logger{
		ticker:            time.NewTicker(frequencyOfWritings),
		pendingToWriteMux: &sync.Mutex{},
		index:             make(map[string]struct{}),
		indexMux:          &sync.Mutex{},
		logWriter:         logWriter,
	}

	go logger.writeBatchInLogFileRecurrently()

	return logger
}

func (l *Logger) Shutdown() {
	l.ticker.Stop()
	l.writeBatch()
}

func (l *Logger) writeBatch() {
	l.pendingToWriteMux.Lock()
	defer l.pendingToWriteMux.Unlock()

	if len(l.pendingToWrite) == 0 {
		return
	}

	content := strings.Join(l.pendingToWrite, newlineSequence)

	_, err := l.logWriter.Write([]byte(content + newlineSequence))
	if err != nil {
		return
	}

	l.pendingToWrite = nil
}

func (l *Logger) addToBatch(s string) {
	l.pendingToWriteMux.Lock()
	l.pendingToWrite = append(l.pendingToWrite, s)
	l.pendingToWriteMux.Unlock()
}

func (l *Logger) OnlyNumbers(s string) error {
	if s == terminateSequence {
		return ErrTerminateSequenceDetected
	}

	if !exactNineDecimalDigitsNumber.MatchString(s) {
		return ErrNonExactDecimalDigitsNumber
	}

	leftZerosStripped := strings.TrimLeft(s, "0")

	l.indexMux.Lock()
	defer l.indexMux.Unlock()

	if _, ok := l.index[leftZerosStripped]; ok {
		return ErrDuplicatedNumber
	}

	l.addToBatch(s)

	l.index[leftZerosStripped] = struct{}{}

	return nil
}

func (l *Logger) writeBatchInLogFileRecurrently() {
	go func() {
		for _, channelIsOpen := <-l.ticker.C; channelIsOpen; {
			l.writeBatch()
		}
	}()
}
