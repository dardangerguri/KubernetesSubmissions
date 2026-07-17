package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type WebhookPayload struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func main() {
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}

	webhookUrl := os.Getenv("WEBHOOK_URL")
	if webhookUrl == "" {
		fmt.Println("Warning: WEBHOOK_URL environment variable is not set.")
	}

	var nc *nats.Conn
	var err error

	for i := 0; i < 10; i++ {
		nc, err = nats.Connect(natsUrl)
		if err == nil {
			break
		}
		fmt.Println("Broadcaster waiting for NATS...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	fmt.Println("Broadcaster connected to NATS")

	_, err = nc.QueueSubscribe("todos", "todos-group", func(m *nats.Msg) {
		messageText := string(m.Data)
		fmt.Printf("Received message: %s\n", messageText)

		if webhookUrl == "" {
			return
		}

		payload := WebhookPayload{
			User:    "bot",
			Message: messageText,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("Error marshalling payload: %v\n", err)
			return
		}

		resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Printf("Error sending to webhook: %v\n", err)
			return
		}
		resp.Body.Close()
		fmt.Println("Message successfully broadcasted forward.")
	})
	if err != nil {
		panic(err)
	}

	select {}
}