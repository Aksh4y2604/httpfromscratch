package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

var KEY_VAL_SEPERATOR = []byte(":")
var RN = []byte("\r\n")
var CRLF = []byte("\r\n\r\n")

var ERROR_INVALID_REQ = fmt.Errorf("Error State")
var ERROR_MAL_HEADER = fmt.Errorf("Malformed Header")

func NewHeaders() Headers {
	return Headers{}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2) // split into key + rest
	if len(parts) != 2 {
		return "", "", ERROR_MAL_HEADER
	}

	if bytes.HasSuffix(parts[0], []byte(" ")) {
		return "", "", ERROR_MAL_HEADER
	}

	key := strings.TrimSpace(string(parts[0]))
	val := strings.TrimSpace(string(parts[1]))

	pattern := `^[A-Za-z0-9!#$%&'*+\-.\^_` + "`" + `|~]+$`

	re := regexp.MustCompile(pattern)

	if !re.MatchString(key) || key == "" || val == "" {
		return "", "", ERROR_MAL_HEADER
	}

	lower_case := strings.ToLower(key)
	return lower_case, val, nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		lineEndIdx := bytes.Index(data[read:], RN)
		if lineEndIdx == -1 {
			break
		}
		if lineEndIdx == 0 {
			done = true
			read += len(RN)
			break
		}

		key, val, err := parseHeader(data[read : read+lineEndIdx])
		if err != nil {
			return 0, false, err
		}

		h.Set(key, val)
		read += lineEndIdx + len(RN)
	}
	return read, done, nil
}
func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}
func (h Headers) Set(key string, val string) {
	if h[key] != "" {
		newVal := h.Get(key) + ", " + val
		h[strings.ToLower(key)] = newVal
	} else {
		h[strings.ToLower(key)] = val
	}
}
func (h Headers) Replace(key string, val string) {
	h[strings.ToLower(key)] = val
}
