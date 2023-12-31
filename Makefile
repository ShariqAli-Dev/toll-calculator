obu:
	@go build -o bin/obu ./obu
	@./bin/obu
receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver
calculator:
	@go build -o bin/distance_calculator ./distance_calculator
	@./bin/distance_calculator

.PHONY: obu invoicer