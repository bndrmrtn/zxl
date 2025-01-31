.PHONY: build run

build:
	@go build -o bin/zexlang

run: build
	@DEBUG=true ./bin/zexlang run examples/$(file).zx --debug

serve: build
	@DEBUG=true ./bin/zexlang serve examples/$(folder) --listenAddr=:3030
