package server

import (
	"context"
	"io"
	"regexp"
)

const (
	terminateKeyword = "terminate"
)

var exactlyNineDecimalDigitsPattern = regexp.MustCompile(`^[0-9]{9}?`)

type Handler struct {
	processor func(s string)
}

func NewHandler(processor func(s string)) *Handler {
	return &Handler{
		processor: processor,
	}
}

func (h Handler) ReadLines(ctx context.Context, incomingData io.Reader) {
	linesRead := h.continuesReader(incomingData)

	for {
		select {
		case <-ctx.Done():
			return

		case line, channelOpen := <-linesRead:
			if !channelOpen {
				return
			}

			if line == terminateKeyword {
				return
			}

			if !exactlyNineDecimalDigitsPattern.MatchString(line) {
				return
			}

			h.processor(line)
		}
	}
}

func (h Handler) continuesReader(incomingData io.Reader) <-chan string {
	lineFeed := make(chan string)
	scanner := newLineScanner(incomingData)

	go func() {
		defer close(lineFeed)

		for {
			ok := scanner.Scan()
			if !ok {
				return
			}

			lineFeed <- scanner.Text()
		}
	}()

	return lineFeed
}
