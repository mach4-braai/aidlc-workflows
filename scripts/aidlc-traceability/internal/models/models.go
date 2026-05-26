package models

import "time"

type ArtifactType string

const (
	ArtifactTypeRequirement ArtifactType = "REQUIREMENT"
	ArtifactTypeStory       ArtifactType = "STORY"
	ArtifactTypeUnit        ArtifactType = "UNIT"
	ArtifactTypeComponent   ArtifactType = "COMPONENT"
	ArtifactTypeCodePlan    ArtifactType = "CODE_PLAN"
	ArtifactTypeCode        ArtifactType = "CODE"
	ArtifactTypeTest        ArtifactType = "TEST"
)

type Artifact struct {
	ID         string       `json:"id"`
	Title      string       `json:"title"`
	Type       ArtifactType `json:"artifact_type"`
	Desc       string       `json:"description"`
	SourceFile string       `json:"source_file"`
	SourceLine int          `json:"source_line"`
	Metadata   map[string]any `json:"metadata"`
}

type Relationship struct {
	SourceID         string `json:"source_id"`
	TargetID         string `json:"target_id"`
	RelationshipType string `json:"relationship_type"`
}

type CoverageGap struct {
	ArtifactID    string       `json:"artifact_id"`
	ArtifactTitle string       `json:"artifact_title"`
	ArtifactType  ArtifactType `json:"artifact_type"`
	GapType       string       `json:"gap_type"`
	Description   string       `json:"description"`
}

type CoverageMetrics struct {
	TotalRequirements       int `json:"total_requirements"`
	TotalStories            int `json:"total_stories"`
	TotalUnits              int `json:"total_units"`
	TotalCodeFiles          int `json:"total_code_files"`
	TotalTests              int `json:"total_tests"`
	RequirementsWithStories int `json:"requirements_with_stories"`
	StoriesWithUnits        int `json:"stories_with_units"`
	UnitsWithCode           int `json:"units_with_code"`
	CodeWithTests           int `json:"code_with_tests"`
}

type TraceabilityReport struct {
	ProjectName   string          `json:"project_name"`
	GeneratedAt   time.Time       `json:"generated_at"`
	Artifacts     []Artifact      `json:"artifacts"`
	Relationships []Relationship  `json:"relationships"`
	Gaps          []CoverageGap   `json:"gaps"`
	Metrics       CoverageMetrics `json:"metrics"`
}
