# AlloraCLI Usage Guide

Welcome to AlloraCLI! This comprehensive guide will help you get started with AI-powered infrastructure management and take you from beginner to advanced user.

## 📋 Table of Contents

1. [Getting Started](#getting-started)
2. [Installation](#installation)
3. [Initial Setup](#initial-setup)
4. [Basic Usage](#basic-usage)
5. [Core Commands](#core-commands)
6. [AI-Powered Features](#ai-powered-features)
7. [Cloud Provider Integration](#cloud-provider-integration)
8. [Advanced Usage](#advanced-usage)
9. [Plugin System](#plugin-system)
10. [Troubleshooting](#troubleshooting)
11. [Best Practices](#best-practices)
12. [Examples](#examples)

---

## 🚀 Getting Started

### What is AlloraCLI?

AlloraCLI is an AI-powered command-line interface that revolutionizes infrastructure management through:
- **Natural Language Processing**: Interact with your infrastructure using plain English
- **Multi-Cloud Support**: Manage AWS, Azure, and GCP from one tool
- **Intelligent Automation**: AI-driven insights and automated operations
- **Real-time Monitoring**: Built-in monitoring and alerting capabilities
- **Security-First Approach**: Comprehensive security analysis and compliance

### Prerequisites

Before you begin, ensure you have:
- **Operating System**: Windows 10+, macOS 10.15+, or Linux
- **Network Access**: Internet connection for AI services and cloud APIs
- **Permissions**: Administrative access for installation (optional)
- **Cloud Accounts**: AWS, Azure, or GCP accounts (for cloud features)

---

## 📦 Installation

### Option 1: Download Pre-built Binaries (Recommended)

Visit the [GitHub Releases page](https://github.com/AlloraAi/AlloraCLI/releases) and download the appropriate binary for your system:

#### Linux/macOS
```bash
# Download the latest release
curl -L https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-linux-amd64 -o allora

# Make it executable
chmod +x allora

# Move to system PATH (optional)
sudo mv allora /usr/local/bin/

# Verify installation
allora version
```

#### Windows
```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-windows-amd64.exe" -OutFile "allora.exe"

# Add to PATH (optional)
# Move allora.exe to a directory in your PATH or add current directory to PATH

# Verify installation
.\allora.exe version
```

### Option 2: Install from Source

```bash
# Clone the repository
git clone https://github.com/AlloraAi/AlloraCLI.git
cd AlloraCLI

# Build from source
go build -o allora ./cmd/allora

# Verify build
./allora version
```

### Option 3: Package Managers (Future Support)

```bash
# Homebrew (macOS/Linux) - Coming Soon
brew install AlloraAi/tap/allora

# Scoop (Windows) - Coming Soon
scoop bucket add AlloraAi https://github.com/AlloraAi/scoop-bucket
scoop install allora

# Chocolatey (Windows) - Coming Soon
choco install allora

# APT (Ubuntu/Debian) - Coming Soon
sudo apt install allora

# YUM (RHEL/CentOS) - Coming Soon
sudo yum install allora
```

---

## ⚙️ Initial Setup

### 1. Initialize AlloraCLI

```bash
# Initialize configuration
allora init

# This creates:
# - ~/.config/alloracli/config.yaml (Linux/macOS)
# - %APPDATA%\alloracli\config.yaml (Windows)
```

### 2. Configure AI Services

```bash
# Set OpenAI API key for AI features
allora config set openai.api_key YOUR_OPENAI_API_KEY

# Optional: Configure other AI providers
allora config set anthropic.api_key YOUR_ANTHROPIC_KEY
allora config set google.api_key YOUR_GOOGLE_AI_KEY
```

### 3. Configure Cloud Providers

#### AWS Configuration
```bash
# Method 1: Direct configuration
allora config set aws.access_key_id YOUR_AWS_ACCESS_KEY
allora config set aws.secret_access_key YOUR_AWS_SECRET_KEY
allora config set aws.region us-west-2

# Method 2: Use AWS CLI profiles
allora config set aws.profile default
allora config set aws.region us-west-2

# Method 3: Use IAM roles (for EC2 instances)
allora config set aws.use_iam_role true
```

#### Azure Configuration
```bash
# Service Principal authentication
allora config set azure.client_id YOUR_CLIENT_ID
allora config set azure.client_secret YOUR_CLIENT_SECRET
allora config set azure.tenant_id YOUR_TENANT_ID
allora config set azure.subscription_id YOUR_SUBSCRIPTION_ID

# Or use Azure CLI authentication
allora config set azure.use_cli_auth true
```

#### Google Cloud Platform
```bash
# Service Account authentication
allora config set gcp.service_account_path /path/to/service-account.json
allora config set gcp.project_id your-project-id

# Or use Application Default Credentials
allora config set gcp.use_adc true
allora config set gcp.project_id your-project-id
```

### 4. Verify Configuration

```bash
# Show current configuration
allora config show

# Test cloud provider connections
allora config test aws
allora config test azure
allora config test gcp

# Run configuration diagnostics
allora config doctor
```

---

## 🎯 Basic Usage

### Help System

```bash
# General help
allora help
allora --help

# Command-specific help
allora help ask
allora help deploy
allora monitor --help

# List all available commands
allora help commands

# Get version information
allora version
allora version --detailed
```

### Configuration Management

```bash
# View all settings
allora config show

# Set a configuration value
allora config set key value

# Get a specific configuration value
allora config get aws.region

# Remove a configuration value
allora config unset aws.secret_access_key

# Reset to defaults
allora config reset

# Export configuration
allora config export > my-config.yaml

# Import configuration
allora config import my-config.yaml
```

---

## 🛠️ Core Commands

### 1. Ask Command - Natural Language Interface

The `ask` command is your gateway to AI-powered infrastructure management:

```bash
# Infrastructure insights
allora ask "What's the current status of my AWS infrastructure?"
allora ask "Show me my most expensive cloud resources"
allora ask "Which services are consuming the most CPU?"

# Cost optimization
allora ask "How can I reduce my cloud costs?"
allora ask "What are my biggest expenses this month?"
allora ask "Suggest cost optimization strategies"

# Performance analysis
allora ask "Why is my application slow?"
allora ask "What's causing high memory usage?"
allora ask "How can I improve database performance?"

# Security insights
allora ask "Are there any security vulnerabilities in my infrastructure?"
allora ask "What security best practices should I implement?"
allora ask "Show me recent security events"

# Capacity planning
allora ask "Do I need to scale my infrastructure?"
allora ask "What's my projected growth for the next quarter?"
allora ask "When should I add more capacity?"
```

### 2. Deploy Command - Application Deployment

```bash
# Basic deployment
allora deploy --app myapp --environment production

# Deployment with specific version
allora deploy --app myapp --version v1.2.3 --environment staging

# Rollback deployment
allora deploy rollback --app myapp --version v1.2.2

# Deployment status
allora deploy status --app myapp

# Deployment history
allora deploy history --app myapp --limit 10

# Deployment validation
allora deploy validate --config deployment.yaml

# Blue-green deployment
allora deploy --strategy blue-green --app myapp

# Canary deployment
allora deploy --strategy canary --app myapp --percentage 10
```

### 3. Monitor Command - Infrastructure Monitoring

```bash
# System overview
allora monitor status
allora monitor dashboard

# Resource-specific monitoring
allora monitor --resource ec2
allora monitor --resource rds
allora monitor --resource kubernetes

# Provider-specific monitoring
allora monitor --provider aws
allora monitor --provider azure
allora monitor --provider gcp

# Real-time monitoring
allora monitor --real-time --resource ec2
allora monitor --follow --service webapp

# Historical data
allora monitor --time-range "last 24h"
allora monitor --from "2025-07-10" --to "2025-07-12"

# Alerts and notifications
allora monitor alerts list
allora monitor alerts create --name "high-cpu" --threshold 80
allora monitor alerts delete alert-123
```

### 4. Troubleshoot Command - Problem Resolution

```bash
# Auto-detect issues
allora troubleshoot auto

# Service-specific troubleshooting
allora troubleshoot --service database
allora troubleshoot --service webserver
allora troubleshoot --service network

# Issue-specific troubleshooting
allora troubleshoot --issue "high latency"
allora troubleshoot --issue "connection timeout"
allora troubleshoot --issue "out of memory"

# Incident management
allora troubleshoot incident create --title "Database slowdown"
allora troubleshoot incident list
allora troubleshoot incident show INC-2025-001

# Diagnostic information
allora troubleshoot diagnose --comprehensive
allora troubleshoot logs --service webapp --tail 100

# Auto-fix capabilities
allora troubleshoot --auto-fix --confirm
allora troubleshoot suggest --issue "disk space"
```

### 5. Security Command - Security Management

```bash
# Security scans
allora security scan
allora security scan --provider aws
allora security scan --comprehensive

# Compliance checks
allora security compliance check --standard iso27001
allora security compliance check --standard pci-dss
allora security compliance report

# Vulnerability assessment
allora security vulnerabilities scan
allora security vulnerabilities report

# Security monitoring
allora security monitor events
allora security monitor threats

# Access management
allora security access review
allora security permissions audit

# Security policies
allora security policies list
allora security policies validate
```

### 6. Cloud Command - Cloud Provider Operations

```bash
# AWS operations
allora cloud aws ec2 list
allora cloud aws s3 buckets
allora cloud aws rds instances

# Azure operations
allora cloud azure vm list
allora cloud azure storage accounts
allora cloud azure sql databases

# GCP operations
allora cloud gcp compute instances
allora cloud gcp storage buckets
allora cloud gcp sql instances

# Cross-cloud operations
allora cloud inventory
allora cloud costs compare
allora cloud migrate plan --from aws --to azure
```

---

## 🤖 AI-Powered Features

### Gemini Interface

Launch the interactive AI interface for conversational infrastructure management:

```bash
# Start Gemini interface
allora gemini

# Start with specific context
allora gemini --context production
allora gemini --provider aws
```

In the Gemini interface, you can:
- Ask complex questions about your infrastructure
- Get real-time insights and recommendations
- Execute commands through natural language
- Visualize data and metrics
- Collaborate with team members

### Natural Language Processing

AlloraCLI understands various ways to express the same request:

```bash
# These all do the same thing:
allora ask "Show me my EC2 instances"
allora ask "List all my AWS servers"
allora ask "What virtual machines do I have on Amazon?"
allora ask "Display my compute instances in AWS"

# Complex queries:
allora ask "Show me EC2 instances in us-west-2 that are running and have high CPU usage"
allora ask "Find databases that haven't been backed up in the last 7 days"
allora ask "List all resources created in the last week that cost more than $100"
```

### Intelligent Suggestions

AlloraCLI provides proactive suggestions based on your infrastructure:

```bash
# Get optimization suggestions
allora suggest optimizations

# Get security recommendations
allora suggest security

# Get cost reduction ideas
allora suggest costs

# Get performance improvements
allora suggest performance
```

---

## ☁️ Cloud Provider Integration

### AWS Integration

```bash
# Setup AWS credentials
allora config set aws.access_key_id YOUR_KEY
allora config set aws.secret_access_key YOUR_SECRET
allora config set aws.region us-west-2

# Common AWS operations
allora ask "List my EC2 instances"
allora ask "Show S3 bucket sizes"
allora ask "What's my AWS bill this month?"

# Specific AWS commands
allora cloud aws ec2 list --region us-east-1
allora cloud aws s3 ls
allora cloud aws iam users
allora cloud aws cloudformation stacks

# AWS monitoring
allora monitor --provider aws --service ec2
allora monitor --provider aws --service rds
```

### Azure Integration

```bash
# Setup Azure credentials
allora config set azure.client_id YOUR_CLIENT_ID
allora config set azure.client_secret YOUR_SECRET
allora config set azure.tenant_id YOUR_TENANT_ID
allora config set azure.subscription_id YOUR_SUBSCRIPTION_ID

# Common Azure operations
allora ask "Show my Azure virtual machines"
allora ask "List storage accounts"
allora ask "What's my Azure spending?"

# Specific Azure commands
allora cloud azure vm list
allora cloud azure storage list
allora cloud azure sql list
allora cloud azure resourcegroups list

# Azure monitoring
allora monitor --provider azure --service vm
allora monitor --provider azure --service sql
```

### Google Cloud Platform Integration

```bash
# Setup GCP credentials
allora config set gcp.service_account_path /path/to/key.json
allora config set gcp.project_id your-project-id

# Common GCP operations
allora ask "Show my GCP compute instances"
allora ask "List Cloud Storage buckets"
allora ask "What's my GCP usage?"

# Specific GCP commands
allora cloud gcp compute list
allora cloud gcp storage buckets
allora cloud gcp sql instances
allora cloud gcp projects list

# GCP monitoring
allora monitor --provider gcp --service compute
allora monitor --provider gcp --service storage
```

---

## 🔧 Advanced Usage

### Configuration Profiles

Manage multiple environments with configuration profiles:

```bash
# Create profiles
allora config profile create production
allora config profile create staging
allora config profile create development

# Switch between profiles
allora config profile use production
allora config profile use staging

# List profiles
allora config profile list

# Copy profile settings
allora config profile copy production staging

# Delete profile
allora config profile delete development
```

### Scripting and Automation

Use AlloraCLI in scripts and automation:

```bash
#!/bin/bash
# Daily infrastructure health check

echo "Running daily infrastructure health check..."

# Check overall status
allora monitor status --json > daily-status.json

# Check for security issues
allora security scan --format json > security-report.json

# Check costs
allora ask "What's my daily spend?" --output json > costs.json

# Send alerts if issues found
if allora troubleshoot auto --check-only; then
    echo "All systems healthy"
else
    echo "Issues detected, running diagnostics..."
    allora troubleshoot diagnose --auto-fix
fi
```

### Output Formats

Control output format for integration with other tools:

```bash
# JSON output
allora monitor status --format json
allora ask "List EC2 instances" --output json

# YAML output
allora config show --format yaml
allora cloud aws ec2 list --format yaml

# Table output (default)
allora monitor status --format table

# CSV output
allora cloud aws ec2 list --format csv

# Custom formatting
allora monitor status --template "{{.Name}}: {{.Status}}"
```

### Environment Variables

Configure AlloraCLI using environment variables:

```bash
# Configuration
export ALLORA_CONFIG_PATH=/custom/path/config.yaml
export ALLORA_LOG_LEVEL=debug
export ALLORA_OUTPUT_FORMAT=json

# Cloud provider credentials
export AWS_ACCESS_KEY_ID=your_key
export AWS_SECRET_ACCESS_KEY=your_secret
export AZURE_CLIENT_ID=your_client_id
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json

# AI services
export OPENAI_API_KEY=your_openai_key
export ANTHROPIC_API_KEY=your_anthropic_key
```

---

## 🔌 Plugin System

### Managing Plugins

```bash
# List available plugins
allora plugin list

# Search for plugins
allora plugin search monitoring
allora plugin search aws

# Install plugins
allora plugin install monitoring-pro
allora plugin install cost-optimizer

# Update plugins
allora plugin update monitoring-pro
allora plugin update --all

# Remove plugins
allora plugin remove cost-optimizer

# Plugin information
allora plugin info monitoring-pro
```

### Using Plugins

```bash
# Run plugin commands
allora plugin run cost-optimizer --provider aws
allora plugin run monitoring-pro dashboard

# Plugin-specific help
allora plugin help cost-optimizer
```

### Developing Plugins

Create custom plugins to extend AlloraCLI functionality:

```bash
# Initialize plugin development
allora plugin init my-plugin

# Build plugin
allora plugin build my-plugin

# Test plugin locally
allora plugin test my-plugin

# Package plugin
allora plugin package my-plugin
```

---

## 🔍 Troubleshooting

### Common Issues

#### Authentication Problems
```bash
# Test connections
allora config test aws
allora config test azure
allora config test gcp

# Check credentials
allora config show aws
allora config get aws.access_key_id

# Reset credentials
allora config unset aws.secret_access_key
allora config set aws.secret_access_key NEW_SECRET
```

#### Network Issues
```bash
# Test connectivity
allora troubleshoot network

# Check proxy settings
allora config set proxy.http http://proxy:8080
allora config set proxy.https https://proxy:8080

# Bypass SSL verification (not recommended for production)
allora config set ssl.verify false
```

#### Performance Issues
```bash
# Enable debug logging
allora config set log.level debug

# Check system resources
allora troubleshoot system

# Clear cache
allora cache clear

# Update to latest version
allora update
```

### Debug Mode

Enable debug mode for detailed troubleshooting:

```bash
# Global debug mode
allora config set debug true

# Command-specific debug
allora --debug ask "Show my instances"
allora --verbose monitor status

# Log output to file
allora --log-file debug.log ask "Show my infrastructure"
```

### Getting Help

```bash
# Built-in help
allora help
allora help troubleshoot

# Check documentation
allora docs open
allora docs search "configuration"

# Community support
allora support contact
allora support forum
```

---

## 💡 Best Practices

### Security Best Practices

1. **Credential Management**
   ```bash
   # Use environment variables for sensitive data
   export AWS_ACCESS_KEY_ID=your_key
   export AWS_SECRET_ACCESS_KEY=your_secret
   
   # Use IAM roles when possible
   allora config set aws.use_iam_role true
   
   # Regularly rotate credentials
   allora security credentials rotate
   ```

2. **Access Control**
   ```bash
   # Follow principle of least privilege
   allora security permissions audit
   
   # Use separate credentials for different environments
   allora config profile create production
   allora config profile create staging
   ```

3. **Monitoring and Alerting**
   ```bash
   # Set up security monitoring
   allora security monitor enable
   
   # Configure alerts
   allora monitor alerts create --name "unauthorized-access" --type security
   ```

### Performance Best Practices

1. **Configuration Optimization**
   ```bash
   # Cache frequently accessed data
   allora config set cache.enabled true
   allora config set cache.ttl 300
   
   # Use specific regions
   allora config set aws.region us-west-2
   ```

2. **Efficient Querying**
   ```bash
   # Use specific filters
   allora ask "Show EC2 instances in production environment"
   
   # Limit results when appropriate
   allora monitor logs --tail 100 --service webapp
   ```

### Operational Best Practices

1. **Regular Maintenance**
   ```bash
   # Update regularly
   allora update check
   allora update install
   
   # Clean up periodically
   allora cache clear
   allora logs cleanup --older-than 30d
   ```

2. **Backup and Recovery**
   ```bash
   # Backup configuration
   allora config export > backup-$(date +%Y%m%d).yaml
   
   # Document your setup
   allora config show > infrastructure-config.yaml
   ```

---

## 📚 Examples

### Example 1: Daily Infrastructure Check

```bash
#!/bin/bash
# daily-check.sh - Daily infrastructure health check

echo "🔍 Daily Infrastructure Health Check - $(date)"
echo "================================================"

# Overall system status
echo "📊 System Status:"
allora monitor status

# Security check
echo "🔒 Security Status:"
allora security scan --quick

# Cost monitoring
echo "💰 Cost Status:"
allora ask "What's my current daily spending rate?"

# Performance check
echo "⚡ Performance Status:"
allora ask "Are there any performance issues I should know about?"

# Recommendations
echo "💡 Recommendations:"
allora suggest optimizations --top 3

echo "✅ Daily check complete"
```

### Example 2: Deployment Pipeline Integration

```bash
#!/bin/bash
# deploy.sh - Deployment script with AlloraCLI

APP_NAME="myapp"
VERSION="$1"
ENVIRONMENT="$2"

if [ -z "$VERSION" ] || [ -z "$ENVIRONMENT" ]; then
    echo "Usage: $0 <version> <environment>"
    exit 1
fi

echo "🚀 Deploying $APP_NAME version $VERSION to $ENVIRONMENT"

# Pre-deployment checks
echo "🔍 Pre-deployment checks..."
allora troubleshoot --environment $ENVIRONMENT --check-only
if [ $? -ne 0 ]; then
    echo "❌ Pre-deployment checks failed"
    exit 1
fi

# Deploy application
echo "📦 Deploying application..."
allora deploy --app $APP_NAME --version $VERSION --environment $ENVIRONMENT

# Post-deployment validation
echo "✅ Post-deployment validation..."
allora monitor health-check --app $APP_NAME --environment $ENVIRONMENT --timeout 300

# Success notification
echo "🎉 Deployment successful!"
allora ask "Send deployment notification for $APP_NAME $VERSION to $ENVIRONMENT"
```

### Example 3: Cost Optimization Script

```bash
#!/bin/bash
# cost-optimization.sh - Weekly cost optimization

echo "💰 Weekly Cost Optimization Report - $(date)"
echo "============================================="

# Current spending
echo "📊 Current Spending:"
allora ask "What's my spending this week compared to last week?"

# Cost analysis
echo "🔍 Cost Analysis:"
allora ask "What are my top 10 most expensive resources?"

# Optimization opportunities
echo "💡 Optimization Opportunities:"
allora suggest costs --detailed

# Unused resources
echo "🗑️ Unused Resources:"
allora ask "What resources haven't been used in the last 30 days?"

# Savings recommendations
echo "💵 Savings Recommendations:"
allora ask "How much could I save by rightsizing my instances?"

# Generate report
allora ask "Generate a cost optimization report" --format pdf --output "cost-report-$(date +%Y%m%d).pdf"

echo "✅ Cost optimization analysis complete"
```

### Example 4: Security Audit

```bash
#!/bin/bash
# security-audit.sh - Monthly security audit

echo "🔒 Monthly Security Audit - $(date)"
echo "=================================="

# Comprehensive security scan
echo "🔍 Running comprehensive security scan..."
allora security scan --comprehensive --format json > security-scan.json

# Compliance check
echo "📋 Checking compliance..."
allora security compliance check --standard iso27001 > compliance-report.txt

# Vulnerability assessment
echo "🚨 Vulnerability assessment..."
allora security vulnerabilities scan --severity high,critical

# Access review
echo "👥 Access review..."
allora security access review --inactive-days 90

# Security recommendations
echo "💡 Security recommendations..."
allora suggest security --priority high

# Generate security report
allora security report generate --format pdf --output "security-audit-$(date +%Y%m%d).pdf"

echo "✅ Security audit complete"
```

---

## 🎓 Learning Resources

### Official Documentation
- [Getting Started Guide](docs/getting-started.md)
- [Configuration Reference](docs/configuration.md)
- [API Documentation](docs/api.md)
- [Plugin Development Guide](docs/plugins.md)

### Community Resources
- [GitHub Repository](https://github.com/AlloraAi/AlloraCLI)
- [Community Forum](https://community.alloracli.com)
- [Discord Server](https://discord.gg/alloracli)
- [YouTube Tutorials](https://youtube.com/@alloracli)

### Training Materials
- [Interactive Tutorial](https://learn.alloracli.com)
- [Video Course](https://courses.alloracli.com)
- [Certification Program](https://certification.alloracli.com)

---

## 🤝 Getting Support

### Self-Help Resources
```bash
# Built-in help
allora help
allora docs search "your topic"

# Configuration diagnostics
allora config doctor

# System diagnostics
allora troubleshoot system
```

### Community Support
- **GitHub Issues**: Report bugs and request features
- **Community Forum**: Ask questions and share knowledge
- **Discord**: Real-time chat with the community
- **Stack Overflow**: Tag questions with `alloracli`

### Professional Support
- **Enterprise Support**: Available for organizations
- **Consulting Services**: Professional implementation and training
- **Custom Development**: Tailored solutions for specific needs

---

## 🔄 Updates and Maintenance

### Updating AlloraCLI

```bash
# Check for updates
allora update check

# Update to latest version
allora update install

# Update to specific version
allora update install --version v1.2.3

# Show update history
allora update history
```

### Maintenance Tasks

```bash
# Clear cache
allora cache clear

# Clean up logs
allora logs cleanup --older-than 30d

# Validate configuration
allora config validate

# Backup configuration
allora config export > backup.yaml
```

---

This comprehensive usage guide should help new users get started with AlloraCLI and progress to advanced usage patterns. For the most up-to-date information, always refer to the official documentation and community resources.

**Happy infrastructure management! 🚀**
