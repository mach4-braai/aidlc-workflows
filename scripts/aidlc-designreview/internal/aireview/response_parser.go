package aireview

import (
	"encoding/json"
	"regexp"
	"strings"
)

var codeBlockRe = regexp.MustCompile("(?s)```(?:json)?\\s*\n(.*?)\n```")

// ExtractJSONFromMarkdown strips a markdown JSON code block, returning raw JSON.
// If no code block is found the input is returned as-is.
func ExtractJSONFromMarkdown(text string) string {
	if m := codeBlockRe.FindStringSubmatch(text); m != nil {
		return strings.TrimSpace(m[1])
	}
	return text
}

// ValidateResponseSchema checks that every required key is present in a JSON string.
func ValidateResponseSchema(response string, required map[string]bool) bool {
	var obj map[string]any
	if err := json.Unmarshal([]byte(response), &obj); err != nil {
		return false
	}
	for key := range required {
		if _, ok := obj[key]; !ok {
			return false
		}
	}
	return true
}
