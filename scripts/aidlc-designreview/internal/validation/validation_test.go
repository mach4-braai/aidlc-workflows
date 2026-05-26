package validation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/validation"
)

func TestClassifyApplicationDesignByContent(t *testing.T) {
	content := "# Application Design\n\n## Architecture\n\nMicroservices with REST APIs."
	artType := validation.ClassifyByContent(content)
	if artType != validation.ArtifactTypeApplicationDesign {
		t.Fatalf("expected APPLICATION_DESIGN, got %s", artType)
	}
}

func TestScanDirectoryFindsAidlcDocs(t *testing.T) {
	tmp := t.TempDir()
	docsDir := filepath.Join(tmp, "aidlc-docs", "construction")
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, "application-design.md"), []byte("# Application Design"), 0644)
	result := validation.ScanDirectory(tmp)
	if !result.HasAidlcDocs {
		t.Fatal("should detect aidlc-docs directory")
	}
}

func TestValidateStructurePassesGoodProject(t *testing.T) {
	tmp := t.TempDir()
	docsDir := filepath.Join(tmp, "aidlc-docs", "construction")
	os.MkdirAll(docsDir, 0755)
	os.WriteFile(filepath.Join(docsDir, "application-design.md"), []byte("# Application Design\n"), 0644)
	result := validation.ValidateStructure(tmp)
	if !result.IsValid {
		t.Fatalf("expected valid, got errors: %v", result.Errors)
	}
}
