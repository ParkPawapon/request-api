.PHONY: run build test test-race fmt vet tidy lint

run:
	go run ./cmd/api

build:
	go build ./cmd/api

test:
	go test ./...

test-race:
	go test -race ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

vet:
	go vet ./...

tidy:
	go mod tidy

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint is not installed; run 'go vet ./...' as the baseline lint gate."; \
		exit 127; \
	fi
