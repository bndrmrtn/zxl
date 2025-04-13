.PHONY: build run rund serve format cache

# build the project
build:
	@go build -o bin/flare

# run the project
run: build
	@SHOW_STACK=true ./bin/flare run testcodes/$(file).fl --debug

# rund the project
rund: build
		@DEBUG=true SHOW_STACK=true ./bin/flare run testcodes/$(file).fl --debug

# start the default http server
serve: build
	@DEBUG=true ./bin/flare serve testcodes/$(folder) --listenAddr=:3030

# format the project's code (dev only, not working properly)
format: build
	@DEBUG=true ./bin/flare format testcodes/format

# clean the project's cache
cache:
	@rm -rf .flcache
