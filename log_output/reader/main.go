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

		messageEnv := os.Getenv("MESSAGE")

		configFileDir, err := os.ReadFile("/config/information.txt")
		fileContent := "file not found"
		if err == nil {
			fileContent = strings.TrimSpace(string(configFileDir))
		}

		pongCount := "0"
		resp, err := http.Get("http://pingpong-app-svc.exercises.svc.cluster.local/pings")
		if err == nil {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				pongCount = strings.TrimSpace(string(body))
			}
		} else {
			fmt.Printf("Error reaching pingpong service: %v\n", err)
		}

		output := fmt.Sprintf("file content: %s\nenv variable: MESSAGE=%s\n%s.Ping / Pongs: %s\n",
			fileContent, messageEnv, strings.TrimSpace(string(logData)), pongCount)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, output)
	})

	port := "8080"
	fmt.Println("Server started on port", port)

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://pingpong-app-svc.exercises.svc.cluster.local/pings")
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "pingpong not ready", http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	http.ListenAndServe(":"+port, nil)
}

