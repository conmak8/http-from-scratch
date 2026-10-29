package http1

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

type Response struct {
	conn        net.Conn
	header      http.Header
	body        bytes.Buffer
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
	//return r.conn.Write(data)

	// Write MUST write to the buffer, NOT the network connection!
	return r.body.Write(data)
}

func (r *Response) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.wroteHeader = true
	r.statusCode = statusCode

	// writeHeader MUST NOT write to the network connection!

	//// status line according to RFC
	//fmt.Fprintf(r.conn, "HTTP/1.1 %d %s\r\n", statusCode, http.StatusText(statusCode))
	//
	//// loop over header and write each key: values
	//for key, values := range r.header {
	//	for _, val := range values {
	//		fmt.Fprintf(r.conn, "%s: %s\r\n", key, val)
	//	}
	//}
	//
	//fmt.Fprintf(r.conn, "\r\n")
}

// only flush method can send bytes over the socket
func (r *Response) Flush(reqMethod string) error {
	// Status-Line
	if _, err := fmt.Fprintf(r.conn, "HTTP/1.1 %d %s\r\n", r.statusCode, http.StatusText(r.statusCode)); err != nil {
		return err
	}

	// Calculate Content-Length if not already provided and method is not HEAD
	if r.header.Get("Content-Length") == "" && reqMethod != http.MethodHead {
		r.header.Set("Content-Length", strconv.Itoa(r.body.Len()))
	}

	// Write Headers
	for key, values := range r.header {
		for _, val := range values {
			if _, err := fmt.Fprintf(r.conn, "%s: %s\r\n", key, val); err != nil {
				return err
			}
		}
	}

	// End of Headers (the standalone CRLF)
	if _, err := fmt.Fprintf(r.conn, "\r\n"); err != nil {
		return err
	}

	// Body (skip body bytes for HEAD requests)
	if reqMethod != http.MethodHead && r.body.Len() > 0 {
		if _, err := r.conn.Write(r.body.Bytes()); err != nil {
			return err
		}
	}

	return nil
}
