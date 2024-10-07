FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o reddit-oauth-service .

EXPOSE 8081

CMD ["./reddit-oauth-service"]
