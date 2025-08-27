package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromscratch/internal/headers"
	"httpfromscratch/internal/request"
	"httpfromscratch/internal/response"
	"httpfromscratch/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

func respond400() []byte {
	return []byte("<html> <head> <title>400 Bad Request</title> </head> <body> <h1>Bad Request</h1> <p>Your request honestly kinda sucked.</p> </body> </html>")
}

func respond500() []byte {
	return []byte("<html> <head> <title>500 Internal Server Error</title> </head> <body> <h1>Internal Server Error</h1> <p>Something went wrong.</p> </body> </html>")
}

func handleYourProblem(h headers.Headers) (response.StatusCode, []byte) {
	body := respond400()
	h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
	return response.BadRequest, body
}

func handleMyProblem(h headers.Headers) (response.StatusCode, []byte) {
	body := respond500()
	h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
	return response.InternalError, body
}

func handleHttpbin(w *response.Writer, h headers.Headers, target string) bool {
	ep := strings.TrimPrefix(target, "/httpbin/")
	epInt, err := strconv.Atoi(ep)
	if err != nil || epInt == 0 {
		return false
	}

	resp, err := http.Get("http://httpbin.org/stream/" + ep)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// prepare chunked response
	h.Delete("Content-length")
	h.Set("Transfer-Encoding", "chunked")
	h.Replace("Content-type", "text/plain")
	h.Set("Trailers", "X-Content-SHA256")
	h.Set("Trailers", "X-Content-Length")

	w.WriteStatusLine(response.OK)
	w.WriteHeaders(h)

	fullBody := []byte{}
	buf := make([]byte, 32)

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			fullBody = append(fullBody, buf[:n]...)
			if _, err := w.WriteChunkedBody(buf[:n]); err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}

	shaSum := fmt.Sprintf("%x", sha256.Sum256(fullBody))
	trailers := headers.NewHeaders()
	trailers.Set("X-Content-SHA256", shaSum)
	trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
	w.WriteTrailers(trailers)

	return true
}

func handleDefault(h headers.Headers) (response.StatusCode, []byte) {
	body := []byte("<html> <head> <title>200 OK</title> </head> <body> <h1>200 OK</h1> <p>Everything is fine.</p> </body> </html>")
	h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
	return response.OK, body
}

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		h.Replace("Content-Type", "text/html")

		var (
			statusCode response.StatusCode
			body       []byte
		)

		switch {
		case req.RequestLine.RequestTarget == "/yourproblem":
			statusCode, body = handleYourProblem(h)

		case req.RequestLine.RequestTarget == "/myproblem":
			statusCode, body = handleMyProblem(h)

		case strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/"):
			if handleHttpbin(w, h, req.RequestLine.RequestTarget) {
				return // handled fully inside
			}
			statusCode, body = response.InternalError, respond500()

		default:
			statusCode, body = handleDefault(h)
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
