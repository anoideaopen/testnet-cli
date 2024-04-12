PROJECT_NAME := "testnet-cli"
PKG := "github.com/anoideaopen/$(PROJECT_NAME)"

OUT_windows_386 := "./output/cli-windows-386.exe"
OUT_windows_amd64 := "./output/cli-windows-amd64.exe"
OUT_darwin_amd64 := "./output/cli-darwin-amd64"
OUT_linux_386 := "./output/cli-linux-386"
OUT_linux_amd64 := "./output/cli-linux-amd64"

PROJECT_DIR=$(PWD)

.PHONY: all
all: build

commitSHA=$(shell git rev-parse HEAD)
dateStr=$(shell date "+%Y-%m-%d_%H:%M:%S")
currentBranch=$(shell git rev-parse --abbrev-ref HEAD)

.PHONY: build
build: ## - build application
	go build -mod=vendor -ldflags "-X main.version=$(currentBranch) -X main.commit=$(commitSHA) -X main.date=$(dateStr)"

.PHONY: build-all
build-all: ## - build application
	GOOS=linux GOARCH=386 go build -o $(OUT_linux_386) -mod=vendor -ldflags "-X main.version=$(currentBranch) -X main.commit=$(commitSHA) -X main.date=$(dateStr)"
	GOOS=linux GOARCH=amd64 go build -o $(OUT_linux_amd64) -mod=vendor -ldflags "-X main.version=$(currentBranch) -X main.commit=$(commitSHA) -X main.date=$(dateStr)"
	GOOS=darwin GOARCH=amd64 go build -o $(OUT_darwin_amd64) -mod=vendor -ldflags "-X main.version=$(currentBranch) -X main.commit=$(commitSHA) -X main.date=$(dateStr)"
	GOOS=windows GOARCH=386 go build -o $(OUT_windows_386) -mod=vendor -ldflags "-X main.version=$(currentBranch) -X main.commit=$(commitSHA) -X main.date=$(dateStr)"
	GOOS=windows GOARCH=amd64 go build -o $(OUT_windows_amd64) -mod=vendor -ldflags "-X main.version=$(currentBranch) -X main.commit=$(commitSHA) -X main.date=$(dateStr)"

.DEFAULT_GOAL := help

.PHONY: help
help: ## - display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
