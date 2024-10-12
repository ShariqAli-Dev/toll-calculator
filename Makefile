build-obu:
	@go build -o ./bin/obu ./cmd/obu
obu: build-obu
	@./bin/obu

build-reciever:
	@go build -o ./bin/reciever ./cmd/reciever
reciever: build-reciever
	@./bin/reciever


build-calc:
	@go build -o ./bin/distance_calculator ./cmd/distance_calculator
calc: build-calc
	@./bin/distance_calculator


build-aggr:
	@go build -o ./bin/aggregator ./cmd/aggregator
aggr: build-aggr
	@./bin/aggregator

kafka:
	@docker compose up -d
kafka-stop:
	@docker compose down

.PHONY: obu