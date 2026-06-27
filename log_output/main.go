package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var id = uuid.New().String()

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/status", statusHandler)

	port := "8080"
	fmt.Println("Server started on port", port)

	http.ListenAndServe(":"+port, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprint(w, `<!DOCTYPE html>
	<html>
	<head>
		<title>Log Output</title>
	</head>
	<body>
		<h1>Log Output App</h1>
		<p>Visit /status for data</p>
	</body>
	</html>`)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w,
		`{"timestamp":"%s","id":"%s"}`,
		time.Now().UTC().Format(time.RFC3339Nano),
		id,
	)
}
