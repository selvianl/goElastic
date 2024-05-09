FROM golang:1.21-alpine

WORKDIR /app
COPY . .

RUN go build -o bin cmd/insider/main.go

ENTRYPOINT ["/app/bin"]