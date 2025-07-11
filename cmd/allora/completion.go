package main

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:
  $ source <(allora completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ allora completion bash > /etc/bash_completion.d/allora
  # macOS:
  $ allora completion bash > /usr/local/etc/bash_completion.d/allora

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ allora completion zsh > "${fpath[1]}/_allora"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ allora completion fish | source

  # To load completions for each session, execute once:
  $ allora completion fish > ~/.config/fish/completions/allora.fish

PowerShell:
  PS> allora completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> allora completion powershell > allora.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}

	return cmd
}
