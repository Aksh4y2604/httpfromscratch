package server

import (
	"fmt"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"io"
	"net"
)

type Server struct {
	closed  bool
	handler Handler
}
type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func handleConn(Server *Server, conn io.ReadWriteCloser) {
	defer conn.Close()

	responseWriter := response.Writer{
		CurrentState: response.WriterStateInit,
		Writer:       conn,
	}

	header := response.GetDefaultHeaders(0)
	req, err := request.RequestFromReader(conn)

	if err != nil {
		responseWriter.WriteStatusLine(response.BadRequest)
		responseWriter.WriteHeaders(header)
		return
	}

	Server.handler(&responseWriter, req)

}

func runServer(Server *Server, listener net.Listener) error {

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			if Server.closed {
				break
			}
			go handleConn(Server, conn)
		}
	}()

	return nil
}

func Serve(port uint16, handler Handler) (*Server, error) {

	server := &Server{}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server.closed = false
	server.handler = handler

	err = runServer(server, listener)

	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
