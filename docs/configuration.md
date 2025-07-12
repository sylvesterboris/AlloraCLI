# Configuration Reference

## Overview
AlloraCLI uses a hierarchical configuration system that supports multiple formats and sources.

## Configuration Files

### Main Configuration File
Default location: `~/.allora/config.yaml`

```yaml
# Cloud Provider Configurations
providers:
  aws:
    profile: default
    region: us-west-2
    access_key_id: ""  # Use environment variables or IAM roles
    secret_access_key: ""
  
  azure:
    subscription_id: ""
    tenant_id: ""
    client_id: ""
    client_secret: ""
    
  gcp:
    project_id: ""
    credentials_file: ""

# AI Configuration
ai:
  openai:
    api_key: ""
    model: "gpt-4"
    max_tokens: 2000
    temperature: 0.7
  
# Monitoring Configuration
monitoring:
  prometheus:
    url: ""
    username: ""
    password: ""
  grafana:
    url: ""
    api_key: ""

# Security Settings
security:
  encryption_key: ""
  audit_log: true
  compliance_mode: "SOC2"

# UI Settings
ui:
  color: true
  interactive: true
  auto_export: false
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ALLORA_CONFIG_DIR` | Configuration directory | `~/.allora` |
| `ALLORA_OPENAI_API_KEY` | OpenAI API key | - |
| `ALLORA_AWS_PROFILE` | AWS profile | `default` |
| `ALLORA_LOG_LEVEL` | Log level | `info` |

## Command-Line Flags

Most configuration options can be overridden via command-line flags:

```bash
allora --config /path/to/config.yaml --log-level debug ask "question"
```

## Configuration Validation

Use `allora config validate` to check your configuration:

```bash
allora config validate
```
