.PHONY: all deps test build

all: deps test build

deps:
	@go mod vendor
	@go mod tidy

test:
	@go vet ./{internal,resources}/...
	@go test -v -race -cover ./{cmd,handler}/...

install: 
	@go install .
