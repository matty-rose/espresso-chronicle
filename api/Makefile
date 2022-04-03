APP_DIR := app
GOFILES := $(shell find . -name "*.go")

.PHONY: fmt
fmt:
	goimports -w $(GOFILES)

.PHONY: test
test:
	@go test ./...

.PHONY: lint
lint:
	@golangci-lint run $(APP_DIR)

.PHONY: build
build:
	@go build -o bin/espresso-chronicle $(APP_DIR)/main.go

.PHONY: run
run:
	@go run $(APP_DIR)/main.go

.PHONY: docker/dev
docker/dev:
	@docker compose up --build -d

.PHONY: docker/emulator
docker/emulator:
	@docker compose up firestore --build -d

.PHONY: docker/logs
docker/logs:
	@docker compose logs -f api
