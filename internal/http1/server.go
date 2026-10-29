package http1

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func ListenAndServe(addr string, handler http.Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer listener.Close()

	if handler == nil {
		handler = http.DefaultServeMux
	}

	for {
		// w8t for connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("There is an error: %v", err)
			continue
		}

		go HandleConnection(conn, handler)

		/*
				Functions that might be useful:
				* listener.Accept()
			  * HandleConnection()
		*/
	}
}

// HandleConnection handle the lifecycle of one client socket
func HandleConnection(conn net.Conn, mux http.Handler) error {
	defer conn.Close()

	// handle the keep alive connection
	// create reader once for whole connection
	reader := bufio.NewReader(conn)

	// loop as to handle multiple consecutive requests(keep-alive)
	for {
		req, err := ParseRequest(reader)
		if err != nil {
			// if client closed the connection will retrun nil
			if err == io.EOF {
				return nil
			}

			// TODO: 400 bad request to conn before returning directly to socket
			errorResponse := "HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"
			conn.Write([]byte(errorResponse))

			return err // as to trigger defer conn.Close()
		}
		res := NewResponse(conn)
		mux.ServeHTTP(res, req)

		// Check if client asked to close, or if HTTP/1.0 didn't ask for keep-alive
		if strings.EqualFold(req.Header.Get("Connection"), "close") {
			return nil // Clean exit -> triggers defer conn.Close()
		}
		if req.Proto == "HTTP/1.0" && !strings.EqualFold(req.Header.Get("Connection"), "keep-alive") {
			return nil // HTTP/1.0 closes by default
		}
	}
	// Temporary verification log:
	// with curl -v -H "X-Custom-Token: abc123" http://localhost:8080/hello
	//log.Printf("Successfully parsed request: Method=%s Path=%s Headers=%v", req.Method, req.URL.Path, req.Header)
}
