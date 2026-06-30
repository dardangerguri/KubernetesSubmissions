package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const filePath = "/tmp/pongs.txt"

func getCounter() int {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Debug - ReadFile Error: %v\n", err)
		return 0
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		fmt.Printf("Debug - Atoi Parse Error: %v\n", err)
		return 0
	}
	return count
}

func saveCounter(count int) {
	_ = os.MkdirAll("/tmp", 0755)
	err := os.WriteFile(filePath, []byte(strconv.Itoa(count)), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		counter := getCounter()
		counter++
		saveCounter(counter)

		fmt.Fprintf(w, "pong %d", counter)
	})

	http.HandleFunc("/pings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprintf(w, "%d", getCounter())
	})

	fmt.Printf("Ping-pong server started in port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err !=nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
