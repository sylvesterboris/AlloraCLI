package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AlloraAi/AlloraCLI/pkg/config"
	"github.com/sirupsen/logrus"
)

// SecurityService interface defines security-related operations
type SecurityService interface {
	ScanVulnerabilities(ctx context.Context, target string) (*ScanResult, error)
	CheckCompliance(ctx context.Context, standard string) (*ComplianceResult, error)
	AuditPermissions(ctx context.Context, resource string) (*AuditResult, error)
	MonitorSecurityEvents(ctx context.Context) (<-chan SecurityEvent, error)
	GenerateSecurityReport(ctx context.Context, options ReportOptions) (*SecurityReport, error)
	ValidateSecurityPolicies(ctx context.Context, policies []Policy) (*ValidationResult, error)
}

// ScanResult represents the result of a security scan
type ScanResult struct {
	ID              string          `json:"id"`
	Target          string          `json:"target"`
	Timestamp       time.Time       `json:"timestamp"`
	Status          string          `json:"status"`
	Summary         ScanSummary     `json:"summary"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Recommendations []string        `json:"recommendations"`
}

// ScanSummary provides a summary of scan results
type ScanSummary struct {
	TotalChecks    int `json:"total_checks"`
	CriticalIssues int `json:"critical_issues"`
	HighIssues     int `json:"high_issues"`
	MediumIssues   int `json:"medium_issues"`
	LowIssues      int `json:"low_issues"`
	InfoIssues     int `json:"info_issues"`
}

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Severity    string            `json:"severity"`
	CVSS        float64           `json:"cvss"`
	CVE         string            `json:"cve,omitempty"`
	Component   string            `json:"component"`
	Version     string            `json:"version"`
	Solution    string            `json:"solution"`
	References  []string          `json:"references"`
	Metadata    map[string]string `json:"metadata"`
}

// ComplianceResult represents compliance check results
type ComplianceResult struct {
	ID        string              `json:"id"`
	Standard  string              `json:"standard"`
	Timestamp time.Time           `json:"timestamp"`
	Status    string              `json:"status"`
	Score     float64             `json:"score"`
	Controls  []ComplianceControl `json:"controls"`
	Summary   ComplianceSummary   `json:"summary"`
}

// ComplianceControl represents a single compliance control
type ComplianceControl struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Evidence    string `json:"evidence"`
	Remediation string `json:"remediation"`
}

// ComplianceSummary provides a summary of compliance results
type ComplianceSummary struct {
	TotalControls   int `json:"total_controls"`
	PassedControls  int `json:"passed_controls"`
	FailedControls  int `json:"failed_controls"`
	WarningControls int `json:"warning_controls"`
}

// AuditResult represents permission audit results
type AuditResult struct {
	ID          string       `json:"id"`
	Resource    string       `json:"resource"`
	Timestamp   time.Time    `json:"timestamp"`
	Permissions []Permission `json:"permissions"`
	Issues      []AuditIssue `json:"issues"`
	Summary     AuditSummary `json:"summary"`
}

// Permission represents a permission entry
type Permission struct {
	Principal  string            `json:"principal"`
	Actions    []string          `json:"actions"`
	Resource   string            `json:"resource"`
	Effect     string            `json:"effect"`
	Conditions map[string]string `json:"conditions"`
}

// AuditIssue represents an audit issue
type AuditIssue struct {
	Type           string `json:"type"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	Resource       string `json:"resource"`
	Principal      string `json:"principal"`
	Recommendation string `json:"recommendation"`
}

// AuditSummary provides a summary of audit results
type AuditSummary struct {
	TotalPermissions int `json:"total_permissions"`
	CriticalIssues   int `json:"critical_issues"`
	HighIssues       int `json:"high_issues"`
	MediumIssues     int `json:"medium_issues"`
	LowIssues        int `json:"low_issues"`
}

// SecurityEvent represents a security event
type SecurityEvent struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Timestamp   time.Time         `json:"timestamp"`
	Source      string            `json:"source"`
	Severity    string            `json:"severity"`
	Description string            `json:"description"`
	Details     map[string]string `json:"details"`
	Actions     []string          `json:"actions"`
}

// SecurityReport represents a comprehensive security report
type SecurityReport struct {
	ID                string             `json:"id"`
	Timestamp         time.Time          `json:"timestamp"`
	Type              string             `json:"type"`
	ExecutiveSummary  ExecutiveSummary   `json:"executive_summary"`
	ScanResults       []ScanResult       `json:"scan_results"`
	ComplianceResults []ComplianceResult `json:"compliance_results"`
	AuditResults      []AuditResult      `json:"audit_results"`
	Recommendations   []string           `json:"recommendations"`
}

