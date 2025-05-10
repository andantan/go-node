package network

import (
	"fmt"
	"time"
)

// TODO: FIRST IMPLEMENT
// The first transport over blockchain server

type ServerOpts struct {
	// This is will be container
	Transports []Transport
}

type Server struct {
	ServerOpts

	rpcCh  chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh:      make(chan RPC),
		quitCh:     make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransports()

	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			// Is there some message from rpc channel
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			// break -> Break select statement NOT for loop
			break free
		case <-ticker.C:
			fmt.Println("Do stuff every x seconds")
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				// Handle
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
