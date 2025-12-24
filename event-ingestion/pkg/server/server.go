package main

import (
	"context"
	api "event-ingestion/pkg/api"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc"
)

type EventServer struct {
	api.UnimplementedEventServiceServer
	jobs sync.Map
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
			Message:  "No event id recieved",
		}, nil

	}

	var jobId = uuid.NewString()
	s.jobs.Store(jobId, "processing")

	go func(id string) {
		time.Sleep(3 * time.Second)
		s.jobs.Store(id, "done")
	}(jobId)

	return &api.EventResponse{
		Accepted: true,
		Message:  "Event recieved successfully for JOBID",
		JobId:    jobId,
	}, nil
}

func (s *EventServer) GetStatus(ctx context.Context, req *api.StatusRequest) (*api.StatusResponse, error) {
	if v, ok := s.jobs.Load(req.JobId); ok {
		if status, ok := v.(string); ok {
			return &api.StatusResponse{Status: status}, nil
		}
	}
	return &api.StatusResponse{Status: "not found"}, nil
}

func main() {

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterEventServiceServer(grpcServer, &EventServer{})

	log.Println("grpc server running on 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
