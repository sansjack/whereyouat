package rpc

import (
	"context"
	"log"
	"net"
	"net/rpc"
	"whereyouat/internal/rpc/services"
)

type Server struct {
	listener  net.Listener
	address   string
	rpcServer *rpc.Server
	ls        *services.LocationService
}

func NewServer(address string, ls *services.LocationService) *Server {
	return &Server{
		address:   address,
		rpcServer: rpc.NewServer(),
		ls:        ls,
	}
}

func (s *Server) Start() error {
	if err := s.rpcServer.Register(s.ls); err != nil {
		log.Printf("Error registering service: %v", err)
		return err
	}

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	s.listener = listener
	log.Printf("RPC server listening on %s", s.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()

			remoteAddr := c.RemoteAddr().String()

			ctx := context.WithValue(context.Background(), services.RemoteAddrKey, remoteAddr)
			s.ls.SetContext(ctx)

			s.rpcServer.ServeConn(c)
		}(conn)
	}
}

func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
