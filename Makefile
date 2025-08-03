.PHONY: up down build-gen run-gen logs

up:
	docker compose up --build -d

down:
	docker compose down -v

build-gen:
	go build -o generator ./cmd/app

run-gen:
	go run ./cmd/app/main.go

logs:
	docker compose logs -f kafka-consumer