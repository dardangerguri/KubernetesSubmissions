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
				img {
					max-width: 400px;
					border-radius: 12px;
					margin-bottom: 20px;
				}
				p {
					color: #555;
				}
			</style>
		</head>
		<body>
			<h1>Todo App</h1>
			<div>
				<img src="/image" alt="Todo Image">
			</div>
			<p>DevOps with Kubernetes</p>
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
