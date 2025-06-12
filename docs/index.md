---
layout: default
title: "MarkIt - Extensible Markup Parser"
description: "Flexible markup parsing with configurable tag protocols. Parse XML, HTML, and custom markup formats with a single, extensible library."
keywords: "markit, xml parser, html parser, go parser, extensible markup, tag protocols"
author: "Khicago Team"
---

# MarkIt

> **Extensible Markup Parser for Go**

A Go library for parsing and rendering markup formats with configurable tag protocols - supporting XML, HTML, and custom markup formats with a unified API.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![Test Coverage](https://img.shields.io/badge/coverage-91.3%25-brightgreen.svg)](https://github.com/khicago/markit)
[![Go Report Card](https://goreportcard.com/badge/github.com/khicago/markit)](https://goreportcard.com/report/github.com/khicago/markit)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Documentation

- [Getting Started](getting-started)
- [API Reference](api-reference)
- [Configuration](configuration)
- [Renderer](renderer)
- [Memory Management](memory-management)
- [FAQ](faq)
- [Examples](examples)
- [Contributing](contributing)

## Quick Start

### Installation

```bash
go get github.com/khicago/markit
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/khicago/markit"
)

func main() {
    // XML content
    xml := `<books>
        <book id="1">
            <title>The Go Programming Language</title>
            <author>Alan Donovan</author>
        </book>
    </books>`
    
    // Parse
    parser := markit.NewParser(xml)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Access parsed content
    fmt.Printf("Parsed %d root elements\n", len(doc.Children))
    
    // Render back to markup
    renderer := markit.NewRenderer()
    output := renderer.Render(doc)
    fmt.Println("Rendered:", output)
}
```

### Parse and Render Workflow

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/khicago/markit"
)

func main() {
    // Original markup
    original := `<article>
        <h1>Hello World</h1>
        <p>This is a <strong>sample</strong> document.</p>
        <img src="photo.jpg" alt="A photo" />
    </article>`
    
    // Parse the document
    parser := markit.NewParser(original)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create renderer with pretty-printing options
    opts := &markit.RenderOptions{
        Indent:         "  ",    // 2-space indentation
        CompactMode:    false,   // Pretty-printed output
        SortAttributes: true,    // Consistent attribute order
        EscapeText:     true,    // Safe HTML output
    }
    renderer := markit.NewRendererWithOptions(opts)
    
    // Render the document
    output, err := renderer.RenderToString(doc)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Pretty-printed output:")
    fmt.Println(output)
    
    // Debug the AST structure
    fmt.Println("\nAST Debug View:")
    fmt.Println(markit.PrettyPrint(doc))
}
```

## Key Features

- **Unified Markup Support**: Parse HTML, XML, and similar markup formats with a consistent API
- **Configurable Parsing**: Adjust behavior for different markup standards
- **Good Performance**: Efficient lexer and parser implementation
- **Error Reporting**: Clear error messages with position information
- **Flexible Rendering**: Convert AST back to markup with formatting options
- **Debug Tools**: Built-in AST visualization for debugging
- **Validation Options**: Optional document validation capabilities
- **Writer Support**: Rendering to io.Writer interfaces

## Project Status

- [Milestone 1 Release Notes](milestone1_release_note) - Current implemented features
- [Milestone 2 Development Plan](milestone2_plan) - Upcoming features and enhancements

## Usage Examples

### HTML5 Document Parsing

```go
html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Sample Page</title>
</head>
<body>
    <h1>Welcome</h1>
    <p>Hello, <strong>world</strong>!</p>
</body>
</html>`

parser := markit.NewParser(html)
doc, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}

// Access document structure
fmt.Printf("Document has %d children\n", len(doc.Children))
```

### Document Rendering and Formatting

```go
// Parse a document
parser := markit.NewParser(`<div><p>Hello</p><p>World</p></div>`)
doc, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}

// Render with different formatting options
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

prettyRenderer := markit.NewRendererWithOptions(&markit.RenderOptions{
    Indent:      "  ",
    CompactMode: false,
    EscapeText:  true,
})
pretty, err := prettyRenderer.RenderToString(doc)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Pretty-printed:")
fmt.Println(pretty)
// Output:
// <div>
//   <p>Hello</p>
//   <p>World</p>
// </div>
```

### Custom Configuration

```go
config := &markit.Config{
    OpenBracket:  "{{",
    CloseBracket: "}}",
    // ... other options
}

parser := markit.NewParserWithConfig(content, config)
doc, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}
```

### AST Traversal and Debugging

```go
// Parse document
doc, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}

// Debug AST structure
fmt.Println("AST Structure:")
fmt.Println(markit.PrettyPrint(doc))

// Traverse nodes
markit.Walk(doc, func(node markit.Node) bool {
    if elem, ok := node.(*markit.Element); ok {
        fmt.Printf("Element: %s\n", elem.Name)
    }
    return true // continue traversal
})
```

## Links

- **GitHub**: [github.com/khicago/markit](https://github.com/khicago/markit)
- **Issues**: [Report bugs or request features](https://github.com/khicago/markit/issues)  
- **Discussions**: [Join the community discussion](https://github.com/khicago/markit/discussions)

## License

MarkIt is released under the MIT License. See [LICENSE](https://github.com/khicago/markit/blob/main/LICENSE) for details.

---

<div class="text-center">
<h3>Ready to get started?</h3>
<a href="https://github.com/khicago/markit" class="btn btn-outline">View on GitHub</a>
</div> 