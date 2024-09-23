FROM golang:1.23.1-alpine
WORKDIR /app
COPY . .
RUN touch .env
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main cmd/api/main.go
EXPOSE 8080
CMD ["/app/bin/main"]