.PHONY: build run serve format

build:
	@go build -o bin/zxlang

run: build
	@DEBUG=false ./bin/zxlang run examples/$(file).zx --debug

serve: build
	@DEBUG=true ./bin/zxlang serve examples/$(folder) --listenAddr=:3030

format: build
	@DEBUG=true ./bin/zxlang format examples/format
