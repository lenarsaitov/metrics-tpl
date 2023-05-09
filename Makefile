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
