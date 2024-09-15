# SERVER_ADDRESS := localhost:9999
# POLL_INTERVAL := 1
# REPORT_INTERVAL := 1
SERVER_PORT := 4343
ADDRESS := localhost:${SERVER_PORT}
TEMP_FILE := $(shell mktemp)


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
	metricstest -test.v -test.run=^TestIteration5$ \
    	-agent-binary-path=cmd/agent/agent \
        -binary-path=cmd/server/server \
        -server-port=$(SERVER_PORT) \
        -source-path=.

.PHONY: 
	runserver runagent all