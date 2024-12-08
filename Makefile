build-gate:
	@go build -o ./bin/gate ./cmd/gateway/
gate: build-gate
	@./bin/gate

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

docker:
	@docker compose up -d
docker-stop:
	@docker compose down
docker-restart: docker-stop docker 

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/types/ptypes.proto

.PHONY: obu
