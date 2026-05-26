package discovery_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/discovery"
)

func TestFindAidlcDocsReturnsNilWhenAbsent(t *testing.T) {
	tmp := t.TempDir()
	result := discovery.FindAidlcDocs(tmp)
	if result != "" {
		t.Fatalf("expected empty, got %q", result)
	}
}

func TestFindAidlcDocsFindsDirectory(t *testing.T) {
	tmp := t.TempDir()
	docsDir := filepath.Join(tmp, "aidlc-docs")
	os.MkdirAll(docsDir, 0755)
	result := discovery.FindAidlcDocs(tmp)
	if result != docsDir {
		t.Fatalf("expected %q, got %q", docsDir, result)
	}
}

func TestDiscoverSourceCodeFindsGoFiles(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main"), 0644)
	files := discovery.DiscoverSourceCode(tmp)
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
}

func TestDiscoverSourceCodeExcludesTestFiles(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmp, "main_test.go"), []byte("package main"), 0644)
	files := discovery.DiscoverSourceCode(tmp)
	for _, f := range files {
		if filepath.Base(f) == "main_test.go" {
			t.Fatal("should not include test files")
		}
	}
}
