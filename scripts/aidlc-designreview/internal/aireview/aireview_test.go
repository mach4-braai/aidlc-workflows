package aireview_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/aireview"
)

func TestSeverityHasFourLevels(t *testing.T) {
	levels := []aireview.Severity{
		aireview.SeverityCritical, aireview.SeverityHigh,
		aireview.SeverityMedium, aireview.SeverityLow,
	}
	if len(levels) != 4 {
		t.Fatal("expected 4 severity levels")
	}
}

func TestExtractJSONFromMarkdownCodeBlock(t *testing.T) {
	input := "```json\n{\"key\": \"value\"}\n```"
	extracted := aireview.ExtractJSONFromMarkdown(input)
	if extracted != `{"key": "value"}` {
		t.Fatalf("expected extracted JSON, got %q", extracted)
	}
}

func TestExtractJSONFromPlainText(t *testing.T) {
	input := `{"key": "value"}`
	extracted := aireview.ExtractJSONFromMarkdown(input)
	if extracted != input {
		t.Fatalf("expected same text, got %q", extracted)
	}
}

func TestValidateResponseSchemaRequiresKeys(t *testing.T) {
	response := `{"findings": [], "summary": "ok"}`
	required := map[string]bool{"findings": true}
	if !aireview.ValidateResponseSchema(response, required) {
		t.Fatal("response with required keys should pass validation")
	}
}

func TestValidateResponseSchemaMissingKeys(t *testing.T) {
	response := `{"wrong_key": []}`
	required := map[string]bool{"findings": true}
	if aireview.ValidateResponseSchema(response, required) {
		t.Fatal("response missing required keys should fail validation")
	}
}

func TestIsRetryableDetectsThrottling(t *testing.T) {
	err := &aireview.BedrockAPIError{Message: "ThrottlingException: rate exceeded"}
	if !aireview.IsRetryable(err) {
		t.Fatal("ThrottlingException should be retryable")
	}
}

func TestIsRetryableFalseForBadInput(t *testing.T) {
	err := &aireview.BedrockAPIError{Message: "ValidationException: invalid input"}
	if aireview.IsRetryable(err) {
		t.Fatal("ValidationException should not be retryable")
	}
}
