package server

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScannerReadLinesEndedWithUnixLikeNewLineSequence(t *testing.T) {
	linesUnixLike := bytes.NewBufferString("a\nb\nc\nd\n")
	scanner := newLineScanner(linesUnixLike)

	require.Equal(t, 4, countLinesRead(scanner))
}

func TestScannerReadLinesEndedWithNonUnixLikeNewLineSequence(t *testing.T) {
	linesNonUnixLike := bytes.NewBufferString("a\r\nb\r\nc\r\n")
	scanner := newLineScanner(linesNonUnixLike)

	require.Equal(t, 3, countLinesRead(scanner))
}

func TestScannerReadOnlyLinesEndingWithNewLine(t *testing.T) {
	linesNonUnixLike := bytes.NewBufferString("a\r\nb\r\nc\r\n1232")
	scanner := newLineScanner(linesNonUnixLike)

	require.Equal(t, 3, countLinesRead(scanner))
}

func countLinesRead(scanner *bufio.Scanner) int {
	var readLines int
	for ok := scanner.Scan(); ok; {
		readLines++
	}

	return readLines
}
