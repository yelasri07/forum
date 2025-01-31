FROM golang:1.22.5-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o main ./cmd

EXPOSE 8080

CMD ["./main"]