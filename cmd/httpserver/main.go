package main

import (
	"fmt"
	"httpServer/internal/http1"
	"log"
	"net/http"
)

const addr = ":8080"

//	func main() {
//		//network: "tcp"
//		//address: ":8080"
//		//listener, err := net.Listen("tcp", ":8080")
//		//if err != nil {
//		//	log.Fatalf("failed to bind: %v", err)
//		//}
//		//defer listener.Close()
//
//		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
//			fmt.Fprintf(w, "hello")
//		})
//		if err := http1.ListenAndServe(addr, nil); err != nil {
//			panic(err)
//		}
//	}
func main() {
	// Creation of custom router by reference
	mux := http1.NewServerMux()

	// Registering routes
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, "Returning list of users\n")
			return
		}
		if r.Method == http.MethodPost {
			fmt.Fprintf(w, "Creating a new user\n")
			return
		}
		// if method not allowed fallback
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed on /users\n", r.Method)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})

	// starting the server using custom Mux instead of nil
	log.Println("Server listening on :8080...")
	if err := http1.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
