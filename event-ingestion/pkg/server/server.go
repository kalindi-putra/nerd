package main

import (
	"context"
	"database/sql"
	api "event-ingestion/pkg/api"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type EventServer struct {
	api.UnimplementedEventServiceServer
	redisClient *redis.Client
	db          *sql.DB
}

func (s *EventServer) IngestEvent(
	ctx context.Context,
	req *api.EventRequest,
) (*api.EventResponse, error) {

	log.Printf("Received EventLoad: id=%s payload=%s timestamp=%d",
		req.EventId,
		req.Payload,
		req.Timestamp)

	if req.EventId == "" {
		return &api.EventResponse{
			Accepted: false,
			Message:  "No event id received",
		}, nil

	}

	jobId := uuid.NewString()

	if _, err := s.db.ExecContext(
		ctx,
		`INSERT INTO events (job_id, event_id, payload, event_timestamp, status) VALUES ($1, $2, $3, $4, $5)`,
		jobId,
		req.EventId,
		req.Payload,
		req.Timestamp,
		"processing",
	); err != nil {
		log.Printf("failed to insert event: %v", err)
		return &api.EventResponse{
			Accepted: false,
			Message:  "Failed to persist event",
		}, nil
	}

	if err := s.redisClient.Set(ctx, jobId, "processing", 24*time.Hour).Err(); err != nil {
		log.Printf("failed to set processing status in redis: %v", err)
	}

	go func(id string) {
		time.Sleep(3 * time.Second)

		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.redisClient.Set(bgCtx, id, "done", 24*time.Hour).Err(); err != nil {
			log.Printf("failed to set done status in redis for job %s: %v", id, err)
		}

		if _, err := s.db.ExecContext(bgCtx, `UPDATE events SET status = $1, updated_at = NOW() WHERE job_id = $2`, "done", id); err != nil {
			log.Printf("failed to update status in postgres for job %s: %v", id, err)
		}
	}(jobId)

	return &api.EventResponse{
		Accepted: true,
		Message:  "Event received successfully",
		JobId:    jobId,
	}, nil
}

func (s *EventServer) GetStatus(ctx context.Context, req *api.StatusRequest) (*api.StatusResponse, error) {
	if status, err := s.redisClient.Get(ctx, req.JobId).Result(); err == nil {
		return &api.StatusResponse{Status: status}, nil
	}

	var status string
	err := s.db.QueryRowContext(ctx, `SELECT status FROM events WHERE job_id = $1`, req.JobId).Scan(&status)
	if err == nil {
		if cacheErr := s.redisClient.Set(ctx, req.JobId, status, 24*time.Hour).Err(); cacheErr != nil {
			log.Printf("failed to backfill redis status for job %s: %v", req.JobId, cacheErr)
		}
		return &api.StatusResponse{Status: status}, nil
	}

	if err != sql.ErrNoRows {
		log.Printf("failed to read status from postgres: %v", err)
	}

	return &api.StatusResponse{Status: "not found"}, nil
}

func main() {
	postgresHost := getEnv("POSTGRES_HOST", "localhost")
	postgresPort := getEnv("POSTGRES_PORT", "5432")
	postgresUser := getEnv("POSTGRES_USER", "postgres")
	postgresPassword := getEnv("POSTGRES_PASSWORD", "postgres")
	postgresDB := getEnv("POSTGRES_DB", "event_ingestion")

	postgresDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDB,
	)

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := createEventsTable(db); err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "localhost:6379"),
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterEventServiceServer(grpcServer, &EventServer{redisClient: redisClient, db: db})

	log.Println("grpc server running on 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func createEventsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			job_id VARCHAR(36) PRIMARY KEY,
			event_id VARCHAR(255) NOT NULL,
			payload TEXT,
			event_timestamp BIGINT,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}
