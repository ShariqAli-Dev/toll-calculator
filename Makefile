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

kafka:
	@docker compose up -d
kafka-stop:
	@docker compose down

prometheus:
	@docker run --name prometheus -d \
		--add-host=host.docker.internal:host-gateway \
		-p 127.0.0.1:9090:9090 \
		-v /home/shariq/projects/courses/fulltimegodev/toll-calculator/.config/prometheus.yml:/etc/prometheus/prometheus.yml \
		prom/prometheus --config.file=/etc/prometheus/prometheus.yml

prometheus-stop: 
	@docker stop prometheus
	@docker rm prometheus

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/types/ptypes.proto

.PHONY: obu
