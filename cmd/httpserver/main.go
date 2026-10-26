package main

import (
	"fmt"
	"httpServer/internal/http1"
	"net/http"
)

const addr = ":8080"

func main() {
	//network: "tcp"
	//address: ":8080"
	//listener, err := net.Listen("tcp", ":8080")
	//if err != nil {
	//	log.Fatalf("failed to bind: %v", err)
	//}
	//defer listener.Close()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})
	if err := http1.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
