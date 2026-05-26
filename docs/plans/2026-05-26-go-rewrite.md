# Go Rewrite of aidlc-workflows Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Rewrite all three Python aidlc-workflows tools (aidlc-designreview, aidlc-traceability, aidlc-evaluator) into Go with full feature parity, equivalent test coverage, single static binary per tool, and GoReleaser distribution — mirroring the Python repo layout as closely as Go best practices allow.

**Architecture:** Three independent Go modules under `scripts/`, each with a Cobra CLI, AWS SDK Go v2 for Bedrock, and a custom thin Bedrock agent loop replacing `strands-agents`. Internal packages mirror the Python layer structure exactly. The evaluator's Docker sandbox uses the Docker SDK for Go rather than shell-outs. Config files, prompt templates, and pattern docs are embedded into the binary via `embed.FS`. Go templates replace Jinja2.

**Tech Stack:**
- Go 1.25
- `github.com/spf13/cobra` — CLI framework (replaces Click/argparse)
- `github.com/aws/aws-sdk-go-v2` — AWS SDK (replaces boto3)
- `github.com/aws/aws-sdk-go-v2/service/bedrockruntime` — Bedrock Converse API (replaces strands-agents)
- `github.com/docker/docker/client` — Docker SDK (evaluator, replaces shell-out)
- `gonum.org/v1/gonum/graph` — directed graph (traceability, replaces networkx)
- `gopkg.in/yaml.v3` — YAML config (replaces pyyaml)
- `html/template` + `text/template` — report templates (replaces Jinja2)
- `github.com/cenkalti/backoff/v4` — exponential backoff (replaces backoff library)
- `github.com/pterm/pterm` — Rich-style terminal output (spinners, colors, progress)
- GoReleaser — cross-platform binary distribution

**Python reference:** `aidlc-workflows-py/` git submodule. Read it for behaviour, not structure.

**Key structural mirror:**
```text
Python                                      Go
scripts/aidlc-designreview/src/design_reviewer/{foundation,validation,parsing,ai_review,reporting,orchestration,cli}/
→ scripts/aidlc-designreview/internal/{foundation,validation,parsing,aireview,reporting,orchestration}/
   + cmd/design-reviewer/main.go

scripts/aidlc-traceability/src/traceability/{models,discovery,parsers,graph,analysis,agent,generators,pipeline,cli}/
→ scripts/aidlc-traceability/internal/{models,discovery,parsers,graph,analysis,agent,generators,pipeline}/
   + cmd/traceability/main.go

scripts/aidlc-evaluator/packages/{shared,execution,cli-harness,ide-harness,qualitative,quantitative,contracttest,reporting,trend-reports}/
→ scripts/aidlc-evaluator/internal/{shared,execution,cliharness,ideharness,qualitative,quantitative,contracttest,reporting,trendreports}/
   + cmd/aidlc-eval/main.go
```

---

## Phase 0: Repository Scaffold

### Task 1: AGENTS.md and top-level docs

**Files:**
- Create: `AGENTS.md`
- Create: `docs/README.md` (index of plans)

**Step 1: Write AGENTS.md** — adapt from `aidlc-workflows-py/AGENTS.md`, replacing Python-specific sections with Go equivalents:

```markdown
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
```

**Step 2: Commit**

```bash
git add AGENTS.md docs/
git commit -m "Add AGENTS.md and docs scaffold"
```

---

### Task 2: GoReleaser config and CI/release workflows

**Files:**
- Create: `.goreleaser.yml`
- Create: `.github/workflows/ci.yml`
- Create: `.github/workflows/release.yml`

**Step 1: Write `.goreleaser.yml`**

```yaml
version: 2

before:
  hooks:
    - go mod tidy
    - go generate ./...

builds:
  - id: traceability
    dir: scripts/aidlc-traceability
    main: ./cmd/traceability
    binary: aidlc-traceability
    env: [CGO_ENABLED=0]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}

  - id: designreview
    dir: scripts/aidlc-designreview
    main: ./cmd/design-reviewer
    binary: aidlc-designreview
    env: [CGO_ENABLED=0]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}

  - id: evaluator
    dir: scripts/aidlc-evaluator
    main: ./cmd/aidlc-eval
    binary: aidlc-evaluator
    env: [CGO_ENABLED=0]
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}

archives:
  - id: all
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    name_template: "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"

checksum:
  name_template: checksums.txt

changelog:
  sort: asc
  filters:
    exclude: ['^docs:', '^test:', '^chore:']
```

**Step 2: Write `.github/workflows/ci.yml`**

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        tool: [aidlc-traceability, aidlc-designreview, aidlc-evaluator]
    steps:
      - uses: actions/checkout@v4
        with:
          submodules: true
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
          cache-dependency-path: scripts/${{ matrix.tool }}/go.sum
      - name: Test ${{ matrix.tool }}
        run: |
          cd scripts/${{ matrix.tool }}
          go vet ./...
          go test -race -count=1 ./...

  goreleaser-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - uses: goreleaser/goreleaser-action@v6
        with:
          args: check
```

**Step 3: Write `.github/workflows/release.yml`**

```yaml
name: Release
on:
  push:
    tags: ['v*']

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
          submodules: true
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - uses: goreleaser/goreleaser-action@v6
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Step 4: Commit**

```bash
git add .goreleaser.yml .github/
git commit -m "Add GoReleaser config and CI/release workflows"
```

---

## Phase 1: aidlc-traceability

**Reference:** `aidlc-workflows-py/scripts/aidlc-traceability/`
**Go module path:** `github.com/mach4-braai/aidlc-workflows/aidlc-traceability`

### Task 3: Module scaffold and models

**Files:**
- Create: `scripts/aidlc-traceability/go.mod`
- Create: `scripts/aidlc-traceability/go.sum`
- Create: `scripts/aidlc-traceability/internal/models/models.go`
- Create: `scripts/aidlc-traceability/internal/models/models_test.go`

**Step 1: Initialize module**

```bash
cd scripts/aidlc-traceability
go mod init github.com/mach4-braai/aidlc-workflows/aidlc-traceability
go get github.com/spf13/cobra@latest
go get github.com/aws/aws-sdk-go-v2/config@latest
go get github.com/aws/aws-sdk-go-v2/service/bedrockruntime@latest
go get github.com/aws/aws-sdk-go-v2/service/sts@latest
go get gonum.org/v1/gonum/graph@latest
go get gonum.org/v1/gonum/graph/simple@latest
go get gopkg.in/yaml.v3@latest
go get github.com/pterm/pterm@latest
go get github.com/cenkalti/backoff/v4@latest
```

**Step 2: Write failing test** — mirrors `tests/test_models.py`

```go
// internal/models/models_test.go
package models_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func TestArtifactIDIsNonEmpty(t *testing.T) {
    a := models.Artifact{ID: "FR-001", Title: "Login", Type: models.ArtifactTypeRequirement}
    if a.ID == "" {
        t.Fatal("artifact ID must not be empty")
    }
}

func TestCoverageMetricsZeroValue(t *testing.T) {
    m := models.CoverageMetrics{}
    if m.TotalRequirements != 0 {
        t.Fatal("zero value should be 0")
    }
}

func TestTraceabilityReportHoldsArtifacts(t *testing.T) {
    r := models.TraceabilityReport{
        Artifacts: []models.Artifact{
            {ID: "FR-001", Title: "Login", Type: models.ArtifactTypeRequirement},
        },
    }
    if len(r.Artifacts) != 1 {
        t.Fatalf("expected 1 artifact, got %d", len(r.Artifacts))
    }
}
```

**Step 3: Run test — expect FAIL**

```bash
go test ./internal/models/...
# Expected: FAIL — models package not defined
```

**Step 4: Write minimal implementation** — mirrors `src/traceability/models.py`

```go
// internal/models/models.go
package models

import "time"

type ArtifactType string

const (
    ArtifactTypeRequirement ArtifactType = "REQUIREMENT"
    ArtifactTypeStory       ArtifactType = "STORY"
    ArtifactTypeUnit        ArtifactType = "UNIT"
    ArtifactTypeComponent   ArtifactType = "COMPONENT"
    ArtifactTypeCodePlan    ArtifactType = "CODE_PLAN"
    ArtifactTypeCode        ArtifactType = "CODE"
    ArtifactTypeTest        ArtifactType = "TEST"
)

type Artifact struct {
    ID         string            `json:"id"`
    Title      string            `json:"title"`
    Type       ArtifactType      `json:"artifact_type"`
    Desc       string            `json:"description"`
    SourceFile string            `json:"source_file"`
    SourceLine int               `json:"source_line"`
    Metadata   map[string]any    `json:"metadata"`
}

type Relationship struct {
    SourceID         string `json:"source_id"`
    TargetID         string `json:"target_id"`
    RelationshipType string `json:"relationship_type"`
}

type CoverageGap struct {
    ArtifactID    string       `json:"artifact_id"`
    ArtifactTitle string       `json:"artifact_title"`
    ArtifactType  ArtifactType `json:"artifact_type"`
    GapType       string       `json:"gap_type"`
    Description   string       `json:"description"`
}

type CoverageMetrics struct {
    TotalRequirements       int `json:"total_requirements"`
    TotalStories            int `json:"total_stories"`
    TotalUnits              int `json:"total_units"`
    TotalCodeFiles          int `json:"total_code_files"`
    TotalTests              int `json:"total_tests"`
    RequirementsWithStories int `json:"requirements_with_stories"`
    StoriesWithUnits        int `json:"stories_with_units"`
    UnitsWithCode           int `json:"units_with_code"`
    CodeWithTests           int `json:"code_with_tests"`
}

type TraceabilityReport struct {
    ProjectName   string         `json:"project_name"`
    GeneratedAt   time.Time      `json:"generated_at"`
    Artifacts     []Artifact     `json:"artifacts"`
    Relationships []Relationship `json:"relationships"`
    Gaps          []CoverageGap  `json:"gaps"`
    Metrics       CoverageMetrics `json:"metrics"`
}
```

