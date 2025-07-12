# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install git (for go get) and libc for bcrypt
RUN apk add --no-cache git libc6-compat

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o server main.go

# ------------------

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY .env .env

EXPOSE 8080

CMD ["./server"]