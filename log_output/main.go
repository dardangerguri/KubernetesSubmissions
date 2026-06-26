package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	id := uuid.New().String()

	fmt.Println("App started with ID:", id)

	for {
		fmt.Printf("%s: %s\n",
			time.Now().UTC().Format(time.RFC3339Nano),
			id,
		)
		time.Sleep(5 * time.Second)
	}
}