**Step 5: Run test — expect PASS**

```bash
go test ./internal/models/... -v
# Expected: PASS
```

**Step 6: Commit**

```bash
git add scripts/aidlc-traceability/
git commit -m "Add aidlc-traceability module scaffold and models"
```

---

### Task 4: Discovery package

**Files:**
- Create: `scripts/aidlc-traceability/internal/discovery/discovery.go`
- Create: `scripts/aidlc-traceability/internal/discovery/discovery_test.go`

**Reference:** `src/traceability/discovery.py`

**Step 1: Write failing tests** — mirrors `tests/test_discovery.py`

```go
// internal/discovery/discovery_test.go
package discovery_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/discovery"
)

func TestFindAidlcDocsReturnsNilWhenAbsent(t *testing.T) {
    tmp := t.TempDir()
    result := discovery.FindAidlcDocs(tmp)
    if result != "" {
        t.Fatalf("expected empty, got %q", result)
    }
}

func TestFindAidlcDocsFindsDirectory(t *testing.T) {
    tmp := t.TempDir()
    docsDir := filepath.Join(tmp, "aidlc-docs")
    os.MkdirAll(docsDir, 0755)
    result := discovery.FindAidlcDocs(tmp)
    if result != docsDir {
        t.Fatalf("expected %q, got %q", docsDir, result)
    }
}

func TestDiscoverSourceCodeFindsGoFiles(t *testing.T) {
    tmp := t.TempDir()
    os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main"), 0644)
    files := discovery.DiscoverSourceCode(tmp)
    if len(files) != 1 {
        t.Fatalf("expected 1 file, got %d", len(files))
    }
}

func TestDiscoverSourceCodeExcludesTestFiles(t *testing.T) {
    tmp := t.TempDir()
    os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main"), 0644)
    os.WriteFile(filepath.Join(tmp, "main_test.go"), []byte("package main"), 0644)
    files := discovery.DiscoverSourceCode(tmp)
    for _, f := range files {
        if filepath.Base(f) == "main_test.go" {
            t.Fatal("should not include test files")
        }
    }
}
```

**Step 2: Run — expect FAIL**

```bash
go test ./internal/discovery/... -v
```

**Step 3: Implement**

```go
// internal/discovery/discovery.go
package discovery

import (
    "os"
    "path/filepath"
    "strings"
)

// FindAidlcDocs walks up from projectRoot looking for aidlc-docs/.
// Returns the full path or empty string if not found.
func FindAidlcDocs(projectRoot string) string {
    candidate := filepath.Join(projectRoot, "aidlc-docs")
    if info, err := os.Stat(candidate); err == nil && info.IsDir() {
        return candidate
    }
    return ""
}

var sourceExtensions = map[string]bool{
    ".go": true, ".py": true, ".java": true, ".ts": true,
    ".cpp": true, ".rs": true, ".kt": true, ".swift": true,
}

// DiscoverSourceCode returns all non-test source files under root.
func DiscoverSourceCode(root string) []string {
    var files []string
    filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
        if err != nil || d.IsDir() {
            return nil
        }
        ext := strings.ToLower(filepath.Ext(path))
        if !sourceExtensions[ext] {
            return nil
        }
        base := filepath.Base(path)
        if strings.HasSuffix(base, "_test.go") || strings.HasSuffix(base, "_test.py") {
            return nil
        }
        files = append(files, path)
        return nil
    })
    return files
}
```

**Step 4: Run — expect PASS**

```bash
go test ./internal/discovery/... -v
```

**Step 5: Commit**

```bash
git add scripts/aidlc-traceability/internal/discovery/
git commit -m "Add traceability discovery package"
```

---

### Task 5: Parsers package

**Files:**
- Create: `scripts/aidlc-traceability/internal/parsers/requirements.go`
- Create: `scripts/aidlc-traceability/internal/parsers/stories.go`
- Create: `scripts/aidlc-traceability/internal/parsers/units.go`
- Create: `scripts/aidlc-traceability/internal/parsers/code_plans.go`
- Create: `scripts/aidlc-traceability/internal/parsers/components.go`
- Create: `scripts/aidlc-traceability/internal/parsers/code.go`
- Create: `scripts/aidlc-traceability/internal/parsers/linker.go`
- Create: `scripts/aidlc-traceability/internal/parsers/parsers_test.go`

**Reference:** `src/traceability/parsers/*.py`

**Step 1: Write failing tests** — mirrors `tests/test_parsers.py`

```go
// internal/parsers/parsers_test.go
package parsers_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/parsers"
)

func TestParseRequirementsExtractsIDs(t *testing.T) {
    content := "# Requirements\n\n## FR-001: User Login\n\nUsers must be able to log in.\n\n## FR-002: Dashboard\n\nShow dashboard.\n"
    tmp := t.TempDir()
    f := filepath.Join(tmp, "requirements.md")
    os.WriteFile(f, []byte(content), 0644)
    artifacts := parsers.ParseRequirements(f)
    if len(artifacts) != 2 {
        t.Fatalf("expected 2 requirements, got %d", len(artifacts))
    }
    if artifacts[0].ID != "FR-001" {
        t.Fatalf("expected FR-001, got %s", artifacts[0].ID)
    }
}

func TestParseStoriesExtractsIDs(t *testing.T) {
    content := "# Stories\n\n## US-1.1: Login Form\n\nAs a user I want to log in.\n"
    tmp := t.TempDir()
    f := filepath.Join(tmp, "stories.md")
    os.WriteFile(f, []byte(content), 0644)
    artifacts := parsers.ParseStories(f)
    if len(artifacts) == 0 {
        t.Fatal("expected at least one story")
    }
    if artifacts[0].Type != models.ArtifactTypeStory {
        t.Fatalf("expected STORY type, got %s", artifacts[0].Type)
    }
}

func TestParseCodeFileReturnsArtifact(t *testing.T) {
    tmp := t.TempDir()
    f := filepath.Join(tmp, "main.go")
    os.WriteFile(f, []byte("package main\n\nfunc main() {}\n"), 0644)
    a := parsers.ParseCodeFile(f, tmp)
    if a == nil {
        t.Fatal("expected artifact, got nil")
    }
    if a.Type != models.ArtifactTypeCode {
        t.Fatalf("expected CODE, got %s", a.Type)
    }
}

func TestInferLinksFindsExplicitReferences(t *testing.T) {
    reqs := []models.Artifact{{ID: "FR-001", Type: models.ArtifactTypeRequirement}}
    stories := []models.Artifact{{ID: "US-1.1", Title: "Implements FR-001", Type: models.ArtifactTypeStory}}
    links := parsers.InferRequirementStoryLinks(reqs, stories)
    if len(links) == 0 {
        t.Fatal("expected at least one inferred link")
    }
}
```

**Step 2: Run — expect FAIL**

```bash
go test ./internal/parsers/... -v
```

**Step 3: Implement** — each parser uses regex to extract headings with ID patterns (mirrors Python regex logic from `requirements.py`, `stories.py`, `units.py` etc.)

Key patterns per parser:
- Requirements: `## (FR|NFR|AR)-\d+:` heading → Artifact
- Stories: `## (US-\d+\.\d+|US-\d+):` heading → Artifact
- Units: `## (UNIT-\d+|[A-Za-z-]+---[A-Za-z-]+):` heading → Artifact
- Components: `## [A-Z][a-zA-Z]+` in `application-components.md` → Artifact
- Code: filepath → Artifact with SourceFile, count LoC
- Linker: scan story Title/Description for requirement ID mentions via regex

Create one `.go` file per parser, plus `linker.go`. Each function signature:
```go
func ParseRequirements(filePath string) []models.Artifact
func ParseStories(filePath string) []models.Artifact
func ParseUnits(filePath string) []models.Artifact
func ParseCodePlans(filePath string) []models.Artifact
func ParseComponents(filePath string) []models.Artifact
func ParseCodeFile(filePath, projectRoot string) *models.Artifact
func InferRequirementStoryLinks(reqs, stories []models.Artifact) []models.Relationship
```

**Step 4: Run — expect PASS**

```bash
go test ./internal/parsers/... -v
```

**Step 5: Commit**

```bash
git add scripts/aidlc-traceability/internal/parsers/
git commit -m "Add traceability parsers (requirements, stories, units, code, linker)"
```

---

### Task 6: Graph and analysis packages

