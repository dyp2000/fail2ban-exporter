#!/usr/bin/env make

.PHONY : help all compile clear

.DEFAULT_GOAL := help

OUT_DIR ?= .
APP ?= fail2ban-exporter
LINUX ?= -linux-amd64
DARWIN ?= -darwin-amd64

help: ## Show this help
	@echo "Make Fail2ban-exporter project"
	@echo "Copyright © 2020 Dennis Y. Parygin (dyp2000@mail.ru)"
	@echo
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[33m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: clear packages release ## Full build project

clear: ## Clear caches, objects fns other
	go clean
	rm -f ${OUT_DIR}/${APP}-${DARWIN}
	rm -f ${OUT_DIR}/${APP}-${LINUX}

packages: ## Install required packages
	go mod vendor -v

debug: ## Compile project with debug info
	cd ./src && go build -v -o ${OUT_DIR}/${APP}

release: ## Compile project Release
	env GOOS=darwin GOARCH=amd64 go build -v -ldflags "-s -w" -o ${OUT_DIR}/${APP}${DARWIN} -i ./cmd/${APP}
	env GOOS=linux GOARCH=amd64 go build -v -ldflags "-s -w" -o ${OUT_DIR}/${APP}${LINUX} -i ./cmd/${APP}
