FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /json_placeholder_client ./cmd/app/
FROM alpine:latest
WORKDIR /
COPY --from=builder /json_placeholder_client .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./json_placeholder_client"]