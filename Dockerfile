# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o osohub-backend ./cmd/main.go
RUN ls -lh /app

# Final image
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/osohub-backend .
COPY .env .env
## Firebase config not needed if not using Firebase

EXPOSE 8080

CMD ["./osohub-backend"]
