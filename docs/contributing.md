---
layout: default
title: "Contributing - Join the MarkIt Community"
description: "Learn how to contribute to MarkIt parser - from setting up development environment to submitting pull requests."
keywords: "markit contributing, go open source, parser development, contribution guide, llm, ai"
author: "Khicago Team"
---

# Contributing to MarkIt

> **Help us build the next-generation markup parser for Go**

We welcome contributions from developers of all skill levels! This guide will help you get started with contributing to MarkIt.

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Code Standards](#code-standards)
- [Testing Guidelines](#testing-guidelines)
- [Submitting Changes](#submitting-changes)
- [Issue Guidelines](#issue-guidelines)
- [Documentation](#documentation)
- [Community](#community)

## Getting Started

### Prerequisites

- **Go 1.19+** - [Download Go](https://golang.org/dl/)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **Make** (optional) - For using Makefile commands

### Fork and Clone

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:

```bash
git clone https://github.com/YOUR_USERNAME/markit.git
cd markit
```

3. **Add upstream remote**:

```bash
git remote add upstream https://github.com/khicago/markit.git
```

4. **Verify your setup**:

```bash
go mod tidy
go test ./...
```

## Development Environment

### Project Structure

```
markit/
├── ast.go              # AST node definitions
├── config.go           # Parser configuration
├── lexer.go            # Lexical analysis
├── parser.go           # Main parser logic
├── visitor.go          # Visitor pattern implementation
├── void_elements.go    # Void elements support
├── examples/           # Example programs
├── docs/              # Documentation
├── tests/             # Test files
└── README.md          # Project overview
```

### Setting Up Your IDE

#### VS Code

Recommended extensions:
- **Go** (by Google)
- **Go Test Explorer**
- **GitLens**

#### GoLand/IntelliJ

Enable Go modules support and configure code style according to our standards.

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v -run TestVoidElements

# Run benchmarks
go test -bench=. -benchmem
```

### Building Examples

```bash
# Build all examples
cd examples
go build -o bin/ ./...

# Run specific example
go run void_elements_demo.go
```

## Code Standards

### Go Style Guide

We follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

#### Key Principles

1. **命名即文档** - Names should be self-documenting
2. **保持简单性** - Prefer simple, clear code over clever solutions
3. **避免技术债** - Clean up as you go, don't leave broken windows

#### Naming Conventions

```go
// ✅ Good - Clear, descriptive names
type ParserConfig struct {
    CaseSensitive     bool
    AllowSelfCloseTags bool
    VoidElements      map[string]bool
}

func (p *Parser) parseElement() (*Element, error) {
    // Implementation
}

// ❌ Bad - Generic, unclear names
type Config struct {
    CS   bool
    AST  bool
    VE   map[string]bool
}

func (p *Parser) parse() (*Node, error) {
    // Implementation
}
```

#### Error Handling

```go
// ✅ Good - Descriptive error messages
func (p *Parser) parseAttribute() (string, string, error) {
    if p.current.Type != TokenIdentifier {
        return "", "", fmt.Errorf("expected attribute name, got %s at position %d", 
            p.current.Type, p.current.Position)
    }
    // ...
}

// ❌ Bad - Generic error messages
func (p *Parser) parseAttribute() (string, string, error) {
    if p.current.Type != TokenIdentifier {
        return "", "", errors.New("parse error")
    }
    // ...
}
```

#### Documentation

```go
// ✅ Good - Clear documentation with examples
// ParseWithConfig creates a new parser with the specified configuration
// and parses the input content into an AST.
//
// Example:
//   config := markit.HTMLConfig()
//   parser := markit.NewParserWithConfig(content, config)
//   doc, err := parser.Parse()
func NewParserWithConfig(content string, config *ParserConfig) *Parser {
    // Implementation
}
```

### Code Formatting

Use `gofmt` and `goimports`:

```bash
# Format code
gofmt -w .

# Organize imports
goimports -w .
```

### Linting

We use `golangci-lint` for comprehensive linting:

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## Testing Guidelines

### Test Structure

```go
func TestParseElement(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        config   *ParserConfig
        expected *Element
        wantErr  bool
    }{
        {
            name:     "simple element",
            input:    "<div>content</div>",
            config:   DefaultConfig(),
            expected: &Element{TagName: "div", /* ... */},
            wantErr:  false,
        },
        {
            name:     "void element",
            input:    "<br>",
            config:   HTMLConfig(),
            expected: &Element{TagName: "br", SelfClosed: true},
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewParserWithConfig(tt.input, tt.config)
            result, err := parser.parseElement()
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected.TagName, result.TagName)
            // More assertions...
        })
    }
}
```

### Test Categories

#### Unit Tests
- Test individual functions and methods
- Mock dependencies when necessary
- Focus on edge cases and error conditions

#### Integration Tests
- Test complete parsing workflows
- Use real-world examples
- Verify end-to-end functionality

#### Benchmark Tests
```go
func BenchmarkParseHTML(b *testing.B) {
    content := generateLargeHTML(1000) // Helper function
    config := HTMLConfig()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        parser := NewParserWithConfig(content, config)
        _, err := parser.Parse()
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Coverage Requirements

- **Minimum coverage**: 90%
- **New features**: 95%+ coverage required
- **Critical paths**: 100% coverage required

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Submitting Changes

### Workflow

1. **Create a feature branch**:
```bash
git checkout -b feature/your-feature-name
```

2. **Make your changes** following our code standards

3. **Add tests** for new functionality

4. **Run the test suite**:
```bash
go test -v -cover ./...
golangci-lint run
```

5. **Commit your changes**:
```bash
git add .
git commit -m "feat: add support for custom void elements"
```

6. **Push to your fork**:
```bash
git push origin feature/your-feature-name
```

7. **Create a Pull Request** on GitHub

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

#### Examples
```
feat(parser): add support for HTML5 void elements

- Implement VoidElements configuration option
- Add HTMLConfig() with predefined void elements
- Update parser logic to handle void elements correctly

Closes #123
```

### Pull Request Guidelines

#### Before Submitting
- [ ] Tests pass locally
- [ ] Code follows style guidelines
- [ ] Documentation is updated
- [ ] CHANGELOG.md is updated (for significant changes)

#### PR Description Template
```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added for new functionality
```

## Issue Guidelines

### Bug Reports

Use the bug report template:

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Parse this content: '...'
2. With this configuration: '...'
3. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- Go version: [e.g. 1.21]
- MarkIt version: [e.g. v1.0.0]
- OS: [e.g. macOS 14.0]

**Additional context**
Any other context about the problem.
```

### Feature Requests

Use the feature request template:

```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Alternative solutions or features you've considered.

**Additional context**
Any other context or screenshots about the feature request.
```

### Security Issues

For security vulnerabilities, please email security@khicago.dev instead of creating a public issue.

## Documentation

### Types of Documentation

1. **Code Documentation** - Inline comments and godoc
2. **User Documentation** - Usage guides and examples
3. **API Documentation** - Complete API reference
4. **Contributing Documentation** - This guide

### Writing Guidelines

- Use clear, concise language
- Include practical examples
- Keep documentation up-to-date with code changes
- Use proper markdown formatting

### Building Documentation

```bash
# Generate godoc
godoc -http=:6060

# Build Jekyll site (for GitHub Pages)
cd docs
bundle install
bundle exec jekyll serve
```

## Community

### Communication Channels

- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - General questions and ideas
- **Email** - security@khicago.dev (security issues only)

### Code of Conduct

We are committed to providing a welcoming and inclusive environment. Please read our [Code of Conduct](CODE_OF_CONDUCT.md).

### Recognition

Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

## Development Tips

### Debugging

```go
// Add debug logging
func (p *Parser) parseElement() (*Element, error) {
    if p.config.Debug {
        log.Printf("Parsing element at position %d", p.position)
    }
    // Implementation
}
```

### Performance Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### Common Pitfalls

1. **Memory Leaks** - Always clean up resources
2. **Race Conditions** - Use proper synchronization
3. **Error Handling** - Don't ignore errors
4. **Testing** - Test edge cases and error conditions

## Getting Help

If you need help:

1. Check existing [documentation](/)
2. Search [GitHub Issues](https://github.com/khicago/markit/issues)
3. Ask in [GitHub Discussions](https://github.com/khicago/markit/discussions)
4. Read the [FAQ](faq)

---

## Thank You! 🙏

Thank you for contributing to MarkIt! Your efforts help make this project better for everyone.

---

<div align="center">

**[🏠 Back to Home](/)** • **[📚 Documentation](/docs)** • **[🐛 Report Issues](https://github.com/khicago/markit/issues)**

</div> 