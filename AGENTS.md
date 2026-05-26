# AGENTS.md

## Project overview
AI-DLC (AI-Driven Development Life Cycle) Go implementation. This repository is the
Go rewrite of github.com/awslabs/aidlc-workflows (Python). Three CLI tools:
`aidlc-designreview`, `aidlc-traceability`, `aidlc-evaluator`.

## Repository structure
aidlc-workflows-py/           # Python submodule (reference only — do not modify)
scripts/
  aidlc-designreview/         # go.mod: AI-powered design review (Go 1.24+)
  aidlc-traceability/         # go.mod: Traceability matrix generator (Go 1.24+)
  aidlc-evaluator/            # go.mod: Evaluation & reporting framework (Go 1.25+)
docs/plans/                   # Implementation plans

## Module paths
- `github.com/mach4-braai/aidlc-workflows/aidlc-traceability`
- `github.com/mach4-braai/aidlc-workflows/aidlc-designreview`
- `github.com/mach4-braai/aidlc-workflows/aidlc-evaluator`

## Build commands
cd scripts/aidlc-traceability && go build ./cmd/traceability/...
cd scripts/aidlc-designreview && go build ./cmd/design-reviewer/...
cd scripts/aidlc-evaluator    && go build ./cmd/aidlc-eval/...

## Test commands
# Unit tests (fast, no external deps)
cd scripts/aidlc-traceability && go test -short ./...
cd scripts/aidlc-designreview && go test -short ./...
cd scripts/aidlc-evaluator    && go test -short ./...

# Full test suite with race detector
cd scripts/aidlc-traceability && go test -race -count=1 -short ./...
cd scripts/aidlc-designreview && go test -race -count=1 -short ./...
cd scripts/aidlc-evaluator    && go test -race -count=1 -short ./...

# Integration tests (requires AWS credentials + Docker daemon)
# Omit -short to include integration tests; they are skipped with -short.

## Release snapshot (all tools)
goreleaser build --snapshot --clean

## Code style
- gofmt/goimports for formatting
- Table-driven tests (t.Run subtests)
- Integration tests gated by testing.Short(): if testing.Short() { t.Skip(...) }
- No comments except for non-obvious WHY

## Key dependencies per module
aidlc-traceability:  cobra, aws-sdk-go-v2, gonum/graph, pterm, backoff/v4, yaml.v3
aidlc-designreview:  cobra, aws-sdk-go-v2, pterm, backoff/v4, yaml.v3
aidlc-evaluator:     cobra, aws-sdk-go-v2, docker/docker, pterm, backoff/v4, yaml.v3

## Embed strategy
Config/prompt/pattern files live in `internal/assets/config/` within each module.
The `internal/assets/assets.go` file embeds the config directory so it is
bundled into the binary at compile time. The `//go:embed` path cannot use `..`.
