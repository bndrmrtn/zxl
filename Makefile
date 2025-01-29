.PHONY: build run

build:
	@go build -o bin/zexlang

run: build
	@./bin/zexlang run examples/$(file).zx --debug
