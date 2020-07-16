.PHONY: build
build:
	go build -v ./cmd/delta

.PHONY: serve
serve:
	go build -v ./cmd/delta
	./delta

.PHONY: dispatch
dispatch:
	go build -v ./cmd/dispatch
	./dispatch

.PHONY: demon
demon:
	go build -v ./cmd/delta
	./delta -deamon

.PHONY: test
test: 
	go test -v -race -timeout 30s ./...


.DEFAULT_GOAL := build