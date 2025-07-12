package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// PrintBanner prints the AlloraCLI banner
func PrintBanner() {
	banner := `
    ▄▄▄       ██▓     ██▓     ▒█████   ██▀███   ▄▄▄      
   ▒████▄    ▓██▒    ▓██▒    ▒██▒  ██▒▓██ ▒ ██▒▒████▄    
   ▒██  ▀█▄  ▒██░    ▒██░    ▒██░  ██▒▓██ ░▄█ ▒▒██  ▀█▄  
   ░██▄▄▄▄██ ▒██░    ▒██░    ▒██   ██░▒██▀▀█▄  ░██▄▄▄▄██ 
   ░▓█   ▓██▒░██████▒░██████▒░ ████▓▒░░██▓ ▒██▒ ▓█   ▓██▒
   ░▒▒   ▓▒█░░ ▒░▓  ░░ ▒░▓  ░░ ▒░▒░▒░ ░ ▒▓ ░▒▓░ ▒▒   ▓▒█░
    ▒   ▒▒ ░░ ░ ▒  ░░ ░ ▒  ░  ░ ▒ ▒░   ░▒ ░ ▒░  ▒   ▒▒ ░
    ░   ▒     ░ ░     ░ ░   ░ ░ ░ ▒    ░░   ░   ░   ▒   
        ░  ░    ░  ░    ░  ░    ░ ░     ░           ░  ░
                                                         
               AI-Powered IT Infrastructure CLI
`
	color.Cyan(banner)
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	return s
}

// DisplayResponse displays a response in the specified format
func DisplayResponse(data interface{}, format string) error {
	switch format {
	case "json":
		return displayJSON(data)
	case "yaml":
		return displayYAML(data)
	case "table":
		return displayTable(data)
	case "text":
		return displayText(data)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// displayJSON displays data in JSON format
func displayJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// displayYAML displays data in YAML format
func displayYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(data)
}

// displayTable displays data in table format
func displayTable(data interface{}) error {
	// This is a simplified table display
	// In a real implementation, you'd need to handle different data types
	fmt.Printf("%+v\n", data)
	return nil
}

// displayText displays data in human-readable text format
func displayText(data interface{}) error {
	fmt.Printf("%+v\n", data)
	return nil
}

// InitializeLogging initializes the logging system
func InitializeLogging(verbose bool) error {
	// Set log level
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})

	return nil
}

// GetUserInput gets user input with a prompt
func GetUserInput(prompt string) (string, error) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read input")
	}
	return strings.TrimSpace(scanner.Text()), scanner.Err()
}

// ConfirmAction asks for user confirmation
func ConfirmAction(message string) bool {
	fmt.Printf("%s (y/N): ", message)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return response == "y" || response == "yes"
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

// CheckForExit checks if the user wants to exit (simplified)
func CheckForExit() bool {
	// This is a simplified implementation
	// In a real implementation, you'd use non-blocking input
	return false
}

// JoinArgs joins command line arguments into a single string
func JoinArgs(args []string) string {
	return strings.Join(args, " ")
}

// ParseKeyValue parses a key=value string
func ParseKeyValue(input string) []string {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		return []string{}
	}
	return parts
}

// FormatBytes formats bytes into human-readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration formats duration into human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

// CreateTable creates a formatted table
func CreateTable(headers []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.AppendBulk(rows)
	table.Render()
}

// Colorize adds color to text based on type
func Colorize(text string, colorType string) string {
	switch colorType {
	case "success":
		return color.GreenString(text)
	case "error":
		return color.RedString(text)
	case "warning":
		return color.YellowString(text)
	case "info":
		return color.BlueString(text)
	case "header":
		return color.CyanString(text)
	default:
		return text
	}
}

// LogInfo logs an info message
func LogInfo(message string) {
	logrus.Info(message)
}

// LogError logs an error message
func LogError(message string) {
	logrus.Error(message)
}

// LogDebug logs a debug message
func LogDebug(message string) {
	logrus.Debug(message)
}

// LogWarning logs a warning message
func LogWarning(message string) {
	logrus.Warning(message)
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveFromSlice removes an item from a string slice
func RemoveFromSlice(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// GetEnvWithDefault gets an environment variable with a default value
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
