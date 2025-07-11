package main

import (
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/security"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newSecurityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: "Security auditing and compliance checking",
		Long:  `Perform security audits, vulnerability assessments, and compliance checks with AI-powered insights.`,
	}

	cmd.AddCommand(newSecurityAuditCmd())
	cmd.AddCommand(newSecurityScanCmd())
	cmd.AddCommand(newSecurityComplianceCmd())
	cmd.AddCommand(newSecurityReportCmd())
	cmd.AddCommand(newSecurityRemediateCmd())

	return cmd
}

func newSecurityAuditCmd() *cobra.Command {
	var target string
	var comprehensive bool
	var format string

	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Perform comprehensive security audit",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityAudit(target, comprehensive, format)
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "target resource or service")
	cmd.Flags().BoolVarP(&comprehensive, "comprehensive", "c", false, "perform comprehensive audit")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newSecurityScanCmd() *cobra.Command {
	var target string
	var scanType string
	var format string

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan for vulnerabilities and security issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityScan(target, scanType, format)
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "target to scan")
	cmd.Flags().StringVarP(&scanType, "type", "s", "all", "scan type (vulnerabilities, misconfigurations, secrets, all)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newSecurityComplianceCmd() *cobra.Command {
	var standard string
	var target string
	var format string

	cmd := &cobra.Command{
		Use:   "compliance",
		Short: "Check compliance against security standards",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityCompliance(standard, target, format)
		},
	}

	cmd.Flags().StringVarP(&standard, "standard", "s", "", "compliance standard (SOC2, ISO27001, PCI-DSS, HIPAA)")
	cmd.Flags().StringVarP(&target, "target", "t", "", "target resource or service")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newSecurityReportCmd() *cobra.Command {
	var reportType string
	var period string
	var format string

	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate security reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityReport(reportType, period, format)
		},
	}

	cmd.Flags().StringVarP(&reportType, "type", "t", "summary", "report type (summary, detailed, executive)")
	cmd.Flags().StringVarP(&period, "period", "p", "30d", "report period (e.g., 7d, 30d, 90d)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml, pdf)")

	return cmd
}

func newSecurityRemediateCmd() *cobra.Command {
	var issue string
	var target string
	var autoFix bool
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "remediate",
		Short: "Remediate security issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecurityRemediate(issue, target, autoFix, dryRun)
		},
	}

	cmd.Flags().StringVarP(&issue, "issue", "i", "", "specific issue to remediate")
	cmd.Flags().StringVarP(&target, "target", "t", "", "target resource or service")
	cmd.Flags().BoolVarP(&autoFix, "auto-fix", "a", false, "automatically fix issues")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "show what would be fixed without making changes")

	return cmd
}

// Implementation functions
func runSecurityAudit(target string, comprehensive bool, format string) error {
	sec, err := security.New()
	if err != nil {
		return fmt.Errorf("failed to initialize security module: %w", err)
	}

	options := security.AuditOptions{
		Target:        target,
		Comprehensive: comprehensive,
	}

	spinner := utils.NewSpinner("Performing security audit...")
	spinner.Start()

	audit, err := sec.PerformAudit(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to perform security audit: %w", err)
	}

	return utils.DisplayResponse(audit, format)
}

func runSecurityScan(target, scanType, format string) error {
	sec, err := security.New()
	if err != nil {
		return fmt.Errorf("failed to initialize security module: %w", err)
	}

	options := security.ScanOptions{
		Target:   target,
		ScanType: scanType,
	}

	spinner := utils.NewSpinner("Scanning for security issues...")
	spinner.Start()

	scan, err := sec.PerformScan(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to perform security scan: %w", err)
	}

	return utils.DisplayResponse(scan, format)
}

func runSecurityCompliance(standard, target, format string) error {
	sec, err := security.New()
	if err != nil {
		return fmt.Errorf("failed to initialize security module: %w", err)
	}

	options := security.ComplianceOptions{
		Standard: standard,
		Target:   target,
	}

	spinner := utils.NewSpinner("Checking compliance...")
	spinner.Start()

	compliance, err := sec.CheckCompliance(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to check compliance: %w", err)
	}

	return utils.DisplayResponse(compliance, format)
}

func runSecurityReport(reportType, period, format string) error {
	sec, err := security.New()
	if err != nil {
		return fmt.Errorf("failed to initialize security module: %w", err)
	}

	options := security.ReportOptions{
		Type:   reportType,
		Period: period,
	}

	spinner := utils.NewSpinner("Generating security report...")
	spinner.Start()

	report, err := sec.GenerateReport(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to generate security report: %w", err)
	}

	return utils.DisplayResponse(report, format)
}

func runSecurityRemediate(issue, target string, autoFix, dryRun bool) error {
	sec, err := security.New()
	if err != nil {
		return fmt.Errorf("failed to initialize security module: %w", err)
	}

	options := security.RemediationOptions{
		Issue:   issue,
		Target:  target,
		AutoFix: autoFix,
		DryRun:  dryRun,
	}

	spinner := utils.NewSpinner("Analyzing remediation options...")
	spinner.Start()

	remediation, err := sec.Remediate(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to remediate security issues: %w", err)
	}

	if dryRun {
		fmt.Println("🔍 Security issues that would be remediated:")
	} else {
		fmt.Println("🔧 Security remediation results:")
	}

	return utils.DisplayResponse(remediation, "text")
}
