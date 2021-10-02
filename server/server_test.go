package server

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServerShutdownsWaitForHandlers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	spy := spyConnectionHandler{holdConnections: 1 * time.Second}

	s := New(spy.countConnectionsAndHoldIt)
	go startNConcurrentClients(3)

	_ = s.Start(ctx)

	require.Greater(t, spy.connectedClients, 0)
	require.Equal(t, spy.connectedClients, spy.handlersDone)
}

func TestServerShutdownsWhenContextIsDone(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	s := New(func(_ context.Context, _ net.Conn) error { return nil })
	err := s.Start(ctx)

	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func Test5MaximumConcurrentClients(t *testing.T) {
	ctx, cancelCTX := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelCTX()

	spy := spyConnectionHandler{holdConnections: 10 * time.Second}

	s := New(spy.countConnectionsAndHoldIt)

	go s.Start(ctx)
	startNConcurrentClients(50)

	<-ctx.Done()

	require.Equal(t, 5, spy.connectedClients)
}


func startNConcurrentClients(n int) {
	for i := 0; i < n; i++ {
		_, _ = net.Dial("tcp", ":4000")
	}
}

type spyConnectionHandler struct {
	connectedClients int
	handlersDone     int
	holdConnections  time.Duration
}

func (s *spyConnectionHandler) countConnectionsAndHoldIt(_ context.Context, _ net.Conn) error {
	s.connectedClients++

	time.Sleep(s.holdConnections)

	s.handlersDone++

	return nil
}
