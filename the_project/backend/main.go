package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Todo struct {
	Text string `json:"text"`
}

var (
	todos	=[]Todo{{Text: "Learn Kubernetes basics"}, {Text: "Deploy application to cluster"}}
	todosMu	sync.Mutex
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			todosMu.Lock()
			json.NewEncoder(w).Encode(todos)
			todosMu.Unlock()

		case http.MethodPost:
			var newTodo Todo
			err := json.NewDecoder(r.Body).Decode(&newTodo)
			if err != nil || newTodo.Text == "" {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			todosMu.Lock()
			todos = append(todos, newTodo)
			todosMu.Unlock()

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newTodo)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Printf("Backend server started in port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err !=nil {
		fmt.Printf("Error starting backend: %v\n", err)
	}
}
