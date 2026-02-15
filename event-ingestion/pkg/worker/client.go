package main

import (
	"context"
	"log"
	"os"
	"time"

	api "event-ingestion/pkg/api"

	"google.golang.org/grpc"
)

func main() {
	// Get server address from environment variable or use default
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "localhost:50051"
	}

	log.Printf("Connecting to server at: %s", serverAddr)

	conn, err := grpc.Dial(
		serverAddr,
		grpc.WithInsecure(), // for local dev
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := api.NewEventServiceClient(conn)

	req := &api.EventRequest{
		EventId:   "evt-123",
		Payload:   `{"type":"user_signup"}`,
		Timestamp: time.Now().Unix(),
	}

	resp, err := client.IngestEvent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response: accepted=%v message=%s",
		resp.Accepted, resp.Message)

	ctx := context.Background()
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Fatal("timeout waiting for message")

		case <-ticker.C:
			statusResp, err := client.GetStatus(ctx, &api.StatusRequest{JobId: resp.JobId})

			if err != nil {
				log.Fatal(err)
			}
			log.Printf("job=%s status=%s", resp.JobId, statusResp.Status)

			if statusResp.Status == "done" {
				return
			}

		}
	}
}
