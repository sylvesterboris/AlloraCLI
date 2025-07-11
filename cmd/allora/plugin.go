package main

import (
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/plugins"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage AlloraCLI plugins",
		Long:  `Install, manage, and develop plugins to extend AlloraCLI functionality.`,
	}

	cmd.AddCommand(newPluginListCmd())
	cmd.AddCommand(newPluginInstallCmd())
	cmd.AddCommand(newPluginRemoveCmd())
	cmd.AddCommand(newPluginUpdateCmd())
	cmd.AddCommand(newPluginSearchCmd())
	cmd.AddCommand(newPluginInfoCmd())

	return cmd
}

func newPluginListCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginList(format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

func newPluginInstallCmd() *cobra.Command {
	var source string
	var version string

	cmd := &cobra.Command{
		Use:   "install [plugin-name]",
		Short: "Install a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginInstall(args[0], source, version)
		},
	}

	cmd.Flags().StringVarP(&source, "source", "s", "", "plugin source (github, local, url)")
	cmd.Flags().StringVarP(&version, "version", "v", "latest", "plugin version")

	return cmd
}

func newPluginRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [plugin-name]",
		Short: "Remove a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginRemove(args[0])
		},
	}

	return cmd
}

func newPluginUpdateCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "update [plugin-name]",
		Short: "Update plugin(s)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := ""
			if len(args) > 0 {
				pluginName = args[0]
			}
			return runPluginUpdate(pluginName, all)
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "update all plugins")

	return cmd
}

func newPluginSearchCmd() *cobra.Command {
	var category string
	var format string

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for plugins",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := ""
			if len(args) > 0 {
				query = args[0]
			}
			return runPluginSearch(query, category, format)
		},
	}

	cmd.Flags().StringVarP(&category, "category", "c", "", "plugin category")
	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

func newPluginInfoCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "info [plugin-name]",
		Short: "Show plugin information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginInfo(args[0], format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

// Implementation functions
func runPluginList(format string) error {
	pluginMgr, err := plugins.New()
	if err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	pluginList, err := pluginMgr.List()
	if err != nil {
		return fmt.Errorf("failed to list plugins: %w", err)
	}

	if len(pluginList) == 0 {
		fmt.Println("No plugins installed.")
		return nil
	}

	return utils.DisplayResponse(pluginList, format)
}

func runPluginInstall(name, source, version string) error {
	pluginMgr, err := plugins.New()
	if err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	options := plugins.InstallOptions{
		Name:    name,
		Source:  source,
		Version: version,
	}

	spinner := utils.NewSpinner(fmt.Sprintf("Installing plugin %s...", name))
	spinner.Start()

	result, err := pluginMgr.Install(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to install plugin: %w", err)
	}

	fmt.Printf("✅ Plugin '%s' installed successfully\n", name)
	if result.RestartRequired {
		fmt.Println("⚠️  Please restart AlloraCLI to activate the plugin")
	}

	return nil
}

func runPluginRemove(name string) error {
	pluginMgr, err := plugins.New()
	if err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	// Confirm removal
	if !utils.ConfirmAction(fmt.Sprintf("Are you sure you want to remove plugin '%s'?", name)) {
		fmt.Println("Plugin removal cancelled.")
		return nil
	}

	spinner := utils.NewSpinner(fmt.Sprintf("Removing plugin %s...", name))
	spinner.Start()

	err = pluginMgr.Remove(name)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to remove plugin: %w", err)
	}

	fmt.Printf("✅ Plugin '%s' removed successfully\n", name)
	return nil
}

func runPluginUpdate(name string, all bool) error {
	pluginMgr, err := plugins.New()
	if err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	if all {
		spinner := utils.NewSpinner("Updating all plugins...")
		spinner.Start()

		results, err := pluginMgr.UpdateAll()
		spinner.Stop()

		if err != nil {
			return fmt.Errorf("failed to update plugins: %w", err)
		}

		fmt.Println("📦 Plugin update results:")
		for _, result := range results {
			if result.Updated {
				fmt.Printf("  ✅ %s updated to %s\n", result.Name, result.Version)
			} else {
				fmt.Printf("  ℹ️  %s already up to date\n", result.Name)
			}
		}
	} else {
		if name == "" {
			return fmt.Errorf("plugin name is required")
		}

		spinner := utils.NewSpinner(fmt.Sprintf("Updating plugin %s...", name))
		spinner.Start()

		result, err := pluginMgr.Update(name)
		spinner.Stop()

		if err != nil {
			return fmt.Errorf("failed to update plugin: %w", err)
		}

		if result.Updated {
			fmt.Printf("✅ Plugin '%s' updated to %s\n", name, result.Version)
		} else {
			fmt.Printf("ℹ️  Plugin '%s' is already up to date\n", name)
		}
	}

	return nil
}

func runPluginSearch(query, category, format string) error {
	pluginMgr, err := plugins.New()
	if err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	options := plugins.SearchOptions{
		Query:    query,
		Category: category,
	}

	spinner := utils.NewSpinner("Searching for plugins...")
	spinner.Start()

	results, err := pluginMgr.Search(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to search plugins: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No plugins found matching your criteria.")
		return nil
	}

	return utils.DisplayResponse(results, format)
}

func runPluginInfo(name, format string) error {
	pluginMgr, err := plugins.New()
	if err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	info, err := pluginMgr.GetInfo(name)
	if err != nil {
		return fmt.Errorf("failed to get plugin info: %w", err)
	}

	return utils.DisplayResponse(info, format)
}
