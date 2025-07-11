package main

import (
	"fmt"
	"time"

	"github.com/AlloraAi/AlloraCLI/pkg/monitor"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newMonitorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Real-time monitoring and alerting",
		Long:  `Monitor your infrastructure in real-time with AI-powered insights and intelligent alerting.`,
	}

	cmd.AddCommand(newMonitorStatusCmd())
	cmd.AddCommand(newMonitorServiceCmd())
	cmd.AddCommand(newMonitorAlertCmd())
	cmd.AddCommand(newMonitorMetricsCmd())
	cmd.AddCommand(newMonitorDashboardCmd())

	return cmd
}

func newMonitorStatusCmd() *cobra.Command {
	var refresh int
	var format string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get overall system status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMonitorStatus(refresh, format)
		},
	}

	cmd.Flags().IntVarP(&refresh, "refresh", "r", 0, "auto-refresh interval in seconds (0 = no refresh)")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

func newMonitorServiceCmd() *cobra.Command {
	var serviceName string
	var detailed bool
	var format string

	cmd := &cobra.Command{
		Use:   "service [service-name]",
		Short: "Monitor specific service health",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				serviceName = args[0]
			}
			return runMonitorService(serviceName, detailed, format)
		},
	}

	cmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "show detailed service information")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

func newMonitorAlertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Manage monitoring alerts",
		Long:  `Create, list, and manage intelligent monitoring alerts.`,
	}

	cmd.AddCommand(newMonitorAlertCreateCmd())
	cmd.AddCommand(newMonitorAlertListCmd())
	cmd.AddCommand(newMonitorAlertDeleteCmd())

	return cmd
}

func newMonitorAlertCreateCmd() *cobra.Command {
	var name, condition, action, severity string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new monitoring alert",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMonitorAlertCreate(name, condition, action, severity, enabled)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "alert name (required)")
	cmd.Flags().StringVarP(&condition, "condition", "c", "", "alert condition (e.g., 'cpu > 80%')")
	cmd.Flags().StringVarP(&action, "action", "a", "", "action to take (e.g., 'notify', 'scale-up')")
	cmd.Flags().StringVarP(&severity, "severity", "s", "medium", "alert severity (low, medium, high, critical)")
	cmd.Flags().BoolVarP(&enabled, "enabled", "e", true, "enable alert")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("condition")
	cmd.MarkFlagRequired("action")

	return cmd
}

func newMonitorAlertListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all monitoring alerts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMonitorAlertList()
		},
	}

	return cmd
}

func newMonitorAlertDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [alert-name]",
		Short: "Delete a monitoring alert",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMonitorAlertDelete(args[0])
		},
	}

	return cmd
}

func newMonitorMetricsCmd() *cobra.Command {
	var metric string
	var duration string
	var format string

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "View system metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMonitorMetrics(metric, duration, format)
		},
	}

	cmd.Flags().StringVarP(&metric, "metric", "m", "", "specific metric to query")
	cmd.Flags().StringVarP(&duration, "duration", "d", "1h", "time range for metrics (e.g., 1h, 24h, 7d)")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml, graph)")

	return cmd
}

func newMonitorDashboardCmd() *cobra.Command {
	var port int
	var host string

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Launch monitoring dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMonitorDashboard(host, port)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 8080, "dashboard port")
	cmd.Flags().StringVarP(&host, "host", "h", "localhost", "dashboard host")

	return cmd
}

// Implementation functions
func runMonitorStatus(refresh int, format string) error {
	mon, err := monitor.New()
	if err != nil {
		return fmt.Errorf("failed to initialize monitor: %w", err)
	}

	if refresh > 0 {
		return runMonitorStatusWithRefresh(mon, refresh, format)
	}

	status, err := mon.GetSystemStatus()
	if err != nil {
		return fmt.Errorf("failed to get system status: %w", err)
	}

	return utils.DisplayResponse(status, format)
}

func runMonitorStatusWithRefresh(mon monitor.Monitor, refresh int, format string) error {
	ticker := time.NewTicker(time.Duration(refresh) * time.Second)
	defer ticker.Stop()

	for {
		utils.ClearScreen()
		status, err := mon.GetSystemStatus()
		if err != nil {
			return fmt.Errorf("failed to get system status: %w", err)
		}

		fmt.Printf("Last updated: %s (refreshing every %ds)\n\n", time.Now().Format("15:04:05"), refresh)
		if err := utils.DisplayResponse(status, format); err != nil {
			return err
		}

		select {
		case <-ticker.C:
			continue
		default:
			// Check for user input to exit
			if utils.CheckForExit() {
				return nil
			}
		}
	}
}

func runMonitorService(serviceName string, detailed bool, format string) error {
	mon, err := monitor.New()
	if err != nil {
		return fmt.Errorf("failed to initialize monitor: %w", err)
	}

	if serviceName == "" {
		// List all services
		services, err := mon.ListServices()
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}
		return utils.DisplayResponse(services, format)
	}

	// Monitor specific service
	service, err := mon.GetServiceStatus(serviceName, detailed)
	if err != nil {
		return fmt.Errorf("failed to get service status: %w", err)
	}

	return utils.DisplayResponse(service, format)
}

func runMonitorAlertCreate(name, condition, action, severity string, enabled bool) error {
	mon, err := monitor.New()
	if err != nil {
		return fmt.Errorf("failed to initialize monitor: %w", err)
	}

	alert := monitor.AlertConfig{
		Name:      name,
		Condition: condition,
		Action:    action,
		Severity:  severity,
		Enabled:   enabled,
		CreatedAt: time.Now(),
	}

	if err := mon.CreateAlert(alert); err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}

	fmt.Printf("✅ Alert '%s' created successfully\n", name)
	return nil
}

func runMonitorAlertList() error {
	mon, err := monitor.New()
	if err != nil {
		return fmt.Errorf("failed to initialize monitor: %w", err)
	}

	alerts, err := mon.ListAlerts()
	if err != nil {
		return fmt.Errorf("failed to list alerts: %w", err)
	}

	if len(alerts) == 0 {
		fmt.Println("No alerts configured.")
		return nil
	}

	fmt.Println("Configured alerts:")
	for _, alert := range alerts {
		fmt.Printf("  • %s (%s) - %s\n", alert.RuleName, alert.Severity, alert.Message)
	}

	return nil
}

func runMonitorAlertDelete(name string) error {
	mon, err := monitor.New()
	if err != nil {
		return fmt.Errorf("failed to initialize monitor: %w", err)
	}

	if err := mon.DeleteAlert(name); err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	fmt.Printf("✅ Alert '%s' deleted successfully\n", name)
	return nil
}

func runMonitorMetrics(metric, duration, format string) error {
	mon, err := monitor.New()
	if err != nil {
		return fmt.Errorf("failed to initialize monitor: %w", err)
	}

	metrics, err := mon.GetMetrics(metric, duration)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}

	return utils.DisplayResponse(metrics, format)
}

func runMonitorDashboard(host string, port int) error {
	mon, err := monitor.New()
	if err != nil {
		return fmt.Errorf("failed to initialize monitor: %w", err)
	}

	fmt.Printf("🚀 Starting monitoring dashboard at http://%s:%d\n", host, port)
	fmt.Println("Press Ctrl+C to stop...")

	return mon.StartDashboard(host, port)
}
