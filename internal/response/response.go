package response

import (
	"fmt"
	"httpfromscratch/internal/headers"
	"io"
	"strconv"
)

type StatusCode int
type WriterState string

const (
	WriterStateInit    WriterState = "Init"
	WriterStateHeaders WriterState = "Headers"
	WriterStateBody    WriterState = "Body"
)

type Writer struct {
	CurrentState WriterState
	Writer       io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {

	if w.CurrentState != WriterStateInit {
		return fmt.Errorf("Writer is not in init state")
	}
	statusLine := []byte{}
	switch statusCode {
	case OK:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case BadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case InternalError:

		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		return fmt.Errorf("Unknown status code: %d", statusCode)
	}

	w.Writer.Write(statusLine)
	w.CurrentState = WriterStateHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.CurrentState != WriterStateHeaders {
		return fmt.Errorf("Writer is not in headers state")
	}

	for key, val := range headers {
		w.Writer.Write([]byte(key + ": " + val + "\r\n"))
	}
	w.Writer.Write([]byte("\r\n"))

	w.CurrentState = WriterStateBody
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.CurrentState != WriterStateBody {
		return 0, fmt.Errorf("Writer is not in body state")
	}

	w.CurrentState = WriterStateInit
	return w.Writer.Write(p)
}

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
