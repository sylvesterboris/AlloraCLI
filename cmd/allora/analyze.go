package main

import (
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/analyze"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newAnalyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "AI-powered analysis and performance optimization",
		Long:  `Analyze logs, performance metrics, and infrastructure data with AI-powered insights.`,
	}

	cmd.AddCommand(newAnalyzeLogsCmd())
	cmd.AddCommand(newAnalyzePerformanceCmd())
	cmd.AddCommand(newAnalyzeCostsCmd())
	cmd.AddCommand(newAnalyzeSecurityCmd())
	cmd.AddCommand(newAnalyzeCapacityCmd())

	return cmd
}

func newAnalyzeLogsCmd() *cobra.Command {
	var logFile string
	var pattern string
	var timeRange string
	var format string

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Analyze log files with AI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyzeLogs(logFile, pattern, timeRange, format)
		},
	}

	cmd.Flags().StringVarP(&logFile, "file", "f", "", "log file path")
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "search pattern or regex")
	cmd.Flags().StringVarP(&timeRange, "time", "t", "24h", "time range (e.g., 1h, 24h, 7d)")
	cmd.Flags().StringVarP(&format, "format", "o", "text", "output format (text, json, yaml)")

	return cmd
}

func newAnalyzePerformanceCmd() *cobra.Command {
	var service string
	var metric string
	var timeRange string
	var format string

	cmd := &cobra.Command{
		Use:   "performance",
		Short: "Analyze performance metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyzePerformance(service, metric, timeRange, format)
		},
	}

	cmd.Flags().StringVarP(&service, "service", "s", "", "service name")
	cmd.Flags().StringVarP(&metric, "metric", "m", "", "specific metric (cpu, memory, disk, network)")
	cmd.Flags().StringVarP(&timeRange, "time", "t", "1h", "time range (e.g., 1h, 24h, 7d)")
	cmd.Flags().StringVarP(&format, "format", "o", "text", "output format (text, json, yaml)")

	return cmd
}

func newAnalyzeCostsCmd() *cobra.Command {
	var period string
	var service string
	var recommendations bool
	var format string

	cmd := &cobra.Command{
		Use:   "costs",
		Short: "Analyze cloud costs and optimization opportunities",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyzeCosts(period, service, recommendations, format)
		},
	}

	cmd.Flags().StringVarP(&period, "period", "p", "30d", "analysis period (e.g., 7d, 30d, 90d)")
	cmd.Flags().StringVarP(&service, "service", "s", "", "specific service or resource")
	cmd.Flags().BoolVarP(&recommendations, "recommendations", "r", true, "include cost optimization recommendations")
	cmd.Flags().StringVarP(&format, "format", "o", "text", "output format (text, json, yaml)")

	return cmd
}

func newAnalyzeSecurityCmd() *cobra.Command {
	var target string
	var deep bool
	var format string

	cmd := &cobra.Command{
		Use:   "security",
		Short: "Analyze security posture and vulnerabilities",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyzeSecurity(target, deep, format)
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "target resource or service")
	cmd.Flags().BoolVarP(&deep, "deep", "d", false, "perform deep security analysis")
	cmd.Flags().StringVarP(&format, "format", "o", "text", "output format (text, json, yaml)")

	return cmd
}

func newAnalyzeCapacityCmd() *cobra.Command {
	var service string
	var forecast string
	var format string

	cmd := &cobra.Command{
		Use:   "capacity",
		Short: "Analyze capacity and forecast future needs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyzeCapacity(service, forecast, format)
		},
	}

	cmd.Flags().StringVarP(&service, "service", "s", "", "service name")
	cmd.Flags().StringVarP(&forecast, "forecast", "f", "30d", "forecast period (e.g., 7d, 30d, 90d)")
	cmd.Flags().StringVarP(&format, "format", "o", "text", "output format (text, json, yaml)")

	return cmd
}

// Implementation functions
func runAnalyzeLogs(logFile, pattern, timeRange, format string) error {
	analyzer, err := analyze.New()
	if err != nil {
		return fmt.Errorf("failed to initialize analyzer: %w", err)
	}

	options := analyze.LogOptions{
		File:      logFile,
		Pattern:   pattern,
		TimeRange: timeRange,
	}

	spinner := utils.NewSpinner("Analyzing logs...")
	spinner.Start()

	analysis, err := analyzer.AnalyzeLogs(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to analyze logs: %w", err)
	}

	return utils.DisplayResponse(analysis, format)
}

func runAnalyzePerformance(service, metric, timeRange, format string) error {
	analyzer, err := analyze.New()
	if err != nil {
		return fmt.Errorf("failed to initialize analyzer: %w", err)
	}

	options := analyze.PerformanceOptions{
		Service:   service,
		Metric:    metric,
		TimeRange: timeRange,
	}

	spinner := utils.NewSpinner("Analyzing performance metrics...")
	spinner.Start()

	analysis, err := analyzer.AnalyzePerformance(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to analyze performance: %w", err)
	}

	return utils.DisplayResponse(analysis, format)
}

func runAnalyzeCosts(period, service string, recommendations bool, format string) error {
	analyzer, err := analyze.New()
	if err != nil {
		return fmt.Errorf("failed to initialize analyzer: %w", err)
	}

	options := analyze.CostOptions{
		Period:          period,
		Service:         service,
		Recommendations: recommendations,
	}

	spinner := utils.NewSpinner("Analyzing costs...")
	spinner.Start()

	analysis, err := analyzer.AnalyzeCosts(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to analyze costs: %w", err)
	}

	return utils.DisplayResponse(analysis, format)
}

func runAnalyzeSecurity(target string, deep bool, format string) error {
	analyzer, err := analyze.New()
	if err != nil {
		return fmt.Errorf("failed to initialize analyzer: %w", err)
	}

	options := analyze.SecurityOptions{
		Target: target,
		Deep:   deep,
	}

	spinner := utils.NewSpinner("Analyzing security...")
	spinner.Start()

	analysis, err := analyzer.AnalyzeSecurity(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to analyze security: %w", err)
	}

	return utils.DisplayResponse(analysis, format)
}

func runAnalyzeCapacity(service, forecast, format string) error {
	analyzer, err := analyze.New()
	if err != nil {
		return fmt.Errorf("failed to initialize analyzer: %w", err)
	}

	options := analyze.CapacityOptions{
		Service:  service,
		Forecast: forecast,
	}

	spinner := utils.NewSpinner("Analyzing capacity...")
	spinner.Start()

	analysis, err := analyzer.AnalyzeCapacity(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to analyze capacity: %w", err)
	}

	return utils.DisplayResponse(analysis, format)
}
