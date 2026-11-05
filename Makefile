APP_NAME := wosaas
CONF ?= ./conf/app.toml
BIN := bin/$(APP_NAME)

.PHONY: all build run test fmt tidy clean

all: build

build:
	go build -o $(BIN) .

run:
	go run . -conf $(CONF)

test:
	go test ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

clean:
	rm -rf bin
