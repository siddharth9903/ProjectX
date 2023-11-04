package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
}

type Server struct {
	ServerOpts
	rpcChan  chan RPC
	quitChan chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcChan:    make(chan RPC),
		quitChan:   make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcChan:
			fmt.Println(rpc)
		case <-s.quitChan:
			break free
		case <-ticker.C:
			fmt.Println("Every 5 seconds")
		}
	}

	fmt.Println("Server shut down")
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcChan <- rpc
			}
		}(tr)
	}
}
