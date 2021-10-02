package server

import (
	"context"
	"net"
	"sync"
)

type connectionHandler func(ctx context.Context, conn net.Conn) error

type Server struct {
	handler             connectionHandler
	incomingConnections chan net.Conn
	concurrencyControl  chan struct{}
	connectionsWG       *sync.WaitGroup
}

const (
	maxConcurrentClients = 5
	network              = "tcp"
	address              = ":4000"
)

// New receives a server.connectionHandler as argument, when connectionHandler
// returns an error it triggers a graceful shutdown, waiting for all the handlers to finish
func New(handler connectionHandler) *Server {
	return &Server{
		handler:             handler,
		connectionsWG:       &sync.WaitGroup{},
		incomingConnections: make(chan net.Conn),
		concurrencyControl:  make(chan struct{}, maxConcurrentClients),
	}
}

func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	go s.acceptConnections(l)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			s.connectionsWG.Wait()
			_ = l.Close()

			return ctx.Err()

		case conn := <-s.incomingConnections:
			go s.handleConnectionWithinWaitGroup(ctx, cancel, conn)
		}
	}
}

func (s *Server) acceptConnections(l net.Listener) {
	for {
		s.concurrencyControl <- struct{}{}

		conn, err := l.Accept()
		if err != nil {
			<-s.concurrencyControl
			continue
		}

		s.incomingConnections <- conn
	}
}

func (s *Server) handleConnectionWithinWaitGroup(ctx context.Context, cancelWG func(), c net.Conn) {
	s.connectionsWG.Add(1)

	defer func() {
		<-s.concurrencyControl
		s.connectionsWG.Done()
		_ = c.Close()
	}()

	if err := s.handler(ctx, c); err != nil {
		cancelWG()
	}
}
