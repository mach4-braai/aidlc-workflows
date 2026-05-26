package generators

import (
	"bytes"
	"text/template"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

var markdownTmpl = template.Must(template.New("md").Parse(`# Traceability Report: {{.ProjectName}}

Generated: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}

## Summary

| Metric | Count |
|--------|-------|
| Requirements | {{.Metrics.TotalRequirements}} |
| Stories | {{.Metrics.TotalStories}} |
| Units | {{.Metrics.TotalUnits}} |
| Code Files | {{.Metrics.TotalCodeFiles}} |
| Tests | {{.Metrics.TotalTests}} |

## Coverage

| Metric | Count |
|--------|-------|
| Requirements with Stories | {{.Metrics.RequirementsWithStories}} |
| Stories with Units | {{.Metrics.StoriesWithUnits}} |
| Units with Code | {{.Metrics.UnitsWithCode}} |
| Code with Tests | {{.Metrics.CodeWithTests}} |

## Artifacts

{{range .Artifacts -}}
- **{{.ID}}** ({{.Type}}): {{.Title}}
{{end}}

## Gaps

{{if .Gaps -}}
{{range .Gaps -}}
- [{{.GapType}}] **{{.ArtifactID}}**: {{.Description}}
{{end -}}
{{else -}}
No coverage gaps detected.
{{end}}
`))

// GenerateMarkdown renders a TraceabilityReport as a markdown string.
func GenerateMarkdown(r models.TraceabilityReport) string {
	var buf bytes.Buffer
	if err := markdownTmpl.Execute(&buf, r); err != nil {
		return ""
	}
	return buf.String()
}
