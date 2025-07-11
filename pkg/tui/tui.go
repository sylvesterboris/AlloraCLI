package tui

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TUIManager manages terminal user interface components
type TUIManager struct {
	app         *tview.Application
	pages       *tview.Pages
	components  map[string]tview.Primitive
	currentPage string
}

// NewTUIManager creates a new TUI manager
func NewTUIManager() *TUIManager {
	app := tview.NewApplication()
	pages := tview.NewPages()
	
	return &TUIManager{
		app:        app,
		pages:      pages,
		components: make(map[string]tview.Primitive),
	}
}

// AddComponent adds a new component to the TUI
func (t *TUIManager) AddComponent(name string, component tview.Primitive) {
	t.components[name] = component
	t.pages.AddPage(name, component, true, false)
}

// ShowComponent shows a specific component
func (t *TUIManager) ShowComponent(name string) error {
	if _, exists := t.components[name]; !exists {
		return fmt.Errorf("component not found: %s", name)
	}
	
	t.pages.SwitchToPage(name)
	t.currentPage = name
	return nil
}

// Run starts the TUI application
func (t *TUIManager) Run() error {
	t.app.SetRoot(t.pages, true)
	return t.app.Run()
}

// Stop stops the TUI application
func (t *TUIManager) Stop() {
	t.app.Stop()
}

// CreateDashboard creates a monitoring dashboard
func (t *TUIManager) CreateDashboard() *tview.Grid {
	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(0, 0).
		SetBorders(true)

	// Header
	header := tview.NewTextView().
		SetText("AlloraCLI - AI-Powered IT Infrastructure Management").
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	header.SetBorder(true).SetTitle("AlloraCLI Dashboard")

	// Main content area
	mainContent := tview.NewFlex().SetDirection(tview.FlexRow)
	
	// Status panel
	statusPanel := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	statusPanel.SetBorder(true).SetTitle("System Status")
	
	// Metrics panel
	metricsPanel := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	metricsPanel.SetBorder(true).SetTitle("Metrics")
	
	// Logs panel
	logsPanel := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)
	logsPanel.SetBorder(true).SetTitle("Logs")
	
	// Add panels to main content
	mainContent.AddItem(statusPanel, 0, 1, false)
	mainContent.AddItem(metricsPanel, 0, 1, false)
	mainContent.AddItem(logsPanel, 0, 1, false)

	// Footer
	footer := tview.NewTextView().
		SetText("Press 'q' to quit, 'Tab' to switch panels").
		SetTextAlign(tview.AlignCenter)
	footer.SetBorder(true).SetTitle("Controls")

	// Add components to grid
	grid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(mainContent, 1, 0, 1, 2, 0, 0, true)
	grid.AddItem(footer, 2, 0, 1, 2, 0, 0, false)

	// Update status periodically
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			t.app.QueueUpdateDraw(func() {
				statusPanel.SetText(fmt.Sprintf(
					"[green]System Status: [white]Healthy\n" +
					"[yellow]CPU Usage: [white]25%%\n" +
					"[blue]Memory Usage: [white]45%%\n" +
					"[cyan]Disk Usage: [white]60%%\n" +
					"[magenta]Network: [white]Active\n" +
					"[red]Alerts: [white]0\n" +
					"[green]Uptime: [white]%s",
					time.Since(time.Now().Add(-24*time.Hour)).String(),
				))
				
				metricsPanel.SetText(fmt.Sprintf(
					"[green]Requests/sec: [white]120\n" +
					"[yellow]Response Time: [white]50ms\n" +
					"[blue]Error Rate: [white]0.1%%\n" +
					"[cyan]Throughput: [white]2.5MB/s\n" +
					"[magenta]Connections: [white]450\n" +
					"[red]Queue Length: [white]15\n" +
					"[green]Last Update: [white]%s",
					time.Now().Format("15:04:05"),
				))
			})
		}
	}()

	return grid
}

