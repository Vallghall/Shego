.PHONY: build build-cli build-repl test run-cli run-repl clean tidy fmt lint

# Output directory
BIN_DIR := bin

# Detect OS for executable extension
ifeq ($(OS),Windows_NT)
    EXE_EXT := .exe
    MKDIR := if not exist $(BIN_DIR) mkdir $(BIN_DIR)
    RMDIR := if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR)
else
    EXE_EXT :=
    MKDIR := mkdir -p $(BIN_DIR)
    RMDIR := rm -rf $(BIN_DIR)
endif

# Executable names
CLI_BIN := $(BIN_DIR)/schego$(EXE_EXT)
REPL_BIN := $(BIN_DIR)/schego-repl$(EXE_EXT)

# Build both executables
build: build-cli build-repl

# Build CLI interpreter
build-cli:
	@$(MKDIR)
	go build -o $(CLI_BIN) ./cmd/schego

# Build REPL
build-repl:
	@$(MKDIR)
	go build -o $(REPL_BIN) ./cmd/schego-repl

# Run all tests
test:
	go test -v ./...

# Run CLI interpreter (requires FILE argument)
# Usage: make run-cli FILE=examples/test.scm
run-cli: build-cli
	$(CLI_BIN) $(FILE)

# Run REPL
run-repl: build-repl
	$(REPL_BIN)

# Clean build artifacts
clean:
	@$(RMDIR)

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...
