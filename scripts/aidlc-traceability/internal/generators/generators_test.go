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
