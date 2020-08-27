#!/usr/bin/env make

.PHONY : help all compile clear

.DEFAULT_GOAL := help

OUT_DIR ?= ../bin
APP ?= fail2ban-exporter

help: ## Show this help
	@echo "Make Fail2ban-exporter project"
	@echo "Copyright © 2020 Dennis Y. Parygin (dyp2000@mail.ru)"
	@echo
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[33m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: clear packages release ## Full build project

debug: ## Compile project with debug info
	cd ./src && go build -v -o ${OUT_DIR}/${APP}

release: ## Compile project Release
	cd ./src && go build -v -ldflags "-s -w" -o ${OUT_DIR}/${APP}

clear: ## Clear caches, objects fns other
	go clean
	rm -f ./bin/*

packages: ## Install required packages
	go get -v -u github.com/prometheus/client_golang/prometheus
	go get -v -u github.com/prometheus/client_golang/prometheus/promauto
	go get -v -u github.com/prometheus/client_golang/prometheus/promhttp
