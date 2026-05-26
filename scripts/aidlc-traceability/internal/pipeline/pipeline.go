package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/analysis"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/discovery"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/generators"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/graph"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/parsers"
)

// Config holds all pipeline inputs.
type Config struct {
	ProjectRoot string
	OutputDir   string
	Format      string // "markdown", "html", "both"
	UseAI       bool
	AWSProfile  string
	AWSRegion   string
	Verbose     bool
}

// Run executes the full traceability pipeline and returns the generated report.
func Run(cfg Config) (models.TraceabilityReport, error) {
	projectName := filepath.Base(cfg.ProjectRoot)
	report := models.TraceabilityReport{
		ProjectName: projectName,
		GeneratedAt: time.Now(),
	}

	docsDir := discovery.FindAidlcDocs(cfg.ProjectRoot)

	var artifacts []models.Artifact
	var relationships []models.Relationship

	if docsDir != "" {
		artifacts, relationships = collectArtifacts(docsDir, cfg.ProjectRoot)
	}

	g := graph.Build(artifacts, relationships)
	gaps := analysis.DetectGaps(artifacts, g)
	metrics := analysis.CalculateMetrics(artifacts, relationships)

	report.Artifacts = artifacts
	report.Relationships = relationships
	report.Gaps = gaps
	report.Metrics = metrics

	if err := writeOutputs(cfg, report); err != nil {
		return report, err
	}
	return report, nil
}

func collectArtifacts(docsDir, projectRoot string) ([]models.Artifact, []models.Relationship) {
	var artifacts []models.Artifact
	var relationships []models.Relationship

	candidates := map[string]func(string) []models.Artifact{
		"requirements.md":          parsers.ParseRequirements,
		"stories.md":               parsers.ParseStories,
		"units.md":                 parsers.ParseUnits,
		"code-plans.md":            parsers.ParseCodePlans,
		"application-components.md": parsers.ParseComponents,
	}

	filepath.WalkDir(docsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if fn, ok := candidates[filepath.Base(path)]; ok {
			artifacts = append(artifacts, fn(path)...)
		}
		return nil
	})

	sourceFiles := discovery.DiscoverSourceCode(projectRoot)
	for _, f := range sourceFiles {
		if a := parsers.ParseCodeFile(f, projectRoot); a != nil {
			artifacts = append(artifacts, *a)
		}
	}

	// Separate by type for linker.
	var reqs, stories, units, components []models.Artifact
	for _, a := range artifacts {
		switch a.Type {
		case models.ArtifactTypeRequirement:
			reqs = append(reqs, a)
		case models.ArtifactTypeStory:
			stories = append(stories, a)
		case models.ArtifactTypeUnit:
			units = append(units, a)
		case models.ArtifactTypeComponent:
			components = append(components, a)
		}
	}

	relationships = append(relationships, parsers.InferRequirementStoryLinks(reqs, stories)...)
	relationships = append(relationships, parsers.InferStoryUnitLinks(stories, units)...)
	relationships = append(relationships, parsers.InferUnitComponentLinks(units, components)...)

	return artifacts, relationships
}

func writeOutputs(cfg Config, report models.TraceabilityReport) error {
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return err
	}

	writeMD := cfg.Format == "markdown" || cfg.Format == "both"
	writeHTML := cfg.Format == "html" || cfg.Format == "both"

	if writeMD {
		content := generators.GenerateMarkdown(report)
		name := fmt.Sprintf("traceability-%s.md", report.GeneratedAt.Format("2006-01-02"))
		if err := os.WriteFile(filepath.Join(cfg.OutputDir, name), []byte(content), 0644); err != nil {
			return err
		}
	}
	if writeHTML {
		content := generators.GenerateHTML(report)
		name := fmt.Sprintf("traceability-%s.html", report.GeneratedAt.Format("2006-01-02"))
		if err := os.WriteFile(filepath.Join(cfg.OutputDir, name), []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}
