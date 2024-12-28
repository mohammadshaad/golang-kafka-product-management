# Variables
GO_TEST=go test

# Targets
.PHONY: all unit integration benchmark clean help

# Default target
all: unit integration benchmark

# Run unit tests
unit:
	$(GO_TEST) ./internal/...

# Run integration tests
integration:
	$(GO_TEST) ./tests/integration/...

# Run benchmarks
benchmark:
	$(GO_TEST) -bench=. ./tests/performance/...

# Clean up temporary test files or artifacts
clean:
	@echo "Cleaning up test artifacts..."
	@rm -rf *.test

# Show help
help:
	@echo "Usage:"
	@echo "  make unit         Run unit tests."
	@echo "  make integration  Run integration tests."
	@echo "  make benchmark    Run benchmarks."
	@echo "  make clean        Clean up temporary files."
	@echo "  make help         Show this help message."
	@echo "  make all          Run all tests and benchmarks."
