package main

import (
	"github.com/AlloraAi/AlloraCLI/pkg/ui"
	"github.com/spf13/cobra"
)

func newGeminiCmd() *cobra.Command {
	var colorEnabled bool
	var exportFile string

	cmd := &cobra.Command{
		Use:   "gemini",
		Short: "Launch Gemini-style AI interface",
		Long: `Launch the Gemini-style AI interface for natural language interactions.
This provides a chat-like experience similar to Google Gemini, allowing you to 
interact with AlloraAi using natural language for infrastructure management tasks.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create and start the Gemini interface
			geminiInterface := ui.NewGeminiInterface(colorEnabled)
			
			// Set export file if provided
			if exportFile != "" {
				defer func() {
					if err := geminiInterface.ExportConversation(exportFile); err != nil {
						cmd.Printf("Warning: Failed to export conversation: %v\n", err)
					} else {
						cmd.Printf("Conversation exported to: %s\n", exportFile)
					}
				}()
			}
			
			// Start the interface
			return geminiInterface.Start()
		},
	}

	// Add flags
	cmd.Flags().BoolVar(&colorEnabled, "color", true, "Enable colorized output")
	cmd.Flags().StringVar(&exportFile, "export", "", "Export conversation to file when exiting")

	return cmd
}
