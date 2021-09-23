package server_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/gfelixc/gigapipe/server"
	"github.com/stretchr/testify/require"
)

func Test5MaximumConcurrentClients(t *testing.T) {
	ctx, cancelCTX := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCTX()

	var connectedClients int

	handler := func(conn net.Conn) {
		connectedClients++

		// do work
		time.Sleep(10 * time.Second)
	}

	setupAndRunServer(ctx, handler)

	for i := 0; i < 50; i++ {
		go func() {
			_, _ = net.Dial("tcp", ":4000")
		}()
	}

	<-ctx.Done()

	require.Equal(t, 5, connectedClients)
}

func setupAndRunServer(ctx context.Context, handler func(conn net.Conn)) {
	s := server.New()
	s.AddHandler(handler)

	go func() {
		_ = s.Start(ctx)
	}()

	time.Sleep(500 * time.Millisecond)
}
