# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o redirix ./cmd/redirix

# Stage 2: Runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/redirix .

ENTRYPOINT ["./redirix"]