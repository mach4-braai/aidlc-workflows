package contracttest_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/contracttest"
)

func TestParseOpenAPISpecExtractsEndpoints(t *testing.T) {
	yamlContent := `
openapi: "3.0.0"
info:
  title: Test API
  version: "1.0"
paths:
  /health:
    get:
      summary: Health check
      responses:
        "200":
          description: OK
`
	spec, err := contracttest.ParseSpec([]byte(yamlContent))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(spec.Endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(spec.Endpoints))
	}
}
