package streaming

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// StreamingResponse represents a streaming response
type StreamingResponse struct {
	ID        string                 `json:"id"`
	Event     string                 `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// StreamingClient handles streaming responses
type StreamingClient struct {
	client *http.Client
}

// NewStreamingClient creates a new streaming client
func NewStreamingClient() *StreamingClient {
	return &StreamingClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// StreamRequest makes a streaming request and returns a channel of responses
func (c *StreamingClient) StreamRequest(ctx context.Context, url string, headers map[string]string) (<-chan *StreamingResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set streaming headers
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create response channel
	responseChan := make(chan *StreamingResponse, 100)

	// Start reading stream
	go func() {
		defer resp.Body.Close()
		defer close(responseChan)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			
			// Parse Server-Sent Events format
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				
				// Try to parse as JSON
				var responseData map[string]interface{}
				if err := json.Unmarshal([]byte(data), &responseData); err == nil {
					response := &StreamingResponse{
						ID:        fmt.Sprintf("stream_%d", time.Now().UnixNano()),
						Event:     "data",
						Data:      responseData,
						Timestamp: time.Now(),
					}
					
					select {
					case responseChan <- response:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return responseChan, nil
}

// StreamWriter handles writing streaming responses
type StreamWriter struct {
	writer io.Writer
}

// NewStreamWriter creates a new stream writer
func NewStreamWriter(writer io.Writer) *StreamWriter {
	return &StreamWriter{writer: writer}
}

// WriteEvent writes a streaming event
func (w *StreamWriter) WriteEvent(event string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	_, err = fmt.Fprintf(w.writer, "event: %s\ndata: %s\n\n", event, string(jsonData))
	return err
}

// WriteData writes streaming data
func (w *StreamWriter) WriteData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	_, err = fmt.Fprintf(w.writer, "data: %s\n\n", string(jsonData))
	return err
}

// WriteMessage writes a plain text message
func (w *StreamWriter) WriteMessage(message string) error {
	_, err := fmt.Fprintf(w.writer, "data: %s\n\n", message)
	return err
}

// Flush flushes the writer if it supports flushing
func (w *StreamWriter) Flush() error {
	if flusher, ok := w.writer.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

// StreamingProgressTracker tracks progress and streams updates
type StreamingProgressTracker struct {
	writer     *StreamWriter
	total      int
	current    int
	lastUpdate time.Time
	interval   time.Duration
}

// NewStreamingProgressTracker creates a new streaming progress tracker
func NewStreamingProgressTracker(writer *StreamWriter, total int) *StreamingProgressTracker {
	return &StreamingProgressTracker{
		writer:     writer,
		total:      total,
		current:    0,
		lastUpdate: time.Now(),
		interval:   100 * time.Millisecond,
	}
}

// Update updates the progress and potentially streams an update
func (t *StreamingProgressTracker) Update(current int, message string) error {
	t.current = current
	
	// Check if we should send an update
	if time.Since(t.lastUpdate) >= t.interval || current == t.total {
		progress := &ProgressUpdate{
			Current:    current,
			Total:      t.total,
			Percentage: float64(current) / float64(t.total) * 100,
			Message:    message,
			Timestamp:  time.Now(),
		}
		
		err := t.writer.WriteEvent("progress", progress)
		if err != nil {
			return err
		}
		
		t.writer.Flush()
		t.lastUpdate = time.Now()
	}
	
	return nil
}

// ProgressUpdate represents a progress update
type ProgressUpdate struct {
	Current    int       `json:"current"`
	Total      int       `json:"total"`
	Percentage float64   `json:"percentage"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
}

// StreamingLogWriter streams log messages
type StreamingLogWriter struct {
	writer *StreamWriter
}

// NewStreamingLogWriter creates a new streaming log writer
func NewStreamingLogWriter(writer *StreamWriter) *StreamingLogWriter {
	return &StreamingLogWriter{writer: writer}
}

// WriteLog writes a log message
func (w *StreamingLogWriter) WriteLog(level, message string) error {
	logEntry := &LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
	}
	
	err := w.writer.WriteEvent("log", logEntry)
	if err != nil {
		return err
	}
	
	return w.writer.Flush()
}

// LogEntry represents a log entry
type LogEntry struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// StreamingMetricsCollector streams metrics data
type StreamingMetricsCollector struct {
	writer   *StreamWriter
	interval time.Duration
	stop     chan struct{}
}

