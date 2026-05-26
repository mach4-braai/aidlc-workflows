package contracttest

import (
	"fmt"
	"net/http"
	"time"
)

// ContractResult holds the outcome of running contract tests.
type ContractResult struct {
	Passed   int
	Failed   int
	Errors   []string
}

// Runner executes HTTP contract tests against a running server.
type Runner struct {
	client *http.Client
}

// NewRunner creates a Runner with a default HTTP client.
func NewRunner() *Runner {
	return &Runner{client: &http.Client{Timeout: 10 * time.Second}}
}

// Run makes requests for each endpoint and validates responses.
func (r *Runner) Run(serverURL string, spec APISpec) (ContractResult, error) {
	var result ContractResult
	for _, ep := range spec.Endpoints {
		url := serverURL + ep.Path
		req, err := http.NewRequest(ep.Method, url, nil)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s %s: %v", ep.Method, ep.Path, err))
			continue
		}
		resp, err := r.client.Do(req)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s %s: %v", ep.Method, ep.Path, err))
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result.Passed++
		} else {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s %s: status %d", ep.Method, ep.Path, resp.StatusCode))
		}
	}
	return result, nil
}
