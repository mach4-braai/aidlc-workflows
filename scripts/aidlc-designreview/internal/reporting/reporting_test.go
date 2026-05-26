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
