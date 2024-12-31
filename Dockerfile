# Dockerfile for orders-service
FROM golang:1.20-alpine
WORKDIR /app
COPY . .
RUN go build -o orders-service
CMD ["./orders-service"]
