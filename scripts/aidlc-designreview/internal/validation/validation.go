package validation

import (
	"os"
	"path/filepath"
	"strings"
)

// ArtifactType identifies the type of a design artifact.
type ArtifactType string

const (
	ArtifactTypeApplicationDesign  ArtifactType = "APPLICATION_DESIGN"
	ArtifactTypeFunctionalDesign   ArtifactType = "FUNCTIONAL_DESIGN"
	ArtifactTypeTechnicalEnv       ArtifactType = "TECHNICAL_ENVIRONMENT"
	ArtifactTypeUnknown            ArtifactType = "UNKNOWN"
)

// ArtifactInfo holds metadata about a discovered artifact file.
type ArtifactInfo struct {
	Path string
	Type ArtifactType
}

// ScanResult holds the output of ScanDirectory.
type ScanResult struct {
	HasAidlcDocs bool
	Files        []string
}

// ValidationResult holds the output of ValidateStructure.
type ValidationResult struct {
	IsValid   bool
	Errors    []string
	Artifacts []ArtifactInfo
}

// ClassifyByContent inspects markdown content to determine the artifact type.
func ClassifyByContent(content string) ArtifactType {
	lower := strings.ToLower(content)
	switch {
	case strings.Contains(lower, "application design") || strings.Contains(lower, "# application"):
		return ArtifactTypeApplicationDesign
	case strings.Contains(lower, "functional design") || strings.Contains(lower, "function") && strings.Contains(lower, "design"):
		return ArtifactTypeFunctionalDesign
	case strings.Contains(lower, "technical environment") || strings.Contains(lower, "tech env"):
		return ArtifactTypeTechnicalEnv
	default:
		return ArtifactTypeUnknown
	}
}

// ClassifyByFilename infers artifact type from a filename.
func ClassifyByFilename(name string) ArtifactType {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "application-design") || strings.Contains(lower, "app-design"):
		return ArtifactTypeApplicationDesign
	case strings.Contains(lower, "functional-design") || strings.Contains(lower, "func-design"):
		return ArtifactTypeFunctionalDesign
	case strings.Contains(lower, "technical-environment") || strings.Contains(lower, "tech-env"):
		return ArtifactTypeTechnicalEnv
	default:
		return ArtifactTypeUnknown
	}
}

// ScanDirectory checks whether a project root contains an aidlc-docs directory.
func ScanDirectory(projectRoot string) ScanResult {
	aidlcDocs := filepath.Join(projectRoot, "aidlc-docs")
	info, err := os.Stat(aidlcDocs)
	if err != nil || !info.IsDir() {
		return ScanResult{}
	}
	var files []string
	filepath.WalkDir(aidlcDocs, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})
	return ScanResult{HasAidlcDocs: true, Files: files}
}

// ValidateStructure checks that the project has a valid aidlc-docs structure.
func ValidateStructure(projectRoot string) ValidationResult {
	scan := ScanDirectory(projectRoot)
	if !scan.HasAidlcDocs {
		return ValidationResult{
			IsValid: false,
			Errors:  []string{"aidlc-docs directory not found"},
		}
	}
	var artifacts []ArtifactInfo
	var errs []string

	for _, f := range scan.Files {
		t := ClassifyByFilename(filepath.Base(f))
		if t == ArtifactTypeUnknown {
			// Try content-based classification.
			if data, err := os.ReadFile(f); err == nil {
				t = ClassifyByContent(string(data))
			}
		}
		artifacts = append(artifacts, ArtifactInfo{Path: f, Type: t})
	}

	hasAppDesign := false
	for _, a := range artifacts {
		if a.Type == ArtifactTypeApplicationDesign {
			hasAppDesign = true
		}
	}
	if !hasAppDesign {
		errs = append(errs, "no application design document found")
	}

	return ValidationResult{
		IsValid:   len(errs) == 0,
		Errors:    errs,
		Artifacts: artifacts,
	}
}
