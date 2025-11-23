# Go Microservices Example (Kafka + Redis + Dedup)


This repository demonstrates a minimal but production‑grade setup with:


- Go microservice with REST → Redis dedup → Kafka producer
- Go consumer microservice reading from Kafka
- Kafka and Redis running on Docker
- Clean folder layout
- Structured logs


## Run Infra
make up
## Run Producer
make run-producer
## Run Consumer
make run-consumer
## Send a test message

curl -X POST localhost:8080/send
-H "Content-Type: application/json"
-d '{"id": "1", "payload": "hello"}'

## Features
- Stateless microservices
- Redis dedup (24h TTL)
- Kafka messaging
- Clean internal package structure
- Minimal dependencies