package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "host=postgres-svc port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Database not available yet")
		return
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS counter (
			id INT PRIMARY KEY,
			value INT NOT NULL
		)
	`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		INSERT INTO counter (id, value)
		VALUES (1, 0)
		ON CONFLICT (id) DO NOTHING
	`)
	if err != nil {
		panic(err)
	}
}

func getCounter() int {
	var value int
	err := db.QueryRow(`
		SELECT value FROM counter WHERE id = 1
	`).Scan(&value)
	if err != nil {
		fmt.Printf("DB read error: %v\n", err)
		return 0
	}
	return value
}

func saveCounter(count int) {
	_, err := db.Exec(`
		UPDATE counter SET value = $1 WHERE id = 1
	`, count)
	if err != nil {
		fmt.Printf("DB write error: %v\n", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	initDB()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK")
			return
		}

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

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, "database not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	fmt.Printf("Ping-pong server started in port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err !=nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
