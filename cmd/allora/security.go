package main

import (
	"context"
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/AlloraAi/AlloraCLI/pkg/security"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newSecurityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: "Security analysis and compliance management",
		Long:  `Comprehensive security analysis, vulnerability scanning, and compliance management with AI-powered recommendations.`,
	}

	cmd.AddCommand(newSecurityScanCmd())
	cmd.AddCommand(newSecurityComplianceCmd())
	cmd.AddCommand(newSecurityAuditCmd())
	cmd.AddCommand(newSecurityReportCmd())
	cmd.AddCommand(newSecurityMonitorCmd())

	return cmd
}

func newSecurityScanCmd() *cobra.Command {
	var target string
	var scanType string
	var format string

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan for security vulnerabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityScan(target, scanType, format)
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "target to scan (IP, domain, or resource)")
	cmd.Flags().StringVarP(&scanType, "type", "T", "comprehensive", "scan type (quick, comprehensive, custom)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newSecurityComplianceCmd() *cobra.Command {
	var standard string
	var format string

	cmd := &cobra.Command{
		Use:   "compliance",
		Short: "Check compliance against security standards",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityCompliance(standard, format)
		},
	}

	cmd.Flags().StringVarP(&standard, "standard", "s", "cis", "compliance standard (cis, pci, sox, hipaa)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newSecurityAuditCmd() *cobra.Command {
	var resource string
	var format string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Audit permissions and access controls",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityAudit(resource, format)
		},
	}

	cmd.Flags().StringVarP(&resource, "resource", "r", "", "resource to audit")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newSecurityReportCmd() *cobra.Command {
	var reportType string
	var format string

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate comprehensive security reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityReport(reportType, format)
		},
	}

	cmd.Flags().StringVarP(&reportType, "type", "t", "summary", "report type (summary, detailed, executive)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml, pdf)")

	return cmd
}

func newSecurityMonitorCmd() *cobra.Command {
	var duration string
	var format string

	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor security events in real-time",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityMonitor(duration, format)
		},
	}

	cmd.Flags().StringVarP(&duration, "duration", "d", "continuous", "monitoring duration (5m, 1h, continuous)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

// Implementation functions
func runSecurityScan(target, scanType, format string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	secService := security.NewSecurityService(cfg)
	ctx := context.Background()

	spinner := utils.NewSpinner("Scanning for security vulnerabilities...")
	spinner.Start()

	result, err := secService.ScanVulnerabilities(ctx, target)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to scan for vulnerabilities: %w", err)
	}

	return utils.DisplayResponse(result, format)
}

func runSecurityCompliance(standard, format string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	secService := security.NewSecurityService(cfg)
	ctx := context.Background()

	spinner := utils.NewSpinner("Checking compliance...")
	spinner.Start()

	result, err := secService.CheckCompliance(ctx, standard)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to check compliance: %w", err)
	}

	return utils.DisplayResponse(result, format)
}

func runSecurityAudit(resource, format string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	secService := security.NewSecurityService(cfg)
	ctx := context.Background()

	spinner := utils.NewSpinner("Auditing permissions...")
	spinner.Start()

	result, err := secService.AuditPermissions(ctx, resource)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to audit permissions: %w", err)
	}

	return utils.DisplayResponse(result, format)
}

func runSecurityReport(reportType, format string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	secService := security.NewSecurityService(cfg)
	ctx := context.Background()

	options := security.ReportOptions{
		Type:           reportType,
		IncludeDetails: true,
		Format:         format,
	}

	spinner := utils.NewSpinner("Generating security report...")
	spinner.Start()

	result, err := secService.GenerateSecurityReport(ctx, options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to generate security report: %w", err)
	}

	return utils.DisplayResponse(result, format)
}

func runSecurityMonitor(duration, format string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	secService := security.NewSecurityService(cfg)
	ctx := context.Background()

	fmt.Println("Starting security monitoring... (Press Ctrl+C to stop)")

	events, err := secService.MonitorSecurityEvents(ctx)
	if err != nil {
		return fmt.Errorf("failed to start security monitoring: %w", err)
	}

	for event := range events {
		if err := utils.DisplayResponse(event, format); err != nil {
			utils.LogError(fmt.Sprintf("Failed to display event: %v", err))
		}
	}

	return nil
}
