package reporting

import (
	"bytes"
	"text/template"
)

var markdownReportTmpl = template.Must(template.New("md").Parse(`# Design Review Report{{if .ProjectName}}: {{.ProjectName}}{{end}}

## Summary

| Severity | Count |
|----------|-------|
| Critical | {{.Summary.CriticalCount}} |
| High     | {{.Summary.HighCount}} |
| Medium   | {{.Summary.MediumCount}} |
| Low      | {{.Summary.LowCount}} |
| **Total**| **{{.Summary.TotalFindings}}** |

## Critique Findings

{{if .ReviewResult.Critique.Findings -}}
{{range .ReviewResult.Critique.Findings -}}
### [{{.Severity}}] {{.Title}} ({{.ID}})

{{.Description}}

{{if .Suggestion}}**Suggestion:** {{.Suggestion}}{{end}}

{{end -}}
{{else -}}
No critique findings.
{{end}}

## Alternative Approaches

{{if .ReviewResult.Alternatives.Suggestions -}}
{{range .ReviewResult.Alternatives.Suggestions -}}
### {{.Title}} ({{.ID}})

{{.Description}}

**Tradeoffs:** {{.Tradeoffs}}

{{end -}}
{{else -}}
No alternative suggestions.
{{end}}

## Gap Analysis

{{if .ReviewResult.GapAnalysis.Findings -}}
{{range .ReviewResult.GapAnalysis.Findings -}}
### [{{.Severity}}] {{.Title}} ({{.ID}})

{{.Description}}

{{if .Suggestion}}**Suggestion:** {{.Suggestion}}{{end}}

{{end -}}
{{else -}}
No gap findings.
{{end}}
`))

// RenderMarkdown renders a ReportData as a markdown string.
func RenderMarkdown(data ReportData) (string, error) {
	var buf bytes.Buffer
	if err := markdownReportTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