// CreateLogViewer creates a log viewer component
func (t *TUIManager) CreateLogViewer() *tview.TextView {
	logView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			t.app.Draw()
		})
	
	logView.SetBorder(true).SetTitle("Log Viewer")
	
	// Simulate log messages
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		
		messages := []string{
			"[green]INFO[white] System initialization complete",
			"[yellow]WARN[white] High CPU usage detected",
			"[blue]DEBUG[white] Processing request #1234",
			"[cyan]INFO[white] Cache cleared successfully",
			"[magenta]INFO[white] Service health check passed",
			"[red]ERROR[white] Failed to connect to database",
			"[green]INFO[white] Connection restored",
		}
		
		i := 0
		for range ticker.C {
			t.app.QueueUpdateDraw(func() {
				timestamp := time.Now().Format("2006-01-02 15:04:05")
				message := messages[i%len(messages)]
				fmt.Fprintf(logView, "[gray]%s[white] %s\n", timestamp, message)
				i++
			})
		}
	}()
	
	return logView
}

// CreateInteractiveMenu creates an interactive menu
func (t *TUIManager) CreateInteractiveMenu() *tview.List {
	menu := tview.NewList().
		AddItem("Dashboard", "View system dashboard", 'd', func() {
			t.ShowComponent("dashboard")
		}).
		AddItem("Logs", "View system logs", 'l', func() {
			t.ShowComponent("logs")
		}).
		AddItem("Monitoring", "View monitoring metrics", 'm', func() {
			t.ShowComponent("monitoring")
		}).
		AddItem("Cloud Resources", "Manage cloud resources", 'c', func() {
			t.ShowComponent("cloud")
		}).
		AddItem("Security", "Security analysis", 's', func() {
			t.ShowComponent("security")
		}).
		AddItem("Settings", "Application settings", 'S', func() {
			t.ShowComponent("settings")
		}).
		AddItem("Exit", "Exit application", 'q', func() {
			t.Stop()
		})
	
	menu.SetBorder(true).SetTitle("AlloraCLI Menu")
	return menu
}

// CreateProgressBar creates a progress bar component
func (t *TUIManager) CreateProgressBar(title string) *tview.TextView {
	progressView := tview.NewTextView().
		SetDynamicColors(true)
	progressView.SetBorder(true).SetTitle(title)
	
	return progressView
}

// UpdateProgressBar updates a progress bar
func (t *TUIManager) UpdateProgressBar(progressView *tview.TextView, current, total int, message string) {
	percent := float64(current) / float64(total) * 100
	barLength := 50
	filled := int(percent / 100 * float64(barLength))
	
	bar := ""
	for i := 0; i < barLength; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	
	t.app.QueueUpdateDraw(func() {
		progressView.SetText(fmt.Sprintf(
			"[cyan]%s\n\n[white]Progress: [green]%s[white] %.1f%%\n\n[yellow]%s",
			message, bar, percent, fmt.Sprintf("%d/%d", current, total),
		))
	})
}

// SetupKeyBindings sets up global key bindings
func (t *TUIManager) SetupKeyBindings() {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			t.Stop()
			return nil
		case tcell.KeyTab:
			// Switch between components
			switch t.currentPage {
			case "menu":
				t.ShowComponent("dashboard")
			case "dashboard":
				t.ShowComponent("logs")
			case "logs":
				t.ShowComponent("menu")
			default:
				t.ShowComponent("menu")
			}
			return nil
		}
		
		switch event.Rune() {
		case 'q':
			t.Stop()
			return nil
		}
		
		return event
	})
}

// ShowInteractiveDemo shows an interactive demo
func (t *TUIManager) ShowInteractiveDemo() error {
	// Create components
	menu := t.CreateInteractiveMenu()
	dashboard := t.CreateDashboard()
	logs := t.CreateLogViewer()
	
	// Add components
	t.AddComponent("menu", menu)
	t.AddComponent("dashboard", dashboard)
	t.AddComponent("logs", logs)
	
	// Setup key bindings
	t.SetupKeyBindings()
	
	// Start with menu
	t.ShowComponent("menu")
	
	// Run the application
	if err := t.Run(); err != nil {
		log.Fatal(err)
	}
	
	return nil
}
