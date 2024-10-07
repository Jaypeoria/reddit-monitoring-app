FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o reddit-fetcher-service .

EXPOSE 8082

CMD ["./reddit-fetcher-service"]
