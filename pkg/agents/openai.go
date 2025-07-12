package agents

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/sashabaranov/go-openai"
)

// OpenAIAgent implements the Agent interface using OpenAI's GPT models
type OpenAIAgent struct {
	*BaseAgent
	client       *openai.Client
	systemPrompt string
}

// NewOpenAIAgent creates a new OpenAI-powered agent
func NewOpenAIAgent(cfg config.Agent, agentType string) (*OpenAIAgent, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	client := openai.NewClient(cfg.APIKey)

	baseAgent := &BaseAgent{
		name:    fmt.Sprintf("openai-%s", agentType),
		config:  cfg,
		context: context.Background(),
	}

	agent := &OpenAIAgent{
		BaseAgent:    baseAgent,
		client:       client,
		systemPrompt: getSystemPrompt(agentType),
	}

	return agent, nil
}

// Query processes a query using OpenAI's GPT model
func (o *OpenAIAgent) Query(ctx context.Context, query *Query) (*Response, error) {
	// Update last activity
	if o.status == nil {
		o.status = &AgentStatus{}
	}
	o.status.LastActivity = time.Now().UTC()
	o.status.State = "processing"

	// Prepare the conversation
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: o.systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: query.Text,
		},
	}

	// Add context if available
	if len(query.Context) > 0 {
		contextStr := formatContext(query.Context)
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf("Additional context: %s", contextStr),
		})
	}

	// Make the API call
	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       o.config.Model,
		Messages:    messages,
		MaxTokens:   o.config.MaxTokens,
		Temperature: float32(o.config.Temperature),
	})

	if err != nil {
		o.status.State = "error"
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		o.status.State = "error"
		return nil, fmt.Errorf("no response from OpenAI")
	}

	o.status.State = "idle"

	// Parse the response for actions and suggestions
	content := resp.Choices[0].Message.Content
	actions := parseActions(content)
	suggestions := parseSuggestions(content)

	return &Response{
		Text:       content,
		Content:    content,
		Type:       "text",
		Confidence: calculateConfidence(resp.Usage),
		Metadata: map[string]interface{}{
			"agent_type":        o.GetType(),
			"model":             o.config.Model,
			"tokens_used":       resp.Usage.TotalTokens,
			"prompt_tokens":     resp.Usage.PromptTokens,
			"completion_tokens": resp.Usage.CompletionTokens,
			"finish_reason":     resp.Choices[0].FinishReason,
		},
		Suggestions: suggestions,
		Actions:     actions,
		Timestamp:   time.Now().UTC(),
	}, nil
}

// GetCapabilities returns the capabilities of the OpenAI agent
func (o *OpenAIAgent) GetCapabilities() []string {
	agentType := o.GetType()
	baseCapabilities := []string{"chat", "analysis", "recommendations"}

	switch agentType {
	case "openai-general":
		return append(baseCapabilities, "general", "help", "documentation")
	case "openai-aws":
		return append(baseCapabilities, "aws", "cloud", "ec2", "s3", "rds", "lambda")
	case "openai-azure":
		return append(baseCapabilities, "azure", "cloud", "vm", "storage", "functions")
	case "openai-gcp":
		return append(baseCapabilities, "gcp", "cloud", "compute", "storage", "functions")
	case "openai-kubernetes":
		return append(baseCapabilities, "kubernetes", "k8s", "pods", "deployments", "services")
	case "openai-monitoring":
		return append(baseCapabilities, "monitoring", "metrics", "alerts", "dashboards")
	case "openai-security":
		return append(baseCapabilities, "security", "compliance", "audit", "vulnerabilities")
	default:
		return baseCapabilities
	}
}

// GetType returns the agent type
func (o *OpenAIAgent) GetType() string {
	return o.name
}

