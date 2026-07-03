package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const imagePath = "/shared/todo_image.jpg"

func getBackendURL() string {
	url := os.Getenv("BACKEND_URL")
	if url == "" {
		return "http://todo-backend-svc"
	}
	return url
}

func getImageURL() string {
	url := os.Getenv("IMAGE_URL")
	if url == "" {
		return "https://picsum.photos/1200"
	}
	return url
}

type Todo struct {
	Text string `json:"text"`
}

func fetchAndCacheImage() error {
	fileInfo, err := os.Stat(imagePath)
	if err == nil {
		if time.Since(fileInfo.ModTime()) < 10 * time.Minute {
			return nil
		}
	}

	fmt.Println("Fetching and caching image...")
	resp, err := http.Get(getImageURL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getTodosFromBackend(backendUrl string) []Todo {
	var list []Todo
	resp, err := http.Get(backendUrl)
	if err != nil {
		fmt.Printf("Error pulling from backend: %v\n", err)
		return []Todo{{Text: "Backend unreachable"}}
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&list)
	return list
}

func isBackendHealthy() bool {
	resp, err := http.Get("http://todo-backend-svc/healthz")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func main() {
	backendUrl := getBackendURL()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if !isBackendHealthy() {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			htmlFailure := `<!DOCTYPE html>
			<html>
			<head>
				<title>Todo App - Failure</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						text-align: center;
						margin-top: 100px;
						background-color: #fff5f5;
					}
					.error-container {
						border: 2px solid #feb2b2;
						background-color: #fff5f5;
						padding: 40px;
						display: inline-block;
						border-radius: 8px;
					}
					h1 {
						color: #9b2c2c;
						font-size: 3rem;
						margin-bottom: 10px;
					}
					p {
						color: #c53030;
						font-size: 1.2rem;
					}
				</style>
			</head>
			<body>
				<div class="error-container">
					<h1>System Failure</h1>
					<p>The Todo App is currently unhealthy. Please wait for recovery.</p>
				</div>
			</body>
			</html>`

			fmt.Fprint(w, htmlFailure)
			return
		}

		currentTodos := getTodosFromBackend(backendUrl)
		var todoItemsHTML string
		for _, todo := range currentTodos {
			todoItemsHTML += fmt.Sprintf(`<div class="todo-item">%s</div>`, todo.Text)
		}

		html := fmt.Sprintf(`<!DOCTYPE html>
		<html>
		<head>
			<title>Todo App</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					text-align: center;
					margin-top: 50px;
				}
				h1 {
					font-size: 2.5rem;
					margin-bottom: 20px;
				}
				h2 {
					font-size: 1.8rem;
					margin-top: 30px;
					margin-bottom: 20px;
				}
				img {
					max-width: 400px;
					border-radius: 12px;
					margin-bottom: 20px;
				}
				.form-container {
					margin-bottom: 30px;
				}
				input[type="text"] {
					width: 400px;
					padding: 12px;
					font-size: 1rem;
					border: 2px solid #28a745;
					border-radius: 6px;
					outline:none;
				}
				button {
					padding: 12px 24px;
					font-size: 1rem;
					background-color: #28a745;
					color: white;
					border: none;
					border-radius: 6px;
					cursor: pointer;
					margin-left: 10px;
				}
				button:hover {
					background-color: #218838;
				}
				.todo-list {
					max-width: 600px;
					margin: 0 auto;
					text-align: left;
				}
				.todo-item {
					background-color: #f8f9fa;
					padding: 15px;
					margin-bottom: 10px;
					border-left: 5px solid #28a745;
					border-radius: 4px;
					font-size: 1.1rem;
				}
			</style>
		</head>
		<body>
			<h1>Todo App</h1>
			<div>
				<img src="/image" alt="Todo Image">
			</div>

			<div class="form-container">
				<form action="/create" method="POST">
					<input type="text" name="todo" placeholder="Enter a new todo (max 140 characters)" maxlength="140" required>
					<button type="submit">Send</button>
				</form>
			</div>

			<h2>Todos</h2>
			<div class="todo-list">
				%s
			</div>
			<form action="/break" method="POST" style="margin-top:20px;">
				<button type="submit" style="background-color:#dc3545;">Break application</button>
			</form>
		</body>
		</html>`, todoItemsHTML)

		fmt.Fprint(w, html)
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		todoText := r.FormValue("todo")
		if todoText != "" {
			todoObj := Todo{Text: todoText}
			jsonData, _ := json.Marshal(todoObj)

			http.Post(getBackendURL()+"/todos", "application/json", bytes.NewBuffer(jsonData))
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/break", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		baseUrl := "http://todo-backend-svc"
		envUrl := os.Getenv("BACKEND_URL")

		if envUrl != "" {
			baseUrl = envUrl
			if len(baseUrl) > 6 && baseUrl[len(baseUrl)-6:] == "/todos" {
				baseUrl = baseUrl[:len(baseUrl)-6]
			}
		}

		http.Post(baseUrl+"/break", "text/plain", nil)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		_ = fetchAndCacheImage()
		http.ServeFile(w, r, imagePath)
	})

	fmt.Printf("Frontend server started in port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting frontend: %v\n", err)
	}
}