# stage 1
FROM golang:1.25.6-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# stage 2
FROM alpine:latest
WORKDIR /root
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["./server"]
