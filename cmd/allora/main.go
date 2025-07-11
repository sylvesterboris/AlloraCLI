package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	if err := newRootCmd().ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var configFile string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "allora",
		Short: "AI-Powered IT Infrastructure Management CLI",
		Long: `AlloraCLI is a powerful command-line interface for AI agents specialized 
in IT infrastructure management and automation. Built by AlloraAi, it provides 
intelligent automation for DevOps and IT operations through natural language 
processing and multi-agent AI systems.`,
		Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Initialize configuration
			if err := config.Initialize(configFile, verbose); err != nil {
				return fmt.Errorf("failed to initialize configuration: %w", err)
			}

			// Initialize logging
			if err := utils.InitializeLogging(verbose); err != nil {
				return fmt.Errorf("failed to initialize logging: %w", err)
			}

			return nil
		},
	}

	// Global flags
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.config/alloracli/config.yaml)")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))

	// Add subcommands
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newAskCmd())
	cmd.AddCommand(newMonitorCmd())
	cmd.AddCommand(newTroubleshootCmd())
	cmd.AddCommand(newDeployCmd())
	cmd.AddCommand(newAnalyzeCmd())
	cmd.AddCommand(newSecurityCmd())
	cmd.AddCommand(newCloudCmd())
	cmd.AddCommand(newPluginCmd())
	cmd.AddCommand(newCompletionCmd())
	cmd.AddCommand(newGeminiCmd())

	// Enable auto-completion
	cmd.CompletionOptions.DisableDefaultCmd = false

	return cmd
}