// ExecutiveSummary provides a high-level summary
type ExecutiveSummary struct {
	OverallRiskScore     float64  `json:"overall_risk_score"`
	CriticalFindings     int      `json:"critical_findings"`
	HighPriorityFindings int      `json:"high_priority_findings"`
	ComplianceScore      float64  `json:"compliance_score"`
	KeyRecommendations   []string `json:"key_recommendations"`
}

// ReportOptions defines options for generating security reports
type ReportOptions struct {
	Type           string   `json:"type"`
	Targets        []string `json:"targets"`
	Standards      []string `json:"standards"`
	IncludeDetails bool     `json:"include_details"`
	Format         string   `json:"format"`
}

// Policy represents a security policy
type Policy struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Rules       []PolicyRule      `json:"rules"`
	Metadata    map[string]string `json:"metadata"`
}

// PolicyRule represents a rule within a policy
type PolicyRule struct {
	ID         string            `json:"id"`
	Condition  string            `json:"condition"`
	Action     string            `json:"action"`
	Parameters map[string]string `json:"parameters"`
	Enabled    bool              `json:"enabled"`
}

// ValidationResult represents policy validation results
type ValidationResult struct {
	ID        string             `json:"id"`
	Timestamp time.Time          `json:"timestamp"`
	Status    string             `json:"status"`
	Policies  []PolicyValidation `json:"policies"`
	Summary   ValidationSummary  `json:"summary"`
}

// PolicyValidation represents validation results for a single policy
type PolicyValidation struct {
	PolicyID string              `json:"policy_id"`
	Status   string              `json:"status"`
	Issues   []ValidationIssue   `json:"issues"`
	Warnings []ValidationWarning `json:"warnings"`
}

