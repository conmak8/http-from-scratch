package http1

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

	return req, nil
}
