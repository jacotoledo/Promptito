.PHONY: build build-gui build-linux test lint clean run docker-build docker-run install service-start service-stop help

BINARY_NAME=promptito
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u)
GOFLAGS=-ldflags="-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"

help: ## Show this help message
	@echo "Promptito Makefile"
	@echo "=================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary (console version)
	@echo "Building ${BINARY_NAME}..."
	go build ${GOFLAGS} -o bin/${BINARY_NAME}.exe ./cmd/promptito
	@echo "Built: bin/${BINARY_NAME}.exe"

build-gui: ## Build Windows GUI version (no console window)
	@echo "Building ${BINARY_NAME} GUI..."
	go build -ldflags="-H=windowsgui ${GOFLAGS}" -o bin/${BINARY_NAME}.exe ./cmd/promptito
	@echo "Built: bin/${BINARY_NAME}.exe (GUI mode)"

build-linux: ## Build for Linux (cross-compile)
	@echo "Building ${BINARY_NAME} for Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${GOFLAGS} -o bin/${BINARY_NAME}-linux-amd64 ./cmd/promptito
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build ${GOFLAGS} -o bin/${BINARY_NAME}-linux-arm64 ./cmd/promptito
	@echo "Built: bin/${BINARY_NAME}-linux-*"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "Cleaned"

run: build ## Build and run
	@echo "Running ${BINARY_NAME}..."
	./bin/${BINARY_NAME}.exe -prompts public/prompts -static public

run-dev: ## Run without building (requires pre-built binary)
	./bin/${BINARY_NAME}.exe -prompts public/prompts -static public

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t promptito:latest .

docker-run: docker-build ## Build and run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --rm promptito:latest

docker-stop: ## Stop Docker container
	docker stop promptito || true

install: ## Install binary to GOPATH
	go install ${GOFLAGS} ./cmd/promptito

service-start: ## Start as Windows service (requires admin)
	@echo "Installing and starting Promptito service..."
	@if exist "%SystemRoot%\System32\sc.exe" ( \
		sc create Promptito binPath= "%cd%\bin\promptito.exe -prompts public/prompts -static public" start= auto && \
		sc start Promptito \
	) else ( \
		echo SC.exe not found. Run as administrator to install service. \
	)

service-stop: ## Stop Windows service
	@echo "Stopping Promptito service..."
	@if exist "%SystemRoot%\System32\sc.exe" ( \
		sc stop Promptito && sc delete Promptito \
	)

deps: ## Download dependencies
	go mod download
	go mod verify

tidy: ## Tidy go.mod
	go mod tidy

all: clean fmt vet test build ## Run all checks and build

.PHONY: build build-gui build-linux test lint clean run docker-build docker-run install service-start service-stop help
