package server

import (
	"context"
	"net"
	"sync"
)

type Server struct {
	concurrencyControl chan struct{}
	handler            func(conn net.Conn)
	connectionsWG      *sync.WaitGroup
}

const (
	maxConcurrentClient = 5
	network             = "tcp"
	address             = ":4000"
)

func New() Server {
	return Server{
		connectionsWG:      &sync.WaitGroup{},
		concurrencyControl: make(chan struct{}, maxConcurrentClient),
	}
}
func (s *Server) AddHandler(fn func(conn net.Conn)) {
	s.handler = fn
}

func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			s.connectionsWG.Wait()
			return ctx.Err()
		case s.concurrencyControl <- struct{}{}:
			conn, err := l.Accept()
			if err != nil {
				return err
			}

			go s.handleConnection(conn)
		}
	}
}

func (s Server) handleConnection(c net.Conn) {
	s.connectionsWG.Add(1)

	defer func() {
		<-s.concurrencyControl
		s.connectionsWG.Done()
	}()

	if s.handler == nil {
		return
	}

	s.handler(c)
}
