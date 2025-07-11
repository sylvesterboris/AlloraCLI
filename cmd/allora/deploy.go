package main

import (
	"fmt"

	"github.com/AlloraAi/AlloraCLI/pkg/deploy"
	"github.com/AlloraAi/AlloraCLI/pkg/utils"
	"github.com/spf13/cobra"
)

func newDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "AI-powered deployment automation",
		Long:  `Deploy and manage infrastructure with AI-powered automation and optimization.`,
	}

	cmd.AddCommand(newDeployInfraCmd())
	cmd.AddCommand(newDeployAppCmd())
	cmd.AddCommand(newDeployStatusCmd())
	cmd.AddCommand(newDeployRollbackCmd())
	cmd.AddCommand(newDeployPlanCmd())

	return cmd
}

func newDeployInfraCmd() *cobra.Command {
	var template string
	var optimize bool
	var dryRun bool
	var vars []string

	cmd := &cobra.Command{
		Use:   "infra",
		Short: "Deploy infrastructure",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployInfra(template, optimize, dryRun, vars)
		},
	}

	cmd.Flags().StringVarP(&template, "template", "t", "", "infrastructure template (terraform, cloudformation, pulumi)")
	cmd.Flags().BoolVarP(&optimize, "optimize", "o", false, "enable AI optimization")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "show what would be deployed without making changes")
	cmd.Flags().StringSliceVarP(&vars, "var", "v", []string{}, "template variables (key=value)")

	return cmd
}

func newDeployAppCmd() *cobra.Command {
	var image string
	var environment string
	var replicas int
	var strategy string

	cmd := &cobra.Command{
		Use:   "app",
		Short: "Deploy application",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployApp(image, environment, replicas, strategy)
		},
	}

	cmd.Flags().StringVarP(&image, "image", "i", "", "container image")
	cmd.Flags().StringVarP(&environment, "env", "e", "production", "deployment environment")
	cmd.Flags().IntVarP(&replicas, "replicas", "r", 1, "number of replicas")
	cmd.Flags().StringVarP(&strategy, "strategy", "s", "rolling", "deployment strategy (rolling, blue-green, canary)")

	return cmd
}

func newDeployStatusCmd() *cobra.Command {
	var deploymentID string
	var format string

	cmd := &cobra.Command{
		Use:   "status [deployment-id]",
		Short: "Check deployment status",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				deploymentID = args[0]
			}
			return runDeployStatus(deploymentID, format)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "table", "output format (table, json, yaml)")

	return cmd
}

func newDeployRollbackCmd() *cobra.Command {
	var deploymentID string
	var version string
	var confirm bool

	cmd := &cobra.Command{
		Use:   "rollback [deployment-id]",
		Short: "Rollback deployment",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				deploymentID = args[0]
			}
			return runDeployRollback(deploymentID, version, confirm)
		},
	}

	cmd.Flags().StringVarP(&version, "version", "v", "", "specific version to rollback to")
	cmd.Flags().BoolVarP(&confirm, "confirm", "y", false, "skip confirmation prompts")

	return cmd
}

func newDeployPlanCmd() *cobra.Command {
	var template string
	var optimize bool
	var format string

	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Generate deployment plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeployPlan(template, optimize, format)
		},
	}

	cmd.Flags().StringVarP(&template, "template", "t", "", "infrastructure template")
	cmd.Flags().BoolVarP(&optimize, "optimize", "o", false, "enable AI optimization")
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format (text, json, yaml)")

	return cmd
}

// Implementation functions
func runDeployInfra(template string, optimize, dryRun bool, vars []string) error {
	deployer, err := deploy.New()
	if err != nil {
		return fmt.Errorf("failed to initialize deployer: %w", err)
	}

	options := deploy.InfraOptions{
		Template: template,
		Optimize: optimize,
		DryRun:   dryRun,
		Variables: parseVariables(vars),
	}

	spinner := utils.NewSpinner("Preparing infrastructure deployment...")
	spinner.Start()

	result, err := deployer.DeployInfrastructure(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to deploy infrastructure: %w", err)
	}

	if dryRun {
		fmt.Println("🔍 Deployment plan:")
	} else {
		fmt.Println("🚀 Deployment result:")
	}

	return utils.DisplayResponse(result, "text")
}

func runDeployApp(image, environment string, replicas int, strategy string) error {
	deployer, err := deploy.New()
	if err != nil {
		return fmt.Errorf("failed to initialize deployer: %w", err)
	}

	options := deploy.AppOptions{
		Image:       image,
		Environment: environment,
		Replicas:    replicas,
		Strategy:    strategy,
	}

	spinner := utils.NewSpinner("Deploying application...")
	spinner.Start()

	result, err := deployer.DeployApplication(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to deploy application: %w", err)
	}

	fmt.Println("🚀 Application deployment result:")
	return utils.DisplayResponse(result, "text")
}

func runDeployStatus(deploymentID, format string) error {
	deployer, err := deploy.New()
	if err != nil {
		return fmt.Errorf("failed to initialize deployer: %w", err)
	}

	if deploymentID == "" {
		// List all deployments
		deployments, err := deployer.ListDeployments()
		if err != nil {
			return fmt.Errorf("failed to list deployments: %w", err)
		}
		return utils.DisplayResponse(deployments, format)
	}

	// Get specific deployment status
	status, err := deployer.GetDeploymentStatus(deploymentID)
	if err != nil {
		return fmt.Errorf("failed to get deployment status: %w", err)
	}

	return utils.DisplayResponse(status, format)
}

func runDeployRollback(deploymentID, version string, confirm bool) error {
	deployer, err := deploy.New()
	if err != nil {
		return fmt.Errorf("failed to initialize deployer: %w", err)
	}

	if deploymentID == "" {
		return fmt.Errorf("deployment ID is required")
	}

	// Get confirmation if not auto-confirmed
	if !confirm {
		if !utils.ConfirmAction(fmt.Sprintf("Are you sure you want to rollback deployment %s?", deploymentID)) {
			fmt.Println("Rollback cancelled.")
			return nil
		}
	}

	spinner := utils.NewSpinner("Rolling back deployment...")
	spinner.Start()

	result, err := deployer.RollbackDeployment(deploymentID, version)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to rollback deployment: %w", err)
	}

	fmt.Println("🔄 Rollback result:")
	return utils.DisplayResponse(result, "text")
}

func runDeployPlan(template string, optimize bool, format string) error {
	deployer, err := deploy.New()
	if err != nil {
		return fmt.Errorf("failed to initialize deployer: %w", err)
	}

	options := deploy.PlanOptions{
		Template: template,
		Optimize: optimize,
	}

	spinner := utils.NewSpinner("Generating deployment plan...")
	spinner.Start()

	plan, err := deployer.GeneratePlan(options)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to generate plan: %w", err)
	}

	return utils.DisplayResponse(plan, format)
}

func parseVariables(vars []string) map[string]string {
	variables := make(map[string]string)
	for _, v := range vars {
		if kv := utils.ParseKeyValue(v); len(kv) == 2 {
			variables[kv[0]] = kv[1]
		}
	}
	return variables
}
