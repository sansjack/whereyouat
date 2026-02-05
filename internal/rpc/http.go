package rpc

import (
	"context"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"whereyouat/internal/rpc/services"
)

type HTTPServer struct {
	listener  net.Listener
	address   string
	rpcServer *rpc.Server
	ls        *services.LocationService
}

func NewHTTPServer(address string, ls *services.LocationService) *HTTPServer {
	return &HTTPServer{
		address:   address,
		rpcServer: rpc.NewServer(),
		ls:        ls,
	}
}

func (s *HTTPServer) Start() error {
	if err := s.rpcServer.Register(s.ls); err != nil {
		return err
	}

	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		remoteAddr := r.RemoteAddr
		
		ctx := context.WithValue(context.Background(), services.RemoteAddrKey, remoteAddr)
		s.ls.SetContext(ctx)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		codec := jsonrpc.NewServerCodec(&httpConn{
			Reader: r.Body,
			Writer: w,
		})
		s.rpcServer.ServeRequest(codec)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	s.listener = listener

	return http.Serve(listener, nil)
}

func (s *HTTPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

type httpConn struct {
	Reader interface{ Read([]byte) (int, error) }
	Writer interface{ Write([]byte) (int, error) }
}

func (c *httpConn) Read(p []byte) (n int, err error) {
	return c.Reader.Read(p)
}

func (c *httpConn) Write(p []byte) (n int, err error) {
	return c.Writer.Write(p)
}

func (c *httpConn) Close() error {
	return nil
}
