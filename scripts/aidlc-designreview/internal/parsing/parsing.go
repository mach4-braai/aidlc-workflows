package parsing

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ApplicationDesignModel holds parsed application design content.
type ApplicationDesignModel struct {
	RawContent  string
	FilePaths   []string
	SourceCount int
	Sections    map[string]string
}

// FunctionalDesignModel holds parsed functional design content.
type FunctionalDesignModel struct {
	RawContent  string
	FilePaths   []string
	UnitNames   []string
	SourceCount int
}

// TechnicalEnvironmentModel holds parsed technical environment content.
type TechnicalEnvironmentModel struct {
	RawContent string
	FilePath   string
}

// DesignData aggregates all parsed design artifacts.
type DesignData struct {
	AppDesign  *ApplicationDesignModel
	FuncDesign *FunctionalDesignModel
	TechEnv    *TechnicalEnvironmentModel
}

var unitHeading = regexp.MustCompile(`(?m)^##\s+(.+)$`)

// ParseApplicationDesign merges multiple content strings into a single model.
func ParseApplicationDesign(contents []string) *ApplicationDesignModel {
	raw := strings.Join(contents, "\n\n")
	sections := extractSections(raw)
	return &ApplicationDesignModel{
		RawContent:  raw,
		SourceCount: len(contents),
		Sections:    sections,
	}
}

// ParseFunctionalDesign merges content and extracts unit names from headings.
func ParseFunctionalDesign(contents []string) *FunctionalDesignModel {
	raw := strings.Join(contents, "\n\n")
	var units []string
	for _, m := range unitHeading.FindAllStringSubmatch(raw, -1) {
		units = append(units, strings.TrimSpace(m[1]))
	}
	return &FunctionalDesignModel{
		RawContent:  raw,
		UnitNames:   units,
		SourceCount: len(contents),
	}
}

// ParseTechnicalEnvironment reads a single file into the model.
func ParseTechnicalEnvironment(filePath string) *TechnicalEnvironmentModel {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return &TechnicalEnvironmentModel{FilePath: filePath}
	}
	return &TechnicalEnvironmentModel{
		RawContent: string(data),
		FilePath:   filePath,
	}
}

// LoadDesignData walks aidlcDocsDir and populates a DesignData.
func LoadDesignData(aidlcDocsDir string) DesignData {
	var data DesignData
	var appContents, funcContents []string

	filepath.WalkDir(aidlcDocsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		lower := strings.ToLower(filepath.Base(path))
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		switch {
		case strings.Contains(lower, "application-design") || strings.Contains(lower, "app-design"):
			appContents = append(appContents, string(content))
		case strings.Contains(lower, "functional-design") || strings.Contains(lower, "func-design"):
			funcContents = append(funcContents, string(content))
		case strings.Contains(lower, "technical-environment") || strings.Contains(lower, "tech-env"):
			data.TechEnv = ParseTechnicalEnvironment(path)
		}
		return nil
	})

	if len(appContents) > 0 {
		data.AppDesign = ParseApplicationDesign(appContents)
	}
	if len(funcContents) > 0 {
		data.FuncDesign = ParseFunctionalDesign(funcContents)
	}
	return data
}

var sectionHeading = regexp.MustCompile(`(?m)^##\s+(.+)$`)

func extractSections(content string) map[string]string {
	sections := make(map[string]string)
	matches := sectionHeading.FindAllStringSubmatchIndex(content, -1)
	for i, m := range matches {
		title := strings.TrimSpace(content[m[2]:m[3]])
		end := len(content)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		sections[title] = strings.TrimSpace(content[m[1]:end])
	}
	return sections
}
