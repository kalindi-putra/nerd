# Deployment Instructions for Event Ingestion Project

This document provides comprehensive instructions for deploying the Event Ingestion project using Docker and Docker Compose.

## Prerequisites

Ensure you have the following installed on your machine:

- **Docker** (version 20.10 or higher)
- **Docker Compose** (version 1.29 or higher)

## Project Structure

The project is organized as follows:

```
event-ingestion
├── Dockerfile                 # Dockerfile for the gRPC server
├── Dockerfile.worker          # Dockerfile for the worker/client
├── .dockerignore             # Files to exclude from Docker build
├── go.mod                    # Go module dependencies
├── pkg/
│   ├── api/                  # Generated gRPC code
│   ├── server/               # gRPC server implementation
│   └── worker/               # Worker client implementation
├── proto/                    # Protocol buffer definitions
└── deploy/
    ├── docker-compose.yml    # Docker Compose configuration
    └── README.md            # This file
```

## Architecture

The application consists of two components:

1. **Event Server**: A gRPC server that receives event ingestion requests and provides job status
2. **Event Worker**: A client that sends events to the server and polls for job completion

## Quick Start

### 1. Navigate to the deployment directory

```bash
cd event-ingestion/deploy
```

### 2. Build and run with Docker Compose

```bash
docker-compose up --build
```

This command will:
- Build Docker images for both server and worker
- Create a Docker network for inter-container communication
- Start the server on port 50051
- Start the worker which connects to the server

### 3. View logs

```bash
# View all logs
docker-compose logs -f

# View server logs only
docker-compose logs -f event-server

# View worker logs only
docker-compose logs -f event-worker
```

### 4. Stop the services

```bash
docker-compose down
```

## Manual Docker Build and Run

### Build Server Image

```bash
cd event-ingestion
docker build -t event-server:latest -f Dockerfile .
```

### Build Worker Image

```bash
docker build -t event-worker:latest -f Dockerfile.worker .
```

### Run Server Container

```bash
docker run -d \
  --name event-server \
  -p 50051:50051 \
  event-server:latest
```

### Run Worker Container

```bash
docker run -d \
  --name event-worker \
  --link event-server \
  -e SERVER_ADDRESS=event-server:50051 \
  event-worker:latest
```

## Configuration

### Environment Variables

#### Worker Configuration

- `SERVER_ADDRESS`: gRPC server address (default: `localhost:50051`)

Example:
```bash
docker run -e SERVER_ADDRESS=event-server:50051 event-worker:latest
```

### Port Configuration

The server listens on port **50051** by default. You can map it to a different host port:

```bash
docker run -p 8080:50051 event-server:latest
```

## Docker Compose Options

### Scale Workers

You can run multiple worker instances:

```bash
docker-compose up --scale event-worker=3
```

### Run in Detached Mode

```bash
docker-compose up -d
```

### Rebuild Images

```bash
docker-compose build --no-cache
```

## Testing the Deployment

### Using the Worker

The worker automatically sends a test event when it starts and polls for completion.

### Manual Testing with grpcurl

If you have `grpcurl` installed:

```bash
# Send an event
grpcurl -plaintext -d '{
  "event_id": "test-001",
  "payload": "{\"type\":\"test\"}",
  "timestamp": 1234567890
}' localhost:50051 event.EventService/IngestEvent

# Check status (replace JOB_ID with actual job ID from response)
grpcurl -plaintext -d '{
  "job_id": "JOB_ID"
}' localhost:50051 event.EventService/GetStatus
```

## Troubleshooting

### Worker Cannot Connect to Server

- Ensure both containers are on the same network
- Check that `SERVER_ADDRESS` environment variable is set correctly
- Verify the server is running: `docker-compose ps`

### Port Already in Use

If port 50051 is already in use, modify the port mapping in `docker-compose.yml`:

```yaml
ports:
  - "50052:50051"  # Maps host port 50052 to container port 50051
```

### View Container Status

```bash
docker-compose ps
```

### Inspect Logs for Errors

```bash
docker-compose logs event-server | grep -i error
docker-compose logs event-worker | grep -i error
```

### Restart Services

```bash
docker-compose restart
```

## Production Considerations

For production deployment, consider:

1. **TLS/SSL**: Enable secure gRPC connections
2. **Health Checks**: Implement proper health check endpoints
3. **Resource Limits**: Set CPU and memory limits in docker-compose.yml
4. **Logging**: Configure centralized logging (e.g., ELK stack)
5. **Monitoring**: Add Prometheus metrics and Grafana dashboards
6. **Secrets Management**: Use Docker secrets or external secret managers
7. **Load Balancing**: Use a load balancer for multiple server instances

### Example Resource Limits

```yaml
services:
  event-server:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

## Cleanup

To remove all containers, networks, and volumes:

```bash
docker-compose down -v
```

To remove images as well:

```bash
docker-compose down --rmi all -v
```

## Support

For issues or questions, please refer to the main project README or open an issue in the repository.
│   │   └── event.pb.go
│   ├── server
│   │   └── server.go
│   └── worker
│       └── client.go
├── go.mod
├── go.sum
├── Dockerfile
├── .dockerignore
├── deploy
│   ├── docker-compose.yml
│   └── README.md
└── README.md
```

## Steps to Containerize and Deploy

1. **Build the Docker Image**
   Navigate to the root directory of the project (where the Dockerfile is located) and run the following command to build the Docker image:

   ```
   docker build -t event-ingestion .
   ```

2. **Run the Application**
   Change to the `deploy` directory and use Docker Compose to start the application:

   ```
   cd deploy
   docker-compose up
   ```

3. **Access the Application**
   Once the containers are running, you can access the gRPC server as specified in the server implementation. Make sure to check the logs for any output or errors.

## Additional Notes

- Ensure that your `.dockerignore` file is properly configured to exclude unnecessary files from the Docker image.
- The `docker-compose.yml` file defines the services required for the application. Review it to understand how the services are configured and any dependencies that may be needed.
- Update the main `README.md` file with any specific instructions or configurations needed for usage.

By following these steps, you should be able to successfully containerize and deploy the Event Ingestion project.