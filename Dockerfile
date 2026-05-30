FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/server ./cmd/app

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/internal/src /app/internal/src

EXPOSE 8080

CMD ["./server"]
