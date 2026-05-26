package shared

import (
	"encoding/json"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Scenario describes an evaluation scenario loaded from scenario.yaml.
type Scenario struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	VisionPath  string `yaml:"vision_path"`
	TechEnvPath string `yaml:"tech_env_path"`
	OpenAPIPath string `yaml:"openapi_path"`
	GoldenPath  string `yaml:"golden_path"`
}

// SandboxConfig describes Docker sandbox parameters.
type SandboxConfig struct {
	Image       string
	WorkDir     string
	MemoryLimit int64
	CPUQuota    int64
	NetworkMode string
}

// LoadScenario reads a scenario definition from a YAML file.
func LoadScenario(path string) (Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, err
	}
	var s Scenario
	return s, yaml.Unmarshal(data, &s)
}

var credentialPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(AKIA[0-9A-Z]{16})`),
	regexp.MustCompile(`(?i)(aws_secret_access_key\s*=\s*\S+)`),
	regexp.MustCompile(`(?i)(password\s*[:=]\s*\S+)`),
	regexp.MustCompile(`(?i)(token\s*[:=]\s*[A-Za-z0-9+/]{20,})`),
	regexp.MustCompile(`(?i)(api[_-]?key\s*[:=]\s*\S+)`),
}

// ScrubCredentials replaces known credential patterns with [REDACTED].
func ScrubCredentials(text string) string {
	for _, re := range credentialPatterns {
		text = re.ReplaceAllString(text, "[REDACTED]")
	}
	return text
}

// AtomicWriteYAML marshals v to YAML and writes it atomically via a temp file.
func AtomicWriteYAML(path string, v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return atomicWrite(path, data)
}

// AtomicWriteJSON marshals v to JSON and writes it atomically via a temp file.
func AtomicWriteJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return atomicWrite(path, data)
}

func atomicWrite(path string, data []byte) error {
	tmp, err := os.CreateTemp("", "aidlc-eval-*")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	tmp.Close()
	return os.Rename(tmp.Name(), path)
}
