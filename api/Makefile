GOFILES := $(shell find . -name "*.go")

.PHONY: fmt
fmt:
	goimports -w $(GOFILES)

.PHONY: test
test:
	@go test ./...

.PHONY: build
build:
	@go build -o bin/espresso-chronicle main.go

.PHONY: run
run:
	@go run main.go
