package main

import (
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage AlloraCLI configuration",
		Long:  `Manage AlloraCLI configuration including agents, cloud providers, and monitoring settings.`,
	}

	cmd.AddCommand(newConfigShowCmd())
	cmd.AddCommand(newConfigAgentCmd())
	cmd.AddCommand(newConfigCloudCmd())
	cmd.AddCommand(newConfigMonitoringCmd())
	cmd.AddCommand(newConfigSecurityCmd())

	return cmd
}

func newConfigShowCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigShow(format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "output format (yaml, json)")

	return cmd
}

func newConfigAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Manage AI agents",
		Long:  `Add, remove, or modify AI agent configurations.`,
	}

	cmd.AddCommand(newConfigAgentAddCmd())
	cmd.AddCommand(newConfigAgentRemoveCmd())
	cmd.AddCommand(newConfigAgentListCmd())

	return cmd
}

func newConfigAgentAddCmd() *cobra.Command {
	var name, agentType, apiKey, model string
	var maxTokens int
	var temperature float64

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new AI agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigAgentAdd(name, agentType, apiKey, model, maxTokens, temperature)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "agent name (required)")
	cmd.Flags().StringVarP(&agentType, "type", "t", "general", "agent type (general, aws, azure, gcp, kubernetes, monitoring)")
	cmd.Flags().StringVarP(&apiKey, "api-key", "k", "", "API key for the agent")
	cmd.Flags().StringVarP(&model, "model", "m", "gpt-4", "AI model to use")
	cmd.Flags().IntVar(&maxTokens, "max-tokens", 2048, "maximum tokens for responses")
	cmd.Flags().Float64Var(&temperature, "temperature", 0.7, "response creativity (0.0-1.0)")

	cmd.MarkFlagRequired("name")

	return cmd
}

func newConfigAgentRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [agent-name]",
		Short: "Remove an AI agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigAgentRemove(args[0])
		},
	}

	return cmd
}

func newConfigAgentListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured AI agents",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigAgentList()
		},
	}

	return cmd
}

func newConfigCloudCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "Manage cloud provider configurations",
		Long:  `Configure cloud providers (AWS, Azure, GCP) for infrastructure management.`,
	}

	cmd.AddCommand(newConfigCloudAWSCmd())
	cmd.AddCommand(newConfigCloudAzureCmd())
	cmd.AddCommand(newConfigCloudGCPCmd())

	return cmd
}

func newConfigCloudAWSCmd() *cobra.Command {
	var region, profile string

	cmd := &cobra.Command{
		Use:   "aws",
		Short: "Configure AWS settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigCloudAWS(region, profile)
		},
	}

	cmd.Flags().StringVarP(&region, "region", "r", "us-west-2", "AWS region")
	cmd.Flags().StringVarP(&profile, "profile", "p", "default", "AWS profile")

	return cmd
}

func newConfigCloudAzureCmd() *cobra.Command {
	var subscriptionID, tenantID string

	cmd := &cobra.Command{
		Use:   "azure",
		Short: "Configure Azure settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigCloudAzure(subscriptionID, tenantID)
		},
	}

	cmd.Flags().StringVarP(&subscriptionID, "subscription-id", "s", "", "Azure subscription ID")
	cmd.Flags().StringVarP(&tenantID, "tenant-id", "t", "", "Azure tenant ID")

	return cmd
}

func newConfigCloudGCPCmd() *cobra.Command {
	var projectID, region string

	cmd := &cobra.Command{
		Use:   "gcp",
		Short: "Configure GCP settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigCloudGCP(projectID, region)
		},
	}

	cmd.Flags().StringVarP(&projectID, "project-id", "p", "", "GCP project ID")
	cmd.Flags().StringVarP(&region, "region", "r", "us-central1", "GCP region")

	return cmd
}

func newConfigMonitoringCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitoring",
		Short: "Configure monitoring tools",
		Long:  `Configure monitoring tools like Prometheus, Grafana, and other observability platforms.`,
	}

	// Add subcommands for specific monitoring tools
	cmd.AddCommand(newConfigMonitoringPrometheusCmd())
	cmd.AddCommand(newConfigMonitoringGrafanaCmd())

	return cmd
}

func newConfigMonitoringPrometheusCmd() *cobra.Command {
	var endpoint string

	cmd := &cobra.Command{
		Use:   "prometheus",
		Short: "Configure Prometheus settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigMonitoringPrometheus(endpoint)
		},
	}

	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:9090", "Prometheus endpoint")

	return cmd
}

