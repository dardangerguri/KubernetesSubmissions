package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("/shared/log.txt")
		if err != nil {
			http.Error(w, "file not ready", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(data))
	})


	port := "8080"
	fmt.Println("Server started on port", port)

	http.ListenAndServe(":"+port, nil)
}

