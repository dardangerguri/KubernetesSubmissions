package main

import (
	"fmt"
	"io"
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
		resp, err := http.Get("http://pingpong-app-svc/pings")
		if err == nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				pongCount = strings.TrimSpace(string(body))
			}
		} else {
			fmt.Printf("Error reaching pingpong service: %v\n", err)
		}

		output := fmt.Sprintf("%s.Ping / Pongs: %s\n", strings.TrimSpace(string(logData)), pongCount)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, output)
	})

	port := "8080"
	fmt.Println("Server started on port", port)

	http.ListenAndServe(":"+port, nil)
}

