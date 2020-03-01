.PHONY: build run client help
.DEFAULT_GOAL := build

PROJECTNAME := $(shell basename "$(PWD)")

## build: build project
build:
	go build -o ./build/$(PROJECTNAME) -v ./cmd/$(PROJECTNAME)
	go build -o ./build/client -v ./cmd/client

## run: run project
run:
	./build/$(PROJECTNAME)

client:
	./build/client

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'