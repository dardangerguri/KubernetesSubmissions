package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const imagePath = "/shared/todo_image.jpg"

func fetchAndCacheImage() error {
	fileInfo, err := os.Stat(imagePath)
	if err == nil {
		if time.Since(fileInfo.ModTime()) < 10 * time.Minute {
			return nil
		}
	}

	fmt.Println("Fetching and caching image...")
	resp, err := http.Get("https://picsum.photos/1200")
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		html := `<!DOCTYPE html>
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
				<input type="text" placeholder="Enter a new todo (max 140 characters)" maxlength="140">
				<button type="button">Send</button>
			</div>

			<h2>Todos</h2>
			<div class="todo-list">
				<div class="todo-item">Learn Kubernetes basics</div>
				<div class="todo-item">Deploy application to cluster</div>
				<div class="todo-item">Configure persistent volumes</div>
			</div>
		</body>
		</html>`

		fmt.Fprint(w, html)
	})

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		err := fetchAndCacheImage()
		if err != nil {
			fmt.Printf("Error fetching image: %v\n", err)
		}
		http.ServeFile(w, r, imagePath)
	})

	fmt.Printf("Server started in port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err !=nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
