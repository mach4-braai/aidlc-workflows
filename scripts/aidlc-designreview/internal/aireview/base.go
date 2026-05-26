package aireview

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	backoff "github.com/cenkalti/backoff/v4"
)

// BaseAgent wraps a Bedrock client with retry logic.
type BaseAgent struct {
	client  *bedrockruntime.Client
	modelID string
}

// NewBaseAgent constructs a BaseAgent loading AWS config from the environment.
func NewBaseAgent(ctx context.Context, modelID, profile, region string) (*BaseAgent, error) {
	opts := []func(*awsconfig.LoadOptions) error{}
	if profile != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(profile))
	}
	if region != "" {
		opts = append(opts, awsconfig.WithRegion(region))
	}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &BaseAgent{
		client:  bedrockruntime.NewFromConfig(cfg),
		modelID: modelID,
	}, nil
}

// InvokeModel sends a prompt to Bedrock with exponential backoff on retryable errors.
func (a *BaseAgent) InvokeModel(ctx context.Context, sysPrompt, userMsg string) (string, TokenUsage, error) {
	var result string
	var usage TokenUsage

	op := func() error {
		resp, err := a.client.Converse(ctx, &bedrockruntime.ConverseInput{
			ModelId: aws.String(a.modelID),
			System: []types.SystemContentBlock{
				&types.SystemContentBlockMemberText{Value: sysPrompt},
			},
			Messages: []types.Message{{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{Value: userMsg},
				},
			}},
		})
		if err != nil {
			apiErr := &BedrockAPIError{Message: err.Error()}
			if IsRetryable(apiErr) {
				return apiErr
			}
			return backoff.Permanent(apiErr)
		}
		result = extractConverseText(resp.Output)
		if resp.Usage != nil {
			usage.InputTokens = int(aws.ToInt32(resp.Usage.InputTokens))
			usage.OutputTokens = int(aws.ToInt32(resp.Usage.OutputTokens))
		}
		return nil
	}

	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 2 * time.Second
	bo.MaxElapsedTime = 60 * time.Second
	if err := backoff.Retry(op, bo); err != nil {
		return "", usage, err
	}
	return result, usage, nil
}

func extractConverseText(output types.ConverseOutput) string {
	msg, ok := output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return ""
	}
	var sb strings.Builder
	for _, block := range msg.Value.Content {
		if t, ok := block.(*types.ContentBlockMemberText); ok {
			sb.WriteString(t.Value)
		}
	}
	return sb.String()
}

// ParseJSON extracts and unmarshals JSON from a potentially markdown-wrapped response.
func ParseJSON(text string, v any) error {
	return json.Unmarshal([]byte(ExtractJSONFromMarkdown(text)), v)
}
