package http1

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrBadRequest          = errors.New("bad request")
	ErrVersionNotSupported = errors.New("http version not supported")
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	// e.g. GET /demo HTTP/1.1 \r\n
	// reading the first line until the delimiter of new line
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// trim the crlf
	line = strings.TrimRight(line, "\r\n")

	// split by space
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("the request is not right : %q", line)
	}

	method := parts[0]
	target := parts[1]
	proto := parts[2]

	// --- NEW VERSION VALIDATION LOGIC ---
	if !strings.HasPrefix(proto, "HTTP/") {
		return nil, ErrBadRequest
	}

	verParts := strings.Split(strings.TrimPrefix(proto, "HTTP/"), ".")
	if len(verParts) != 2 {
		return nil, ErrBadRequest
	}

	major, err1 := strconv.Atoi(verParts[0])
	minor, err2 := strconv.Atoi(verParts[1])
	if err1 != nil || err2 != nil {
		return nil, ErrBadRequest
	}

	// Only HTTP/1.0 and HTTP/1.1 are supported
	if !(major == 1 && (minor == 0 || minor == 1)) {
		return nil, ErrVersionNotSupported
	}
	// ------------------------------------

	// convert (/demo) into the *url.URL that struct requires
	parsedURL, err := url.ParseRequestURI(target)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// build request struct
	req := &http.Request{
		Method: method,
		URL:    parsedURL,
		Proto:  proto,
		Header: make(http.Header),
	}

	for {
		// read one line
		nextLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// trim line breaks
		nextLine = strings.TrimRight(nextLine, "\r\n")

		// detect end of headers
		if nextLine == "" {
			break
		}

		// split key value
		parts := strings.SplitN(nextLine, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("the header is malformed: %q", nextLine)
		}

		// trim possible white spaces
		header := strings.TrimSpace(parts[0])
		headerValue := strings.TrimSpace(parts[1])

		// store in map
		req.Header.Add(header, headerValue)
	}

	req.Body = http.NoBody
	return req, nil
}
