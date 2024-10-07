FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o reddit-statistics-service .

EXPOSE 8083

CMD ["./reddit-statistics-service"]
