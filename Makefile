BINARY  ?= payme
CMD     ?= ./cmd/payme
BIN_DIR ?= bin

.PHONY: all build test lint vet tidy vendor deps clean
.PHONY: release-patch release-minor release-major release

all: build

build:
	go build -o $(BIN_DIR)/$(BINARY) $(CMD)

test:
	go test ./...

lint:
	golangci-lint run --allow-parallel-runners

vet:
	go vet ./...

tidy:
	go mod tidy

vendor: tidy
	go mod vendor

# Update, tidy and vendor all dependencies.
deps: vendor

clean:
	rm -rf $(BIN_DIR)

release-patch release-minor release-major:
	$(MAKE) release VERSION=$(shell go run github.com/mdomke/git-semver/v6@latest -target $(subst release-,,$@))

release:
	git tag -s -a $(VERSION) -m "Release $(VERSION)"
	@echo "Now run:"
	@echo "- git push --tags"
	@echo "- gh release create --generate-notes $(VERSION)"
