GOLANGCI_LINT ?= golangci-lint-v2

.PHONY: default
default: lint test

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run

.PHONY: test
test:
	go test ./... -race -cover

.PHONY: fmt
fmt:
	gofmt -l -w .
	goimports -l -w .

.PHONY: vet
vet:
	go vet ./...