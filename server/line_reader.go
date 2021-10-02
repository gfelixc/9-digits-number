package server

import (
	"bufio"
	"bytes"
	"io"
)

// lineReader returns a string channel to delivery the content
// of each line clean of CRLF sequences. Channel is closed
// when Reader reach its end or an error occurs
func lineReader(incomingData io.Reader) <-chan string {
	linesFeed := make(chan string)
	scanner := bufio.NewScanner(incomingData)
	scanner.Split(scanOnlyLinesEndedWithNewlineSequence)

	go func() {
		for scanner.Scan() {
			linesFeed <- scanner.Text()
		}

		close(linesFeed)
	}()

	return linesFeed
}

// scanOnlyLinesEndedWithNewlineSequence is a bufio.SplitFunc implementation
// intended to scan lines only if ends with a newline sequence
func scanOnlyLinesEndedWithNewlineSequence(data []byte, _ bool) (advance int, token []byte, err error) {
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, dropCR(data[0:i]), nil
	}

	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}

	return data
}
