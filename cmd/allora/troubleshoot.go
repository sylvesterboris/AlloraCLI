package main

import (
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/troubleshoot"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newTroubleshootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "troubleshoot",
		Short: "AI-powered troubleshooting and incident response",
		Long:  `Diagnose and resolve infrastructure issues with AI-powered troubleshooting and automated remediation.`,
	}

	cmd.AddCommand(newTroubleshootIncidentCmd())
	cmd.AddCommand(newTroubleshootSuggestCmd())
	cmd.AddCommand(newTroubleshootAutofixCmd())
	cmd.AddCommand(newTroubleshootDiagnoseCmd())
	cmd.AddCommand(newTroubleshootHistoryCmd())

	return cmd
}

func newTroubleshootIncidentCmd() *cobra.Command {
	var logs string
	var service string
	var severity string
	var format string

	cmd := &cobra.Command{
		Use:   "incident",
		Short: "Analyze incidents and provide solutions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTroubleshootIncident(logs, service, severity, format)
		},
	}

	cmd.Flags().StringVarP(&logs, "logs", "l", "", "path to log file or log content")
	cmd.Flags().StringVarP(&service, "service", "s", "", "service name related to the incident")
	cmd.Flags().StringVarP(&severity, "severity", "v", "medium", "incident severity (low, medium, high, critical)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newTroubleshootSuggestCmd() *cobra.Command {
	var service string
	var issue string
	var context string
	var format string

	cmd := &cobra.Command{
		Use:   "suggest",
		Short: "Get AI-powered remediation suggestions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTroubleshootSuggest(service, issue, context, format)
		},
	}

	cmd.Flags().StringVarP(&service, "service", "s", "", "service name")
	cmd.Flags().StringVarP(&issue, "issue", "i", "", "issue description")
	cmd.Flags().StringVarP(&context, "context", "c", "", "additional context")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newTroubleshootAutofixCmd() *cobra.Command {
	var severity string
	var dryRun bool
	var confirm bool

	cmd := &cobra.Command{
		Use:   "autofix",
		Short: "Automatically fix common issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTroubleshootAutofix(severity, dryRun, confirm)
		},
	}

	cmd.Flags().StringVarP(&severity, "severity", "s", "medium", "maximum severity to auto-fix (low, medium, high)")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "show what would be fixed without making changes")
	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompts")

	return cmd
}

func newTroubleshootDiagnoseCmd() *cobra.Command {
	var target string
	var deep bool
	var format string

	cmd := &cobra.Command{
		Use:   "diagnose [target]",
		Short: "Run comprehensive system diagnostics",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				target = args[0]
			}
			return runTroubleshootDiagnose(target, deep, format)
		},
	}

	cmd.Flags().BoolVarP(&deep, "deep", "d", false, "perform deep diagnostics")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

func newTroubleshootHistoryCmd() *cobra.Command {
	var limit int
	var format string

	cmd := &cobra.Command{
		Use:   "history",
		Short: "View troubleshooting history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTroubleshootHistory(limit, format)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "number of recent entries to show")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

// Implementation functions
func runTroubleshootIncident(logs, service, severity, format string) error {
	ts, err := troubleshoot.New()
	if err != nil {
		return fmt.Errorf("failed to initialize troubleshooter: %w", err)
	}

	incident := troubleshoot.Incident{
		Logs:     logs,
		Service:  service,
		Severity: severity,
	}

	spinner := utils.NewSpinner("Analyzing incident...")
	spinner.Start()

	analysis, err := ts.AnalyzeIncident(incident)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to analyze incident: %w", err)
	}

	return utils.DisplayResponse(analysis, format)
}

func runTroubleshootSuggest(service, issue, context, format string) error {
	ts, err := troubleshoot.New()
	if err != nil {
		return fmt.Errorf("failed to initialize troubleshooter: %w", err)
	}

	request := troubleshoot.SuggestionRequest{
		Service: service,
		Issue:   issue,
		Context: context,
	}

	spinner := utils.NewSpinner("Generating remediation suggestions...")
	spinner.Start()

	suggestions, err := ts.GetSuggestions(request)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to get suggestions: %w", err)
	}

	return utils.DisplayResponse(suggestions, format)
}

func runTroubleshootAutofix(severity string, dryRun, confirm bool) error {
	ts, err := troubleshoot.New()
	if err != nil {
		return fmt.Errorf("failed to initialize troubleshooter: %w", err)
	}

	options := troubleshoot.AutofixOptions{
		Severity: severity,
		DryRun:   dryRun,
		Confirm:  confirm,
	}

	spinner := utils.NewSpinner("Scanning for issues to auto-fix...")
	spinner.Start()

	results, err := ts.AutoFix(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to run autofix: %w", err)
	}

	if dryRun {
		fmt.Println("🔍 Issues that would be fixed:")
	} else {
		fmt.Println("🔧 Auto-fix results:")
	}

	for _, result := range results {
		status := "✅ Fixed"
		if result.Error != "" {
			status = fmt.Sprintf("❌ Error: %s", result.Error)
		} else if dryRun {
			status = "🔄 Would fix"
		}
		fmt.Printf("  • %s - %s\n", result.Issue, status)
	}

	return nil
}

func runTroubleshootDiagnose(target string, deep bool, format string) error {
	ts, err := troubleshoot.New()
	if err != nil {
		return fmt.Errorf("failed to initialize troubleshooter: %w", err)
	}

	options := troubleshoot.DiagnosticOptions{
		Target: target,
		Deep:   deep,
	}

	spinner := utils.NewSpinner("Running diagnostics...")
	spinner.Start()

	diagnostics, err := ts.RunDiagnostics(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to run diagnostics: %w", err)
	}

	return utils.DisplayResponse(diagnostics, format)
}

func runTroubleshootHistory(limit int, format string) error {
	ts, err := troubleshoot.New()
	if err != nil {
		return fmt.Errorf("failed to initialize troubleshooter: %w", err)
	}

	history, err := ts.GetHistory(limit)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	if len(history) == 0 {
		fmt.Println("No troubleshooting history found.")
		return nil
	}

	return utils.DisplayResponse(history, format)
}
