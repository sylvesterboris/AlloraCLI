package ui

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlloraAi/AlloraCLI/pkg/agents"
	"github.com/fatih/color"
)

// AnimatedLogo represents the animated ASCII art logo
type AnimatedLogo struct {
	frames []string
	colors []color.Attribute
}

// NewAnimatedLogo creates a new animated logo
func NewAnimatedLogo() *AnimatedLogo {
	frames := []string{
		"AlloraAi - AI-Powered Infrastructure",
		"AlloraAi - Your AI Assistant Ready",
	}

	colors := []color.Attribute{
		color.FgBlue,
		color.FgCyan,
		color.FgGreen,
		color.FgMagenta,
		color.FgYellow,
	}

	return &AnimatedLogo{
		frames: frames,
		colors: colors,
	}
}

// Start begins the animation
func (a *AnimatedLogo) Start() {
	for i := 0; i < 3; i++ {
		for frameIndex, frame := range a.frames {
			// Clear screen
			fmt.Print("\033[2J\033[H")
			
			// Set color
			color.Set(a.colors[frameIndex%len(a.colors)])
			
			// Print frame with ASCII art
			fmt.Println("╔═══════════════════════════════════════════════════════════════════════════╗")
			fmt.Println("║                                                                           ║")
			fmt.Printf("║                          %s                        ║\n", frame)
			fmt.Println("║                                                                           ║")
			fmt.Println("╚═══════════════════════════════════════════════════════════════════════════╝")
			
			// Reset color
			color.Unset()
			
			// Wait
			time.Sleep(time.Millisecond * 800)
		}
	}
}

// Message represents a conversation message
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// GeminiInterface represents the Gemini-style interface
type GeminiInterface struct {
	colorEnabled bool
	conversation []Message
	agents       *agents.AgentManager
}

// NewGeminiInterface creates a new Gemini interface
func NewGeminiInterface(colorEnabled bool) *GeminiInterface {
	return &GeminiInterface{
		colorEnabled: colorEnabled,
		conversation: make([]Message, 0),
		agents:       agents.NewAgentManager(),
	}
}

// displayWelcome shows the welcome screen
func (g *GeminiInterface) displayWelcome() {
	// Clear screen
	fmt.Print("\033[2J\033[H")
	
	// Create and start animated logo
	logo := NewAnimatedLogo()
	logo.Start()
	
	// Clear screen again
	fmt.Print("\033[2J\033[H")
	
	// Display welcome message
	if g.colorEnabled {
		color.Set(color.FgCyan, color.Bold)
	}
	
	fmt.Println("╔═══════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                           ║")
	fmt.Println("║                            AlloraAi Assistant                            ║")
	fmt.Println("║                       AI-Powered Infrastructure                          ║")
	fmt.Println("║                                                                           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════════════════╝")
	
	if g.colorEnabled {
		color.Unset()
	}
	
	// Welcome text
	fmt.Println()
	if g.colorEnabled {
		color.Set(color.FgGreen, color.Bold)
	}
	fmt.Println("🎉 Welcome to AlloraAi - Your AI-Powered Infrastructure Assistant!")
	if g.colorEnabled {
		color.Unset()
	}
	
	// Animated typing effect for description
	g.typeText("I'm here to help you manage your cloud infrastructure, deploy applications, monitor systems, and troubleshoot issues using natural language.")
	
	fmt.Println()
	
	// Show quick tips
	g.showQuickTips()
}

// typeText creates a typing animation effect
func (g *GeminiInterface) typeText(text string) {
	if g.colorEnabled {
		color.Set(color.FgWhite)
	}
	
	for _, char := range text {
		fmt.Print(string(char))
		// Random delay between 10-50ms for realistic typing
		delay := time.Duration(rand.Intn(40)+10) * time.Millisecond
		time.Sleep(delay)
	}
	
	if g.colorEnabled {
		color.Unset()
	}
	
	fmt.Println()
}

// showQuickTips displays quick tips for using the interface
func (g *GeminiInterface) showQuickTips() {
	if g.colorEnabled {
		color.Set(color.FgYellow, color.Bold)
	}
	
	fmt.Println("💡 Quick Tips:")
	
	if g.colorEnabled {
		color.Unset()
		color.Set(color.FgWhite)
	}
	
	tips := []string{
		"• Ask me anything about your infrastructure: \"Monitor my AWS EC2 instances\"",
		"• Get help with deployments: \"Deploy my app to Kubernetes\"",
		"• Troubleshoot issues: \"Why is my service responding slowly?\"",
		"• Type /help for available commands",
		"• Type /examples to see sample queries",
	}
	
	for _, tip := range tips {
		fmt.Println(tip)
		time.Sleep(time.Millisecond * 200)
	}
	
	if g.colorEnabled {
		color.Unset()
	}
	
	fmt.Println()
}

