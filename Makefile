.PHONY: run build build-mac

run:
	@go run .

build:
	@go build .

.ONESHELL:
build-mac:
	@export CGO_ENABLED=0
	@export GOOS=darwin
	@export GOARCH=arm64
	@go build .

.ONESHELL:
build-mac-amd:
	@export CGO_ENABLED=0
	@export GOOS=darwin
	@export GOARCH=amd64
	@go build .