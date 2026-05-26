package aireview

// Severity classifies the impact level of a finding.
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
)

// CritiqueFinding is a single design issue found by the critique agent.
type CritiqueFinding struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Severity    Severity `json:"severity"`
	Description string   `json:"description"`
	Suggestion  string   `json:"suggestion"`
	Category    string   `json:"category"`
}

// AlternativeSuggestion is an alternative design approach.
type AlternativeSuggestion struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Tradeoffs   string `json:"tradeoffs"`
}

// GapFinding is a missing design concern identified by the gap agent.
type GapFinding struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Severity    Severity `json:"severity"`
	Description string   `json:"description"`
	Suggestion  string   `json:"suggestion"`
}

// TokenUsage tracks Bedrock token consumption.
type TokenUsage struct {
	InputTokens  int
	OutputTokens int
}

// CritiqueResult holds the output of the critique agent.
type CritiqueResult struct {
	Findings []CritiqueFinding
	Summary  string
	Usage    TokenUsage
}

// AlternativesResult holds the output of the alternatives agent.
type AlternativesResult struct {
	Suggestions []AlternativeSuggestion
	Summary     string
	Usage       TokenUsage
}

// GapAnalysisResult holds the output of the gap agent.
type GapAnalysisResult struct {
	Findings []GapFinding
	Summary  string
	Usage    TokenUsage
}

// ReviewResult aggregates outputs from all three agents.
type ReviewResult struct {
	Critique     CritiqueResult
	Alternatives AlternativesResult
	GapAnalysis  GapAnalysisResult
	TotalUsage   TokenUsage
}

// BedrockAPIError wraps a Bedrock API error with the raw message.
type BedrockAPIError struct {
	Message string
}

func (e *BedrockAPIError) Error() string { return e.Message }
