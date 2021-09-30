package logger_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/gfelixc/gigapipe/logger"
	"github.com/stretchr/testify/require"
)

var logFilename = "numbers.log"

func tearDown() {
	err := os.Remove(logFilename)
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	if err != nil {
		println("unable to remove log file: ", err.Error())
	}
}

func Test200KNumbersPerSecond(t *testing.T) {
	t.Cleanup(tearDown)

	timeout := time.After(1 * time.Second)
	l := logger.New()

	go func() {
		for i := 0; i <= 200_000; i++ {
			_ = l.OnlyNumbers(fmt.Sprintf("%09d", i))
		}
	}()

	<-timeout

	require.NoError(t, l.Flush())

	linesInLogFile, err := countLinesInLogFile()
	require.NoError(t, err)

	require.Equal(t, 200_000, linesInLogFile)
}

func TestInstanceLoggerCreatesALogFile(t *testing.T) {
	t.Cleanup(tearDown)

	require.False(t, logFileExists())

	logger.New()

	require.True(t, logFileExists())
}

func TestInstanceLoggerClearsAnExistentLogFile(t *testing.T) {
	t.Cleanup(tearDown)

	file, err := createLogFilePreFilled()
	require.NoError(t, err)

	logger.New()

	stat, err := file.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(0), stat.Size())
}

func TestOnlyNumbersMayBeWrittenToTheLogFile(t *testing.T) {
	t.Cleanup(tearDown)

	l := logger.New()

	t.Run("Alphanumeric", func(t *testing.T) {
		err := l.OnlyNumbers("Aca2321")
		require.ErrorIs(t, err, logger.ErrNonExactDecimalDigitsNumber)
	})

	t.Run("SpecialChars", func(t *testing.T) {
		err := l.OnlyNumbers("123 12-3")
		require.ErrorIs(t, err, logger.ErrNonExactDecimalDigitsNumber)
	})

	t.Run("OnlyDecimals", func(t *testing.T) {
		err := l.OnlyNumbers("123456789")
		require.NoError(t, err)
	})
}

func TestEachNumberMustBeFollowedByAServerNativeNewlineSequence(t *testing.T) {
	t.Cleanup(tearDown)

	l := logger.New()

	err := l.OnlyNumbers("123456789")
	require.NoError(t, err)

	content, err := readLogFileContent()
	require.NoError(t, err)

	require.Equal(t, "123456789\n", string(content))
}

func TestNoDuplicateNumbersMayBeWrittenToTheLogFile(t *testing.T) {
	t.Cleanup(tearDown)

	l := logger.New()

	err := l.OnlyNumbers("123456789")
	require.NoError(t, err)

	err = l.OnlyNumbers("123456789")
	require.ErrorIs(t, err, logger.ErrDuplicatedNumber)

	content, err := readLogFileContent()
	require.NoError(t, err)

	require.Equal(t, "123456789\n", string(content))
}

func readLogFileContent() ([]byte, error) {
	f, err := os.Open(logFilename)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func createLogFilePreFilled() (*os.File, error) {
	file, err := os.Create(logFilename)
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString("lorem ipsum")
	if err != nil {
		return nil, err
	}

	err = file.Sync()
	if err != nil {
		return nil, err
	}

	return file, nil
}

func logFileExists() bool {
	if _, err := os.Stat(logFilename); os.IsNotExist(err) {
		return false
	}

	return true
}

func countLinesInLogFile() (int, error) {
	content, err := readLogFileContent()
	if err != nil {
		return 0, err
	}

	lines := bytes.Count(content, []byte("\n"))
	return lines, nil
}
