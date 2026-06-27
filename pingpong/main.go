package main

import (
	"fmt"
	"net/http"
	"os"
)

var counter = 0

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		counter++
		fmt.Fprintf(w, "pong %d", counter)
	})

	fmt.Printf("Ping-pong server started in port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err !=nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
