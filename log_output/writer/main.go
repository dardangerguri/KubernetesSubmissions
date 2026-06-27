package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	id := uuid.New().String()
	file := "/shared/log.txt"

	for {
		f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("error opening file:", err)
			continue
		}

		fmt.Fprintf(f, "%s: %s\n", time.Now().UTC().Format(time.RFC3339Nano), id)
		f.Close()

		time.Sleep(5 * time.Second)
	}
}
