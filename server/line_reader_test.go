package server

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLineReaderIgnoresTheLastLineDueALackOfNewLineSequence(t *testing.T) {
	linesNonUnixLike := bytes.NewBufferString("a\nb")
	scanner := lineReader(linesNonUnixLike)
	lines := getAnSliceOfStrings(scanner)

	require.Equal(t, []string{"a"}, lines)
}

func TestLineReaderCleansCRLFSequences(t *testing.T) {
	linesNonUnixLike := bytes.NewBufferString("a\r\nb\nc\r\n")
	scanner := lineReader(linesNonUnixLike)
	lines := getAnSliceOfStrings(scanner)

	require.Equal(t, []string{"a", "b", "c"}, lines)
}

func getAnSliceOfStrings(scanner <-chan string) []string {
	var lines []string
	for v := range scanner {
		lines = append(lines, v)
	}
	return lines
}

