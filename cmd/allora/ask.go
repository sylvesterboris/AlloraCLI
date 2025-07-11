package main

import (
	"context"
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/agents"
	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newAskCmd() *cobra.Command {
	var agentName string
	var format string
	var interactive bool

	cmd := &cobra.Command{
		Use:   "ask [query]",
		Short: "Ask AI agents about IT infrastructure questions",
		Long: `Ask natural language questions to AI agents about your IT infrastructure.
The agent will analyze your query and provide intelligent responses, suggestions,
and actionable insights based on your infrastructure context.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAsk(args, agentName, format, interactive)
		},
	}

	cmd.Flags().StringVarP(&agentName, "agent", "a", "", "specific agent to use (default: first available)")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "interactive mode for follow-up questions")

	return cmd
}

func runAsk(args []string, agentName, format string, interactive bool) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate that we have at least one agent configured
	if len(cfg.Agents) == 0 {
		return fmt.Errorf("no agents configured. Run 'allora init' to set up your first agent")
	}

	// Select agent
	var selectedAgent config.Agent
	if agentName != "" {
		agent, exists := cfg.Agents[agentName]
		if !exists {
			return fmt.Errorf("agent '%s' not found", agentName)
		}
		selectedAgent = agent
	} else {
		// Use the first available agent
		for name, agent := range cfg.Agents {
			agentName = name
			selectedAgent = agent
			break
		}
	}

	// Initialize agent
	aiAgent, err := agents.NewAgent(selectedAgent)
	if err != nil {
		return fmt.Errorf("failed to initialize agent: %w", err)
	}

	// Join all arguments into a single query
	query := utils.JoinArgs(args)

	if interactive {
		return runInteractiveAsk(aiAgent, query, format)
	}

	return runSingleAsk(aiAgent, query, format)
}

func runSingleAsk(agent agents.Agent, query, format string) error {
	// Show spinner while processing
	spinner := utils.NewSpinner("Processing your question...")
	spinner.Start()

	// Process the query
	ctx := context.Background()
	agentQuery := &agents.Query{
		Text:    query,
		Context: make(map[string]interface{}),
	}
	response, err := agent.Query(ctx, agentQuery)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to process query: %w", err)
	}

	// Format and display response
	return utils.DisplayResponse(response, format)
}

func runInteractiveAsk(agent agents.Agent, initialQuery, format string) error {
	fmt.Println("🤖 Interactive mode - Type 'exit' to quit, 'help' for commands")
	fmt.Println()

	// Process initial query if provided
	if initialQuery != "" {
		fmt.Printf("You: %s\n", initialQuery)
		if err := runSingleAsk(agent, initialQuery, format); err != nil {
			return err
		}
		fmt.Println()
	}

	// Interactive loop
	for {
		query, err := utils.GetUserInput("You: ")
		if err != nil {
			return err
		}

		// Handle special commands
		switch query {
		case "exit", "quit":
			fmt.Println("👋 Goodbye!")
			return nil
		case "help":
			printInteractiveHelp()
			continue
		case "clear":
			utils.ClearScreen()
			continue
		}

		// Process the query
		if err := runSingleAsk(agent, query, format); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println()
	}
}

func printInteractiveHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  exit, quit  - Exit interactive mode")
	fmt.Println("  help        - Show this help message")
	fmt.Println("  clear       - Clear the screen")
	fmt.Println()
	fmt.Println("Example questions:")
	fmt.Println("  \"What's the CPU usage of my servers?\"")
	fmt.Println("  \"How do I optimize my database performance?\"")
	fmt.Println("  \"Check the health of my Kubernetes cluster\"")
	fmt.Println("  \"What security vulnerabilities should I be aware of?\"")
	fmt.Println()
}
