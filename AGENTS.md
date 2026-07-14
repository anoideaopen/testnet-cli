# testnet-cli

Go CLI for Hyperledger Fabric testnet. Uses Cobra, Viper, Zap logging, and a [custom fork](https://github.com/anoideaopen/fabric-sdk-go) of fabric-sdk-go (`go.mod` replace directive).

## Build

```sh
make build           # go build -mod=vendor (current platform)
make build-all       # cross-compile to linux/darwin/windows, 386/amd64
```

`vendor/` is gitignored. Populate first: `go mod vendor`. Ldflags inject version/commit/date.

## CI pipeline (mandatory order)

1. **check-cyrillic-comments** — grep fails on any Cyrillic chars in source files
2. **validate-go** — `go mod tidy` + `git diff --exit-code` on `go.mod`; then `go fmt ./...` + `go fix ./...` + `git diff --exit-code`
3. **golangci-lint** — v2.9.0, config `.golangci.yml` (v2 format), lints non-test files only
4. **`go test -count 1 ./...`** — **no tests exist**; this passes instantly
5. **`go test ./... -coverprofile=./coverage.out`** — threshold all 0%

## Config

YAML via `--config` / `CLI_CONFIG`. Env prefix `CLI_` overrides file. Panics for on-chain commands without config.

Key flags: `-k` key type (`ed25519`|`secp256k1`|`gost`), `-r` response type, `-w` wait for batch event, `-n` request count, `-t` rate limit.

## Quirks

- **zap**, not logrus: README says logrus, actual code uses `go.uber.org/zap`.
- **osv-scanner-config.toml**: referenced in CI workflow but does not exist on disk.
- **Dead commands**: 9 cobra subcommands (`fetchBatch`, `chaincodeVersion`, `validateBlock`, etc.) are fully implemented but commented out of `root.go` registration.
- **Linter allowances**: revive/staticcheck permit dot-imports for `onsi/gomega` and `onsi/ginkgo/v2`, but neither package is used anywhere.
