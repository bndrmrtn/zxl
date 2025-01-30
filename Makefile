.PHONY: build run

build:
	@go build -o bin/zexlang

run: build
	@DEBUG=true ./bin/zexlang run examples/$(file).zx --debug
