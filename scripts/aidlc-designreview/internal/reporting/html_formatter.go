package reporting

import (
	"bytes"
	"html/template"
)

var htmlReportTmpl = template.Must(template.New("html").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Design Review Report{{if .ProjectName}}: {{.ProjectName}}{{end}}</title>
  <style>
    body { font-family: sans-serif; max-width: 1200px; margin: 0 auto; padding: 2rem; }
    table { border-collapse: collapse; width: 100%; margin-bottom: 1.5rem; }
    th, td { border: 1px solid #ddd; padding: 0.5rem 1rem; text-align: left; }
    th { background: #f5f5f5; }
    .CRITICAL { color: #c0392b; font-weight: bold; }
    .HIGH { color: #e67e22; font-weight: bold; }
    .MEDIUM { color: #f39c12; }
    .LOW { color: #27ae60; }
  </style>
</head>
<body>
  <h1>Design Review Report{{if .ProjectName}}: {{.ProjectName}}{{end}}</h1>

  <h2>Summary</h2>
  <table>
    <tr><th>Severity</th><th>Count</th></tr>
    <tr><td class="CRITICAL">Critical</td><td>{{.Summary.CriticalCount}}</td></tr>
    <tr><td class="HIGH">High</td><td>{{.Summary.HighCount}}</td></tr>
    <tr><td class="MEDIUM">Medium</td><td>{{.Summary.MediumCount}}</td></tr>
    <tr><td class="LOW">Low</td><td>{{.Summary.LowCount}}</td></tr>
    <tr><td><strong>Total</strong></td><td><strong>{{.Summary.TotalFindings}}</strong></td></tr>
  </table>

  {{if .ReviewResult.Critique.Findings}}
  <h2>Critique Findings</h2>
  {{range .ReviewResult.Critique.Findings}}
  <div class="finding">
    <h3><span class="{{.Severity}}">[{{.Severity}}]</span> {{.Title}} ({{.ID}})</h3>
    <p>{{.Description}}</p>
    {{if .Suggestion}}<p><strong>Suggestion:</strong> {{.Suggestion}}</p>{{end}}
  </div>
  {{end}}
  {{end}}

  {{if .ReviewResult.Alternatives.Suggestions}}
  <h2>Alternative Approaches</h2>
  {{range .ReviewResult.Alternatives.Suggestions}}
  <div class="suggestion">
    <h3>{{.Title}} ({{.ID}})</h3>
    <p>{{.Description}}</p>
    <p><strong>Tradeoffs:</strong> {{.Tradeoffs}}</p>
  </div>
  {{end}}
  {{end}}

  {{if .ReviewResult.GapAnalysis.Findings}}
  <h2>Gap Analysis</h2>
  {{range .ReviewResult.GapAnalysis.Findings}}
  <div class="finding">
    <h3><span class="{{.Severity}}">[{{.Severity}}]</span> {{.Title}} ({{.ID}})</h3>
    <p>{{.Description}}</p>
    {{if .Suggestion}}<p><strong>Suggestion:</strong> {{.Suggestion}}</p>{{end}}
  </div>
  {{end}}
  {{end}}
</body>
</html>
`))

// RenderHTML renders a ReportData as an HTML string.
func RenderHTML(data ReportData) (string, error) {
	var buf bytes.Buffer
	if err := htmlReportTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
