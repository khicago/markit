# MarkIt

**Extensible Markup Parser & Renderer for Go**

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![Test Coverage](https://img.shields.io/badge/coverage-91.3%25-brightgreen.svg)](https://github.com/khicago/markit)
[![Go Report Card](https://goreportcard.com/badge/github.com/khicago/markit)](https://goreportcard.com/report/github.com/khicago/markit)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![GitHub Release](https://img.shields.io/github/release/khicago/markit.svg)](https://github.com/khicago/markit/releases)
[![GitHub Stars](https://img.shields.io/github/stars/khicago/markit.svg)](https://github.com/khicago/markit/stargazers)

> **Markup parsing and rendering with configurable tag protocols** - Parse, transform, and render XML, HTML, and custom markup formats with a single, extensible library.

## Why MarkIt?

Most parsers are designed for specific markup languages with fixed behavior. MarkIt provides a flexible approach with its **Tag Protocol** system and **Rendering Engine**, allowing you to work with various tag-based syntaxes.

```go
// Parse markup content
parser := markit.NewParser(input)
ast, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}

// Render with formatting options
renderer := markit.NewRendererWithOptions(&markit.RenderOptions{
    Indent:         "  ",
    SortAttributes: true,
})
output, err := renderer.RenderToString(ast)
if err != nil {
    log.Fatal(err)
}
```

## Documentation

For comprehensive documentation, please visit:

- **[Complete Documentation](docs/)** - API reference and guides
- **[Contributing Guide](docs/contributing.md)** - How to contribute 
- **[Changelog](docs/CHANGELOG.md)** - Version history
- **[FAQ](docs/faq.md)** - Frequently asked questions
- **[Memory Management](docs/memory-management.md)** - Optimization strategies

## Installation & Quick Start

```bash
go get github.com/khicago/markit
```

```go
package main

import (
    "fmt"
    "log"
    "github.com/khicago/markit"
)

func main() {
    content := `<root>
        <item id="1">Hello World</item>
        <!-- Comments work too -->
    </root>`
    
    parser := markit.NewParser(content)
    ast, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    renderer := markit.NewRendererWithOptions(&markit.RenderOptions{
        Indent:         "  ",
        SortAttributes: true,
    })
    
    output, err := renderer.RenderToString(ast)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(output)
}
```

## Core Features

- **Universal Markup Support**: Parse XML, HTML, and custom formats with a unified API
- **Configurable Parsing**: Adjust behavior for different markup standards
- **Smart Lexical Analysis**: Efficient and accurate tokenization
- **Professional Rendering Engine**: Configurable formatting options for multiple styles
- **Error Reporting**: Clear error messages with position information
- **Memory Efficiency**: Optimized for performance with minimal allocations
- **Thread Safety**: Guidelines for concurrent parsing and rendering
- **Visitor Pattern**: Easy AST traversal and transformation
- **Validation Options**: Well-formedness and encoding validation
- **Void Elements Support**: Complete HTML5 void elements and custom configurations

## Usage Examples

### HTML5 Document Parsing

```go
// Use HTML configuration for HTML5 documents
config := markit.HTMLConfig()
parser := markit.NewParserWithConfig(htmlContent, config)
doc, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}

// Access document structure
title := doc.Root.FindDescendantByTag("title")
if title != nil && len(title.Children) > 0 {
    fmt.Println("Title:", title.Children[0].String())
}
```

### Document Rendering with Formatting Options

```go
// Parse a document
parser := markit.NewParser(`<div><p>Hello</p><p>World</p></div>`)
doc, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}

// Compact rendering
compactRenderer := markit.NewRendererWithOptions(&markit.RenderOptions{
    CompactMode: true,
    EscapeText:  true,
})
compact, err := compactRenderer.RenderToString(doc)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Compact:", compact)
// Output: <div><p>Hello</p><p>World</p></div>

// Pretty-printed rendering
prettyRenderer := markit.NewRendererWithOptions(&markit.RenderOptions{
    Indent:      "  ",
    CompactMode: false,
    EscapeText:  true,
})
pretty, err := prettyRenderer.RenderToString(doc)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Pretty-printed:\n", pretty)
// Output:
// <div>
//   <p>Hello</p>
//   <p>World</p>
// </div>
```

### Custom Void Elements

```go
// Define custom void elements for UI components
config := markit.DefaultConfig()
config.SetVoidElements([]string{
    "ui-icon",
    "ui-spacer", 
    "ui-divider",
})

content := `<ui-card>
    <ui-avatar user="john" size="large">
    <ui-spacer height="20">
    <h2>User Profile</h2>
    <ui-divider>
    <p>User information content</p>
</ui-card>`

parser := markit.NewParserWithConfig(content, config)
ast, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}
```

### Memory-Efficient Processing for Large Documents

```go
func processLargeDocument(reader io.Reader) error {
    scanner := bufio.NewScanner(reader)
    buffer := bytes.Buffer{}
    count := 0
    
    for scanner.Scan() {
        line := scanner.Text()
        buffer.WriteString(line)
        buffer.WriteString("\n")
        count++
        
        // Process in chunks of 1000 lines
        if count >= 1000 {
            if err := processChunk(buffer.String()); err != nil {
                return err
            }
            buffer.Reset()
            count = 0
        }
    }
    
    // Process any remaining content
    if buffer.Len() > 0 {
        if err := processChunk(buffer.String()); err != nil {
            return err
        }
    }
    
    return scanner.Err()
}

func processChunk(content string) error {
    parser := markit.NewParser(content)
    doc, err := parser.Parse()
    if err != nil {
        return err
    }
    // Process document...
    return nil
}
```

## Performance Benchmarks

| Feature | MarkIt | Standard XML | HTML Parser |
|---------|--------|--------------|-------------|
| **Parsing Speed** | Fast | Medium | Medium |
| **Rendering Speed** | Fast | Slow | Slow |
| **Memory Usage** | Minimal | Moderate | Moderate |
| **Flexibility** | High | XML Only | HTML Only |
| **Streaming** | Supported | Limited | Limited |

```bash
# Run benchmarks
go test -bench=. -benchmem
```

## Architecture

### Tag Protocol System

MarkIt's approach is based on configurable tag protocols that define how tags are recognized:

```go
// Core protocols define how tags are identified
coreProtocols := []markit.CoreProtocol{
    {
        Name:     "standard-tag",
        OpenSeq:  "<",
        CloseSeq: ">",
        SelfClose: "/",
    },
    {
        Name:     "comment",
        OpenSeq:  "<!--",
        CloseSeq: "-->",
    },
}
```

### Rendering Engine

```go
type RenderOptions struct {
    Indent             string            // Indentation string
    EscapeText         bool              // Text escaping
    PreserveSpace      bool              // Whitespace handling
    CompactMode        bool              // Compact vs pretty-print
    SortAttributes     bool              // Attribute ordering
    EmptyElementStyle  EmptyElementStyle // Element closing style
    IncludeDeclaration bool              // Declaration inclusion
}
```

## Use Cases

- **Web Development**: Parse HTML, extract metadata, transform markup
- **Document Processing**: Convert between formats, extract structured data
- **Template Engines**: Custom template syntax, macro expansion
- **API Integration**: Parse XML responses, transform data formats
- **Enterprise Applications**: Large document processing with streaming

## Testing & Quality

MarkIt maintains **91.3% test coverage** with comprehensive test suites:

```bash
# Run tests with coverage
go test -v -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Contributing

We welcome contributions! Here's how to get started:

```bash
# Clone the repository
git clone https://github.com/khicago/markit.git
cd markit

# Install dependencies
go mod download

# Run tests
go test -v ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with ❤️ by the Go community
- Special thanks to all [contributors](https://github.com/khicago/markit/contributors)

---

<div align="center">

**[⭐ Star us on GitHub](https://github.com/khicago/markit)** • **[📖 Read the Docs](https://pkg.go.dev/github.com/khicago/markit)** • **[💬 Join Discussions](https://github.com/khicago/markit/discussions)**

Made with ❤️ for the Go community

</div> 