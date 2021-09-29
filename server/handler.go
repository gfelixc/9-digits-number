package server

import (
	"context"
	"errors"
	"io"

	"github.com/gfelixc/gigapipe/logger"
)

const (
	terminateKeyword = "terminate"
)

type Handler struct {
	processor func(s string) error
}

func NewHandler(processor func(s string) error) *Handler {
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

			err := h.processor(line)
			if errors.Is(err, logger.ErrNonExactDecimalDigitsNumber) {
				return
			}

			if err != nil {
				println("unable to process line", err.Error())
			}
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
