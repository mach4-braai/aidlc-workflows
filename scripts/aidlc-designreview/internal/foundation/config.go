package foundation

import (
	"io/fs"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/mach4-braai/aidlc-workflows/aidlc-designreview/internal/assets"
)

var configFS = assets.FS

// Config is the top-level configuration model.
type Config struct {
	Models  ModelConfig    `yaml:"models"`
	AWS     AWSConfig      `yaml:"aws"`
	Review  ReviewSettings `yaml:"review"`
}

type ModelConfig struct {
	DefaultModel     string `yaml:"default_model"`
	Critique         string `yaml:"critique"`
	Alternatives     string `yaml:"alternatives"`
	Gap              string `yaml:"gap"`
}

type AWSConfig struct {
	Region           string `yaml:"region"`
	ProfileName      string `yaml:"profile_name"`
	GuardrailID      string `yaml:"guardrail_id"`
	GuardrailVersion string `yaml:"guardrail_version"`
}

type ReviewSettings struct {
	SeverityThreshold    string  `yaml:"severity_threshold"`
	EnableAlternatives   bool    `yaml:"enable_alternatives"`
	EnableGapAnalysis    bool    `yaml:"enable_gap_analysis"`
}

// defaultConfig returns sane defaults.
func defaultConfig() Config {
	return Config{
		AWS: AWSConfig{
			Region:      "us-east-1",
			ProfileName: "default",
		},
		Review: ReviewSettings{
			SeverityThreshold:  "medium",
			EnableAlternatives: true,
			EnableGapAnalysis:  true,
		},
	}
}

// LoadDefaultConfig reads the embedded default-config.yaml.
func LoadDefaultConfig() (Config, error) {
	data, err := configFS.ReadFile("config/default-config.yaml")
	if err != nil {
		return defaultConfig(), nil
	}
	cfg := defaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// LoadConfig reads a user config from disk and merges it over defaults.
func LoadConfig(path string) (Config, error) {
	cfg, err := LoadDefaultConfig()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	return cfg, yaml.Unmarshal(data, &cfg)
}

// PatternLibrary holds all loaded design patterns.
type PatternLibrary struct {
	Patterns map[string]string
}

// LoadPatternLibrary reads all pattern markdown files from the embedded FS.
func LoadPatternLibrary() (*PatternLibrary, error) {
	lib := &PatternLibrary{Patterns: make(map[string]string)}
	err := fs.WalkDir(configFS, "config/patterns", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		data, err := configFS.ReadFile(path)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(d.Name(), ".md")
		lib.Patterns[name] = string(data)
		return nil
	})
	return lib, err
}

// PromptManager loads prompt templates from the embedded FS.
type PromptManager struct {
	templates map[string]*template.Template
}

// LoadPromptManager reads all prompt markdown files from the embedded FS.
func LoadPromptManager() (*PromptManager, error) {
	pm := &PromptManager{templates: make(map[string]*template.Template)}
	err := fs.WalkDir(configFS, "config/prompts", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		data, err := configFS.ReadFile(path)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(d.Name(), ".md")
		// Strip version suffix (e.g. critique-v1 → critique)
		if idx := strings.LastIndex(name, "-v"); idx != -1 {
			name = name[:idx]
		}
		tmpl, err := template.New(name).Option("missingkey=zero").Parse(string(data))
		if err != nil {
			return err
		}
		pm.templates[name] = tmpl
		return nil
	})
	return pm, err
}

// BuildAgentPrompt renders the named prompt template with vars substituted.
func (pm *PromptManager) BuildAgentPrompt(name string, vars map[string]string) (string, error) {
	tmpl, ok := pm.templates[name]
	if !ok {
		// Fall back: return a basic prompt if template not found.
		return "Analyze the following design content and provide feedback.", nil
	}
	var sb strings.Builder
	if err := tmpl.Execute(&sb, vars); err != nil {
		return "", err
	}
	result := sb.String()
	if result == "" {
		result = "Analyze the following design content and provide feedback."
	}
	return result, nil
}
