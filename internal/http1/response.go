package http1

import (
	"fmt"
	"net"
	"net/http"
)

type Response struct {
	conn        net.Conn
	header      http.Header
	statusCode  int
	wroteHeader bool
}

func NewResponse(conn net.Conn) *Response {
	return &Response{
		conn:        conn,              // sendig data to whom
		header:      make(http.Header), // store headers without nil map panic
		statusCode:  http.StatusOK,     // HTTP status ? if handler doesnt set one?
		wroteHeader: false,             // are any bytes already send?
	}
}

func (r *Response) Header() http.Header {
	return r.header
}

func (r *Response) Write(data []byte) (n int, err error) {
	// make sure http header have been transmitted
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}

	//stream raw body bytes
	return r.conn.Write(data)
}

func (r *Response) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.wroteHeader = true
	r.statusCode = statusCode

	// status line according to RFC
	fmt.Fprintf(r.conn, "HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))

	// loop over header and write each key: values
	for key, values := range r.header {
		for _, val := range values {
			fmt.Fprintf(r.conn, "%s: %s\r\n", key, val)
		}
	}

	fmt.Fprintf(r.conn, "\r\n")
}
