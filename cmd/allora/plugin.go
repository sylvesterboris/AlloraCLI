package main

import (
	"context"
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/AlloraAi/AlloraCLI/pkg/plugins"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Plugin management and execution",
		Long:  `Manage and execute plugins to extend AlloraCLI functionality.`,
	}

	cmd.AddCommand(newPluginListCmd())
	cmd.AddCommand(newPluginInstallCmd())
	cmd.AddCommand(newPluginUninstallCmd())
	cmd.AddCommand(newPluginUpdateCmd())
	cmd.AddCommand(newPluginSearchCmd())
	cmd.AddCommand(newPluginRunCmd())

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

	cmd.Flags().StringVarP(&source, "source", "s", "", "plugin source (URL or registry)")
	cmd.Flags().StringVarP(&version, "version", "v", "latest", "plugin version")

	return cmd
}

func newPluginUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall [plugin-name]",
		Short: "Uninstall a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginUninstall(args[0])
		},
	}

	return cmd
}

func newPluginUpdateCmd() *cobra.Command {
	var all bool

	cmd := &cobra.Command{
		Use:   "update [plugin-name]",
		Short: "Update a plugin",
		RunE: func(cmd *cobra.Command, args []string) error {
			if all {
				return runPluginUpdateAll()
			}
			if len(args) == 0 {
				return fmt.Errorf("plugin name required when --all is not specified")
			}
			return runPluginUpdate(args[0])
		},
	}

	cmd.Flags().BoolVarP(&all, "all", "a", false, "update all plugins")

	return cmd
}

func newPluginSearchCmd() *cobra.Command {
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
			return runPluginSearch(query, format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

func newPluginRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [plugin-name] [plugin-args...]",
		Short: "Run a plugin",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginName := args[0]
			pluginArgs := args[1:]
			return runPluginRun(pluginName, pluginArgs)
		},
	}

	return cmd
}

// Implementation functions
func runPluginList(format string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	pluginService, err := plugins.NewPluginService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin service: %w", err)
	}

	ctx := context.Background()
	pluginList, err := pluginService.ListPlugins(ctx)
	if err != nil {
		return fmt.Errorf("failed to list plugins: %w", err)
	}

	return utils.DisplayResponse(pluginList, format)
}

func runPluginInstall(name, source, version string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	pluginService, err := plugins.NewPluginService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin service: %w", err)
	}

	ctx := context.Background()
	
	if source == "" {
		source = fmt.Sprintf("https://registry.alloraai.com/plugins/%s", name)
	}

	spinner := utils.NewSpinner(fmt.Sprintf("Installing plugin %s...", name))
	spinner.Start()

	err = pluginService.InstallPlugin(ctx, name, source)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to install plugin: %w", err)
	}

	fmt.Printf("✅ Plugin %s installed successfully!\n", name)
	return nil
}

func runPluginUninstall(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	pluginService, err := plugins.NewPluginService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin service: %w", err)
	}

	ctx := context.Background()

	spinner := utils.NewSpinner(fmt.Sprintf("Uninstalling plugin %s...", name))
	spinner.Start()

	err = pluginService.UninstallPlugin(ctx, name)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to uninstall plugin: %w", err)
	}

	fmt.Printf("✅ Plugin %s uninstalled successfully!\n", name)
	return nil
}

func runPluginUpdate(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	pluginService, err := plugins.NewPluginService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin service: %w", err)
	}

	ctx := context.Background()

	spinner := utils.NewSpinner(fmt.Sprintf("Updating plugin %s...", name))
	spinner.Start()

	err = pluginService.UpdatePlugin(ctx, name)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to update plugin: %w", err)
	}

	fmt.Printf("✅ Plugin %s updated successfully!\n", name)
	return nil
}

func runPluginUpdateAll() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	pluginService, err := plugins.NewPluginService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin service: %w", err)
	}

	ctx := context.Background()

	// List all plugins first
	pluginList, err := pluginService.ListPlugins(ctx)
	if err != nil {
		return fmt.Errorf("failed to list plugins: %w", err)
	}

	for _, plugin := range pluginList {
		if !plugin.Enabled {
			continue
		}

		spinner := utils.NewSpinner(fmt.Sprintf("Updating plugin %s...", plugin.Name))
		spinner.Start()

		err = pluginService.UpdatePlugin(ctx, plugin.Name)
		spinner.Stop()

		if err != nil {
			utils.LogError(fmt.Sprintf("Failed to update plugin %s: %v", plugin.Name, err))
			continue
		}

		fmt.Printf("✅ Plugin %s updated successfully!\n", plugin.Name)
	}

	return nil
}

func runPluginSearch(query, format string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	pluginService, err := plugins.NewPluginService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin service: %w", err)
	}

	ctx := context.Background()

	spinner := utils.NewSpinner("Searching for plugins...")
	spinner.Start()

	results, err := pluginService.SearchPlugins(ctx, query)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to search plugins: %w", err)
	}

	return utils.DisplayResponse(results, format)
}

func runPluginRun(name string, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	pluginService, err := plugins.NewPluginService(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize plugin service: %w", err)
	}

	ctx := context.Background()

	result, err := pluginService.ExecutePlugin(ctx, name, args)
	if err != nil {
		return fmt.Errorf("failed to execute plugin: %w", err)
	}

	if result.ExitCode != 0 {
		if result.Error != "" {
			fmt.Printf("Plugin error: %s\n", result.Error)
		}
		return fmt.Errorf("plugin exited with code %d", result.ExitCode)
	}

	if result.Output != "" {
		fmt.Print(result.Output)
	}

	return nil
}
