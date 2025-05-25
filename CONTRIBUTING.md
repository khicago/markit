# Contributing to MarkIt

We love your input! We want to make contributing to MarkIt as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Pull Requests

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git

### Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/markit.git
cd markit

# Add upstream remote
git remote add upstream https://github.com/khicago/markit.git

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...
```

## Code Style

We use standard Go formatting and linting tools:

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run all checks
make check
```

### Code Quality Standards

- **100% test coverage** is required for all new code
- All code must pass `go vet` and `golangci-lint`
- Follow Go best practices and idioms
- Write clear, self-documenting code
- Add comprehensive tests for edge cases

## Testing

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v -run TestSpecificFunction

# Run benchmarks
go test -bench=. -benchmem ./...
```

### Writing Tests

- Write table-driven tests when appropriate
- Test edge cases and error conditions
- Use descriptive test names
- Include benchmarks for performance-critical code

Example test structure:

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "<tag>content</tag>",
            expected: "content",
            wantErr:  false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParseFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("ParseFunction() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Documentation

### Code Documentation

- All public functions and types must have Go doc comments
- Comments should explain the "why", not just the "what"
- Include examples in doc comments when helpful

### README Updates

- Update README.md if you change public APIs
- Add examples for new features
- Update performance benchmarks if applicable

## Commit Messages

Write clear, descriptive commit messages:

```
feat: add support for custom bracket protocols

- Implement configurable open/close sequences
- Add protocol registration system
- Include comprehensive tests

Fixes #123
```

### Commit Message Format

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

## Issue Reporting

### Bug Reports

Use the bug report template and include:

- Go version
- Operating system
- Minimal reproduction case
- Expected vs actual behavior
- Stack trace if applicable

### Feature Requests

Use the feature request template and include:

- Clear description of the problem
- Proposed solution
- Alternative solutions considered
- Additional context

## Performance Considerations

- Profile code changes that might affect performance
- Include benchmarks for new features
- Avoid breaking changes to public APIs
- Consider memory allocation patterns

## Release Process

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release PR
4. Tag release after merge
5. GitHub Actions will handle the rest

## Community

- Be respectful and inclusive
- Help others learn and grow
- Share knowledge and best practices
- Follow the [Go Code of Conduct](https://golang.org/conduct)

## Questions?

Feel free to open an issue or start a discussion if you have questions about contributing!

## License

By contributing, you agree that your contributions will be licensed under the MIT License. 