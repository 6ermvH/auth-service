# Используем официальный образ Go
FROM golang:1.22
WORKDIR /app
COPY internal/ .
RUN go mod download
RUN go build -o auth-service cmd/main.go
EXPOSE 8080
CMD ["./auth-service"]

