package server

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/gfelixc/gigapipe/logger"
	"github.com/stretchr/testify/require"
)

type components struct {
	handler       *ReadAndLogLines
	logger        *logger.Logger
	logWriterMock *bytes.Buffer
}

func setup() components {
	w := new(bytes.Buffer)
	l := logger.New(w)
	h := NewReadAndLogLines(l)
	return components{
		handler:       h,
		logger:        l,
		logWriterMock: w,
	}
}

func TestFinishWhenContextDone(t *testing.T) {
	c := setup()
	conn := connectionMocked()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := c.handler.HandleConnection(ctx, conn)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestReturnNoErrorWhenConnectionIsClosed(t *testing.T) {
	c := setup()
	conn := connectionThatWillBeClosedIn100ms()
	err := c.handler.HandleConnection(context.Background(), conn)
	require.NoError(t, err)
}

func TestContinueReadingWhenErrDuplicatedIsReceived(t *testing.T) {
	connA, connB := net.Pipe()

	input := "123456789\n123456789\n012345678\n"
	expected := "123456789\n012345678\n"

	wait100msWriteDataAndCloseConnections([]byte(input), connA, connB)

	c := setup()
	_ = c.handler.HandleConnection(context.Background(), connA)
	c.logger.Shutdown()

	require.Equal(t, expected, c.logWriterMock.String())
}

func TestReturnsNilErrWhenErrNonExactDecimalDigitsNumber(t *testing.T) {
	connA, connB := net.Pipe()

	input := "123456789\n126789\n012345678\n"

	wait100msWriteData([]byte(input), connB)

	c := setup()
	err := c.handler.HandleConnection(context.Background(), connA)
	require.NoError(t, err)
}


func TestReturnsErrWhenErrTerminateSequenceDetected(t *testing.T) {
	connA, connB := net.Pipe()

	input := "terminate\n"

	wait100msWriteData([]byte(input), connB)

	c := setup()
	err := c.handler.HandleConnection(context.Background(), connA)
	require.ErrorIs(t, err, logger.ErrTerminateSequenceDetected)
}

func wait100msWriteData(data []byte, connB net.Conn)  {
	go func() {
		time.Sleep(100 * time.Millisecond)

		_, _ = connB.Write(data)
	}()
}

func wait100msWriteDataAndCloseConnections(data []byte, connA net.Conn, connB net.Conn)  {
	go func() {
		time.Sleep(100 * time.Millisecond)

		_, _ = connB.Write(data)
		_ = connB.Close()
		_ = connA.Close()
	}()
}

func connectionMocked() net.Conn {
	_, conn := net.Pipe()
	return conn
}

func connectionThatWillBeClosedIn100ms() net.Conn {
	conn := connectionMocked()

	go func() {
		time.Sleep(100 * time.Millisecond)

		_ = conn.Close()
	}()

	return conn
}
