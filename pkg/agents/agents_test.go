package agents

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAgentManager(t *testing.T) {
	manager := NewAgentManager()

	// Test adding an agent
	agent := &MockAgent{
		name:      "test-agent",
		agentType: "general",
	}

	err := manager.AddAgent(agent)
	if err != nil {
		t.Errorf("AddAgent() failed: %v", err)
	}

	// Test getting an agent
	retrievedAgent, err := manager.GetAgent("test-agent")
	if err != nil {
		t.Errorf("GetAgent() failed: %v", err)
	}

	if retrievedAgent.GetName() != "test-agent" {
		t.Errorf("Expected agent name 'test-agent', got '%s'", retrievedAgent.GetName())
	}

	// Test listing agents
	agents := manager.ListAgents()
	if len(agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agents))
	}

	// Test removing an agent
	err = manager.RemoveAgent("test-agent")
	if err != nil {
		t.Errorf("RemoveAgent() failed: %v", err)
	}

	agents = manager.ListAgents()
	if len(agents) != 0 {
		t.Errorf("Expected 0 agents after removal, got %d", len(agents))
	}
}

func TestAgentQuery(t *testing.T) {
	agent := &MockAgent{
		name:      "test-agent",
		agentType: "general",
	}

	ctx := context.Background()
	query := &Query{
		Text: "What is the status of my servers?",
		Context: map[string]interface{}{
			"user_id": "test-user",
		},
	}

	response, err := agent.Query(ctx, query)
	if err != nil {
		t.Errorf("Query() failed: %v", err)
	}

	if response.Text == "" {
		t.Error("Expected non-empty response text")
	}

	if response.Confidence < 0 || response.Confidence > 1 {
		t.Errorf("Expected confidence between 0 and 1, got %f", response.Confidence)
	}
}

func TestAgentCapabilities(t *testing.T) {
	agent := &MockAgent{
		name:      "monitoring-agent",
		agentType: "monitoring",
	}

	capabilities := agent.GetCapabilities()
	if len(capabilities) == 0 {
		t.Error("Expected non-empty capabilities list")
	}

	// Check if the agent has expected capabilities
	hasMonitoring := false
	for _, cap := range capabilities {
		if cap == "monitoring" {
			hasMonitoring = true
			break
		}
	}

	if !hasMonitoring {
		t.Error("Expected monitoring agent to have 'monitoring' capability")
	}
}

func TestAgentStatus(t *testing.T) {
	agent := &MockAgent{
		name:      "test-agent",
		agentType: "general",
	}

	status := agent.GetStatus()
	if status.State != "idle" {
		t.Errorf("Expected initial state 'idle', got '%s'", status.State)
	}

	if status.LastActivity.IsZero() {
		t.Error("Expected non-zero last activity time")
	}
}

func TestAgentConfiguration(t *testing.T) {
	agent := &MockAgent{
		name:      "test-agent",
		agentType: "general",
	}

	// Test getting initial configuration
	config := agent.GetConfiguration()
	if config.Model == "" {
		t.Error("Expected non-empty model name")
	}

	// Test updating configuration
	newConfig := &AgentConfig{
		Model:       "gpt-4",
		MaxTokens:   2000,
		Temperature: 0.5,
	}

	err := agent.UpdateConfiguration(newConfig)
	if err != nil {
		t.Errorf("UpdateConfiguration() failed: %v", err)
	}

	updatedConfig := agent.GetConfiguration()
	if updatedConfig.Model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", updatedConfig.Model)
	}
}

func BenchmarkAgentQuery(b *testing.B) {
	agent := &MockAgent{
		name:      "benchmark-agent",
		agentType: "general",
	}

	ctx := context.Background()
	query := &Query{
		Text: "Test query for benchmarking",
		Context: map[string]interface{}{
			"benchmark": true,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := agent.Query(ctx, query)
		if err != nil {
			b.Errorf("Query() failed: %v", err)
		}
	}
}

// MockAgent is a test implementation of the Agent interface
type MockAgent struct {
	name      string
	agentType string
	config    *AgentConfig
	status    *AgentStatus
}

func (m *MockAgent) GetName() string {
	return m.name
}

func (m *MockAgent) GetType() string {
	return m.agentType
}

func (m *MockAgent) Query(ctx context.Context, query *Query) (*Response, error) {
	// Simulate some processing time
	time.Sleep(10 * time.Millisecond)

	return &Response{
		Text:       fmt.Sprintf("Mock response to: %s", query.Text),
		Confidence: 0.85,
		Metadata: map[string]interface{}{
			"agent_type": m.agentType,
			"timestamp":  time.Now().UTC(),
		},
	}, nil
}

func (m *MockAgent) GetCapabilities() []string {
	switch m.agentType {
	case "monitoring":
		return []string{"monitoring", "alerting", "metrics"}
	case "security":
		return []string{"security", "compliance", "audit"}
	case "cloud":
		return []string{"cloud", "aws", "azure", "gcp"}
	default:
		return []string{"general", "chat", "help"}
	}
}

func (m *MockAgent) GetStatus() *AgentStatus {
	if m.status == nil {
		m.status = &AgentStatus{
			State:        "idle",
			LastActivity: time.Now().UTC(),
			Health:       "healthy",
		}
	}
	return m.status
}

func (m *MockAgent) GetConfiguration() *AgentConfig {
	if m.config == nil {
		m.config = &AgentConfig{
			Model:       "gpt-3.5-turbo",
			MaxTokens:   1000,
			Temperature: 0.7,
		}
	}
	return m.config
}

func (m *MockAgent) UpdateConfiguration(config *AgentConfig) error {
	m.config = config
	return nil
}

func (m *MockAgent) Start() error {
	if m.status == nil {
		m.status = &AgentStatus{}
	}
	m.status.State = "running"
	m.status.LastActivity = time.Now().UTC()
	return nil
}

func (m *MockAgent) Stop() error {
	if m.status == nil {
		m.status = &AgentStatus{}
	}
	m.status.State = "stopped"
	return nil
}

func (m *MockAgent) IsHealthy() bool {
	return m.GetStatus().Health == "healthy"
}
