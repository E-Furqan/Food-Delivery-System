FROM golang:1.23.1 as builder

WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application code
COPY . .

# Build the Go application
RUN go build -o user_service ./main.go

# Final stage: create a minimal image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/user_service .

# Set environment variables
ENV DATABASE_HOST=db
ENV DATABASE_PORT=5432
ENV DATABASE_USER=furqan
ENV DATABASE_PASSWORD=furqan
ENV DATABASE_NAME=User

# Run the application
CMD ["./user_service"]
