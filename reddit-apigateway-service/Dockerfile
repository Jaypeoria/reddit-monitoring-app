FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o reddit-apigateway-service .

EXPOSE 8080

CMD ["./reddit-apigateway-service"]