// ValidationIssue represents a validation issue
type ValidationIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	RuleID      string `json:"rule_id"`
	Solution    string `json:"solution"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	RuleID      string `json:"rule_id"`
	Suggestion  string `json:"suggestion"`
}

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	TotalPolicies   int `json:"total_policies"`
	ValidPolicies   int `json:"valid_policies"`
	InvalidPolicies int `json:"invalid_policies"`
	TotalIssues     int `json:"total_issues"`
	TotalWarnings   int `json:"total_warnings"`
}

// DefaultSecurityService provides a default implementation
type DefaultSecurityService struct {
	config *config.Config
}

// NewSecurityService creates a new security service
func NewSecurityService(cfg *config.Config) SecurityService {
	return &DefaultSecurityService{
		config: cfg,
	}
}

// ScanVulnerabilities performs a vulnerability scan
func (s *DefaultSecurityService) ScanVulnerabilities(ctx context.Context, target string) (*ScanResult, error) {
	// Mock implementation - in real implementation, this would integrate with security scanners
	return &ScanResult{
		ID:        "scan-001",
		Target:    target,
		Timestamp: time.Now(),
		Status:    "completed",
		Summary: ScanSummary{
			TotalChecks:    100,
			CriticalIssues: 0,
			HighIssues:     2,
			MediumIssues:   5,
			LowIssues:      10,
			InfoIssues:     15,
		},
		Vulnerabilities: []Vulnerability{
			{
				ID:          "vuln-001",
				Title:       "Outdated TLS Configuration",
				Description: "TLS 1.0 and 1.1 are deprecated and should be disabled",
				Severity:    "high",
				CVSS:        7.5,
				Component:   "web-server",
				Version:     "1.0",
				Solution:    "Upgrade to TLS 1.2 or higher",
				References:  []string{"https://example.com/tls-security"},
			},
		},
		Recommendations: []string{
			"Update TLS configuration to use only TLS 1.2 and above",
			"Implement security headers for web applications",
			"Enable log monitoring for security events",
		},
	}, nil
}

// CheckCompliance performs compliance checks
func (s *DefaultSecurityService) CheckCompliance(ctx context.Context, standard string) (*ComplianceResult, error) {
	// Mock implementation
	return &ComplianceResult{
		ID:        "compliance-001",
		Standard:  standard,
		Timestamp: time.Now(),
		Status:    "completed",
		Score:     85.5,
		Controls: []ComplianceControl{
			{
				ID:          "control-001",
				Title:       "Access Control",
				Description: "Ensure proper access controls are in place",
				Status:      "passed",
				Evidence:    "Access controls properly configured",
				Remediation: "",
			},
		},
		Summary: ComplianceSummary{
			TotalControls:   20,
			PassedControls:  17,
			FailedControls:  2,
			WarningControls: 1,
		},
	}, nil
}

// AuditPermissions performs permission audits
func (s *DefaultSecurityService) AuditPermissions(ctx context.Context, resource string) (*AuditResult, error) {
	// Mock implementation
	return &AuditResult{
		ID:        "audit-001",
		Resource:  resource,
		Timestamp: time.Now(),
		Permissions: []Permission{
			{
				Principal: "user:admin",
				Actions:   []string{"read", "write", "delete"},
				Resource:  resource,
				Effect:    "allow",
			},
		},
		Issues: []AuditIssue{
			{
				Type:           "excessive_permissions",
				Severity:       "medium",
				Description:    "User has more permissions than required",
				Resource:       resource,
				Principal:      "user:admin",
				Recommendation: "Review and reduce permissions to minimum required",
			},
		},
		Summary: AuditSummary{
			TotalPermissions: 10,
			CriticalIssues:   0,
			HighIssues:       0,
			MediumIssues:     1,
			LowIssues:        2,
		},
	}, nil
}

// MonitorSecurityEvents monitors security events
func (s *DefaultSecurityService) MonitorSecurityEvents(ctx context.Context) (<-chan SecurityEvent, error) {
	events := make(chan SecurityEvent, 100)

	// Mock implementation - would integrate with SIEM systems
	go func() {
		defer close(events)
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				event := SecurityEvent{
					ID:          fmt.Sprintf("event-%d", time.Now().Unix()),
					Type:        "login_attempt",
					Timestamp:   time.Now(),
					Source:      "auth-service",
					Severity:    "info",
					Description: "User login attempt",
					Details: map[string]string{
						"user": "admin",
						"ip":   "192.168.1.100",
					},
					Actions: []string{"logged"},
				}
				select {
				case events <- event:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return events, nil
}

// GenerateSecurityReport generates a comprehensive security report
func (s *DefaultSecurityService) GenerateSecurityReport(ctx context.Context, options ReportOptions) (*SecurityReport, error) {
	// Mock implementation
	return &SecurityReport{
		ID:        "report-001",
		Timestamp: time.Now(),
		Type:      options.Type,
		ExecutiveSummary: ExecutiveSummary{
			OverallRiskScore:     7.5,
			CriticalFindings:     0,
			HighPriorityFindings: 2,
			ComplianceScore:      85.5,
			KeyRecommendations: []string{
				"Update TLS configuration",
				"Implement security monitoring",
				"Review access permissions",
			},
		},
		ScanResults:       []ScanResult{},
		ComplianceResults: []ComplianceResult{},
		AuditResults:      []AuditResult{},
		Recommendations: []string{
			"Implement continuous security monitoring",
			"Regular security training for staff",
			"Establish incident response procedures",
		},
	}, nil
}

// ValidateSecurityPolicies validates security policies
func (s *DefaultSecurityService) ValidateSecurityPolicies(ctx context.Context, policies []Policy) (*ValidationResult, error) {
	// Mock implementation
	return &ValidationResult{
		ID:        "validation-001",
		Timestamp: time.Now(),
		Status:    "completed",
		Policies: []PolicyValidation{
			{
				PolicyID: "policy-001",
				Status:   "valid",
				Issues:   []ValidationIssue{},
				Warnings: []ValidationWarning{},
			},
		},
		Summary: ValidationSummary{
			TotalPolicies:   len(policies),
			ValidPolicies:   len(policies),
			InvalidPolicies: 0,
			TotalIssues:     0,
			TotalWarnings:   0,
		},
	}, nil
}

// Enhanced Security Manager with Encryption and Audit Logging
type SecurityManager struct {
	config     *SecurityConfig
	auditor    *AuditLogger
	encryptor  *Encryptor
	keyManager *KeyManager
	logger     *logrus.Logger
	mu         sync.RWMutex
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	Encryption     bool   `json:"encryption" yaml:"encryption"`
	AuditLogging   bool   `json:"audit_logging" yaml:"audit_logging"`
	KeyManagement  string `json:"key_management" yaml:"key_management"`
	ComplianceMode string `json:"compliance_mode" yaml:"compliance_mode"`
	AuditLogPath   string `json:"audit_log_path" yaml:"audit_log_path"`
	KeyStorePath   string `json:"key_store_path" yaml:"key_store_path"`
	RotationPeriod int    `json:"rotation_period" yaml:"rotation_period"`
}

// AuditEvent represents an audit event
type AuditEvent struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	EventType  string                 `json:"event_type"`
	User       string                 `json:"user"`
	Resource   string                 `json:"resource"`
	Action     string                 `json:"action"`
	Result     string                 `json:"result"`
	Details    map[string]interface{} `json:"details"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	SessionID  string                 `json:"session_id"`
	Severity   string                 `json:"severity"`
	Compliance []string               `json:"compliance"`
}

