package server

import (
	"bytes"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func handleConn(Server *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	header := response.GetDefaultHeaders(0)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.BadRequest)
		response.WriteHeaders(conn, header)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	handlerError := Server.handler(buf, req)

	var body []byte = nil
	var status response.StatusCode = response.OK

	if handlerError != nil {
		status = handlerError.StatusCode
		body = []byte(handlerError.Message)
	} else {
		body = buf.Bytes()
	}
	header.Replace("content-length", fmt.Sprintf("%d", len(body)))
	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, header)
	conn.Write(body)

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
