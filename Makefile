PACKAGES := $(shell go list ./...)
name := $(shell basename ${PWD})

# Default make target
all: help

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## init: initialize project (e.g., make init module=github.com/user/project)
.PHONY: init
init:
	@if [ ! -f go.mod ]; then \
		go mod init ${module}; \
	fi
	go install github.com/cosmtrek/air@latest

## vet: vet the code
.PHONY: vet
vet:
	go vet $(PACKAGES)

## build: build a binary
.PHONY: build
build:
	go build -o ./app -v

## docker-build: build project into a Docker container image
.PHONY: docker-build
docker-build:
	GOPROXY=direct docker buildx build -t ${name} .

## docker-run: run project in a Docker container
.PHONY: docker-run
docker-run:
	docker run -it --rm -p 8080:8080 ${name}

## start: build and run local project
.PHONY: start
start: build
	air
