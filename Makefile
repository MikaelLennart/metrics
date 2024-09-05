runserver: 
	go run ./cmd/server/main.go

runagent: 
	go run ./cmd/agent/main.go

all: runserver runagent

.PHONY: 
	runserver runagent all