func newConfigMonitoringGrafanaCmd() *cobra.Command {
	var endpoint string

	cmd := &cobra.Command{
		Use:   "grafana",
		Short: "Configure Grafana settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigMonitoringGrafana(endpoint)
		},
	}

	cmd.Flags().StringVarP(&endpoint, "endpoint", "e", "http://localhost:3000", "Grafana endpoint")

	return cmd
}

func newConfigSecurityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: "Configure security settings",
		Long:  `Configure security settings including encryption, audit logging, and access control.`,
	}

	// Add security configuration subcommands
	cmd.AddCommand(newConfigSecurityEncryptionCmd())
	cmd.AddCommand(newConfigSecurityAuditCmd())

	return cmd
}

func newConfigSecurityEncryptionCmd() *cobra.Command {
	var enable bool

	cmd := &cobra.Command{
		Use:   "encryption",
		Short: "Configure encryption settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSecurityEncryption(enable)
		},
	}

	cmd.Flags().BoolVarP(&enable, "enable", "e", true, "enable encryption")

	return cmd
}

func newConfigSecurityAuditCmd() *cobra.Command {
	var enable bool

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Configure audit logging",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigSecurityAudit(enable)
		},
	}

	cmd.Flags().BoolVarP(&enable, "enable", "e", true, "enable audit logging")

	return cmd
}

// Implementation functions
func runConfigShow(format string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	return config.Display(cfg, format)
}

func runConfigAgentAdd(name, agentType, apiKey, model string, maxTokens int, temperature float64) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate agent type
	validTypes := []string{"general", "aws", "azure", "gcp", "kubernetes", "monitoring"}
	if !contains(validTypes, agentType) {
		return fmt.Errorf("invalid agent type: %s. Valid types: %v", agentType, validTypes)
	}

	// Add agent to configuration
	cfg.Agents[name] = config.Agent{
		Type:        agentType,
		APIKey:      apiKey,
		Model:       model,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	// Save configuration
	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ Agent '%s' added successfully\n", name)
	return nil
}

func runConfigAgentRemove(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if _, exists := cfg.Agents[name]; !exists {
		return fmt.Errorf("agent '%s' not found", name)
	}

	delete(cfg.Agents, name)

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ Agent '%s' removed successfully\n", name)
	return nil
}

func runConfigAgentList() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(cfg.Agents) == 0 {
		fmt.Println("No agents configured. Run 'allora config agent add' to add one.")
		return nil
	}

	fmt.Println("Configured agents:")
	for name, agent := range cfg.Agents {
		fmt.Printf("  • %s (type: %s, model: %s)\n", name, agent.Type, agent.Model)
	}

	return nil
}

func runConfigCloudAWS(region, profile string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg.CloudProviders.AWS = config.AWSConfig{
		Region:  region,
		Profile: profile,
	}

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ AWS configuration updated (region: %s, profile: %s)\n", region, profile)
	return nil
}

func runConfigCloudAzure(subscriptionID, tenantID string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg.CloudProviders.Azure = config.AzureConfig{
		SubscriptionID: subscriptionID,
		TenantID:       tenantID,
	}

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ Azure configuration updated\n")
	return nil
}

func runConfigCloudGCP(projectID, region string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg.CloudProviders.GCP = config.GCPConfig{
		ProjectID: projectID,
		Region:    region,
	}

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ GCP configuration updated (project: %s, region: %s)\n", projectID, region)
	return nil
}

func runConfigMonitoringPrometheus(endpoint string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg.Monitoring.Prometheus = config.PrometheusConfig{
		Endpoint: endpoint,
	}

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ Prometheus configuration updated (endpoint: %s)\n", endpoint)
	return nil
}

func runConfigMonitoringGrafana(endpoint string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg.Monitoring.Grafana = config.GrafanaConfig{
		Endpoint: endpoint,
	}

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("✅ Grafana configuration updated (endpoint: %s)\n", endpoint)
	return nil
}

func runConfigSecurityEncryption(enable bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg.Security.Encryption = enable

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	status := "enabled"
	if !enable {
		status = "disabled"
	}
	fmt.Printf("✅ Encryption %s\n", status)
	return nil
}

func runConfigSecurityAudit(enable bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	cfg.Security.AuditLogging = enable

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	status := "enabled"
	if !enable {
		status = "disabled"
	}
	fmt.Printf("✅ Audit logging %s\n", status)
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
