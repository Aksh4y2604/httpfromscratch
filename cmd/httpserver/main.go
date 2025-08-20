package main

import (
	"fmt"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"httpfromscratch/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {

		h := response.GetDefaultHeaders(0)
		h.Replace("Content-Type", "text/html")
		body := []byte{}
		statusCode := response.OK

		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			statusCode = response.BadRequest
			body = []byte("<html> <head> <title>400 Bad Request</title> </head> <body> <h1>Bad Request</h1> <p>Your request honestly kinda sucked.</p> </body> </html>")
			h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
		case "/myproblem":
			statusCode = response.InternalError
			body = []byte("<html> <head> <title>500 Internal Server Error</title> </head> <body> <h1>Internal Server Error</h1> <p>Something went wrong.</p> </body> </html>")
			h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
		default:

			statusCode = response.OK
			body = []byte("<html> <head> <title>200 OK</title> </head> <body> <h1>200 OK</h1> <p>Everything is fine.</p> </body> </html>")
			h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
		}

		w.WriteStatusLine(statusCode)
		w.WriteHeaders(h)
		w.WriteBody(body)
	})

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
