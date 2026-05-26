package generators

import (
	"bytes"
	"html/template"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

var htmlTmpl = template.Must(template.New("html").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Traceability Report: {{.ProjectName}}</title>
  <style>
    body { font-family: sans-serif; max-width: 1200px; margin: 0 auto; padding: 2rem; }
    table { border-collapse: collapse; width: 100%; margin-bottom: 1.5rem; }
    th, td { border: 1px solid #ddd; padding: 0.5rem 1rem; text-align: left; }
    th { background: #f5f5f5; }
    .gap { color: #c0392b; }
  </style>
</head>
<body>
  <h1>Traceability Report: {{.ProjectName}}</h1>
  <p>Generated: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</p>

  <h2>Summary</h2>
  <table>
    <tr><th>Metric</th><th>Count</th></tr>
    <tr><td>Requirements</td><td>{{.Metrics.TotalRequirements}}</td></tr>
    <tr><td>Stories</td><td>{{.Metrics.TotalStories}}</td></tr>
    <tr><td>Units</td><td>{{.Metrics.TotalUnits}}</td></tr>
    <tr><td>Code Files</td><td>{{.Metrics.TotalCodeFiles}}</td></tr>
    <tr><td>Tests</td><td>{{.Metrics.TotalTests}}</td></tr>
  </table>

  <h2>Artifacts</h2>
  <table>
    <tr><th>ID</th><th>Type</th><th>Title</th></tr>
    {{range .Artifacts -}}
    <tr><td>{{.ID}}</td><td>{{.Type}}</td><td>{{.Title}}</td></tr>
    {{end}}
  </table>

  {{if .Gaps}}
  <h2>Coverage Gaps</h2>
  <table>
    <tr><th>Artifact</th><th>Gap Type</th><th>Description</th></tr>
    {{range .Gaps -}}
    <tr class="gap"><td>{{.ArtifactID}}</td><td>{{.GapType}}</td><td>{{.Description}}</td></tr>
    {{end}}
  </table>
  {{end}}
</body>
</html>
`))

// GenerateHTML renders a TraceabilityReport as an HTML string.
func GenerateHTML(r models.TraceabilityReport) string {
	var buf bytes.Buffer
	if err := htmlTmpl.Execute(&buf, r); err != nil {
		return ""
	}
	return buf.String()
}
