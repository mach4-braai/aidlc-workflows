package reporting

import (
	"bytes"
	"fmt"
	"text/template"
)

// RunResult holds the aggregated result of an evaluation run.
type RunResult struct {
	RunID           string
	ScenarioName    string
	QualityPercent  float64
	TotalLOC        int
	ContractPassed  int
	ContractFailed  int
	Notes           string
}

var mdTmpl = template.Must(template.New("run").Parse(`# Evaluation Run: {{.RunID}}

**Scenario:** {{.ScenarioName}}

## Results

| Metric | Value |
|--------|-------|
| Quality Score | {{printf "%.1f" .QualityPercent}}% |
| Total LOC | {{.TotalLOC}} |
| Contract Tests Passed | {{.ContractPassed}} |
| Contract Tests Failed | {{.ContractFailed}} |

{{if .Notes}}## Notes

{{.Notes}}{{end}}
`))

// RenderMarkdown renders a RunResult as markdown.
func RenderMarkdown(r RunResult) string {
	var buf bytes.Buffer
	if err := mdTmpl.Execute(&buf, r); err != nil {
		return fmt.Sprintf("# Run: %s\n(render error: %v)", r.RunID, err)
	}
	return buf.String()
}