// AuditLogger handles audit logging
type AuditLogger struct {
	config *SecurityConfig
	logger *logrus.Logger
	file   *os.File
	mu     sync.Mutex
}

// Encryptor handles encryption operations
type Encryptor struct {
	keyManager *KeyManager
	logger     *logrus.Logger
}

// KeyManager manages encryption keys
type KeyManager struct {
	config   *SecurityConfig
	keys     map[string][]byte
	keyStore string
	logger   *logrus.Logger
	mu       sync.RWMutex
}

// NewSecurityManager creates a new security manager
func NewSecurityManager(config *SecurityConfig) (*SecurityManager, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Initialize key manager
	keyManager, err := NewKeyManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize key manager: %w", err)
	}

	// Initialize encryptor
	encryptor := NewEncryptor(keyManager)

	// Initialize audit logger
	auditor, err := NewAuditLogger(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize audit logger: %w", err)
	}

	return &SecurityManager{
		config:     config,
		auditor:    auditor,
		encryptor:  encryptor,
		keyManager: keyManager,
		logger:     logger,
	}, nil
}

// NewKeyManager creates a new key manager
func NewKeyManager(config *SecurityConfig) (*KeyManager, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	km := &KeyManager{
		config:   config,
		keys:     make(map[string][]byte),
		keyStore: config.KeyStorePath,
		logger:   logger,
	}

	// Load existing keys
	if err := km.loadKeys(); err != nil {
		return nil, fmt.Errorf("failed to load keys: %w", err)
	}

	return km, nil
}

