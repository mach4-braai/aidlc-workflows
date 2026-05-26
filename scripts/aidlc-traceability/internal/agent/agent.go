package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

const defaultModelID = "us.anthropic.claude-sonnet-4-20250514-v1:0"

type agentResponse struct {
	Relationships []struct {
		SourceID         string `json:"source_id"`
		TargetID         string `json:"target_id"`
		RelationshipType string `json:"relationship_type"`
	} `json:"relationships"`
	Insights string `json:"insights"`
}

// ParseAgentJSON extracts relationships from a Bedrock agent response.
// Only relationships where both IDs are in validIDs are returned.
func ParseAgentJSON(responseText string, validIDs map[string]bool) ([]models.Relationship, []string) {
	var resp agentResponse
	if err := json.Unmarshal([]byte(responseText), &resp); err != nil {
		return nil, nil
	}
	var rels []models.Relationship
	for _, r := range resp.Relationships {
		if !validIDs[r.SourceID] || !validIDs[r.TargetID] {
			continue
		}
		rels = append(rels, models.Relationship{
			SourceID:         r.SourceID,
			TargetID:         r.TargetID,
			RelationshipType: r.RelationshipType,
		})
	}
	var insights []string
	if resp.Insights != "" {
		insights = []string{resp.Insights}
	}
	return rels, insights
}

// RunReqStoryAnalysis calls Bedrock to infer requirement→story relationships.
func RunReqStoryAnalysis(ctx context.Context, reqs, stories []models.Artifact, profile, region string) ([]models.Relationship, []string) {
	return runAnalysis(ctx, reqs, stories, "traces_to", profile, region)
}

// RunStoryUnitAnalysis calls Bedrock to infer story→unit relationships.
func RunStoryUnitAnalysis(ctx context.Context, stories, units []models.Artifact, profile, region string) ([]models.Relationship, []string) {
	return runAnalysis(ctx, stories, units, "implemented_by", profile, region)
}

// RunUnitComponentAnalysis calls Bedrock to infer unit→component relationships.
func RunUnitComponentAnalysis(ctx context.Context, units, components []models.Artifact, profile, region string) ([]models.Relationship, []string) {
	return runAnalysis(ctx, units, components, "realized_by", profile, region)
}

// RunComponentCodeAnalysis calls Bedrock to infer component→code relationships.
func RunComponentCodeAnalysis(ctx context.Context, components, codeFiles []models.Artifact, profile, region string) ([]models.Relationship, []string) {
	return runAnalysis(ctx, components, codeFiles, "implemented_in", profile, region)
}

func runAnalysis(ctx context.Context, sources, targets []models.Artifact, relType, profile, region string) ([]models.Relationship, []string) {
	client, err := newClient(ctx, profile, region)
	if err != nil {
		return nil, nil
	}

	validIDs := make(map[string]bool, len(sources)+len(targets))
	for _, a := range sources {
		validIDs[a.ID] = true
	}
	for _, a := range targets {
		validIDs[a.ID] = true
	}

	prompt := buildPrompt(sources, targets, relType)
	text, err := converseLoop(ctx, client, defaultModelID, systemPrompt(), prompt)
	if err != nil {
		return nil, nil
	}
	return ParseAgentJSON(text, validIDs)
}

func newClient(ctx context.Context, profile, region string) (*bedrockruntime.Client, error) {
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
	return bedrockruntime.NewFromConfig(cfg), nil
}

func converseLoop(ctx context.Context, client *bedrockruntime.Client, modelID, sysPrompt, userMsg string) (string, error) {
	messages := []types.Message{{
		Role: types.ConversationRoleUser,
		Content: []types.ContentBlock{
			&types.ContentBlockMemberText{Value: userMsg},
		},
	}}
	for {
		resp, err := client.Converse(ctx, &bedrockruntime.ConverseInput{
			ModelId: aws.String(modelID),
			System: []types.SystemContentBlock{
				&types.SystemContentBlockMemberText{Value: sysPrompt},
			},
			Messages: messages,
		})
		if err != nil {
			return "", err
		}
		if resp.StopReason == types.StopReasonEndTurn {
			return extractText(resp.Output), nil
		}
		messages = append(messages, types.Message{
			Role:    types.ConversationRoleAssistant,
			Content: resp.Output.(*types.ConverseOutputMemberMessage).Value.Content,
		})
	}
}

func extractText(output types.ConverseOutput) string {
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

func systemPrompt() string {
	return `You are a traceability analyst. Given two lists of software artifacts, identify relationships between them.
Respond ONLY with valid JSON in this format:
{"relationships": [{"source_id": "...", "target_id": "...", "relationship_type": "..."}], "insights": "..."}`
}

func buildPrompt(sources, targets []models.Artifact, relType string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Identify '%s' relationships between these artifacts.\n\nSources:\n", relType))
	for _, a := range sources {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", a.ID, a.Title))
	}
	sb.WriteString("\nTargets:\n")
	for _, a := range targets {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", a.ID, a.Title))
	}
	return sb.String()
}
