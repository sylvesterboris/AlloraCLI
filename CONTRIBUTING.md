# Contributing to AlloraCLI

Thank you for your interest in contributing to AlloraCLI! We welcome contributions from the community and are excited to collaborate with you.

## 🌟 Ways to Contribute

- 🐛 **Bug Reports**: Help us identify and fix issues
- 💡 **Feature Requests**: Suggest new features or improvements
- 📖 **Documentation**: Improve our documentation
- 🧪 **Testing**: Add or improve tests
- 🔧 **Code**: Submit pull requests with bug fixes or new features
- 🎨 **Design**: Contribute UI/UX improvements
- 💬 **Community**: Help other users in discussions

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Basic understanding of Go programming language
- Familiarity with CLI tools and infrastructure concepts

### Development Setup

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/AlloraCLI.git
   cd AlloraCLI
   ```

3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/AlloraAi/AlloraCLI.git
   ```

4. **Install dependencies**:
   ```bash
   go mod download
   ```

5. **Run tests**:
   ```bash
   make test
   ```

6. **Build the project**:
   ```bash
   make build
   ```

## 🔧 Development Workflow

### Making Changes

1. **Create a new branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Follow the coding standards below
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   make test
   make lint
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request** on GitHub

### Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(agents): add support for Azure OpenAI
fix(config): resolve configuration parsing error
docs(readme): update installation instructions
test(cloud): add AWS provider integration tests
```

## 📝 Coding Standards

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Run `golint` and address any issues
- Use meaningful variable and function names
- Add comments for exported functions and types

### Code Structure

```
AlloraCLI/
├── cmd/                    # CLI commands
│   └── allora/
│       ├── main.go        # Main entry point
│       └── *.go           # Individual command files
├── pkg/                   # Core packages
│   ├── agents/            # AI agent implementations
│   ├── cloud/             # Cloud provider integrations
│   ├── config/            # Configuration management
│   ├── monitor/           # Monitoring integrations
│   └── utils/             # Utility functions
├── docs/                  # Documentation
├── test/                  # Test files
└── scripts/               # Build and deployment scripts
```

### Testing Guidelines

- Write unit tests for new functions
- Use table-driven tests where appropriate
- Include integration tests for major features
- Aim for good test coverage (>80%)
- Use mocks for external dependencies

**Test Example:**
```go
func TestAgentManager_AddAgent(t *testing.T) {
    tests := []struct {
        name    string
        agent   Agent
        wantErr bool
    }{
        {
            name:    "valid agent",
            agent:   &mockAgent{name: "test-agent"},
            wantErr: false,
        },
        {
            name:    "nil agent",
            agent:   nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            manager := NewAgentManager()
            err := manager.AddAgent(tt.agent)
            if (err != nil) != tt.wantErr {
                t.Errorf("AddAgent() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 📋 Pull Request Process

1. **Check Prerequisites**:
   - [ ] Branch is up to date with main
   - [ ] All tests pass
   - [ ] Code is properly formatted
   - [ ] Documentation is updated

2. **Create Pull Request**:
   - Use a clear, descriptive title
   - Fill out the PR template completely
   - Link related issues
   - Add screenshots for UI changes

3. **Code Review**:
   - Address reviewer feedback promptly
   - Keep discussions respectful and constructive
   - Make requested changes in new commits

4. **Merge**:
   - Maintainers will merge after approval
   - Your branch will be deleted after merge

## 🐛 Reporting Issues

### Bug Reports

When reporting bugs, please include:

1. **Clear description** of the issue
2. **Steps to reproduce** the problem
3. **Expected behavior** vs. actual behavior
4. **Environment details**:
   - OS (Windows, macOS, Linux)
   - Go version
   - AlloraCLI version
5. **Logs or error messages**
6. **Screenshots** if applicable

### Feature Requests

When requesting features, please include:

1. **Clear description** of the feature
2. **Use case** or problem it solves
3. **Proposed solution** (if any)
4. **Alternative solutions** considered
5. **Impact** on existing functionality

## 📚 Documentation

### Writing Documentation

- Use clear, concise language
- Include code examples where helpful
- Keep documentation up to date with code changes
- Use proper Markdown formatting

### Documentation Structure

```
docs/
├── getting-started.md      # Getting started guide
├── configuration.md        # Configuration reference
├── api.md                 # API documentation
├── plugins.md             # Plugin development guide
├── troubleshooting.md     # Common issues and solutions
└── assets/                # Images and other assets
```

## 🏗️ Project Structure

### Adding New Commands

1. Create a new file in `cmd/allora/`
2. Define the command using Cobra
3. Add the command to the root command
4. Update documentation
5. Add tests

### Adding New Agents

1. Create implementation in `pkg/agents/`
2. Implement the `Agent` interface
3. Add configuration support
4. Add unit tests
5. Update documentation

### Adding Cloud Providers

1. Create implementation in `pkg/cloud/`
2. Implement standard interfaces
3. Add configuration support
4. Add integration tests
5. Update documentation

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Provide constructive feedback
- Help others learn and grow

### Communication

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and discussions
- **Discord**: Real-time chat and support
- **Email**: security@alloracli.com for security issues

## 🏆 Recognition

We appreciate all contributions! Contributors will be:

- Added to the [contributors list](https://github.com/AlloraAi/AlloraCLI/graphs/contributors)
- Mentioned in release notes for significant contributions
- Eligible for special recognition in our community

## 📞 Getting Help

If you need help contributing:

1. Check the [documentation](docs/)
2. Search existing [issues](https://github.com/AlloraAi/AlloraCLI/issues)
3. Ask in [GitHub Discussions](https://github.com/AlloraAi/AlloraCLI/discussions)
4. Join our [Discord server](https://discord.gg/alloracli)

## 🔒 Security

If you discover a security vulnerability, please email security@alloracli.com instead of creating a public issue.

---

Thank you for contributing to AlloraCLI! Your help makes this project better for everyone. 🙏
