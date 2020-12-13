GOPATH			?= $(shell go env GOPATH)
GOBIN			?= $(firstword $(subst :, ,${GOPATH}))/bin
GO				?= $(shell which go)

FILES_TO_FMT	?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

.PHONY: deps
deps:
	@$(GO) mod tidy
	@$(GO) mod verify

.PHONY: format
format:
	@echo ">> formatting code"
	@gofmt -s -w $(FILES_TO_FMT)
	@goimports -w $(FILES_TO_FMT)

.PHONY: lint
lint:
	@golangci-lint run --timeout 5m

.PHONY: check-git
check-git:
ifneq ($(GIT),)
	@test -x $(GIT) || (echo >&2 "No git executable binary found at $(GIT)."; exit 1)
else
	@echo >&2 "No git binary found."; exit 1
endif