// NewEncryptor creates a new encryptor
func NewEncryptor(keyManager *KeyManager) *Encryptor {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	return &Encryptor{
		keyManager: keyManager,
		logger:     logger,
	}
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(config *SecurityConfig) (*AuditLogger, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	auditor := &AuditLogger{
		config: config,
		logger: logger,
	}

	// Open audit log file if audit logging is enabled
	if config.AuditLogging {
		if err := auditor.openLogFile(); err != nil {
			return nil, fmt.Errorf("failed to open audit log file: %w", err)
		}
	}

	return auditor, nil
}

// Encrypt encrypts data using AES-GCM
func (e *Encryptor) Encrypt(data []byte, keyName string) ([]byte, error) {
	key, err := e.keyManager.GetKey(keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-GCM
func (e *Encryptor) Decrypt(data []byte, keyName string) ([]byte, error) {
	key, err := e.keyManager.GetKey(keyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get decryption key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// GetKey retrieves a key by name
func (km *KeyManager) GetKey(name string) ([]byte, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	key, exists := km.keys[name]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", name)
	}

	return key, nil
}

// GenerateKey generates a new key
func (km *KeyManager) GenerateKey(name string) ([]byte, error) {
	km.mu.Lock()
	defer km.mu.Unlock()

	key := make([]byte, 32) // 256-bit key
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	km.keys[name] = key

	// Save key to store
	if err := km.saveKey(name, key); err != nil {
		return nil, fmt.Errorf("failed to save key: %w", err)
	}

	km.logger.Infof("Generated new key: %s", name)
	return key, nil
}

// LogEvent logs an audit event
func (al *AuditLogger) LogEvent(event *AuditEvent) error {
	if !al.config.AuditLogging {
		return nil
	}

	al.mu.Lock()
	defer al.mu.Unlock()

	// Set event ID and timestamp if not set
	if event.ID == "" {
		event.ID = fmt.Sprintf("audit_%d", time.Now().UnixNano())
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Marshal event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal audit event: %w", err)
	}

	// Write to log file
	if al.file != nil {
		if _, err := al.file.WriteString(string(data) + "\n"); err != nil {
			return fmt.Errorf("failed to write audit log: %w", err)
		}
		al.file.Sync()
	}

	// Log to standard logger
	al.logger.WithFields(logrus.Fields{
		"event_id":   event.ID,
		"event_type": event.EventType,
		"user":       event.User,
		"resource":   event.Resource,
		"action":     event.Action,
		"result":     event.Result,
		"severity":   event.Severity,
	}).Info("Audit event logged")

	return nil
}

// openLogFile opens the audit log file
func (al *AuditLogger) openLogFile() error {
	if al.config.AuditLogPath == "" {
		return fmt.Errorf("audit log path not configured")
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(al.config.AuditLogPath), 0755); err != nil {
		return fmt.Errorf("failed to create audit log directory: %w", err)
	}

	file, err := os.OpenFile(al.config.AuditLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return fmt.Errorf("failed to open audit log file: %w", err)
	}

	al.file = file
	return nil
}

// loadKeys loads keys from the key store
func (km *KeyManager) loadKeys() error {
	if km.keyStore == "" {
		return nil
	}

	// Create key store directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(km.keyStore), 0700); err != nil {
		return fmt.Errorf("failed to create key store directory: %w", err)
	}

	// Load keys from file if it exists
	if _, err := os.Stat(km.keyStore); os.IsNotExist(err) {
		// Generate default key if no key store exists
		if _, err := km.GenerateKey("default"); err != nil {
			return fmt.Errorf("failed to generate default key: %w", err)
		}
		return nil
	}

	data, err := os.ReadFile(km.keyStore)
	if err != nil {
		return fmt.Errorf("failed to read key store: %w", err)
	}

	var keyData map[string]string
	if err := json.Unmarshal(data, &keyData); err != nil {
		return fmt.Errorf("failed to unmarshal key data: %w", err)
	}

	// Decode keys
	for name, encodedKey := range keyData {
		key, err := base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			km.logger.Warnf("Failed to decode key %s: %v", name, err)
			continue
		}
		km.keys[name] = key
	}

	km.logger.Infof("Loaded %d keys from key store", len(km.keys))
	return nil
}

// saveKey saves a key to the key store
func (km *KeyManager) saveKey(name string, key []byte) error {
	if km.keyStore == "" {
		return nil
	}

	// Load existing keys
	var keyData map[string]string
	if data, err := os.ReadFile(km.keyStore); err == nil {
		json.Unmarshal(data, &keyData)
	}

	if keyData == nil {
		keyData = make(map[string]string)
	}

	// Add new key
	keyData[name] = base64.StdEncoding.EncodeToString(key)

	// Save to file
	data, err := json.MarshalIndent(keyData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal key data: %w", err)
	}

	if err := os.WriteFile(km.keyStore, data, 0600); err != nil {
		return fmt.Errorf("failed to write key store: %w", err)
	}

	return nil
}

// LogSecurityEvent logs a security-related event
func (sm *SecurityManager) LogSecurityEvent(eventType, user, resource, action, result string, details map[string]interface{}) error {
	event := &AuditEvent{
		EventType: eventType,
		User:      user,
		Resource:  resource,
		Action:    action,
		Result:    result,
		Details:   details,
		Severity:  sm.determineSeverity(eventType, result),
	}

	return sm.auditor.LogEvent(event)
}

// determineSeverity determines the severity of an event
func (sm *SecurityManager) determineSeverity(eventType, result string) string {
	if result == "failure" || result == "error" {
		return "high"
	}

	switch eventType {
	case "authentication", "authorization":
		return "medium"
	case "data_access", "configuration_change":
		return "low"
	default:
		return "low"
	}
}

// EncryptSensitiveData encrypts sensitive configuration data
func (sm *SecurityManager) EncryptSensitiveData(data map[string]interface{}) (map[string]interface{}, error) {
	if !sm.config.Encryption {
		return data, nil
	}

	result := make(map[string]interface{})
	for k, v := range data {
		if sm.isSensitiveField(k) {
			// Convert to JSON and encrypt
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal sensitive data: %w", err)
			}

			encrypted, err := sm.encryptor.Encrypt(jsonData, "default")
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt sensitive data: %w", err)
			}

			result[k] = base64.StdEncoding.EncodeToString(encrypted)
		} else {
			result[k] = v
		}
	}

	return result, nil
}

// isSensitiveField checks if a field contains sensitive data
func (sm *SecurityManager) isSensitiveField(field string) bool {
	sensitiveFields := []string{
		"password", "secret", "key", "token", "credential",
		"api_key", "access_key", "private_key", "client_secret",
	}

	fieldLower := strings.ToLower(field)
	for _, sensitive := range sensitiveFields {
		if strings.Contains(fieldLower, sensitive) {
			return true
		}
	}

	return false
}
