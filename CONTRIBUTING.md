# Contributing to Supabase Management Tools

First off, thank you for considering contributing to this project! It's people like you that make this tool better for everyone.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Workflow](#development-workflow)
- [Style Guidelines](#style-guidelines)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to uphold. Please be respectful and constructive in all interactions.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional but recommended)
- Docker (for testing supascale.sh)
- golangci-lint (for linting)

### Setting Up Your Development Environment

1. **Fork the repository** on GitHub

2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/supactl.git
   cd supactl
   ```

3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/qubitquilt/supactl.git
   ```

4. **Install dependencies**:
   ```bash
   make deps
   ```

5. **Verify your setup** by running tests:
   ```bash
   make test
   ```

## How Can I Contribute?

### Reporting Bugs

Before creating a bug report, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the problem
- **Expected behavior**
- **Actual behavior**
- **Screenshots** (if applicable)
- **Environment details** (OS, Go version, etc.)

Use the bug report template:

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce:
1. Run command '...'
2. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g., Ubuntu 22.04]
- Go version: [e.g., 1.22]
- supactl version: [e.g., 1.0.0]
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear title** describing the enhancement
- **Provide a detailed description** of the proposed feature
- **Explain why this enhancement would be useful**
- **List potential implementation approaches** if you have ideas

### Contributing Code

1. **Find or create an issue** describing what you'll work on
2. **Comment on the issue** to let others know you're working on it
3. **Follow the development workflow** below
4. **Submit a pull request** when ready

## Development Workflow

### 1. Create a Branch

Always create a new branch for your work:

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation changes
- `refactor/description` - Code refactoring
- `test/description` - Test additions/modifications

### 2. Make Your Changes

- Write clean, readable code
- Follow Go best practices
- Add tests for new functionality
- Update documentation as needed
- Keep commits atomic and focused

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run linter
make lint

# Format code
make fmt

# Check for issues
make vet

# Build to verify compilation
make build
```

### 4. Commit Your Changes

Follow the commit guidelines (see below) and commit your changes:

```bash
git add .
git commit -m "feat: add new command for instance management"
```

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Style Guidelines

### Go Code Style

We follow standard Go conventions:

- **Use `gofmt`** - All code must be formatted with `gofmt` (run `make fmt`)
- **Follow Go best practices** - See [Effective Go](https://golang.org/doc/effective_go.html)
- **Keep functions small** - Aim for single responsibility
- **Write clear variable names** - Prefer clarity over brevity
- **Add comments** - Document exported functions and complex logic
- **Handle errors properly** - Never ignore errors without good reason

Example:

```go
// CreateInstance creates a new Supabase instance with the given name.
// It returns the created instance or an error if creation fails.
func (c *Client) CreateInstance(name string) (*Instance, error) {
    if name == "" {
        return nil, fmt.Errorf("instance name cannot be empty")
    }

    // Implementation...
}
```

### Documentation Style

- Use clear, concise language
- Include code examples where helpful
- Keep line length reasonable (80-100 chars)
- Use proper Markdown formatting

## Commit Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat` - A new feature
- `fix` - A bug fix
- `docs` - Documentation changes
- `style` - Code style changes (formatting, etc.)
- `refactor` - Code refactoring
- `test` - Adding or updating tests
- `chore` - Maintenance tasks
- `perf` - Performance improvements

### Examples

```bash
# Feature
git commit -m "feat(cli): add logs command to view instance logs"

# Bug fix
git commit -m "fix(auth): correct permission error when saving config"

# Documentation
git commit -m "docs: update installation instructions"

# Refactoring
git commit -m "refactor(api): simplify error handling in client"
```

### Best Practices

- Use present tense ("add feature" not "added feature")
- Use imperative mood ("move cursor to..." not "moves cursor to...")
- Keep the subject line under 50 characters
- Capitalize the subject line
- Don't end the subject line with a period
- Use the body to explain what and why (not how)

## Pull Request Process

### Before Submitting

1. ✅ All tests pass (`make test`)
2. ✅ Code is formatted (`make fmt`)
3. ✅ Linter passes (`make lint`)
4. ✅ Documentation is updated
5. ✅ Commits follow guidelines
6. ✅ Branch is up to date with main

### Pull Request Template

```markdown
## Description
Brief description of the changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe how you tested your changes.

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-reviewed the code
- [ ] Commented complex parts
- [ ] Updated documentation
- [ ] Added tests
- [ ] All tests pass
- [ ] No new warnings
```

### Review Process

1. **Automated checks** must pass (tests, linting)
2. **At least one maintainer** must review and approve
3. **Address review feedback** promptly
4. **Squash commits** if requested
5. **Update branch** if main has changed

### After Approval

- Maintainers will merge your PR
- Your contribution will be included in the next release
- You'll be added to the contributors list

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
make test-coverage

# Run specific package tests
go test ./internal/api -v

# Run specific test
go test ./internal/api -run TestCreateInstance

# Run benchmarks
make bench
```

### Writing Tests

- Place tests in `*_test.go` files
- Use table-driven tests when appropriate
- Test both success and error cases
- Mock external dependencies
- Aim for high coverage on critical paths

Example:

```go
func TestCreateInstance(t *testing.T) {
    tests := []struct {
        name        string
        projectName string
        wantErr     bool
    }{
        {
            name:        "valid project name",
            projectName: "my-project",
            wantErr:     false,
        },
        {
            name:        "invalid empty name",
            projectName: "",
            wantErr:     true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Project-Specific Notes

### supactl (Go CLI)

- Commands live in `cmd/`
- API client in `internal/api/`
- Auth logic in `internal/auth/`
- Use Cobra for command structure
- Use survey for interactive prompts

### supascale.sh (Bash Script)

- Keep bash compatibility (avoid bashisms where possible)
- Use `jq` for JSON manipulation
- Test with shellcheck if modifying
- Maintain existing code style

## Getting Help

- **Questions?** Open a GitHub discussion
- **Stuck?** Comment on your issue/PR
- **Found a bug?** Create an issue

## Recognition

Contributors are recognized in:
- GitHub contributors list
- Release notes
- README acknowledgments (for significant contributions)

Thank you for contributing! 🎉
