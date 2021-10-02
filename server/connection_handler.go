package server

import (
	"context"
	"errors"
	"net"

	"github.com/gfelixc/gigapipe/logger"
)

type ReadAndLogLines struct {
	logger *logger.LoggerInstrumented
}

func NewReadAndLogLines(logger *logger.LoggerInstrumented) *ReadAndLogLines {
	return &ReadAndLogLines{logger: logger}
}

func (h *ReadAndLogLines) HandleConnection(ctx context.Context, c net.Conn) error {
	defer c.Close()
	linesRead := lineReader(c)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case line, channelOpen := <-linesRead:
			if !channelOpen {
				return nil
			}

			err := h.logger.OnlyNumbers(line)

			if errors.Is(err, logger.ErrDuplicatedNumber) {
				continue
			}

			if errors.Is(err, logger.ErrNonExactDecimalDigitsNumber) {
				return nil
			}

			if errors.Is(err, logger.ErrTerminateSequenceDetected) {
				return err
			}

			if err != nil {
				return err
			}
		}
	}
}