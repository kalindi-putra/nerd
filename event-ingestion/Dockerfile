# Multi-stage build for event-ingestion service
FROM golang:1.25.5 AS builder

# Set build argument for which component to build (server or worker)
ARG COMPONENT=server

WORKDIR /app

# Copy go module files
COPY go.mod go.sum* ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the specified component
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -o app ./pkg/${COMPONENT}/${COMPONENT}.go

# Final stage - use distroless for minimal image size
FROM gcr.io/distroless/base-debian11

WORKDIR /

# Copy the binary from builder
COPY --from=builder /app/app /app

# Expose the gRPC port (only relevant for server)
EXPOSE 50051

# Run the application
CMD ["/app"]