// NewStreamingMetricsCollector creates a new streaming metrics collector
func NewStreamingMetricsCollector(writer *StreamWriter, interval time.Duration) *StreamingMetricsCollector {
	return &StreamingMetricsCollector{
		writer:   writer,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start starts streaming metrics
func (c *StreamingMetricsCollector) Start(ctx context.Context) error {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.stop:
			return nil
		case <-ticker.C:
			metrics := c.collectMetrics()
			err := c.writer.WriteEvent("metrics", metrics)
			if err != nil {
				return fmt.Errorf("failed to write metrics: %w", err)
			}
			c.writer.Flush()
		}
	}
}

// Stop stops streaming metrics
func (c *StreamingMetricsCollector) Stop() {
	close(c.stop)
}

// collectMetrics collects current metrics
func (c *StreamingMetricsCollector) collectMetrics() map[string]interface{} {
	// This would integrate with actual metrics collection
	// For now, return mock data
	return map[string]interface{}{
		"cpu_usage":    25.5,
		"memory_usage": 45.2,
		"disk_usage":   60.1,
		"network_io":   1024.0,
		"timestamp":    time.Now(),
	}
}

// StreamingCommandExecutor executes commands and streams output
type StreamingCommandExecutor struct {
	writer *StreamWriter
}

// NewStreamingCommandExecutor creates a new streaming command executor
func NewStreamingCommandExecutor(writer *StreamWriter) *StreamingCommandExecutor {
	return &StreamingCommandExecutor{writer: writer}
}

// ExecuteCommand executes a command and streams its output
func (e *StreamingCommandExecutor) ExecuteCommand(ctx context.Context, command string, args []string) error {
	// Start command execution
	err := e.writer.WriteEvent("command_start", map[string]interface{}{
		"command": command,
		"args":    args,
	})
	if err != nil {
		return err
	}
	
	// Stream command output (this would integrate with actual command execution)
	// For now, simulate streaming output
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
				output := map[string]interface{}{
					"line":   fmt.Sprintf("Command output line %d", i+1),
					"stream": "stdout",
				}
				e.writer.WriteEvent("command_output", output)
				e.writer.Flush()
			}
		}
		
		// Command completed
		e.writer.WriteEvent("command_complete", map[string]interface{}{
			"exit_code": 0,
			"duration":  "5.2s",
		})
		e.writer.Flush()
	}()
	
	return nil
}

// StreamingHTTPHandler creates HTTP handlers for streaming responses
type StreamingHTTPHandler struct {
	streamWriter *StreamWriter
}

// NewStreamingHTTPHandler creates a new streaming HTTP handler
func NewStreamingHTTPHandler() *StreamingHTTPHandler {
	return &StreamingHTTPHandler{}
}

// ServeHTTP handles HTTP requests with streaming responses
func (h *StreamingHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set headers for Server-Sent Events
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
	
	// Create stream writer
	streamWriter := NewStreamWriter(w)
	
	// Handle different endpoints
	switch r.URL.Path {
	case "/stream/logs":
		h.streamLogs(r.Context(), streamWriter)
	case "/stream/metrics":
		h.streamMetrics(r.Context(), streamWriter)
	case "/stream/progress":
		h.streamProgress(r.Context(), streamWriter)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// streamLogs streams log messages
func (h *StreamingHTTPHandler) streamLogs(ctx context.Context, writer *StreamWriter) {
	logWriter := NewStreamingLogWriter(writer)
	
	// Stream sample log messages
	logs := []struct {
		level   string
		message string
	}{
		{"INFO", "Application started"},
		{"DEBUG", "Processing request"},
		{"WARN", "High CPU usage detected"},
		{"ERROR", "Database connection failed"},
		{"INFO", "Connection restored"},
	}
	
	for _, log := range logs {
		select {
		case <-ctx.Done():
			return
		default:
			logWriter.WriteLog(log.level, log.message)
			time.Sleep(1 * time.Second)
		}
	}
}

// streamMetrics streams metrics data
func (h *StreamingHTTPHandler) streamMetrics(ctx context.Context, writer *StreamWriter) {
	collector := NewStreamingMetricsCollector(writer, 2*time.Second)
	collector.Start(ctx)
}

// streamProgress streams progress updates
func (h *StreamingHTTPHandler) streamProgress(ctx context.Context, writer *StreamWriter) {
	tracker := NewStreamingProgressTracker(writer, 100)
	
	for i := 0; i <= 100; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			tracker.Update(i, fmt.Sprintf("Processing item %d", i))
			time.Sleep(100 * time.Millisecond)
		}
	}
}