**Files:**
- Create: `scripts/aidlc-traceability/internal/graph/graph.go`
- Create: `scripts/aidlc-traceability/internal/graph/graph_test.go`
- Create: `scripts/aidlc-traceability/internal/analysis/analysis.go`
- Create: `scripts/aidlc-traceability/internal/analysis/analysis_test.go`

**Reference:** `src/traceability/graph.py` and `analysis.py`

**Step 1: Write failing tests** — mirrors `tests/test_graph.py`

```go
// internal/graph/graph_test.go
package graph_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/graph"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func TestBuildGraphHasNodes(t *testing.T) {
    artifacts := []models.Artifact{
        {ID: "FR-001", Type: models.ArtifactTypeRequirement},
        {ID: "US-1.1", Type: models.ArtifactTypeStory},
    }
    rels := []models.Relationship{
        {SourceID: "FR-001", TargetID: "US-1.1", RelationshipType: "traces_to"},
    }
    g := graph.Build(artifacts, rels)
    if graph.NodeCount(g) != 2 {
        t.Fatalf("expected 2 nodes, got %d", graph.NodeCount(g))
    }
    if graph.EdgeCount(g) != 1 {
        t.Fatalf("expected 1 edge, got %d", graph.EdgeCount(g))
    }
}

func TestBuildGraphHandlesOrphanedArtifacts(t *testing.T) {
    artifacts := []models.Artifact{
        {ID: "FR-001", Type: models.ArtifactTypeRequirement},
    }
    g := graph.Build(artifacts, nil)
    if graph.NodeCount(g) != 1 {
        t.Fatalf("expected 1 node, got %d", graph.NodeCount(g))
    }
}
```

```go
// internal/analysis/analysis_test.go
package analysis_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/analysis"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/graph"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func TestDetectGapsFindsOrphanedRequirement(t *testing.T) {
    artifacts := []models.Artifact{
        {ID: "FR-001", Type: models.ArtifactTypeRequirement},
    }
    g := graph.Build(artifacts, nil)
    gaps := analysis.DetectGaps(artifacts, g)
    if len(gaps) == 0 {
        t.Fatal("expected gap for orphaned requirement")
    }
}

func TestCalculateMetricsCountsCorrectly(t *testing.T) {
    artifacts := []models.Artifact{
        {ID: "FR-001", Type: models.ArtifactTypeRequirement},
        {ID: "US-1.1", Type: models.ArtifactTypeStory},
    }
    rels := []models.Relationship{
        {SourceID: "FR-001", TargetID: "US-1.1"},
    }
    m := analysis.CalculateMetrics(artifacts, rels)
    if m.TotalRequirements != 1 {
        t.Fatalf("expected 1 requirement, got %d", m.TotalRequirements)
    }
    if m.RequirementsWithStories != 1 {
        t.Fatalf("expected 1 requirement with story, got %d", m.RequirementsWithStories)
    }
}
```

**Step 2: Run — expect FAIL**

```bash
go test ./internal/graph/... ./internal/analysis/... -v
```

**Step 3: Implement**

`graph.go` — wraps `gonum.org/v1/gonum/graph/simple.DirectedGraph`. Assigns stable int64 IDs to artifact string IDs via a lookup map. Exposes `Build()`, `NodeCount()`, `EdgeCount()`, `HasSuccessor(g, id)`, `HasPredecessor(g, id)`.

`analysis.go` — `DetectGaps()` iterates artifacts, checks for outgoing/incoming edges per type (requirements should have stories, stories should have units, etc.). `CalculateMetrics()` counts artifacts by type, then counts those with at least one outgoing edge to the expected next type.

**Step 4: Run — expect PASS**

```bash
go test ./internal/graph/... ./internal/analysis/... -v
```

**Step 5: Commit**

```bash
git add scripts/aidlc-traceability/internal/graph/ scripts/aidlc-traceability/internal/analysis/
git commit -m "Add traceability graph and analysis packages"
```

---

### Task 7: Bedrock agent package (thin agent loop)

**Files:**
- Create: `scripts/aidlc-traceability/internal/agent/agent.go`
- Create: `scripts/aidlc-traceability/internal/agent/agent_test.go`

**Reference:** `src/traceability/agent.py`

**Context:** Python uses `strands-agents` which wraps Bedrock Converse API. Go version calls `bedrockruntime.Converse` directly with a simple agentic loop: send message → if response has tool_use → execute tool → append tool_result → repeat until stop_reason is "end_turn".

**Step 1: Write failing tests**

```go
// internal/agent/agent_test.go
package agent_test

import (
    "context"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/agent"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func TestParseAgentJSONExtractsRelationships(t *testing.T) {
    validIDs := map[string]bool{"FR-001": true, "US-1.1": true}
    responseText := `{"relationships": [{"source_id": "FR-001", "target_id": "US-1.1", "relationship_type": "traces_to"}], "insights": "good match"}`
    rels, insights := agent.ParseAgentJSON(responseText, validIDs)
    if len(rels) != 1 {
        t.Fatalf("expected 1 relationship, got %d", len(rels))
    }
    if len(insights) == 0 {
        t.Fatal("expected insights")
    }
}

func TestParseAgentJSONSkipsInvalidIDs(t *testing.T) {
    validIDs := map[string]bool{"FR-001": true}
    responseText := `{"relationships": [{"source_id": "INVALID", "target_id": "US-1.1", "relationship_type": "traces_to"}]}`
    rels, _ := agent.ParseAgentJSON(responseText, validIDs)
    if len(rels) != 0 {
        t.Fatal("expected 0 relationships for invalid IDs")
    }
}

func TestParseAgentJSONHandlesMalformedJSON(t *testing.T) {
    validIDs := map[string]bool{"FR-001": true}
    rels, _ := agent.ParseAgentJSON("not json", validIDs)
    if len(rels) != 0 {
        t.Fatal("expected 0 relationships for malformed JSON")
    }
}

// Integration test — skipped unless AWS credentials present
func TestRunReqStoryAnalysis_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    // Requires: AWS_PROFILE or default credentials + Bedrock access
    ctx := context.Background()
    reqs := []models.Artifact{{ID: "FR-001", Title: "User Login", Type: models.ArtifactTypeRequirement}}
    stories := []models.Artifact{{ID: "US-1.1", Title: "Implements FR-001 login form", Type: models.ArtifactTypeStory}}
    rels, _ := agent.RunReqStoryAnalysis(ctx, reqs, stories, "", "us-east-1")
    _ = rels // just check it doesn't panic
}
```

**Step 2: Run — expect FAIL**

```bash
go test ./internal/agent/... -v -short
```

**Step 3: Implement `agent.go`**

Key exported symbols:
- `ParseAgentJSON(responseText string, validIDs map[string]bool) ([]models.Relationship, []string)` — JSON extraction + validation (pure, no Bedrock)
- `RunReqStoryAnalysis(ctx context.Context, reqs, stories []models.Artifact, profile, region string) ([]models.Relationship, []string)` — calls Bedrock Converse
- `RunStoryUnitAnalysis(...)`, `RunUnitComponentAnalysis(...)`, `RunComponentCodeAnalysis(...)` — same pattern

Internal `converseLoop` function:
```go
func converseLoop(ctx context.Context, client *bedrockruntime.Client, modelID string, systemPrompt string, userMsg string, tools []types.Tool) (string, error) {
    messages := []types.Message{{
        Role: types.ConversationRoleUser,
        Content: []types.ContentBlock{&types.ContentBlockMemberText{Value: types.ContentBlock(&types.ContentBlockMemberText{Value: userMsg})}},
    }}
    for {
        resp, err := client.Converse(ctx, &bedrockruntime.ConverseInput{
            ModelId: aws.String(modelID),
            System:  []types.SystemContentBlock{&types.SystemContentBlockMemberText{Value: systemPrompt}},
            Messages: messages,
        })
        if err != nil { return "", err }
        if resp.StopReason == types.StopReasonEndTurn {
            return extractText(resp.Output), nil
        }
        // handle tool_use blocks if present
        messages = append(messages, toolResultMessages(resp)...)
    }
}
```

**Step 4: Run — expect PASS (short mode)**

```bash
go test ./internal/agent/... -v -short
```

**Step 5: Commit**

```bash
git add scripts/aidlc-traceability/internal/agent/
git commit -m "Add thin Bedrock agent loop replacing strands-agents"
```

---

### Task 8: Generators (markdown + HTML)

**Files:**
- Create: `scripts/aidlc-traceability/internal/generators/markdown.go`
- Create: `scripts/aidlc-traceability/internal/generators/html.go`
- Create: `scripts/aidlc-traceability/internal/generators/generators_test.go`

**Reference:** `src/traceability/generators/{markdown,html}.py`

**Step 1: Write failing tests** — mirrors `tests/test_generators.py`

