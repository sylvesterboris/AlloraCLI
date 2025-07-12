# Frequently Asked Questions

## General Questions

### What is AlloraCLI?
AlloraCLI is an AI-powered command-line interface for infrastructure management that supports multiple cloud providers and uses natural language processing for DevOps automation.

### Which cloud providers are supported?
Currently supported: AWS, Azure, Google Cloud Platform. More providers are planned for future releases.

### Do I need programming knowledge to use AlloraCLI?
No! AlloraCLI is designed to work with natural language queries. You can ask questions in plain English.

## Installation & Setup

### What are the system requirements?
- Go 1.21 or higher
- Supported OS: Windows, macOS, Linux
- Internet connection for AI services

### How do I update AlloraCLI?
Download the latest release or rebuild from source:
```bash
go install github.com/AlloraAi/AlloraCLI/cmd/allora@latest
```

## Usage

### Can I use AlloraCLI without AI features?
Yes, many commands work without AI integration, but you'll miss the intelligent automation features.

### Is my data secure?
Yes, AlloraCLI implements enterprise-grade security with encryption and audit logging. See our [Security Policy](../SECURITY.md).

### Can I use AlloraCLI in CI/CD pipelines?
Absolutely! AlloraCLI is designed for both interactive and automated use.

## Troubleshooting

### Why am I getting authentication errors?
Check your API keys and cloud provider credentials. See the [Configuration Guide](configuration.md).

### The AI responses seem incorrect. What should I do?
1. Verify your query is clear and specific
2. Check your AI model configuration
3. Report issues on GitHub

## Contributing

### How can I contribute?
See our [Contributing Guide](../CONTRIBUTING.md) for detailed information.

### Can I create custom plugins?
Yes! See the [Plugin Development Guide](plugins.md).

## Support

### Where can I get help?
- GitHub Issues for bugs and feature requests
- Discord for community support
- Email for enterprise support

### Is commercial support available?
Yes, we offer enterprise support plans. Contact sales@alloracli.com for details.
