GOBIN=./cmd

## lint: Run golangci-lint for project
.PHONY:
lint:
	golangci-lint run

## build: Build go binary
build:
	go build -o $(GOBIN)/main $(GOBIN)/main.go

## run: Run go server
run:
	./local.sh

## get: Run go get missing dependencies
get:
	go get ./...

## test: Run all tests in project
test:
	go test -v -race -cover -bench=. ./...

## deploy: Run commands to deploy apprepo to container
deploy:
	get
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
          -ldflags='-w -s -extldflags "-static"' -a \
          -o $(GOBIN)/main $(GOBIN)/main.go

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.DEFAULT_GOAL := help
