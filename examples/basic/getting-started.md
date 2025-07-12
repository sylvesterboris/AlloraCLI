# Getting Started Examples

## Installation and Setup

### 1. Download and Install

```bash
# Linux/macOS
curl -L https://github.com/AlloraAi/AlloraCLI/releases/latest/download/allora-linux-amd64 -o allora
chmod +x allora
sudo mv allora /usr/local/bin/

# Verify installation
allora version
```

### 2. Initialize Configuration

```bash
# Initialize AlloraCLI
allora init

# This creates ~/.allora/config.yaml with default settings
```

### 3. Configure Your First Provider (AWS)

```bash
# Set AWS credentials
allora config set aws.profile default
allora config set aws.region us-west-2

# Or use environment variables
export AWS_PROFILE=default
export AWS_REGION=us-west-2
```

### 4. Test Your Setup

```bash
# Ask a simple question
allora ask "What AWS resources do I have in us-west-2?"

# List your EC2 instances
allora cloud aws list-resources ec2

# Check connection status
allora config validate
```

## First Tasks

### Infrastructure Discovery

```bash
# Discover all resources across providers
allora ask "Show me all my cloud resources"

# Get cost breakdown
allora ask "What are my current cloud costs?"

# Security audit
allora security audit --provider aws
```

### Basic Monitoring

```bash
# Check system health
allora monitor health

# Get metrics for a specific resource
allora monitor metrics --resource i-1234567890abcdef0

# Set up basic alerting
allora ask "How do I set up CPU alerts for my EC2 instances?"
```

### Getting Help

```bash
# General help
allora help

# Command-specific help
allora ask --help
allora deploy --help

# Interactive mode
allora gemini
```

## Next Steps

1. Read the [Configuration Guide](../../docs/configuration.md)
2. Explore [Common Commands](commands.md)
3. Set up additional [Cloud Providers](../cloud/)
4. Try [Advanced Use Cases](../advanced/)
