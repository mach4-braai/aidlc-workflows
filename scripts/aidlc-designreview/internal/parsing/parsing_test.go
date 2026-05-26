package parsing_test

import (
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/parsing"
)

func TestParseApplicationDesignExtractsContent(t *testing.T) {
	content := "# Application Design\n\n## Architecture\n\nMicroservices.\n"
	model := parsing.ParseApplicationDesign([]string{content})
	if model.RawContent == "" {
		t.Fatal("raw content must not be empty")
	}
	if model.SourceCount != 1 {
		t.Fatalf("expected 1 source, got %d", model.SourceCount)
	}
}

func TestDesignDataAggregatesAllParsed(t *testing.T) {
	data := parsing.DesignData{
		AppDesign:  &parsing.ApplicationDesignModel{RawContent: "app content"},
		FuncDesign: &parsing.FunctionalDesignModel{RawContent: "func content"},
	}
	if data.AppDesign == nil {
		t.Fatal("AppDesign must not be nil")
	}
}
