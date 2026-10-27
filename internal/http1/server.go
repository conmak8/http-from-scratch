package http1

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
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
	//panic("HandleConnection: Not Implemented")

	defer conn.Close()

	reader := bufio.NewReader(conn)
	req, err := ParseRequest(reader)
	if err != nil {
		return err
	}

	// Temporary verification log:
	// with curl -v -H "X-Custom-Token: abc123" http://localhost:8080/hello
	//log.Printf("Successfully parsed request: Method=%s Path=%s Headers=%v", req.Method, req.URL.Path, req.Header)

	res := NewResponse(conn)
	mux.ServeHTTP(res, req)

	return nil

}
