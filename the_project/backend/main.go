package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Todo struct {
	Text string `json:"text"`
}

var db *sql.DB

func initDB() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	for {
		err = db.Ping()
		if err == nil {
			break
		}
		fmt.Println("waiting for postgres...")
		time.Sleep(2 * time.Second)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			text TEXT NOT NULL
		)
	`)
	if err != nil {
		panic(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	initDB()

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			rows, err := db.Query("SELECT text FROM todos")
			if err != nil {
				http.Error(w, "DB error", 500)
				return
			}
			defer rows.Close()

			todos:= []Todo{}
			for rows.Next() {
				var t Todo
				rows.Scan(&t.Text)
				todos = append(todos, t)
			}

			json.NewEncoder(w).Encode(todos)

		case http.MethodPost:
			var newTodo Todo
			err := json.NewDecoder(r.Body).Decode(&newTodo)
			if err != nil || newTodo.Text == "" {
				http.Error(w, "DB request", http.StatusBadRequest)
				return
			}

			_, err = db.Exec("INSERT INTO todos (text) VALUES ($1)", newTodo.Text)
			if err != nil {
				http.Error(w, "DB error", 500)
				return
			}

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
