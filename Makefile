APP=che
GOENV := CGO_ENABLED=0 GOOS=linux
GOFLAGS := -ldflags '-w -s'

.PHONY: build
## build: build the application
build: clean
	@echo "Building..."
	@$(GOENV) go build ${GOFLAGS} -o ${APP} main.go

.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@go clean

.PHONY: test
## test: runs go test with default values
test:
	go test -v -count=1 -race ./...

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'