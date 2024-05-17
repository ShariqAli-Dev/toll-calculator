obu:
	@go build -o ./bin/obu ./obu
	@./bin/obu
reciever:
	@go build -o ./bin/reciever ./data_reciever
	@./bin/reciever
calculator:
	@go build -o ./bin/distance_calculator ./distance_calculator
	@./bin/distance_calculator

.PHONY: obu invoicer