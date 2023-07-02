# Stage 1: Build the Go binary
FROM golang:1.20-alpine3.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api

# Stage 2: Create a minimal image to run the Go binary
FROM alpine:latest

WORKDIR /app

RUN adduser -D api && chown api:api /app

RUN apk --no-cache add ca-certificates

USER api

COPY --from=builder --chown=api:api /app/migrations migrations
COPY --from=builder --chown=api:api /app/api .

CMD ["./api"]
