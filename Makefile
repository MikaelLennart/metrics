runserver: 
	go run ./cmd/server/main.go

runagent: 
	go run ./cmd/agent/main.go

buildserver:
	go build -o ./cmd/server/server ./cmd/server/main.go

buildagent:
	go build -o ./cmd/agent/agent ./cmd/agent/main.go

clearbin:
	rm -f ./cmd/server/server && rm -f ./cmd/agent/agent

buildall: buildagent buildserver

all: runserver runagent

runtests:
	metricstest -test.v -test.run=^TestIteration2[AB]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent

.PHONY: 
	runserver runagent all