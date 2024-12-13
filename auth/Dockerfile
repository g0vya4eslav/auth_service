FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth_service main.go
EXPOSE 50051 9000
CMD ["./auth_service"]






















