set shell := ["/usr/bin/env", "bash", "-c"]

@build:
    go build

@test:
    go test ./...

@fmt:
    go fmt ./...

@mocks:
    ./generate_mocks.sh

@lint:
    go vet
    gosec ./...
    golangci-lint run

# Run before merging PR
@checks: fmt mocks test lint
    # Fail if files not up-to-date (particularly mocks)
    [[ $(git diff --name-only | wc --lines) -lt 1 ]] || (echo "Uncommitted files found" && exit 1)

# Requires a git tag first
@release: checks
    goreleaser release
