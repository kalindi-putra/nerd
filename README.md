# Event Ingestion Service

A high-performance event ingestion service built with Go and gRPC. The service provides asynchronous event processing with job tracking and status monitoring.

## Features

- **gRPC API**: High-performance protocol buffer-based communication
- **Asynchronous Processing**: Events are queued and processed asynchronously
- **Job Tracking**: Monitor the status of ingested events via job ID
- **Containerized**: Full Docker and Docker Compose support
- **Scalable**: Designed for horizontal scaling

## Architecture

The system consists of two main components:

1. **Event Server**: Receives events via gRPC and processes them asynchronously
2. **Event Worker/Client**: Sends events to the server and polls for completion

## Prerequisites

- **Docker** and **Docker Compose** (for containerized deployment)
- **Go 1.25.5+** (for local development)
- **Protocol Buffers compiler** (if modifying .proto files)

## Quick Start with Docker

The fastest way to get started is using Docker Compose:

```bash
# Navigate to the project directory
cd event-ingestion

# Build and start all services
make up

# View logs
make logs

# Stop services
make down
```

### Available Make Commands

```bash
make build        # Build Docker images
make up           # Start services (detached)
make up-fg        # Start services (foreground)
make down         # Stop services
make logs         # View all logs
make logs-server  # View server logs only
make logs-worker  # View worker logs only
make restart      # Restart services
make rebuild      # Rebuild and restart
make clean        # Remove everything
make ps           # Show running containers
```

## Detailed Deployment

For detailed deployment instructions, see [deploy/README.md](deploy/README.md).

### Using Docker Compose

```bash
cd deploy
docker-compose up --build
```

This will:
- Build the server and worker images
- Start the gRPC server on port 50051
- Start the worker client that sends test events

### Manual Docker Build

```bash
# Build server
docker build -t event-server:latest -f Dockerfile .

# Build worker
docker build -t event-worker:latest -f Dockerfile.worker .

# Run server
docker run -d -p 50051:50051 --name event-server event-server:latest

# Run worker (connected to server)
docker run -d --link event-server -e SERVER_ADDRESS=event-server:50051 --name event-worker event-worker:latest
```

## Local Development

### Prerequisites

1. Install Go 1.25.5 or higher
2. Install Protocol Buffers compiler (if modifying protos)

### Setup

```bash
# Clone the repository
git clone <repository-url>
cd event-ingestion

# Download dependencies
go mod download

# Run the server
go run pkg/server/server.go

# In another terminal, run the worker
go run pkg/worker/client.go
```

### Regenerate Protocol Buffers (if needed)

```bash
protoc --go_out=. --go-grpc_out=. proto/event.proto
```

## API Reference

### IngestEvent

Ingests an event and returns a job ID for tracking.

**Request:**
```protobuf
message EventRequest {
  string event_id = 1;
  string payload = 2;
  int64 timestamp = 3;
}
```

**Response:**
```protobuf
message EventResponse {
  bool accepted = 1;
  string message = 2;
  string job_id = 3;
}
```

### GetStatus

Checks the status of a job by job ID.

**Request:**
```protobuf
message StatusRequest {
  string job_id = 1;
}
```

**Response:**
```protobuf
message StatusResponse {
  string status = 1;  // "processing", "done", or "not found"
}
```

## Configuration

### Environment Variables

#### Worker/Client

- `SERVER_ADDRESS`: Address of the gRPC server (default: `localhost:50051`)

### Server Configuration

The server listens on port `50051` by default. This can be modified in `pkg/server/server.go`.

## Testing

### Manual Testing with grpcurl

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Send an event
grpcurl -plaintext -d '{
  "event_id": "test-001",
  "payload": "{\"type\":\"user_signup\"}",
  "timestamp": 1234567890
}' localhost:50051 event.EventService/IngestEvent

# Check job status
grpcurl -plaintext -d '{
  "job_id": "<job-id-from-previous-response>"
}' localhost:50051 event.EventService/GetStatus
```

## Project Structure

```
event-ingestion/
├── Dockerfile              # Server container image
├── Dockerfile.worker       # Worker container image
├── Makefile               # Common deployment commands
├── go.mod                 # Go module definition
├── .dockerignore          # Docker build exclusions
├── proto/
│   └── event.proto        # Protocol buffer definitions
├── pkg/
│   ├── api/              # Generated gRPC code
│   │   ├── event.pb.go
│   │   └── event_grpc.pb.go
│   ├── server/           # gRPC server implementation
│   │   └── server.go
│   └── worker/           # Worker client implementation
│       └── client.go
└── deploy/
    ├── docker-compose.yml  # Docker Compose configuration
    └── README.md          # Deployment guide
```

## Troubleshooting

### Server won't start
- Check if port 50051 is already in use
- View logs: `make logs-server` or `docker-compose logs event-server`

### Worker cannot connect
- Ensure server is running: `make ps`
- Check network configuration in docker-compose.yml
- Verify SERVER_ADDRESS environment variable

### Build failures
- Run `make clean` to remove old images
- Rebuild with `make rebuild`

## Production Considerations

Before deploying to production:

1. **Enable TLS**: Configure secure gRPC connections
2. **Add Authentication**: Implement proper auth mechanisms
3. **Configure Logging**: Set up centralized logging
4. **Add Monitoring**: Implement metrics and health checks
5. **Resource Limits**: Set appropriate CPU/memory limits
6. **Load Balancing**: Use a load balancer for multiple instances
7. **Persistent Storage**: Add database for job persistence
8. **Error Handling**: Implement retry logic and dead letter queues

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

[Specify your license here]

## Support

For questions or issues, please open an issue in the repository.
