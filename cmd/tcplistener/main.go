package main

import (
	"fmt"
	"httpfromscratch/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		rq, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
		}
		fmt.Printf("Request Line: \n - Method: %s \n - Target: %s \n - Version: %s \n", rq.RequestLine.Method, rq.RequestLine.RequestTarget, rq.RequestLine.HttpVersion)
	}
}
