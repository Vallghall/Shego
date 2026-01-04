.PHONY: build build-cli build-repl test run-cli run-repl clean tidy

# Output directory
BIN_DIR := bin

# Build both executables
build: build-cli build-repl

# Build CLI interpreter
build-cli:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/schego ./cmd/schego

# Build REPL
build-repl:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/schego-repl ./cmd/schego-repl

# Run all tests
test:
	go test -v ./...

# Run CLI interpreter (requires FILE argument)
run-cli: build-cli
	./$(BIN_DIR)/schego $(FILE)

# Run REPL
run-repl: build-repl
	./$(BIN_DIR)/schego-repl

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)

# Tidy dependencies
tidy:
	go mod tidy

