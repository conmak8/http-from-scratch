package http1

import (
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
	panic("Response.Header: Not Implemented")
}

func (r *Response) Write(data []byte) (n int, err error) {
	panic("Response.Write: Not Implemented")
}

func (r *Response) WriteHeader(statusCode int) {
	panic("Response.WriteHeader: Not Implemented")
}
