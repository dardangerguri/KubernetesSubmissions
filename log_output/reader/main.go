package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logData, err := os.ReadFile("/shared/log.txt")
		if err != nil {
			http.Error(w, "file not ready", http.StatusInternalServerError)
			return
		}

		pongCount := "0"
		pongData, err := os.ReadFile("/shared/pongs.txt")
		if err == nil {
			pongCount = strings.TrimSpace(string(pongData))
		}

		output := fmt.Sprintf("%s.Ping / Pongs: %s", strings.TrimSpace(string(logData)), pongCount)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, output)
	})

	port := "8080"
	fmt.Println("Server started on port", port)

	http.ListenAndServe(":"+port, nil)
}

