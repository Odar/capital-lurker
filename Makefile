APP?=lurker
BIN?=./bin/$(APP)
GO?=go
LOCAL_BIN:=$(CURDIR)/bin

.PHONY: build
build:
	$(GO) build -o $(BIN) ./cmd/$(APP)

.PHONY: run
run: build
	$(BIN)

.PHONY: test
test:
	$(GO) test ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: .deps
.deps:
	$(info #Install dependencies...)
	go mod tidy

# install project dependencies
.PHONY: deps
deps: .deps

.PHONY: .generate
.generate:
	${GO} generate ./...

.PHONY: generate
generate: .bin-deps .generate

.PHONY: migrate
migrate: migrate-deps migrate-build
	bin/migrate --db=capital | xargs bin/goose

.PHONY: migrate-build
migrate-build:
	$(GO) build -o ./bin/migrate ./tools/migrate

.PHONY: migrate-deps
migrate-deps:
	GO111MODULE=off GOBIN=$(LOCAL_BIN) go get -u -v github.com/pressly/goose/cmd/goose