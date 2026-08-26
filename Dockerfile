# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/main.go

# Run stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /app/docs ./docs
COPY .env.example .env

EXPOSE 8080

CMD ["./api"]
