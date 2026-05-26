package qualitative

import (
	"math"
	"regexp"
	"strings"
)

// ScoreResult holds the result of a qualitative document comparison.
type ScoreResult struct {
	Percent float64
	Details string
}

var headingRe = regexp.MustCompile(`(?m)^#{1,3}\s+(.+)$`)

// Score compares golden and actual documents by section overlap.
// Returns a percentage representing how well actual covers golden's content.
func Score(golden, actual string) ScoreResult {
	if golden == "" {
		return ScoreResult{Percent: 100.0, Details: "empty golden, trivially satisfied"}
	}
	if actual == "" {
		return ScoreResult{Percent: 0.0, Details: "actual document is empty"}
	}
	if golden == actual {
		return ScoreResult{Percent: 100.0, Details: "documents are identical"}
	}

	goldenWords := tokenize(golden)
	actualWords := tokenize(actual)

	if len(goldenWords) == 0 {
		return ScoreResult{Percent: 100.0}
	}

	// Jaccard-like token overlap.
	goldenSet := toSet(goldenWords)
	actualSet := toSet(actualWords)
	intersection := 0
	for w := range goldenSet {
		if actualSet[w] {
			intersection++
		}
	}
	union := len(goldenSet) + len(actualSet) - intersection
	if union == 0 {
		return ScoreResult{Percent: 0.0}
	}

	pct := math.Round(float64(intersection)/float64(union)*100*10) / 10
	return ScoreResult{Percent: pct, Details: "token overlap scoring"}
}

func tokenize(text string) []string {
	lower := strings.ToLower(text)
	// Split on non-alphanumeric.
	re := regexp.MustCompile(`[^a-z0-9]+`)
	parts := re.Split(lower, -1)
	var tokens []string
	for _, p := range parts {
		if len(p) > 2 {
			tokens = append(tokens, p)
		}
	}
	return tokens
}

func toSet(words []string) map[string]bool {
	s := make(map[string]bool, len(words))
	for _, w := range words {
		s[w] = true
	}
	return s
}
