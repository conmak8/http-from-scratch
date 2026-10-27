package http1

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	// e.g. GET /demo HTTP/1.1
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

}
