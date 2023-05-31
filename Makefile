## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run/server: run server
run/server:
	go run ./cmd/server/main.go

## run/agent: run agent
run/agent:
	go run ./cmd/agent/main.go

## test: run all tests
test:
	go test -cover ./...

## test: run all autotests
autotest:
	go build -o cmd/server/server cmd/server/*.go
	go build -o cmd/agent/agent cmd/agent/*.go
	./metricstest -test.v -test.run=^TestIteration7$  -agent-binary-path=cmd/agent/agent -binary-path=cmd/server/server -source-path=.
