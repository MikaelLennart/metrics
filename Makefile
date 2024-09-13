runserver: 
	go run ./cmd/server/main.go

runagent: 
	go run ./cmd/agent/main.go

runall: runserver runagent

buildserver:
	go build -o ./cmd/server/server ./cmd/server/main.go

buildagent:
	go build -o ./cmd/agent/agent ./cmd/agent/main.go

clearbin:
	rm -f ./cmd/server/server && rm -f ./cmd/agent/agent

buildall: buildagent buildserver

all: runserver runagent

runtests:
	SERVER_PORT := 4343
    ADDRESS := "localhost:${SERVER_PORT}"
    TEMP_FILE := $(shell mktemp)
    metricstest -test.v -test.run=^TestIteration4$ \
    	-agent-binary-path=cmd/agent/agent \
        -binary-path=cmd/server/server \
        -server-port=$(SERVER_PORT) \
        -source-path=.

.PHONY: 
	runserver runagent all