// displayPrompt shows the input prompt
func (g *GeminiInterface) displayPrompt() {
	if g.colorEnabled {
		color.Set(color.FgCyan, color.Bold)
	}
	
	fmt.Print("💬 You: ")
	
	if g.colorEnabled {
		color.Unset()
	}
}

// displayThinking shows the thinking animation
func (g *GeminiInterface) displayThinking() {
	if g.colorEnabled {
		color.Set(color.FgMagenta)
	}
	
	thinkingChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	
	for i := 0; i < 20; i++ {
		fmt.Printf("\r🤖 AlloraAi: %s Thinking...", thinkingChars[i%len(thinkingChars)])
		time.Sleep(time.Millisecond * 100)
	}
	
	fmt.Print("\r🤖 AlloraAi: ")
	
	if g.colorEnabled {
		color.Unset()
	}
}

// displayResponse shows the AI response with typing effect
func (g *GeminiInterface) displayResponse(response string) {
	g.displayThinking()
	
	if g.colorEnabled {
		color.Set(color.FgGreen)
	}
	
	// Type the response
	for _, char := range response {
		fmt.Print(string(char))
		// Faster typing for responses
		delay := time.Duration(rand.Intn(20)+5) * time.Millisecond
		time.Sleep(delay)
	}
	
	if g.colorEnabled {
		color.Unset()
	}
	
	fmt.Println()
	fmt.Println()
}

// displayError shows error messages
func (g *GeminiInterface) displayError(errMsg string) {
	if g.colorEnabled {
		color.Set(color.FgRed, color.Bold)
	}
	
	fmt.Printf("❌ Error: %s\n", errMsg)
	
	if g.colorEnabled {
		color.Unset()
	}
}

// handleUserInput processes user input and generates responses
func (g *GeminiInterface) handleUserInput(input string) error {
	// Add user message to conversation
	g.addToConversation("user", input)
	
	// Create context for AI processing
	ctx := context.Background()
	
	// Process the input with AI agents
	response, err := g.agents.ProcessQuery(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to process query: %w", err)
	}
	
	// Display response
	g.displayResponse(response)
	
	// Add AI response to conversation
	g.addToConversation("assistant", response)
	
	return nil
}

// addToConversation adds a message to the conversation history
func (g *GeminiInterface) addToConversation(role, content string) {
	message := Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	}
	g.conversation = append(g.conversation, message)
}

// clearConversation clears the conversation history
func (g *GeminiInterface) clearConversation() {
	g.conversation = make([]Message, 0)
	fmt.Println("🗑️ Conversation history cleared!")
}

// ExportConversation exports the conversation to a file
func (g *GeminiInterface) ExportConversation(filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write conversation as JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	exportData := map[string]interface{}{
		"exported_at": time.Now().Format(time.RFC3339),
		"conversation": g.conversation,
	}
	
	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("failed to encode conversation: %w", err)
	}

	return nil
}

// LoadConversation loads a conversation from a file
func (g *GeminiInterface) LoadConversation(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var exportData map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&exportData); err != nil {
		return fmt.Errorf("failed to decode conversation: %w", err)
	}

	// Clear current conversation
	g.conversation = []Message{}

	// Load conversation messages
	if conv, ok := exportData["conversation"].([]interface{}); ok {
		for _, msgInterface := range conv {
			if msgMap, ok := msgInterface.(map[string]interface{}); ok {
				msg := Message{
					Role:    msgMap["role"].(string),
					Content: msgMap["content"].(string),
				}
				if timestamp, ok := msgMap["timestamp"].(string); ok {
					if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
						msg.Timestamp = t
					}
				}
				g.conversation = append(g.conversation, msg)
			}
		}
	}

	return nil
}

// GetConversationSummary returns a summary of the current conversation
func (g *GeminiInterface) GetConversationSummary() string {
	if len(g.conversation) == 0 {
		return "No conversation history"
	}

	userMessages := 0
	assistantMessages := 0
	totalChars := 0

	for _, msg := range g.conversation {
		if msg.Role == "user" {
			userMessages++
		} else if msg.Role == "assistant" {
			assistantMessages++
		}
		totalChars += len(msg.Content)
	}

	return fmt.Sprintf("Messages: %d user, %d assistant | Characters: %d | Started: %s",
		userMessages, assistantMessages, totalChars, g.conversation[0].Timestamp.Format("15:04:05"))
}