```go
// internal/generators/generators_test.go
package generators_test

import (
    "strings"
    "testing"
    "time"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/generators"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

func sampleReport() models.TraceabilityReport {
    return models.TraceabilityReport{
        ProjectName: "test-project",
        GeneratedAt: time.Now(),
        Artifacts: []models.Artifact{
            {ID: "FR-001", Title: "Login", Type: models.ArtifactTypeRequirement},
        },
        Metrics: models.CoverageMetrics{TotalRequirements: 1},
    }
}

func TestGenerateMarkdownContainsProjectName(t *testing.T) {
    out := generators.GenerateMarkdown(sampleReport())
    if !strings.Contains(out, "test-project") {
        t.Fatal("markdown output must contain project name")
    }
}

func TestGenerateMarkdownContainsArtifactID(t *testing.T) {
    out := generators.GenerateMarkdown(sampleReport())
    if !strings.Contains(out, "FR-001") {
        t.Fatal("markdown output must contain artifact ID")
    }
}

func TestGenerateHTMLIsValidHTML(t *testing.T) {
    out := generators.GenerateHTML(sampleReport())
    if !strings.Contains(out, "<html") {
        t.Fatal("HTML output must contain <html> tag")
    }
    if !strings.Contains(out, "</html>") {
        t.Fatal("HTML output must close <html> tag")
    }
}
```

**Step 2: Run — expect FAIL**

```bash
go test ./internal/generators/... -v
```

**Step 3: Implement** — use `text/template` for markdown, `html/template` for HTML. Templates are defined as string constants within the Go files (not external files — embed is for config/prompts). Render to string, return.

**Step 4: Run — expect PASS**

```bash
go test ./internal/generators/... -v
```

**Step 5: Commit**

```bash
git add scripts/aidlc-traceability/internal/generators/
git commit -m "Add traceability markdown and HTML generators"
```

---

### Task 9: Pipeline and CLI

**Files:**
- Create: `scripts/aidlc-traceability/internal/pipeline/pipeline.go`
- Create: `scripts/aidlc-traceability/internal/pipeline/pipeline_test.go`
- Create: `scripts/aidlc-traceability/cmd/traceability/main.go`
- Create: `scripts/aidlc-traceability/cmd/traceability/root.go`

**Reference:** `src/traceability/pipeline.py` and `cli.py`

**Step 1: Write failing pipeline tests** — mirrors `tests/test_cli_pipeline.py`

```go
// internal/pipeline/pipeline_test.go
package pipeline_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/pipeline"
)

func TestRunPipelineOnEmptyProject(t *testing.T) {
    tmp := t.TempDir()
    report, err := pipeline.Run(pipeline.Config{
        ProjectRoot: tmp,
        UseAI:       false,
        OutputDir:   tmp,
        Format:      "markdown",
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if report.ProjectName == "" {
        t.Fatal("report must have a project name")
    }
}

func TestRunPipelineWritesOutputFile(t *testing.T) {
    tmp := t.TempDir()
    docsDir := filepath.Join(tmp, "aidlc-docs")
    os.MkdirAll(docsDir, 0755)
    os.WriteFile(filepath.Join(docsDir, "requirements.md"), []byte("# Requirements\n\n## FR-001: Login\n\nUsers must log in.\n"), 0644)

    _, err := pipeline.Run(pipeline.Config{
        ProjectRoot: tmp,
        UseAI:       false,
        OutputDir:   tmp,
        Format:      "markdown",
    })
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    entries, _ := os.ReadDir(tmp)
    hasMarkdown := false
    for _, e := range entries {
        if filepath.Ext(e.Name()) == ".md" {
            hasMarkdown = true
        }
    }
    if !hasMarkdown {
        t.Fatal("expected markdown output file")
    }
}
```

**Step 2: Run — expect FAIL**

```bash
go test ./internal/pipeline/... -v
```

**Step 3: Implement pipeline.go** — orchestrates: `discovery → parsers → linker → graph → analysis → [agent if UseAI] → generators → write files`

```go
// Config holds all pipeline inputs
type Config struct {
    ProjectRoot string
    OutputDir   string
    Format      string // "markdown", "html", "both"
    UseAI       bool
    AWSProfile  string
    AWSRegion   string
    Verbose     bool
}

func Run(cfg Config) (models.TraceabilityReport, error) { ... }
```

**Step 4: Implement cmd/traceability/main.go + root.go** — Cobra `generate` subcommand wiring all Config fields to CLI flags. Mirror Python: `traceability generate --input <path> [--output <path>] [--format markdown|html|both] [--no-ai] [--profile <profile>] [--region <region>]`

**Step 5: Run tests — expect PASS**

```bash
go test ./... -v
```

**Step 6: Build and smoke test**

```bash
go build ./cmd/traceability/...
./traceability generate --help
```

Expected: help text with all flags listed.

**Step 7: Commit**

```bash
git add scripts/aidlc-traceability/internal/pipeline/ scripts/aidlc-traceability/cmd/
git commit -m "Add traceability pipeline orchestration and Cobra CLI"
```

---

## Phase 2: aidlc-designreview

**Reference:** `aidlc-workflows-py/scripts/aidlc-designreview/`
**Go module path:** `github.com/mach4-braai/aidlc-workflows/aidlc-designreview`

### Task 10: Module scaffold, config, and foundation

**Files:**
- Create: `scripts/aidlc-designreview/go.mod`
- Create: `scripts/aidlc-designreview/internal/foundation/config.go`
- Create: `scripts/aidlc-designreview/internal/foundation/config_test.go`
- Create: `scripts/aidlc-designreview/config/default-config.yaml` (copy from Python)
- Create: `scripts/aidlc-designreview/config/prompts/` (copy .md files from Python)
- Create: `scripts/aidlc-designreview/config/patterns/` (copy .md files from Python)

**Reference:** `src/design_reviewer/foundation/{config_manager.py,config_models.py,exceptions.py,pattern_library.py,prompt_manager.py}`

**Step 1: Initialize module**

```bash
cd scripts/aidlc-designreview
go mod init github.com/mach4-braai/aidlc-workflows/aidlc-designreview
go get github.com/spf13/cobra@latest
go get github.com/aws/aws-sdk-go-v2/config@latest
go get github.com/aws/aws-sdk-go-v2/service/bedrockruntime@latest
go get gopkg.in/yaml.v3@latest
go get github.com/pterm/pterm@latest
go get github.com/cenkalti/backoff/v4@latest
```

**Step 2: Write failing tests** — mirrors `tests/unit1_foundation/`

```go
// internal/foundation/config_test.go
package foundation_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/foundation"
)

func TestLoadConfigFromYAML(t *testing.T) {
    yaml := `
models:
  critique: us.anthropic.claude-sonnet-4-20250514-v1:0
  alternatives: us.anthropic.claude-sonnet-4-20250514-v1:0
  gap: us.anthropic.claude-sonnet-4-20250514-v1:0
aws:
  region: us-east-1
`
    tmp := t.TempDir()
    f := filepath.Join(tmp, "config.yaml")
    os.WriteFile(f, []byte(yaml), 0644)
    cfg, err := foundation.LoadConfig(f)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cfg.AWS.Region != "us-east-1" {
        t.Fatalf("expected us-east-1, got %s", cfg.AWS.Region)
    }
}

func TestLoadConfigUsesDefaults(t *testing.T) {
    cfg, err := foundation.LoadDefaultConfig()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cfg.AWS.Region == "" {
        t.Fatal("default region must be set")
    }
}

func TestPatternLibraryLoadsPatterns(t *testing.T) {
    lib, err := foundation.LoadPatternLibrary()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(lib.Patterns) == 0 {
        t.Fatal("expected at least one design pattern")
    }
}

func TestPromptManagerBuildsPrompt(t *testing.T) {
    pm, err := foundation.LoadPromptManager()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    prompt, err := pm.BuildAgentPrompt("critique", map[string]string{"design_content": "test content"})
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if prompt == "" {
        t.Fatal("prompt must not be empty")
    }
}
```

**Step 3: Run — expect FAIL**

```bash
go test ./internal/foundation/... -v
```

**Step 4: Implement foundation package**

Key types:
```go
type Config struct {
    Models ModelConfig `yaml:"models"`
    AWS    AWSConfig   `yaml:"aws"`
    Review ReviewSettings `yaml:"review"`
}
type AWSConfig struct {
    Region      string `yaml:"region"`
    ProfileName string `yaml:"profile_name"`
    GuardrailID string `yaml:"guardrail_id"`
    GuardrailVersion string `yaml:"guardrail_version"`
}
type ModelConfig struct {
    Critique     string `yaml:"critique"`
    Alternatives string `yaml:"alternatives"`
    Gap          string `yaml:"gap"`
}
```

Config, patterns, and prompts are embedded via `//go:embed`:
```go
//go:embed ../../config
var configFS embed.FS
```

`LoadDefaultConfig()` reads from embedded FS. `LoadConfig(path)` reads from disk and merges with defaults.

`PatternLibrary` reads all `.md` files from `config/patterns/` embedded FS. `PromptManager` reads `.md` files from `config/prompts/` and uses `text/template` for variable substitution.

**Step 5: Run — expect PASS**

```bash
go test ./internal/foundation/... -v
```

**Step 6: Commit**

```bash
git add scripts/aidlc-designreview/
git commit -m "Add designreview module scaffold, config, and foundation"
```

---

### Task 11: Validation and parsing layers

**Files:**
- Create: `scripts/aidlc-designreview/internal/validation/` (models.go, scanner.go, classifier.go, discoverer.go, loader.go, validator.go)
- Create: `scripts/aidlc-designreview/internal/validation/validation_test.go`
- Create: `scripts/aidlc-designreview/internal/parsing/` (models.go, app_design.go, func_design.go, tech_env.go, base.go)
- Create: `scripts/aidlc-designreview/internal/parsing/parsing_test.go`

**Reference:** `src/design_reviewer/validation/` and `src/design_reviewer/parsing/`

**Step 1: Write failing tests** — mirrors `tests/unit2_validation/` and `tests/unit3_parsing/`

```go
// internal/validation/validation_test.go
package validation_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/validation"
)

func TestClassifyApplicationDesignByContent(t *testing.T) {
    content := "# Application Design\n\n## Architecture\n\nMicroservices with REST APIs."
    artType := validation.ClassifyByContent(content)
    if artType != validation.ArtifactTypeApplicationDesign {
        t.Fatalf("expected APPLICATION_DESIGN, got %s", artType)
    }
}

func TestScanDirectoryFindsAidlcDocs(t *testing.T) {
    tmp := t.TempDir()
    docsDir := filepath.Join(tmp, "aidlc-docs", "construction")
    os.MkdirAll(docsDir, 0755)
    os.WriteFile(filepath.Join(docsDir, "application-design.md"), []byte("# Application Design"), 0644)
    result := validation.ScanDirectory(tmp)
    if !result.HasAidlcDocs {
        t.Fatal("should detect aidlc-docs directory")
    }
}

func TestValidateStructurePassesGoodProject(t *testing.T) {
    tmp := t.TempDir()
    docsDir := filepath.Join(tmp, "aidlc-docs", "construction")
    os.MkdirAll(docsDir, 0755)
    os.WriteFile(filepath.Join(docsDir, "application-design.md"), []byte("# Application Design\n"), 0644)
    result := validation.ValidateStructure(tmp)
    if !result.IsValid {
        t.Fatalf("expected valid, got errors: %v", result.Errors)
    }
}
```

```go
// internal/parsing/parsing_test.go
package parsing_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/parsing"
)

func TestParseApplicationDesignExtractsContent(t *testing.T) {
    content := "# Application Design\n\n## Architecture\n\nMicroservices.\n"
    model := parsing.ParseApplicationDesign([]string{content})
    if model.RawContent == "" {
        t.Fatal("raw content must not be empty")
    }
    if model.SourceCount != 1 {
        t.Fatalf("expected 1 source, got %d", model.SourceCount)
    }
}

func TestDesignDataAggregatesAllParsed(t *testing.T) {
    data := parsing.DesignData{
        AppDesign:  &parsing.ApplicationDesignModel{RawContent: "app content"},
        FuncDesign: &parsing.FunctionalDesignModel{RawContent: "func content"},
    }
    if data.AppDesign == nil {
        t.Fatal("AppDesign must not be nil")
    }
}
```

**Step 2: Implement both packages** following the Python layer structure exactly.

Key validation types:
```go
type ArtifactType string
const (
    ArtifactTypeApplicationDesign  ArtifactType = "APPLICATION_DESIGN"
    ArtifactTypeFunctionalDesign   ArtifactType = "FUNCTIONAL_DESIGN"
    ArtifactTypeTechnicalEnv       ArtifactType = "TECHNICAL_ENVIRONMENT"
)
type ScanResult struct { HasAidlcDocs bool; Files []string }
type ValidationResult struct { IsValid bool; Errors []string; Artifacts []ArtifactInfo }
```

Key parsing types:
```go
type ApplicationDesignModel struct { RawContent string; FilePaths []string; SourceCount int }
type FunctionalDesignModel  struct { RawContent string; FilePaths []string; UnitNames []string; SourceCount int }
type TechnicalEnvironmentModel struct { RawContent string; FilePath string }
type DesignData struct {
    AppDesign  *ApplicationDesignModel
    FuncDesign *FunctionalDesignModel
    TechEnv    *TechnicalEnvironmentModel
}
```

**Step 3: Run tests — expect PASS**

```bash
go test ./internal/validation/... ./internal/parsing/... -v
```

**Step 4: Commit**

```bash
git add scripts/aidlc-designreview/internal/validation/ scripts/aidlc-designreview/internal/parsing/
git commit -m "Add designreview validation and parsing layers"
```

---

### Task 12: AI review agents (critique, alternatives, gap)

**Files:**
- Create: `scripts/aidlc-designreview/internal/aireview/models.go`
- Create: `scripts/aidlc-designreview/internal/aireview/base.go`
- Create: `scripts/aidlc-designreview/internal/aireview/critique.go`
- Create: `scripts/aidlc-designreview/internal/aireview/alternatives.go`
- Create: `scripts/aidlc-designreview/internal/aireview/gap.go`
- Create: `scripts/aidlc-designreview/internal/aireview/orchestrator.go`
- Create: `scripts/aidlc-designreview/internal/aireview/response_parser.go`
- Create: `scripts/aidlc-designreview/internal/aireview/retry.go`
- Create: `scripts/aidlc-designreview/internal/aireview/aireview_test.go`

**Reference:** `src/design_reviewer/ai_review/`

**Step 1: Write failing tests** — mirrors `tests/unit4_ai_review/`

```go
// internal/aireview/aireview_test.go
package aireview_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/aireview"
)

func TestSeverityHasFourLevels(t *testing.T) {
    levels := []aireview.Severity{
        aireview.SeverityCritical, aireview.SeverityHigh,
        aireview.SeverityMedium,  aireview.SeverityLow,
    }
    if len(levels) != 4 {
        t.Fatal("expected 4 severity levels")
    }
}

func TestExtractJSONFromMarkdownCodeBlock(t *testing.T) {
    input := "```json\n{\"key\": \"value\"}\n```"
    extracted := aireview.ExtractJSONFromMarkdown(input)
    if extracted != `{"key": "value"}` {
        t.Fatalf("expected extracted JSON, got %q", extracted)
    }
}

func TestExtractJSONFromPlainText(t *testing.T) {
    input := `{"key": "value"}`
    extracted := aireview.ExtractJSONFromMarkdown(input)
    if extracted != input {
        t.Fatalf("expected same text, got %q", extracted)
    }
}

func TestValidateResponseSchemaRequiresKeys(t *testing.T) {
    response := `{"findings": [], "summary": "ok"}`
    required := map[string]bool{"findings": true}
    if !aireview.ValidateResponseSchema(response, required) {
        t.Fatal("response with required keys should pass validation")
    }
}

func TestValidateResponseSchemaMissingKeys(t *testing.T) {
    response := `{"wrong_key": []}`
    required := map[string]bool{"findings": true}
    if aireview.ValidateResponseSchema(response, required) {
        t.Fatal("response missing required keys should fail validation")
    }
}

func TestIsRetryableDetectsThrottling(t *testing.T) {
    err := &aireview.BedrockAPIError{Message: "ThrottlingException: rate exceeded"}
    if !aireview.IsRetryable(err) {
        t.Fatal("ThrottlingException should be retryable")
    }
}

func TestIsRetryableFalseForBadInput(t *testing.T) {
    err := &aireview.BedrockAPIError{Message: "ValidationException: invalid input"}
    if aireview.IsRetryable(err) {
        t.Fatal("ValidationException should not be retryable")
    }
}
```

**Step 2: Run — expect FAIL**

```bash
go test ./internal/aireview/... -v
```

**Step 3: Implement**

`models.go` defines:
- `Severity` string enum: `CRITICAL | HIGH | MEDIUM | LOW`
- `CritiqueFinding` struct (mirrors Python Pydantic model exactly)
- `AlternativeSuggestion` struct
- `GapFinding` struct
- `CritiqueResult`, `AlternativesResult`, `GapAnalysisResult` aggregation structs
- `ReviewResult` top-level container
- `TokenUsage` struct
- `BedrockAPIError` implementing `error`

`base.go` defines `BaseAgent` struct with:
- `invokeModel(ctx, prompt) (string, TokenUsage, error)` using `backoff.RetryNotify`
- `ExtractJSONFromMarkdown(text) string` (exported for testing)
- `ValidateResponseSchema(text string, required map[string]bool) bool` (exported for testing)

`retry.go` — `IsRetryable(err error) bool` checks error message for throttling/timeout keywords.

`critique.go`, `alternatives.go`, `gap.go` — each embeds `BaseAgent`, implements `Execute(ctx, data DesignData) (Result, error)`.

`orchestrator.go` — `AIOrchestrator.Run(ctx, data DesignData) (ReviewResult, error)` calls agents in sequence, collects token usage.

`response_parser.go` — parses Bedrock JSON responses into typed structs.

**Step 4: Run — expect PASS**

```bash
go test ./internal/aireview/... -v
```

**Step 5: Commit**

```bash
git add scripts/aidlc-designreview/internal/aireview/
git commit -m "Add designreview AI review agents and orchestrator"
```

---

### Task 13: Reporting layer (templates + formatters)

**Files:**
- Create: `scripts/aidlc-designreview/internal/reporting/models.go`
- Create: `scripts/aidlc-designreview/internal/reporting/report_builder.go`
- Create: `scripts/aidlc-designreview/internal/reporting/markdown_formatter.go`
- Create: `scripts/aidlc-designreview/internal/reporting/html_formatter.go`
- Create: `scripts/aidlc-designreview/internal/reporting/templates/html_report.html`
- Create: `scripts/aidlc-designreview/internal/reporting/templates/markdown_report.tmpl`
- Create: `scripts/aidlc-designreview/internal/reporting/reporting_test.go`

**Reference:** `src/design_reviewer/reporting/` + `templates/`

**Step 1: Write failing tests** — mirrors `tests/unit5_reporting/`

```go
// internal/reporting/reporting_test.go
package reporting_test

import (
    "strings"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/aireview"
    "github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/reporting"
)

func sampleReviewResult() aireview.ReviewResult {
    return aireview.ReviewResult{
        Critique: aireview.CritiqueResult{
            Findings: []aireview.CritiqueFinding{
                {ID: "C-001", Title: "Missing auth", Severity: aireview.SeverityHigh, Description: "No auth layer"},
            },
        },
    }
}

func TestBuildReportHasSeverityCount(t *testing.T) {
    result := sampleReviewResult()
    data := reporting.BuildReport(result)
    if data.Summary.HighCount != 1 {
        t.Fatalf("expected 1 HIGH finding, got %d", data.Summary.HighCount)
    }
}

func TestMarkdownFormatterOutputContainsFindings(t *testing.T) {
    data := reporting.BuildReport(sampleReviewResult())
    out, err := reporting.RenderMarkdown(data)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out, "Missing auth") {
        t.Fatal("markdown must contain finding title")
    }
}

func TestHTMLFormatterOutputIsHTML(t *testing.T) {
    data := reporting.BuildReport(sampleReviewResult())
    out, err := reporting.RenderHTML(data)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(out, "<html") {
        t.Fatal("HTML output must start with html tag")
    }
}
```

**Step 2: Run — expect FAIL**

```bash
go test ./internal/reporting/... -v
```

**Step 3: Implement**

Templates are embedded via `//go:embed templates/*`. Convert Jinja2 `.jinja2` templates to Go `html/template` (HTML) and `text/template` (Markdown). All Jinja2 filters become Go template functions registered with `FuncMap`.

`ReportData` struct mirrors Python's `ReportData`. `BuildReport()` computes `ReviewSummary` (severity counts, agent statuses). `RenderMarkdown()` and `RenderHTML()` execute the embedded templates.

**Step 4: Run — expect PASS**

```bash
go test ./internal/reporting/... -v
```

**Step 5: Commit**

```bash
git add scripts/aidlc-designreview/internal/reporting/
git commit -m "Add designreview reporting layer with Go templates"
```

---

### Task 14: Orchestration and CLI

**Files:**
- Create: `scripts/aidlc-designreview/internal/orchestration/orchestrator.go`
- Create: `scripts/aidlc-designreview/internal/orchestration/orchestrator_test.go`
- Create: `scripts/aidlc-designreview/cmd/design-reviewer/main.go`
- Create: `scripts/aidlc-designreview/cmd/design-reviewer/root.go`

**Reference:** `src/design_reviewer/orchestration/orchestrator.py` and `cli/`

**Step 1: Write failing tests** — mirrors `tests/unit5_orchestration/` and `tests/unit5_cli/`

```go
// internal/orchestration/orchestrator_test.go
package orchestration_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/orchestration"
)

func TestReviewOrchestratorFailsOnMissingDir(t *testing.T) {
    orch := orchestration.NewReviewOrchestrator(orchestration.Config{})
    _, err := orch.Run("/nonexistent/path")
    if err == nil {
        t.Fatal("expected error for missing directory")
    }
}

func TestReviewOrchestratorRunsOnMinimalProject(t *testing.T) {
    tmp := t.TempDir()
    docsDir := filepath.Join(tmp, "aidlc-docs", "construction")
    os.MkdirAll(docsDir, 0755)
    os.WriteFile(filepath.Join(docsDir, "application-design.md"),
        []byte("# Application Design\n\n## Architecture\n\nMonolith.\n"), 0644)

    orch := orchestration.NewReviewOrchestrator(orchestration.Config{
        MockAI: true, // skips actual Bedrock calls
    })
    result, err := orch.Run(tmp)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result == nil {
        t.Fatal("expected non-nil result")
    }
}
```

**Step 2: Implement orchestrator** — 6-stage pipeline matching Python exactly. `Config.MockAI bool` controls whether Bedrock is called (allows unit testing without credentials). Each stage updates a `pterm` spinner.

**Step 3: Implement CLI** — Cobra root command:
```bash
design-reviewer --aidlc-docs <path> [--output <path>] [--config <path>] [--format markdown|html|both]
```

**Step 4: Build and smoke test**

```bash
go build ./cmd/design-reviewer/...
./design-reviewer --help
```

**Step 5: Run all tests**

```bash
go test ./... -v
```

**Step 6: Commit**

```bash
git add scripts/aidlc-designreview/internal/orchestration/ scripts/aidlc-designreview/cmd/
git commit -m "Add designreview orchestration pipeline and Cobra CLI"
```

---

## Phase 3: aidlc-evaluator

**Reference:** `aidlc-workflows-py/scripts/aidlc-evaluator/`
**Go module path:** `github.com/mach4-braai/aidlc-workflows/aidlc-evaluator`
**Structure:** Python uv workspace (10 packages) → single Go module with 10 `internal/` packages

### Task 15: Module scaffold and shared package

**Files:**
- Create: `scripts/aidlc-evaluator/go.mod`
- Create: `scripts/aidlc-evaluator/internal/shared/` (scenario.go, io.go, sandbox.go, credential_scrubber.go)
- Create: `scripts/aidlc-evaluator/internal/shared/shared_test.go`
- Create: `scripts/aidlc-evaluator/config/default.yaml` (copy from Python)
- Create: `scripts/aidlc-evaluator/test_cases/` (copy test case directories from Python)

**Step 1: Initialize module**

```bash
cd scripts/aidlc-evaluator
go mod init github.com/mach4-braai/aidlc-workflows/aidlc-evaluator
go get github.com/spf13/cobra@latest
go get github.com/aws/aws-sdk-go-v2/config@latest
go get github.com/aws/aws-sdk-go-v2/service/bedrockruntime@latest
go get github.com/aws/aws-sdk-go-v2/service/s3@latest
go get github.com/docker/docker/client@latest
go get gopkg.in/yaml.v3@latest
go get github.com/pterm/pterm@latest
go get github.com/cenkalti/backoff/v4@latest
```

**Step 2: Write failing tests** — mirrors `packages/shared/tests/`

```go
// internal/shared/shared_test.go
package shared_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/shared"
)

func TestScrubCredentialsRemovesAWSKey(t *testing.T) {
    input := "token: AKIAIOSFODNN7EXAMPLE and more text"
    output := shared.ScrubCredentials(input)
    if output == input {
        t.Fatal("should scrub AWS access key")
    }
    if contains(output, "AKIAIOSFODNN7EXAMPLE") {
        t.Fatal("scrubbed output must not contain raw key")
    }
}

func TestLoadScenarioFromYAML(t *testing.T) {
    yaml := "name: sci-calc\ndescription: Scientific calculator\n"
    tmp := t.TempDir()
    f := filepath.Join(tmp, "scenario.yaml")
    os.WriteFile(f, []byte(yaml), 0644)
    s, err := shared.LoadScenario(f)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if s.Name != "sci-calc" {
        t.Fatalf("expected sci-calc, got %s", s.Name)
    }
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsHelper(s, sub)) }
func containsHelper(s, sub string) bool {
    for i := 0; i <= len(s)-len(sub); i++ {
        if s[i:i+len(sub)] == sub { return true }
    }
    return false
}
```

**Step 3: Implement shared package** — mirrors `packages/shared/src/shared/`:
- `Scenario` struct with Name, Description, VisionPath, TechEnvPath, OpenAPIPath, GoldenPath
- `LoadScenario(path string) (Scenario, error)` — reads scenario.yaml
- `ScrubCredentials(text string) string` — regex replacement for AWS keys, tokens, passwords
- `AtomicWriteYAML(path string, v any) error` — write-to-temp then rename
- `SandboxConfig` struct

**Step 4: Run — expect PASS**

```bash
go test ./internal/shared/... -v
```

**Step 5: Commit**

```bash
git add scripts/aidlc-evaluator/
git commit -m "Add aidlc-evaluator module scaffold and shared package"
```

---

### Task 16: Execution package (runner + Bedrock agents)

**Files:**
- Create: `scripts/aidlc-evaluator/internal/execution/` (config.go, runner.go, agents/executor.go, agents/simulator.go, tools/file_ops.go, tools/rule_loader.go, tools/run_command.go, metrics.go, post_run.go, progress.go)
- Create: `scripts/aidlc-evaluator/internal/execution/execution_test.go`

**Reference:** `packages/execution/src/aidlc_runner/`

**Step 1: Write failing tests** — mirrors `packages/execution/tests/`

```go
// internal/execution/execution_test.go
package execution_test

import (
    "os/exec"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/execution"
)

func TestRunCommandExecutesBashEcho(t *testing.T) {
    if _, err := exec.LookPath("echo"); err != nil {
        t.Skip("echo not available")
    }
    result := execution.RunCommand("echo hello")
    if result.ExitCode != 0 {
        t.Fatalf("expected exit 0, got %d", result.ExitCode)
    }
    if result.Stdout != "hello\n" {
        t.Fatalf("expected 'hello\\n', got %q", result.Stdout)
    }
}

func TestRunnerConfigValidation(t *testing.T) {
    cfg := execution.RunnerConfig{}
    if err := cfg.Validate(); err == nil {
        t.Fatal("empty config should fail validation")
    }
}

func TestMetricsCollectorTracksTokens(t *testing.T) {
    m := execution.NewMetricsCollector()
    m.AddTokens("executor", 100, 50)
    stats := m.Summary()
    if stats.TotalInputTokens != 100 {
        t.Fatalf("expected 100 input tokens, got %d", stats.TotalInputTokens)
    }
}
```

**Step 2: Implement execution package**

Key types:
```go
type RunnerConfig struct {
    VisionPath    string
    TechEnvPath   string
    AIDLCConfig   AIDLCConfig
    ExecutorModel string
    SimulatorModel string
    OutputDir     string
    AWSProfile    string
    AWSRegion     string
}

type CommandResult struct {
    ExitCode int
    Stdout   string
    Stderr   string
    Duration time.Duration
}

func RunCommand(cmd string) CommandResult
```

The `executor` and `simulator` agents use the same Bedrock Converse API + tool loop pattern established in the traceability agent package. Key tools for the executor agent: `file_ops` (CRUD files), `rule_loader` (read AIDLC rules from local path or ZIP URL), `run_command` (shell execution with sandboxing).

**Step 3: Run — expect PASS**

```bash
go test ./internal/execution/... -v
```

**Step 4: Commit**

```bash
git add scripts/aidlc-evaluator/internal/execution/
git commit -m "Add evaluator execution package with runner and Bedrock agents"
```

---

### Task 17: Docker sandbox (replaces shell-out to build.sh)

**Files:**
- Create: `scripts/aidlc-evaluator/internal/execution/sandbox/sandbox.go`
- Create: `scripts/aidlc-evaluator/internal/execution/sandbox/sandbox_test.go`
- Create: `scripts/aidlc-evaluator/docker/sandbox/Dockerfile` (copy from Python)

**Reference:** `docker/sandbox/build.sh` (Python shells out to this)

**Step 1: Write failing tests**

```go
// internal/execution/sandbox/sandbox_test.go
package sandbox_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/execution/sandbox"
)

func TestSandboxConfigDefaultImage(t *testing.T) {
    cfg := sandbox.DefaultConfig()
    if cfg.Image == "" {
        t.Fatal("default config must specify an image")
    }
}

// Integration test: requires Docker daemon
func TestBuildAndRunSandbox_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping docker integration test")
    }
    cfg := sandbox.DefaultConfig()
    err := sandbox.Build(cfg)
    if err != nil {
        t.Fatalf("build failed: %v", err)
    }
    result, err := sandbox.Run(cfg, "echo hello")
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if result.ExitCode != 0 {
        t.Fatalf("expected exit 0, got %d: %s", result.ExitCode, result.Stderr)
    }
}
```

**Step 2: Implement** using `github.com/docker/docker/client`:

```go
type Config struct {
    Image       string
    Dockerfile  string
    WorkDir     string
    MemoryLimit int64  // bytes
    CPUQuota    int64
    NetworkMode string
}

func DefaultConfig() Config { ... }

// Build builds the sandbox Docker image using the Docker SDK.
func Build(cfg Config) error {
    cli, _ := client.NewClientWithOpts(client.FromEnv)
    // docker.ImageBuild() with build context tar stream
    ...
}

// Run executes cmd inside a fresh sandbox container and returns results.
func Run(cfg Config, cmd string) (CommandResult, error) {
    cli, _ := client.NewClientWithOpts(client.FromEnv)
    // ContainerCreate → ContainerStart → ContainerWait → ContainerLogs
    ...
}
```

**Step 3: Run (short mode)**

```bash
go test ./internal/execution/sandbox/... -v -short
```

**Step 4: Commit**

```bash
git add scripts/aidlc-evaluator/internal/execution/sandbox/ scripts/aidlc-evaluator/docker/
git commit -m "Replace sandbox shell-out with Docker SDK for Go"
```

---

### Task 18: CLI and IDE harness adapters

**Files:**
- Create: `scripts/aidlc-evaluator/internal/cliharness/` (adapter.go, orchestrator.go, normalizer.go, prompt_template.go, registry.go, adapters/claude_code.go, adapters/kiro_cli.go)
- Create: `scripts/aidlc-evaluator/internal/cliharness/cliharness_test.go`
- Create: `scripts/aidlc-evaluator/internal/ideharness/` (adapter.go, orchestrator.go, normalizer.go, prompt_template.go, registry.go, adapters/{cursor,cline,kiro,copilot,windsurf}.go)
- Create: `scripts/aidlc-evaluator/internal/ideharness/ideharness_test.go`

**Reference:** `packages/cli-harness/` and `packages/ide-harness/`

**Step 1: Write failing tests** — mirrors `packages/cli-harness/tests/` and `packages/ide-harness/tests/`

```go
// internal/cliharness/cliharness_test.go
package cliharness_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/cliharness"
)

func TestRegistryListsKnownAdapters(t *testing.T) {
    r := cliharness.NewRegistry()
    adapters := r.List()
    if len(adapters) < 2 {
        t.Fatalf("expected at least 2 adapters, got %d", len(adapters))
    }
}

func TestNormalizerStripsANSIEscapeCodes(t *testing.T) {
    input := "\x1b[32mGreen text\x1b[0m"
    output := cliharness.Normalize(input)
    if output != "Green text" {
        t.Fatalf("expected stripped text, got %q", output)
    }
}
```

**Step 2: Implement** — `Adapter` interface with `Run(scenario, prompt string) (string, error)`. Each adapter wraps `exec.Cmd` to spawn the CLI tool, feeds prompts via stdin or flags, captures stdout. `Normalizer` strips ANSI codes, extracts aidlc-docs paths. `Registry` maps adapter name → factory.

**Step 3: Run — expect PASS**

```bash
go test ./internal/cliharness/... ./internal/ideharness/... -v
```

**Step 4: Commit**

```bash
git add scripts/aidlc-evaluator/internal/cliharness/ scripts/aidlc-evaluator/internal/ideharness/
git commit -m "Add CLI and IDE harness adapter packages"
```

---

### Task 19: Qualitative and quantitative scoring

**Files:**
- Create: `scripts/aidlc-evaluator/internal/qualitative/` (models.go, document.go, comparator.go, scorer.go)
- Create: `scripts/aidlc-evaluator/internal/qualitative/qualitative_test.go`
- Create: `scripts/aidlc-evaluator/internal/quantitative/` (models.go, analyzers.go, scanner.go)
- Create: `scripts/aidlc-evaluator/internal/quantitative/quantitative_test.go`

**Reference:** `packages/qualitative/` and `packages/quantitative/`

**Step 1: Write failing tests** — mirrors `packages/qualitative/tests/` and `packages/quantitative/tests/`

```go
// internal/qualitative/qualitative_test.go
package qualitative_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/qualitative"
)

func TestScoreReturnsOneForIdenticalDocs(t *testing.T) {
    content := "# Requirements\n\n## FR-001: Login\n\nUsers must log in.\n"
    score := qualitative.Score(content, content)
    if score.Percent < 95.0 {
        t.Fatalf("identical docs should score ≥ 95%%, got %.1f", score.Percent)
    }
}

func TestScoreReturnsZeroForEmptyActual(t *testing.T) {
    golden := "# Requirements\n\n## FR-001: Login\n\nUsers must log in.\n"
    score := qualitative.Score(golden, "")
    if score.Percent > 10.0 {
        t.Fatalf("empty actual should score < 10%%, got %.1f", score.Percent)
    }
}
```

```go
// internal/quantitative/quantitative_test.go
package quantitative_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/quantitative"
)

func TestScanCountsLinesOfCode(t *testing.T) {
    tmp := t.TempDir()
    os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"), 0644)
    result := quantitative.Scan(tmp)
    if result.TotalLOC == 0 {
        t.Fatal("expected non-zero LOC")
    }
}
```

**Step 2: Implement both packages**

Qualitative: `Score(golden, actual string) ScoreResult` — compare document sections by heading, compute overlap percentage. For the AI-based comparison, call Bedrock Converse with a scoring prompt (same pattern as traceability agent).

Quantitative: `Scan(dir string) ScanResult` — count LoC (non-blank, non-comment), file count, extension breakdown.

**Step 3: Run — expect PASS**

```bash
go test ./internal/qualitative/... ./internal/quantitative/... -v
```

**Step 4: Commit**

```bash
git add scripts/aidlc-evaluator/internal/qualitative/ scripts/aidlc-evaluator/internal/quantitative/
git commit -m "Add qualitative and quantitative scoring packages"
```

---

### Task 20: Contract testing package

**Files:**
- Create: `scripts/aidlc-evaluator/internal/contracttest/` (runner.go, server.go, spec.go)
- Create: `scripts/aidlc-evaluator/internal/contracttest/contracttest_test.go`

**Reference:** `packages/contracttest/`

**Step 1: Write failing tests** — mirrors `packages/contracttest/tests/`

```go
// internal/contracttest/contracttest_test.go
package contracttest_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/contracttest"
)

func TestParseOpenAPISpecExtractsEndpoints(t *testing.T) {
    yamlContent := `
openapi: "3.0.0"
info:
  title: Test API
  version: "1.0"
