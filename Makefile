.PHONY: build build-windows build-all run rundebug serve format cache proto

FILE ?= main
FOLDER ?= .

# Linux build
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=1 GOOS=linux go build -o bin/flare

# Windows build
build-windows:
	@echo "Building for Windows..."
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o bin/flare.exe

# Production build
build-all: build-linux build-windows

# Development build
build:
	@go build -o bin/flare

# Run a flare file without debug mode, but stack trace enabled
run:
	@$(MAKE) --silent build
	@echo "Running: testcodes/$(FILE).fl"
	@SHOW_STACK=true ./bin/flare run testcodes/$(FILE).fl --debug

# Run a flare file with debug mode enabled
rundebug: build
	@echo "Running: testcodes/$(FILE).fl"
	@DEBUG=true SHOW_STACK=true ./bin/flare run testcodes/$(FILE).fl --debug

# Start the HTTP server on the specified folder
serve: build
	@echo "Starting server on folder: testcodes/$(FOLDER)"
	DEBUG=true ./bin/flare serve testcodes/$(FOLDER) --listenAddr=:3030 --dev

# Format project code (development only; not fully functional)
format: build
	@echo "Formatting code..."
	@DEBUG=true ./bin/flare format testcodes/$(FOLDER)

# Clean project cache
cache:
	@echo "Cleaning cache..."
	@rm -rf .flcache
	@rm -rf bin
	@rm -rf .flmod
	@rm -rf debug

proto:
	@protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative internal/extensions/extension.proto
