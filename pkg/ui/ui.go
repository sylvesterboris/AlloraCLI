package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
)

// UIManager manages user interface components
type UIManager struct {
	colorEnabled bool
	verboseMode  bool
}

// NewUIManager creates a new UI manager
func NewUIManager(colorEnabled, verboseMode bool) *UIManager {
	return &UIManager{
		colorEnabled: colorEnabled,
		verboseMode:  verboseMode,
	}
}

// Colors
var (
	InfoColor    = color.New(color.FgCyan)
	SuccessColor = color.New(color.FgGreen)
	WarningColor = color.New(color.FgYellow)
	ErrorColor   = color.New(color.FgRed)
	HeaderColor  = color.New(color.FgBlue, color.Bold)
)

// PrintHeader prints a formatted header
func (ui *UIManager) PrintHeader(text string) {
	if ui.colorEnabled {
		HeaderColor.Printf("\n=== %s ===\n\n", text)
	} else {
		fmt.Printf("\n=== %s ===\n\n", text)
	}
}

// PrintInfo prints info message
func (ui *UIManager) PrintInfo(format string, args ...interface{}) {
	if ui.colorEnabled {
		InfoColor.Printf("ℹ  "+format+"\n", args...)
	} else {
		fmt.Printf("INFO: "+format+"\n", args...)
	}
}

// PrintSuccess prints success message
func (ui *UIManager) PrintSuccess(format string, args ...interface{}) {
	if ui.colorEnabled {
		SuccessColor.Printf("✓  "+format+"\n", args...)
	} else {
		fmt.Printf("SUCCESS: "+format+"\n", args...)
	}
}

// PrintWarning prints warning message
func (ui *UIManager) PrintWarning(format string, args ...interface{}) {
	if ui.colorEnabled {
		WarningColor.Printf("⚠  "+format+"\n", args...)
	} else {
		fmt.Printf("WARNING: "+format+"\n", args...)
	}
}

// PrintError prints error message
func (ui *UIManager) PrintError(format string, args ...interface{}) {
	if ui.colorEnabled {
		ErrorColor.Printf("✗  "+format+"\n", args...)
	} else {
		fmt.Printf("ERROR: "+format+"\n", args...)
	}
}

// PrintVerbose prints verbose message
func (ui *UIManager) PrintVerbose(format string, args ...interface{}) {
	if ui.verboseMode {
		if ui.colorEnabled {
			color.New(color.FgWhite, color.Faint).Printf("  "+format+"\n", args...)
		} else {
			fmt.Printf("  "+format+"\n", args...)
		}
	}
}

// CreateProgressBar creates a progress bar
func (ui *UIManager) CreateProgressBar(max int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(max,
		progressbar.OptionEnableColorCodes(ui.colorEnabled),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetRenderBlankState(true),
	)
}

// CreateSpinnerProgressBar creates a spinner progress bar for indeterminate progress
func (ui *UIManager) CreateSpinnerProgressBar(description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(-1,
		progressbar.OptionEnableColorCodes(ui.colorEnabled),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetRenderBlankState(true),
	)
}

// InteractivePrompt creates an interactive prompt
func (ui *UIManager) InteractivePrompt(label string, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}
	
	return prompt.Run()
}

// InteractiveSelect creates an interactive selection menu
func (ui *UIManager) InteractiveSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	
	_, result, err := prompt.Run()
	return result, err
}

// InteractiveConfirm creates an interactive confirmation
func (ui *UIManager) InteractiveConfirm(label string, defaultValue bool) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Default:   func() string {
			if defaultValue {
				return "y"
			}
			return "n"
		}(),
	}
	
	result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	
	return strings.ToLower(result) == "y" || strings.ToLower(result) == "yes", nil
}

// InteractivePassword creates an interactive password input
func (ui *UIManager) InteractivePassword(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}
	
	return prompt.Run()
}

// InteractiveMultiSelect creates an interactive multi-selection menu
func (ui *UIManager) InteractiveMultiSelect(label string, items []string) ([]string, error) {
	var selected []string
	
	for {
		remaining := []string{}
		for _, item := range items {
			found := false
			for _, sel := range selected {
				if item == sel {
					found = true
					break
				}
			}
			if !found {
				remaining = append(remaining, item)
			}
		}
		
		if len(remaining) == 0 {
			break
		}
		
		options := append(remaining, "Done")
		
		prompt := promptui.Select{
			Label: fmt.Sprintf("%s (selected: %d)", label, len(selected)),
			Items: options,
		}
		
		_, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		
		if result == "Done" {
			break
		}
		
		selected = append(selected, result)
	}
	
	return selected, nil
}

// DisplayTable displays a formatted table
func (ui *UIManager) DisplayTable(headers []string, rows [][]string) {
	if len(rows) == 0 {
		ui.PrintInfo("No data to display")
		return
	}
	
	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}
	
	// Print header
	ui.printTableRow(headers, colWidths, true)
	ui.printTableSeparator(colWidths)
	
	// Print rows
	for _, row := range rows {
		ui.printTableRow(row, colWidths, false)
	}
	
	fmt.Println()
}

// printTableRow prints a single table row
func (ui *UIManager) printTableRow(row []string, widths []int, isHeader bool) {
	fmt.Print("| ")
	for i, cell := range row {
		if i < len(widths) {
			if isHeader && ui.colorEnabled {
				HeaderColor.Printf("%-*s", widths[i], cell)
			} else {
				fmt.Printf("%-*s", widths[i], cell)
			}
			fmt.Print(" | ")
		}
	}
	fmt.Println()
}

