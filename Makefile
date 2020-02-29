.PHONY: build run help
.DEFAULT_GOAL := build

PROJECTNAME := $(shell basename "$(PWD)")

## build: build project
build:
	go build -v ./cmd/$(PROJECTNAME) -o ./build/$(PROJECTNAME)

## run: run project
run:
	./build/$(PROJECTNAME)

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'