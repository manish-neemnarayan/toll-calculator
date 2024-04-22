.PHONY: obu receiver

obu:
	go build -o bin/obu obu/main.go
	./bin/obu

receiver:
	go build -o bin/receiver ./dataReceiver 
	./bin/receiver