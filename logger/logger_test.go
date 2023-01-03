package logger_test

import (
	"bytes"
	"testing"

	"github.com/gfelixc/9-digits-number/logger"
	"github.com/stretchr/testify/require"
)

const (
	validNumber        = "123456789"
	validNumberLogLine = "123456789\n"
)

func TestShutdownDumpsInMemoryNumbersToLogWriter(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf)
	_ = l.OnlyNumbers(validNumber)
	l.Shutdown()

	require.Equal(t, validNumberLogLine, buf.String())
}

func TestOnlyNumbersMayBeWrittenToTheLogFile(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf)
	defer func() {
		l.Shutdown()
	}()

	t.Run("Alphanumeric", func(t *testing.T) {
		err := l.OnlyNumbers("Aca2321")
		require.ErrorIs(t, err, logger.ErrNonExactDecimalDigitsNumber)
	})

	t.Run("SpecialChars", func(t *testing.T) {
		err := l.OnlyNumbers("123 12-3")
		require.ErrorIs(t, err, logger.ErrNonExactDecimalDigitsNumber)
	})

	t.Run("OnlyDecimals", func(t *testing.T) {
		err := l.OnlyNumbers(validNumber)
		require.NoError(t, err)
	})
}

func TestEachNumberMustBeFollowedByAServerNativeNewlineSequence(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf)
	_ = l.OnlyNumbers(validNumber)
	l.Shutdown()

	require.Equal(t, validNumberLogLine, buf.String())
}

func TestTerminateSequenceReturnsErrTerminateSequenceDetected(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf)
	defer func() {
		l.Shutdown()
	}()

	_ = l.OnlyNumbers(validNumber)

	err := l.OnlyNumbers("terminate")
	require.ErrorIs(t, err, logger.ErrTerminateSequenceDetected)
}

func TestDuplicatedNumbersReturnsErrDuplicatedNumber(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf)
	defer func() {
		l.Shutdown()
	}()

	_ = l.OnlyNumbers(validNumber)

	err := l.OnlyNumbers(validNumber)
	require.ErrorIs(t, err, logger.ErrDuplicatedNumber)
}

func TestNoDuplicatedNumbersMayBeWrittenToTheLogFile(t *testing.T) {
	var buf bytes.Buffer

	l := logger.New(&buf)
	_ = l.OnlyNumbers(validNumber)
	_ = l.OnlyNumbers(validNumber)
	l.Shutdown()

	require.Equal(t, validNumberLogLine, buf.String())
}
