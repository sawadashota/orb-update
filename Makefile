SHELL := /bin/bash
.DEFAULT_GOAL := help

init: ## prepare to develop
	go mod download
	go mod verify

docker: ## build docker image
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build
	docker build -t sawadashota/orb-update:latest .
	rm orb-update

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -h -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

# https://postd.cc/auto-documented-makefile/
help: ## Show help message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

