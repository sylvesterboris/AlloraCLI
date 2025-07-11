package main

import (
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/cloud"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newCloudCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud",
		Short: "Cloud provider operations and management",
		Long:  `Manage cloud resources across AWS, Azure, and GCP with AI-powered insights and automation.`,
	}

	cmd.AddCommand(newCloudResourcesCmd())
	cmd.AddCommand(newCloudCostsCmd())
	cmd.AddCommand(newCloudOptimizeCmd())
	cmd.AddCommand(newCloudMigrateCmd())
	cmd.AddCommand(newCloudBackupCmd())

	return cmd
}

func newCloudResourcesCmd() *cobra.Command {
	var provider string
	var resourceType string
	var format string

	cmd := &cobra.Command{
		Use:   "resources",
		Short: "Manage cloud resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCloudResources(provider, resourceType, format)
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "cloud provider (aws, azure, gcp)")
	cmd.Flags().StringVarP(&resourceType, "type", "t", "", "resource type (ec2, s3, rds, etc.)")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

func newCloudCostsCmd() *cobra.Command {
	var provider string
	var period string
	var breakdown bool
	var format string

	cmd := &cobra.Command{
		Use:   "costs",
		Short: "Analyze cloud costs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCloudCosts(provider, period, breakdown, format)
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "cloud provider (aws, azure, gcp)")
	cmd.Flags().StringVarP(&period, "period", "d", "30d", "analysis period (e.g., 7d, 30d, 90d)")
	cmd.Flags().BoolVarP(&breakdown, "breakdown", "b", false, "show cost breakdown by service")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

func newCloudOptimizeCmd() *cobra.Command {
	var provider string
	var resourceType string
	var autoApply bool
	var format string

	cmd := &cobra.Command{
		Use:   "optimize",
		Short: "Optimize cloud resources with AI recommendations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCloudOptimize(provider, resourceType, autoApply, format)
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "cloud provider (aws, azure, gcp)")
	cmd.Flags().StringVarP(&resourceType, "type", "t", "", "resource type to optimize")
	cmd.Flags().BoolVarP(&autoApply, "auto-apply", "a", false, "automatically apply optimization recommendations")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newCloudMigrateCmd() *cobra.Command {
	var source string
	var target string
	var plan bool
	var format string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate resources between cloud providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCloudMigrate(source, target, plan, format)
		},
	}

	cmd.Flags().StringVarP(&source, "source", "s", "", "source cloud provider")
	cmd.Flags().StringVarP(&target, "target", "t", "", "target cloud provider")
	cmd.Flags().BoolVarP(&plan, "plan", "p", false, "generate migration plan only")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newCloudBackupCmd() *cobra.Command {
	var provider string
	var resourceType string
	var schedule string
	var format string

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Manage cloud resource backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCloudBackup(provider, resourceType, schedule, format)
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "cloud provider (aws, azure, gcp)")
	cmd.Flags().StringVarP(&resourceType, "type", "t", "", "resource type to backup")
	cmd.Flags().StringVarP(&schedule, "schedule", "s", "", "backup schedule (daily, weekly, monthly)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

// Implementation functions
func runCloudResources(provider, resourceType, format string) error {
	cloudMgr, err := cloud.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cloud manager: %w", err)
	}

	options := cloud.ResourceOptions{
		Provider: provider,
		Type:     resourceType,
	}

	spinner := utils.NewSpinner("Fetching cloud resources...")
	spinner.Start()

	resources, err := cloudMgr.ListResources(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to list cloud resources: %w", err)
	}

	return utils.DisplayResponse(resources, format)
}

func runCloudCosts(provider, period string, breakdown bool, format string) error {
	cloudMgr, err := cloud.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cloud manager: %w", err)
	}

	options := cloud.CostOptions{
		Provider:  provider,
		Period:    period,
		Breakdown: breakdown,
	}

	spinner := utils.NewSpinner("Analyzing cloud costs...")
	spinner.Start()

	costs, err := cloudMgr.AnalyzeCosts(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to analyze cloud costs: %w", err)
	}

	return utils.DisplayResponse(costs, format)
}

func runCloudOptimize(provider, resourceType string, autoApply bool, format string) error {
	cloudMgr, err := cloud.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cloud manager: %w", err)
	}

	options := cloud.OptimizeOptions{
		Provider:  provider,
		Type:      resourceType,
		AutoApply: autoApply,
	}

	spinner := utils.NewSpinner("Generating optimization recommendations...")
	spinner.Start()

	optimization, err := cloudMgr.OptimizeResources(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to optimize cloud resources: %w", err)
	}

	return utils.DisplayResponse(optimization, format)
}

func runCloudMigrate(source, target string, plan bool, format string) error {
	cloudMgr, err := cloud.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cloud manager: %w", err)
	}

	options := cloud.MigrateOptions{
		Source: source,
		Target: target,
		Plan:   plan,
	}

	spinner := utils.NewSpinner("Preparing migration...")
	spinner.Start()

	migration, err := cloudMgr.MigrateResources(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to migrate cloud resources: %w", err)
	}

	return utils.DisplayResponse(migration, format)
}

func runCloudBackup(provider, resourceType, schedule, format string) error {
	cloudMgr, err := cloud.New()
	if err != nil {
		return fmt.Errorf("failed to initialize cloud manager: %w", err)
	}

	options := cloud.BackupOptions{
		Provider: provider,
		Type:     resourceType,
		Schedule: schedule,
	}

	spinner := utils.NewSpinner("Configuring backup...")
	spinner.Start()

	backup, err := cloudMgr.ManageBackups(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to manage cloud backups: %w", err)
	}

	return utils.DisplayResponse(backup, format)
}
