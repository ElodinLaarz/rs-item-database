.PHONY: all check fmt lint test frontend-check

# Default target
all: check frontend-check

# formatting check (fails if changes needed)
check-fmt:
	@echo "Checking Go formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Go code is not formatted. Run 'make fmt' to fix."; \
		gofmt -l .; \
		exit 1; \
	fi

# Apply formatting
fmt:
	@echo "Formatting Go code..."
	gofmt -w .

# Run Go Vet
vet:
	@echo "Running Go Vet..."
	go vet ./...

# Run Go Tests
test:
	@echo "Running Go Tests..."
	go test -v ./...

# Run all backend checks
check: check-fmt vet test

# Run frontend checks (Build + Typecheck)
frontend-check:
	@echo "Checking Frontend..."
	cd frontend && npm install && npm run build
