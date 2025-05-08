FROM golang:1.24.1 AS builder

# Install C dependencies needed for CGO and H3
RUN apt-get update && apt-get install -y gcc libc6-dev

WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN GOOS=linux go build main.go

# Final stage
FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y \
	libc6 \
	ca-certificates \
	&& rm -rf /var/lib/apt/lists/*  # Clean up to reduce image size

# Copy only the built binary from the builder stage
COPY --from=builder /app/main /app/main

COPY . .

EXPOSE 8080

#CMD ["/app/main"]
