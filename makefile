BINARY=git-release
VERSION=$(shell git describe)
GOENV=CGO_ENABLED=0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"
TEST_FLAGS=-cover -count 1

$(BINARY):
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY) $(SRC) ## Build binary (default)

.PHONY: check
check: test lint ## Test and lint

.PHONY: test
unit-test: ## Run go unit tests
	go test $(TEST_FLAGS) ./...

.PHONY: test
test: ## Run full test suite
	go test $(TEST_FLAGS) -tags=integration ./...

.PHONY: lint
lint: ## Run go vet and staticcheck against codebase
	go vet ./...
	staticcheck ./...

.PHONY: clean
clean: ## Clean workspace
	rm -rf $(BINARY)
	go clean

.PHONY: help
help:
	@echo "Available targets:"
	@if [ -t 1 ]; then \
		awk -F ':|##' '/^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST) | grep -v '^help:'; \
	else \
		awk -F ':|##' '/^[a-zA-Z0-9_-]+:.*?##/ { printf "  %-20s %s\n", $$1, $$NF }' $(MAKEFILE_LIST) | grep -v '^help:'; \
	fi