// getSystemPrompt returns the system prompt based on agent type
func getSystemPrompt(agentType string) string {
	basePrompt := "You are an AI assistant specialized in IT infrastructure management and DevOps operations. "

	switch agentType {
	case "general":
		return basePrompt + "You help with general IT infrastructure questions, provide guidance on best practices, and assist with troubleshooting. Always provide actionable advice and consider security implications."
	case "aws":
		return basePrompt + "You are an AWS expert. Help with AWS services, architecture, cost optimization, security, and best practices. Provide specific AWS CLI commands and configuration examples when relevant."
	case "azure":
		return basePrompt + "You are an Azure expert. Help with Azure services, architecture, cost optimization, security, and best practices. Provide specific Azure CLI commands and configuration examples when relevant."
	case "gcp":
		return basePrompt + "You are a Google Cloud Platform expert. Help with GCP services, architecture, cost optimization, security, and best practices. Provide specific gcloud commands and configuration examples when relevant."
	case "kubernetes":
		return basePrompt + "You are a Kubernetes expert. Help with K8s deployments, services, troubleshooting, scaling, and best practices. Provide kubectl commands and YAML configurations when relevant."
	case "monitoring":
		return basePrompt + "You are a monitoring and observability expert. Help with metrics, logging, alerting, dashboards, and performance monitoring. Focus on Prometheus, Grafana, and other monitoring tools."
	case "security":
		return basePrompt + "You are a security expert. Help with security analysis, compliance, vulnerability assessment, and security best practices. Always prioritize security in your recommendations."
	default:
		return basePrompt + "You provide general IT infrastructure assistance and guidance."
	}
}

// formatContext converts the context map to a readable string
func formatContext(context map[string]interface{}) string {
	var parts []string
	for key, value := range context {
		parts = append(parts, fmt.Sprintf("%s: %v", key, value))
	}
	return strings.Join(parts, ", ")
}

// parseActions extracts actionable items from the response
func parseActions(content string) []Action {
	var actions []Action

	// Look for command patterns
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for AWS CLI commands
		if strings.Contains(line, "aws ") && strings.HasPrefix(line, "aws ") {
			actions = append(actions, Action{
				Type:        "aws-command",
				Description: "Execute AWS CLI command",
				Command:     line,
				Parameters:  map[string]interface{}{"type": "aws-cli"},
				Risk:        "medium",
			})
		}

		// Look for kubectl commands
		if strings.Contains(line, "kubectl ") && strings.HasPrefix(line, "kubectl ") {
			actions = append(actions, Action{
				Type:        "kubectl-command",
				Description: "Execute kubectl command",
				Command:     line,
				Parameters:  map[string]interface{}{"type": "kubectl"},
				Risk:        "medium",
			})
		}

		// Look for Azure CLI commands
		if strings.Contains(line, "az ") && strings.HasPrefix(line, "az ") {
			actions = append(actions, Action{
				Type:        "azure-command",
				Description: "Execute Azure CLI command",
				Command:     line,
				Parameters:  map[string]interface{}{"type": "azure-cli"},
				Risk:        "medium",
			})
		}

		// Look for gcloud commands
		if strings.Contains(line, "gcloud ") && strings.HasPrefix(line, "gcloud ") {
			actions = append(actions, Action{
				Type:        "gcp-command",
				Description: "Execute gcloud command",
				Command:     line,
				Parameters:  map[string]interface{}{"type": "gcloud"},
				Risk:        "medium",
			})
		}
	}

	return actions
}

// parseSuggestions extracts suggestions from the response
func parseSuggestions(content string) []string {
	var suggestions []string

	// Look for common suggestion patterns
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for lines that start with suggestion indicators
		if strings.HasPrefix(line, "•") || strings.HasPrefix(line, "-") ||
			strings.HasPrefix(line, "1.") || strings.HasPrefix(line, "2.") ||
			strings.HasPrefix(line, "Consider") || strings.HasPrefix(line, "Recommend") ||
			strings.HasPrefix(line, "Suggestion:") {
			suggestions = append(suggestions, line)
		}
	}

	// If no specific suggestions found, generate generic ones
	if len(suggestions) == 0 {
		suggestions = []string{
			"Run diagnostic commands to gather more information",
			"Check logs for error messages",
			"Verify configuration settings",
			"Consider monitoring and alerting setup",
		}
	}

	return suggestions
}

// calculateConfidence calculates confidence based on token usage
func calculateConfidence(usage openai.Usage) float64 {
	// Higher token usage generally indicates more comprehensive responses
	if usage.TotalTokens > 1000 {
		return 0.9
	} else if usage.TotalTokens > 500 {
		return 0.8
	} else if usage.TotalTokens > 200 {
		return 0.7
	} else {
		return 0.6
	}
}
