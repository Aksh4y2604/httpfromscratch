package response

import (
	"fmt"
	"httpfromscratch/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	OK            StatusCode = 200
	BadRequest    StatusCode = 400
	InternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case OK:
		w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case BadRequest:
		w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case InternalError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	default:
		return fmt.Errorf("Unknown status code: %d", statusCode)
	}
	return nil
}
func GetDefaultHeaders(contentLen int) headers.Headers {

	headers := headers.NewHeaders()
	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		w.Write([]byte(key + ": " + val + "\r\n"))
	}
	w.Write([]byte("\r\n"))
	return nil
}
