# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Final stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata             # Recommended, if your app makes HTTPS requests
COPY --from=builder /app/main ./server
COPY .env .
EXPOSE 8080
ENTRYPOINT ["./server"]

