# Contributing to repos-sync

## Prerequisites

- Linux
- Go 1.27+ (`go version`)
- git

## Build, Test, Gate

```sh
make build          # ./bin/repos-sync
make test           # go test ./...
go test -race ./... # fuller gate when you have cycles
make vet            # go vet ./...
make lint           # gofmt -l + go vet
```

CI runs `.github/workflows/ci.yml` (build/vet/test/gofmt + golangci-lint),
currently blocked by GitHub Actions billing on this account; the local gate
is authoritative until it is resolved.

## Commits

Conventional Commits (`feat:`, `fix:`, `docs:`, `test:`, `build:`, `chore:`),
one logical change per commit.

## Rules

1. Stdlib-first: no new dependency without justification in the PR.
2. No `sh -c`; always `exec.Command` with an argument array.
3. Every read-only claim must hold under fixtures: repos-sync never writes to
   the repositories it inspects. Do not break that guarantee.
4. Tests generate fixtures (`t.TempDir()` + real `git`), never fixture files.
5. `schema_version` is `1`; do not rename/remove existing JSON keys. Adding
   keys is minor; removing or renaming is breaking.

## Layout

`cmd/repos-sync` entry · `internal/git` git wrappers · `internal/scan` root
walking · `internal/model` report + sorting · `internal/output` renderers ·
`test/integration` end-to-end binary semantics.
