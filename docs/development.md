# Development Guide

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Git
- Make (optional but recommended)

### Setup Development Environment

1. **Clone the repository**
```bash
git clone https://github.com/AlloraAi/AlloraCLI.git
cd AlloraCLI
```

2. **Install dependencies**
```bash
go mod download
go mod tidy
```

3. **Build the project**
```bash
make build
# or
go build -o bin/allora ./cmd/allora
```

4. **Run tests**
```bash
make test
# or
go test ./...
```

5. **Install pre-commit hooks**
```bash
# Install pre-commit (if not already installed)
pip install pre-commit
pre-commit install
```

## Project Structure

```
AlloraCLI/
├── cmd/                    # Application entry points
│   └── allora/            # Main CLI application
├── pkg/                   # Core packages
│   ├── agents/           # AI agent system
│   ├── cloud/            # Cloud provider integrations
│   ├── config/           # Configuration management
│   ├── monitor/          # Monitoring and metrics
│   ├── security/         # Security features
│   ├── ui/               # User interface components
│   └── utils/            # Utility functions
├── plugins/              # Plugin system
├── docs/                 # Documentation
├── examples/             # Usage examples
├── test/                 # Test utilities and integration tests
├── scripts/              # Build and deployment scripts
├── .github/              # GitHub workflows and templates
├── Makefile              # Build automation
├── go.mod                # Go module definition
└── README.md             # Project overview
```

## Development Workflow

### 1. Feature Development

```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make changes
# Write tests
# Update documentation

# Run tests
make test

# Check code quality
make lint

# Commit changes
git commit -m "feat: add your feature description"

# Push and create PR
git push origin feature/your-feature-name
```

### 2. Testing

#### Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Integration Tests
```bash
# Run integration tests (requires setup)
go test ./test/integration/...
```

#### Benchmarks
```bash
# Run benchmarks
go test -bench=. ./...
```

### 3. Code Quality

#### Linting
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

#### Formatting
```bash
# Format code
go fmt ./...

# Import organization
goimports -w .
```

#### Security Scanning
```bash
# Install gosec
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Run security scan
gosec ./...
```

## Build System

### Makefile Targets

```bash
make help          # Show available targets
make build         # Build the application
make test          # Run tests
make lint          # Run linters
make clean         # Clean build artifacts
make install       # Install to $GOPATH/bin
make release       # Create release builds
```

### Cross-Platform Building

```bash
# Build for all platforms
make build-all

# Build for specific platform
GOOS=linux GOARCH=amd64 make build
GOOS=windows GOARCH=amd64 make build
GOOS=darwin GOARCH=amd64 make build
```

## Testing Guidelines

### Writing Tests

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **End-to-End Tests**: Test complete workflows

#### Example Test Structure

```go
func TestAgentQuery(t *testing.T) {
    // Arrange
    agent := &MockAgent{}
    query := &agents.Query{Text: "test query"}
    
    // Act
    response, err := agent.Query(context.Background(), query)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Contains(t, response.Text, "expected content")
}
```

### Mocking

Use interfaces for testability:

```go
// Good: interface allows mocking
type CloudProvider interface {
    ListResources(ctx context.Context) ([]*Resource, error)
}

// Test with mock
func TestResourceListing(t *testing.T) {
    mockProvider := &MockCloudProvider{}
    mockProvider.On("ListResources").Return([]*Resource{}, nil)
    
    // Test implementation
}
```

## Contributing Guidelines

### Code Style

1. Follow [Effective Go](https://golang.org/doc/effective_go)
2. Use meaningful variable and function names
3. Write comprehensive comments for public APIs
4. Keep functions small and focused

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new AI agent capability
fix: resolve authentication issue
docs: update API documentation
test: add integration tests for cloud providers
```

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Update documentation
6. Submit pull request with clear description

## Debugging

### Debug Mode

```bash
# Enable debug logging
ALLORA_LOG_LEVEL=debug ./bin/allora <command>

# Or use flag
./bin/allora --log-level debug <command>
```

### Profiling

```bash
# CPU profiling
go tool pprof ./bin/allora cpu.prof

# Memory profiling
go tool pprof ./bin/allora mem.prof
```

### Common Issues

1. **Import cycles**: Restructure packages to avoid circular dependencies
2. **Race conditions**: Use `go test -race` to detect
3. **Memory leaks**: Use `go tool pprof` to analyze

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes (backward compatible)

### Creating a Release

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create and push tag: `git tag v1.0.0`
4. GitHub Actions will automatically create release

## Getting Help

- Check existing [GitHub Issues](https://github.com/AlloraAi/AlloraCLI/issues)
- Join our [Discord](https://discord.gg/alloracli)
- Read the [Contributing Guide](../CONTRIBUTING.md)
- Email: developers@alloracli.com
