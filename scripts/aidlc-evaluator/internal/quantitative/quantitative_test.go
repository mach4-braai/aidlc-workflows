package quantitative_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/quantitative"
)

func TestScanCountsLinesOfCode(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "main.go"), []byte("package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"), 0644)
	result := quantitative.Scan(tmp)
	if result.TotalLOC == 0 {
		t.Fatal("expected non-zero LOC")
	}
}