// printTableSeparator prints table separator
func (ui *UIManager) printTableSeparator(widths []int) {
	fmt.Print("|")
	for _, width := range widths {
		fmt.Print(strings.Repeat("-", width+2))
		fmt.Print("|")
	}
	fmt.Println()
}

// DisplayKeyValue displays key-value pairs
func (ui *UIManager) DisplayKeyValue(data map[string]interface{}) {
	for key, value := range data {
		if ui.colorEnabled {
			InfoColor.Printf("  %s: ", key)
			fmt.Printf("%v\n", value)
		} else {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

// DisplayList displays a formatted list
func (ui *UIManager) DisplayList(items []string, title string) {
	if title != "" {
		ui.PrintHeader(title)
	}
	
	for i, item := range items {
		if ui.colorEnabled {
			InfoColor.Printf("  %d. %s\n", i+1, item)
		} else {
			fmt.Printf("  %d. %s\n", i+1, item)
		}
	}
	fmt.Println()
}

// ShowSpinner shows a spinner with a message
func (ui *UIManager) ShowSpinner(message string, duration time.Duration) {
	spinner := ui.CreateSpinnerProgressBar(message)
	
	start := time.Now()
	for time.Since(start) < duration {
		spinner.Add(1)
		time.Sleep(100 * time.Millisecond)
	}
	
	spinner.Finish()
}

// InteractiveWizard creates an interactive wizard
func (ui *UIManager) InteractiveWizard(title string, steps []WizardStep) (map[string]interface{}, error) {
	ui.PrintHeader(title)
	
	results := make(map[string]interface{})
	
	for i, step := range steps {
		ui.PrintInfo("Step %d/%d: %s", i+1, len(steps), step.Title)
		
		if step.Description != "" {
			fmt.Printf("  %s\n\n", step.Description)
		}
		
		var result interface{}
		var err error
		
		switch step.Type {
		case WizardStepTypeInput:
			result, err = ui.InteractivePrompt(step.Prompt, step.Default)
		case WizardStepTypeSelect:
			result, err = ui.InteractiveSelect(step.Prompt, step.Options)
		case WizardStepTypeConfirm:
			result, err = ui.InteractiveConfirm(step.Prompt, step.Default == "true")
		case WizardStepTypePassword:
			result, err = ui.InteractivePassword(step.Prompt)
		case WizardStepTypeMultiSelect:
			result, err = ui.InteractiveMultiSelect(step.Prompt, step.Options)
		}
		
		if err != nil {
			return nil, fmt.Errorf("step %d failed: %w", i+1, err)
		}
		
		results[step.Key] = result
		fmt.Println()
	}
	
	return results, nil
}

// WizardStep represents a step in an interactive wizard
type WizardStep struct {
	Title       string
	Description string
	Type        WizardStepType
	Key         string
	Prompt      string
	Default     string
	Options     []string
	Required    bool
}

// WizardStepType represents the type of wizard step
type WizardStepType string

const (
	WizardStepTypeInput       WizardStepType = "input"
	WizardStepTypeSelect      WizardStepType = "select"
	WizardStepTypeConfirm     WizardStepType = "confirm"
	WizardStepTypePassword    WizardStepType = "password"
	WizardStepTypeMultiSelect WizardStepType = "multiselect"
)

// PrintBanner prints the application banner
func (ui *UIManager) PrintBanner() {
	banner := `
   ___   _ _                   ___ _    ___ 
  / _ \ | | |                 / __| |  |_ _|
 / /_\ \| | | ___  _ __ __ _ | |  | |   | | 
 |  _  || | |/ _ \| '__/ _` || |  | |   | | 
 | | | || | | (_) | | | (_| || |__| |__|_|_|
 \_| |_/|_|_|\___/|_|  \__,_| \____\____(_)
                                           
 AI-Powered IT Infrastructure Management
`
	
	if ui.colorEnabled {
		HeaderColor.Println(banner)
	} else {
		fmt.Println(banner)
	}
}

// ClearScreen clears the terminal screen
func (ui *UIManager) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// MoveCursor moves the cursor to a specific position
func (ui *UIManager) MoveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

// SaveCursor saves the current cursor position
func (ui *UIManager) SaveCursor() {
	fmt.Print("\033[s")
}

// RestoreCursor restores the saved cursor position
func (ui *UIManager) RestoreCursor() {
	fmt.Print("\033[u")
}

// GetTerminalSize returns the terminal size
func (ui *UIManager) GetTerminalSize() (int, int) {
	// This would need platform-specific implementation
	// For now, return default values
	return 80, 24
}

// PressEnterToContinue waits for user to press Enter
func (ui *UIManager) PressEnterToContinue(message string) {
	if message == "" {
		message = "Press Enter to continue..."
	}
	
	fmt.Print(message)
	fmt.Scanln()
}

// DisplayMenu displays a menu and returns the selected option
func (ui *UIManager) DisplayMenu(title string, options []string) (int, error) {
	ui.PrintHeader(title)
	
	for i, option := range options {
		fmt.Printf("  %d. %s\n", i+1, option)
	}
	
	var choice int
	fmt.Print("\nEnter your choice: ")
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return 0, err
	}
	
	if choice < 1 || choice > len(options) {
		return 0, fmt.Errorf("invalid choice: %d", choice)
	}
	
	return choice - 1, nil
}

// IsTerminalInteractive checks if the terminal is interactive
func (ui *UIManager) IsTerminalInteractive() bool {
	// Check if stdin is a terminal
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	
	return (stat.Mode() & os.ModeCharDevice) != 0
}
