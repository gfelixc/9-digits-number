package server_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/gfelixc/gigapipe/server"
	"github.com/stretchr/testify/require"
)

type spyProcessor struct {
	linesRead int
}

func (sp *spyProcessor) cancelAfterNLinesRead(n int, cancel func()) func(s string) error {
	return func(s string) error {
		sp.linesRead++
		if sp.linesRead == n {
			cancel()
		}

		return nil
	}
}

func (sp *spyProcessor) increaseLinesRead(s string) error {
	sp.linesRead++

	return nil
}

func TestStopReadingWhenATerminateSequenceIsReceived(t *testing.T) {
	spy := spyProcessor{}
	handler := server.NewHandler(spy.increaseLinesRead)

	dataWithTerminateSequence := bytes.NewBuffer([]byte("098765432\n123456789\n0000000001\nterminate\n098765678\n"))
	handler.ReadLines(context.TODO(), dataWithTerminateSequence)

	require.Equal(t, 3, spy.linesRead)
}

func TestCancellingHandlerUsingTheContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	spy := spyProcessor{}
	handler := server.NewHandler(spy.cancelAfterNLinesRead(5, cancel))

	reader, writer := io.Pipe()
	writeLinesExactlyNineDecimalDigits(writer)
	handler.ReadLines(ctx, reader)

	require.GreaterOrEqual(t, spy.linesRead, 5)
}

func writeLinesExactlyNineDecimalDigits(w *io.PipeWriter) {
	go func() {
		for {
			_, _ = w.Write([]byte("123456789\n"))
		}
	}()
}