// displayMenu shows the interactive menu
func (g *GeminiInterface) displayMenu() {
	if g.colorEnabled {
		color.Set(color.FgCyan, color.Bold)
	}
	
	fmt.Println("\n╭─────────────────────────────────────────────────────────────────────────────╮")
	fmt.Println("│                            AlloraAi Interactive Menu                       │")
	fmt.Println("├─────────────────────────────────────────────────────────────────────────────┤")
	fmt.Println("│ /help      - Show available commands                                        │")
	fmt.Println("│ /clear     - Clear conversation history                                     │")
	fmt.Println("│ /export    - Export conversation to file                                   │")
	fmt.Println("│ /load      - Load conversation from file                                   │")
	fmt.Println("│ /summary   - Show conversation summary                                     │")
	fmt.Println("│ /examples  - Show example queries                                          │")
	fmt.Println("│ /quit      - Exit the interface                                           │")
	fmt.Println("╰─────────────────────────────────────────────────────────────────────────────╯")
	
	if g.colorEnabled {
		color.Unset()
	}
}

// displayExamples shows example queries
func (g *GeminiInterface) displayExamples() {
	examples := []string{
		"Monitor my AWS EC2 instances and alert if CPU usage exceeds 80%",
		"Deploy a new Kubernetes cluster on Azure with 3 nodes",
		"Analyze my GCP billing and suggest cost optimizations",
		"Set up automated backups for my database servers",
		"Create a monitoring dashboard for my microservices",
		"Troubleshoot network connectivity issues in my infrastructure",
		"Generate a security audit report for my cloud resources",
		"Scale my application based on traffic patterns",
	}

	if g.colorEnabled {
		color.Set(color.FgGreen, color.Bold)
	}
	
	fmt.Println("\n╭─────────────────────────────────────────────────────────────────────────────╮")
	fmt.Println("│                           Example Queries                                  │")
	fmt.Println("├─────────────────────────────────────────────────────────────────────────────┤")
	
	for i, example := range examples {
		fmt.Printf("│ %d. %-71s │\n", i+1, example)
	}
	
	fmt.Println("╰─────────────────────────────────────────────────────────────────────────────╯")
	
	if g.colorEnabled {
		color.Unset()
	}
}

// handleSpecialCommands processes special commands like /help, /clear, etc.
func (g *GeminiInterface) handleSpecialCommands(input string) bool {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "/help":
		g.displayMenu()
		return true
	case "/clear":
		g.clearConversation()
		g.displayWelcome()
		return true
	case "/export":
		filename := fmt.Sprintf("conversation_%s.json", time.Now().Format("20060102_150405"))
		if err := g.ExportConversation(filename); err != nil {
			g.displayError(fmt.Sprintf("Failed to export conversation: %v", err))
		} else {
			fmt.Printf("✅ Conversation exported to: %s\n", filename)
		}
		return true
	case "/load":
		fmt.Print("Enter filename to load: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			filename := strings.TrimSpace(scanner.Text())
			if filename != "" {
				if err := g.LoadConversation(filename); err != nil {
					g.displayError(fmt.Sprintf("Failed to load conversation: %v", err))
				} else {
					fmt.Printf("✅ Conversation loaded from: %s\n", filename)
				}
			}
		}
		return true
	case "/summary":
		summary := g.GetConversationSummary()
		fmt.Printf("📊 %s\n", summary)
		return true
	case "/examples":
		g.displayExamples()
		return true
	case "/quit", "/exit":
		g.displayGoodbye()
		return true
	default:
		return false
	}
}

// displayGoodbye shows the goodbye message
func (g *GeminiInterface) displayGoodbye() {
	if g.colorEnabled {
		color.Set(color.FgMagenta, color.Bold)
	}
	
	fmt.Println("\n╭─────────────────────────────────────────────────────────────────────────────╮")
	fmt.Println("│                          Thank you for using AlloraAi!                     │")
	fmt.Println("│                        Your AI Infrastructure Assistant                     │")
	fmt.Println("╰─────────────────────────────────────────────────────────────────────────────╯")
	
	if g.colorEnabled {
		color.Unset()
	}
}

// Start begins the Gemini interface
func (g *GeminiInterface) Start() error {
	// Display welcome message
	g.displayWelcome()
	
	// Display menu
	g.displayMenu()
	
	// Initialize scanner for user input
	scanner := bufio.NewScanner(os.Stdin)
	
	// Main interaction loop
	for {
		// Display prompt
		g.displayPrompt()
		
		// Read user input
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		
		// Skip empty input
		if input == "" {
			continue
		}
		
		// Handle special commands
		if g.handleSpecialCommands(input) {
			// Check if user wants to quit
			if strings.ToLower(input) == "/quit" || strings.ToLower(input) == "/exit" {
				break
			}
			continue
		}
		
		// Process regular user input
		if err := g.handleUserInput(input); err != nil {
			g.displayError(fmt.Sprintf("Error processing input: %v", err))
		}
	}
	
	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("input error: %w", err)
	}
	
	return nil
}
