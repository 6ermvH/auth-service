# Используем официальный образ Go
FROM golang:1.22
WORKDIR /app
COPY internal/go.mod internal/go.sum ./
RUN go mod download
COPY internal/ .
RUN go build -o auth-service main.go
EXPOSE 8080
CMD ["./auth-service"]

