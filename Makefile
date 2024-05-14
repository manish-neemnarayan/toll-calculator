.PHONY: obu receiver calc invoicer proto

obu:
	go build -o bin/obu obu/main.go
	./bin/obu
	
receiver:
	go build -o bin/receiver ./dataReceiver 
	./bin/receiver

calc:
	go build -o bin/calculate ./distanceCalculator 
	./bin/calculate

agg:
	go build -o bin/aggregator ./aggregator 
	./bin/aggregator

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto