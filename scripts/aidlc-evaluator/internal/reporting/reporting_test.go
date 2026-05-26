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
