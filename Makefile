.PHONY: obu receiver calc invoicer

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