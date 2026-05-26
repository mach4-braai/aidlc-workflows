package aireview

import "strings"

var retryableKeywords = []string{
	"ThrottlingException", "ServiceUnavailableException",
	"RequestTimeout", "TooManyRequestsException", "InternalServerException",
}

// IsRetryable reports whether a BedrockAPIError should trigger a retry.
func IsRetryable(err error) bool {
	msg := err.Error()
	for _, kw := range retryableKeywords {
		if strings.Contains(msg, kw) {
			return true
		}
	}
	return false
}
