package contracttest

import "gopkg.in/yaml.v3"

// Endpoint represents a single API endpoint from an OpenAPI spec.
type Endpoint struct {
	Path   string
	Method string
	Summary string
}

// APISpec holds the parsed OpenAPI specification.
type APISpec struct {
	Title     string
	Version   string
	Endpoints []Endpoint
}

type openAPIDoc struct {
	Info  struct {
		Title   string `yaml:"title"`
		Version string `yaml:"version"`
	} `yaml:"info"`
	Paths map[string]map[string]struct {
		Summary string `yaml:"summary"`
	} `yaml:"paths"`
}

// ParseSpec parses an OpenAPI YAML document into an APISpec.
func ParseSpec(data []byte) (APISpec, error) {
	var doc openAPIDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return APISpec{}, err
	}
	spec := APISpec{
		Title:   doc.Info.Title,
		Version: doc.Info.Version,
	}
	for path, methods := range doc.Paths {
		for method, op := range methods {
			spec.Endpoints = append(spec.Endpoints, Endpoint{
				Path:    path,
				Method:  method,
				Summary: op.Summary,
			})
		}
	}
	return spec, nil
}
