package reporting

import "github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/aireview"

// ReviewSummary holds counts of findings by severity.
type ReviewSummary struct {
	CriticalCount int
	HighCount     int
	MediumCount   int
	LowCount      int
	TotalFindings int
}

// ReportData aggregates all data needed to render a report.
type ReportData struct {
	Summary      ReviewSummary
	ReviewResult aireview.ReviewResult
	ProjectName  string
}

// BuildReport computes summary counts from a ReviewResult.
func BuildReport(result aireview.ReviewResult) ReportData {
	var s ReviewSummary
	for _, f := range result.Critique.Findings {
		switch f.Severity {
		case aireview.SeverityCritical:
			s.CriticalCount++
		case aireview.SeverityHigh:
			s.HighCount++
		case aireview.SeverityMedium:
			s.MediumCount++
		case aireview.SeverityLow:
			s.LowCount++
		}
	}
	for _, f := range result.GapAnalysis.Findings {
		switch f.Severity {
		case aireview.SeverityCritical:
			s.CriticalCount++
		case aireview.SeverityHigh:
			s.HighCount++
		case aireview.SeverityMedium:
			s.MediumCount++
		case aireview.SeverityLow:
			s.LowCount++
		}
	}
	s.TotalFindings = s.CriticalCount + s.HighCount + s.MediumCount + s.LowCount
	return ReportData{Summary: s, ReviewResult: result}
}