paths:
  /health:
    get:
      summary: Health check
      responses:
        "200":
          description: OK
`
    spec, err := contracttest.ParseSpec([]byte(yamlContent))
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(spec.Endpoints) != 1 {
        t.Fatalf("expected 1 endpoint, got %d", len(spec.Endpoints))
    }
}
```

**Step 2: Implement** — `ParseSpec(yaml []byte) (APISpec, error)` parses OpenAPI YAML. `Runner.Run(serverURL string, spec APISpec) (ContractResult, error)` makes HTTP requests and validates responses against spec.

**Step 3: Run — expect PASS**

```bash
go test ./internal/contracttest/... -v
```

**Step 4: Commit**

```bash
git add scripts/aidlc-evaluator/internal/contracttest/
git commit -m "Add contract testing package"
```

---

### Task 21: Reporting and trend reports

**Files:**
- Create: `scripts/aidlc-evaluator/internal/reporting/` (models.go, collector.go, baseline.go, render_html.go, render_md.go)
- Create: `scripts/aidlc-evaluator/internal/reporting/reporting_test.go`
- Create: `scripts/aidlc-evaluator/internal/trendreports/` (models.go, collector.go, fetcher.go, gate.go, sparkline.go, render_html.go, render_md.go, render_yaml.go)
- Create: `scripts/aidlc-evaluator/internal/trendreports/trendreports_test.go`

**Reference:** `packages/reporting/` and `packages/trend-reports/`

**Step 1: Write failing tests** — mirrors both test suites

```go
// internal/reporting/reporting_test.go
package reporting_test

import (
    "strings"
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/reporting"
)

func TestRenderMarkdownContainsRunID(t *testing.T) {
    run := reporting.RunResult{RunID: "2026-05-26T10:00:00"}
    md := reporting.RenderMarkdown(run)
    if !strings.Contains(md, "2026-05-26") {
        t.Fatal("markdown must contain run ID")
    }
}
```

```go
// internal/trendreports/trendreports_test.go
package trendreports_test

import (
    "testing"
    "github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/trendreports"
)

func TestSparklineRendersCorrectWidth(t *testing.T) {
    values := []float64{0.5, 0.6, 0.7, 0.8, 0.75}
    line := trendreports.Sparkline(values)
    if len([]rune(line)) != len(values) {
        t.Fatalf("sparkline width %d != data points %d", len([]rune(line)), len(values))
    }
}

func TestGatePassesWhenAboveThreshold(t *testing.T) {
    passed := trendreports.CheckGate(0.85, 0.80)
    if !passed {
        t.Fatal("0.85 should pass 0.80 threshold")
    }
}

func TestGateFailsWhenBelowThreshold(t *testing.T) {
    passed := trendreports.CheckGate(0.75, 0.80)
    if passed {
        t.Fatal("0.75 should fail 0.80 threshold")
    }
}
```

**Step 2: Implement both packages**

`Sparkline(values []float64) string` — maps values to Unicode block characters (▁▂▃▄▅▆▇█), same as Python `sparkline.py`.

**Step 3: Run — expect PASS**

```bash
go test ./internal/reporting/... ./internal/trendreports/... -v
```

**Step 4: Commit**

```bash
git add scripts/aidlc-evaluator/internal/reporting/ scripts/aidlc-evaluator/internal/trendreports/
git commit -m "Add evaluator reporting and trend-reports packages"
```

---

### Task 22: Evaluator CLI (master dispatcher)

**Files:**
- Create: `scripts/aidlc-evaluator/cmd/aidlc-eval/main.go`
- Create: `scripts/aidlc-evaluator/cmd/aidlc-eval/root.go`
- Create: `scripts/aidlc-evaluator/cmd/aidlc-eval/full.go`
- Create: `scripts/aidlc-evaluator/cmd/aidlc-eval/cli_cmd.go`
- Create: `scripts/aidlc-evaluator/cmd/aidlc-eval/ide_cmd.go`
- Create: `scripts/aidlc-evaluator/cmd/aidlc-eval/batch.go`
- Create: `scripts/aidlc-evaluator/cmd/aidlc-eval/trend.go`

**Reference:** `run.py` (Python dispatcher) and `scripts/run_*.py`

**Step 1: Write failing integration test**

```go
// cmd/aidlc-eval/main_test.go
package main_test

import (
    "os/exec"
    "strings"
    "testing"
)

func TestCLIHelp(t *testing.T) {
    cmd := exec.Command("go", "run", "./cmd/aidlc-eval", "--help")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("help failed: %v\n%s", err, out)
    }
    if !strings.Contains(string(out), "full") {
        t.Fatal("help must mention 'full' subcommand")
    }
}
```

**Step 2: Implement Cobra subcommands** — mirror Python mode dispatch:

```text
aidlc-eval full    [--vision] [--tech-env] [--golden] [--openapi] [--config] [--aws-profile] [--aws-region]
aidlc-eval cli     [--cli {claude-code,kiro-cli}] [--scenario] [--config]
aidlc-eval ide     [--ide {cursor,cline,kiro,windsurf,copilot}] [--scenario] [--config]
aidlc-eval batch   [--models] [--scenario] [--config]
aidlc-eval trend   [--baseline] [--runs-dir] [--output]
aidlc-eval test    (runs go test ./...)
```

**Step 3: Build and smoke test**

```bash
go build ./cmd/aidlc-eval/...
./aidlc-eval --help
./aidlc-eval full --help
```

**Step 4: Run all evaluator tests**

```bash
go test ./... -v
```

**Step 5: Commit**

```bash
git add scripts/aidlc-evaluator/cmd/
git commit -m "Add aidlc-evaluator Cobra CLI dispatcher (all subcommands)"
```

---

## Phase 4: Cross-cutting completion

### Task 23: GoReleaser smoke build

Verify all three tools build cleanly for the release matrix.

**Step 1:**

```bash
# From repo root
goreleaser build --snapshot --clean
```

Expected: `dist/` contains 6 binaries (3 tools × linux_amd64 + darwin_arm64 minimum).

**Step 2:** If any build fails, fix the issue — typically missing `go.sum` entries or import path mismatches.

**Step 3: Commit any fixes**

```bash
git add .
git commit -m "Fix goreleaser build issues"
```

---

### Task 24: Final test sweep

Run all tests across all three modules with the race detector.

**Step 1:**

```bash
cd scripts/aidlc-traceability && go test -race -count=1 ./... && cd -
cd scripts/aidlc-designreview && go test -race -count=1 ./... && cd -
cd scripts/aidlc-evaluator    && go test -race -count=1 ./... && cd -
```

Expected: all PASS, no data races.

**Step 2:** Fix any failing tests before proceeding.

**Step 3: Commit**

```bash
git commit -m "All tests passing: traceability, designreview, evaluator" --allow-empty
```

---

### Task 25: Update AGENTS.md with final structure

Now that all three tools exist, update `AGENTS.md` with the actual package tree, exact build commands, and Go version.

**Step 1:** Update `AGENTS.md` to reflect:
- Actual `go.mod` module paths confirmed
- Exact test commands
- GoReleaser snapshot command
- Integration test flag (`-short` to skip)

**Step 2: Commit**

```bash
git add AGENTS.md
git commit -m "Update AGENTS.md with final Go structure and commands"
```

---

## Key Reference Files

| Python source | Go target |
|---|---|
| `aidlc-traceability/src/traceability/models.py` | `internal/models/models.go` |
| `aidlc-traceability/src/traceability/agent.py` | `internal/agent/agent.go` |
| `aidlc-traceability/src/traceability/pipeline.py` | `internal/pipeline/pipeline.go` |
| `aidlc-designreview/src/design_reviewer/ai_review/base.py` | `internal/aireview/base.go` |
| `aidlc-designreview/src/design_reviewer/reporting/templates/*.jinja2` | `internal/reporting/templates/*.html` + `*.tmpl` |
| `aidlc-evaluator/packages/execution/src/aidlc_runner/runner.py` | `internal/execution/runner.go` |
| `aidlc-evaluator/docker/sandbox/build.sh` | `internal/execution/sandbox/sandbox.go` (Docker SDK) |
| `aidlc-evaluator/packages/trend-reports/src/trend_reports/sparkline.py` | `internal/trendreports/sparkline.go` |

## Dependency Decisions

| Python lib | Go equivalent | Reason |
|---|---|---|
| `click` / `argparse` | `github.com/spf13/cobra` | Standard Go CLI framework |
| `boto3` | `github.com/aws/aws-sdk-go-v2` | Official AWS Go SDK |
| `strands-agents` | Custom `converseLoop()` | No Go port; strands is thin wrapper |
| `networkx` | `gonum.org/v1/gonum/graph/simple` | Most complete Go graph library |
| `pydantic` | Go structs + `encoding/json` | Go is typed; validation in constructors |
| `jinja2` | `html/template` + `text/template` | Stdlib; templates embedded via `embed.FS` |
| `rich` | `github.com/pterm/pterm` | Best Rich equivalent for Go |
| `backoff` | `github.com/cenkalti/backoff/v4` | Direct equivalent |
| `docker` shell-out | `github.com/docker/docker/client` | Docker SDK for Go |
| `pyyaml` | `gopkg.in/yaml.v3` | Standard Go YAML library |
