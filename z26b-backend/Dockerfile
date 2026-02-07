# Stage 1: Build Go binary
FROM golang:1.24 AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y git ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Create minimal runtime image
FROM ubuntu:22.04

WORKDIR /app

# Install ca-certificates and tzdata for HTTPS requests and timezone
RUN sed -i '/security.ubuntu.com/d' /etc/apt/sources.list && \
    apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

# Set timezone
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
