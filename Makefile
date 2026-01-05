.PHONY: build build-cli build-repl test run-cli run-repl clean tidy

# Output directory
BIN_DIR := bin

# Build both executables
build: build-cli build-repl

# Build CLI interpreter
build-cli:
	@if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	go build -o $(BIN_DIR)/schego.exe ./cmd/schego

# Build REPL
build-repl:
	@if not exist $(BIN_DIR) mkdir $(BIN_DIR)
	go build -o $(BIN_DIR)/schego-repl.exe ./cmd/schego-repl

# Run all tests
test:
	go test -v ./...

# Run CLI interpreter (requires FILE argument)
run-cli: build-cli
	./$(BIN_DIR)/schego.exe $(FILE)

# Run REPL
run-repl: build-repl
	./$(BIN_DIR)/schego-repl.exe

# Clean build artifacts
clean:
	@if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR)

# Tidy dependencies
tidy:
	go mod tidy
