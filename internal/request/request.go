package request

import (
	"bytes"
	"fmt"
	"httpfromscratch/internal/headers"
	"io"
	"strconv"
)

type parserState string

const (
	StateInit   parserState = "init"
	StateDone   parserState = "done"
	StateError  parserState = "error"
	StateHeader parserState = "header"
	StateBody   parserState = "body"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     headers.Headers
	Body        []byte
}

var ERROR_MALFORMED_REQ_LINE = fmt.Errorf("Bad request line")
var ERROR_ERROR_STATE = fmt.Errorf("Error State")
var ERROR_BODY_PARSE = fmt.Errorf("Error parsing body")
var SEPERATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQ_LINE
	}
	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ERROR_MALFORMED_REQ_LINE
	}
	var consumedLen int = len(startLine) + len(SEPERATOR)
	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, consumedLen, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
	if r.Headers == nil {
		r.Headers = headers.NewHeaders()
	}
outer:
	for {
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateHeader
		case StateHeader:
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}
			read += n
			if done {
				r.state = StateBody
			} else {
				break outer
			}

		case StateBody:
			contentLengthStr := r.Headers.Get("content-length")
			if contentLengthStr == "" {
				r.state = StateDone
				break outer
			}

			contentLen, err := strconv.Atoi(contentLengthStr)
			if err != nil {
				return 0, err
			}

			// Not enough data yet to read full body
			if len(data[read:]) < contentLen {
				break outer
			}
			if len(data[read:]) > contentLen {
				return 0, ERROR_BODY_PARSE
			}

			if contentLen == 0 {
				r.Body = []byte{}
			} else {
				r.Body = data[read : read+contentLen]
				read += contentLen
			}

			r.state = StateDone
		case StateDone:
			break outer
		case StateError:
			return 0, ERROR_ERROR_STATE
		}
	}
	return read, nil
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}
		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}
	return request, nil
}

