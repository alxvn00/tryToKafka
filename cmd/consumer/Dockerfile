FROM golang:1.23

RUN apt-get update && apt-get install -y librdkafka-dev librdkafka1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o consumer ./cmd/consumer

COPY .env .env

CMD ["./consumer"]