# AGENTS.md

## Project overview
AI-DLC (AI-Driven Development Life Cycle) Go implementation. This repository is the
Go rewrite of github.com/awslabs/aidlc-workflows (Python). Three CLI tools:
`aidlc-designreview`, `aidlc-traceability`, `aidlc-evaluator`.

## Repository structure
aidlc-workflows-py/     # Python submodule (reference only — do not modify)
scripts/
  aidlc-designreview/   # go.mod: AI-powered design review
  aidlc-traceability/   # go.mod: Traceability matrix generator
  aidlc-evaluator/      # go.mod: Evaluation & reporting framework
docs/plans/             # Implementation plans

## Setup commands
# Build all tools
cd scripts/aidlc-traceability && go build ./cmd/traceability/...
cd scripts/aidlc-designreview && go build ./cmd/design-reviewer/...
cd scripts/aidlc-evaluator    && go build ./cmd/aidlc-eval/...

# Test all tools
cd scripts/aidlc-traceability && go test ./...
cd scripts/aidlc-designreview && go test ./...
cd scripts/aidlc-evaluator    && go test ./...

# Release snapshot (all tools)
goreleaser build --snapshot --clean

## Code style
- gofmt/goimports for formatting
- Table-driven tests (t.Run subtests)
- Integration tests use build tag: //go:build integration
- No comments except for non-obvious WHY
