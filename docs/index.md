---
layout: default
title: "MarkIt - Next-Generation Extensible Markup Parser"
description: "Revolutionary markup parsing with configurable tag bracket protocols. Parse XML, HTML, and any custom markup format with a single, extensible parser."
keywords: "markit, xml parser, html parser, go parser, extensible markup, tag protocols, llm, ai"
author: "Khicago Team"
---

# MarkIt

> **The Next-Generation Extensible Markup Parser for Go**

Revolutionary markup parsing with configurable tag bracket protocols - Parse XML, HTML, and any custom markup format with a single, extensible parser.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![Test Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen.svg)](https://github.com/khicago/markit)
[![Go Report Card](https://goreportcard.com/badge/github.com/khicago/markit)](https://goreportcard.com/report/github.com/khicago/markit)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 🚀 Quick Start

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
}
```

## 🌟 Key Features

### Universal Markup Support
Parse **XML**, **HTML**, and **custom markup** formats with one unified API.

### Configurable Protocols  
Define your own bracket sequences: `<>`, `{{}}`, `[]`, or any custom delimiters.

### High Performance
- Zero-copy parsing for maximum speed
- Memory-efficient token streaming
- 100% test coverage

### Developer-Friendly
- Rich error reporting with line/column positions
- Type-safe AST with visitor pattern
- Comprehensive documentation

## 📚 Usage Examples

### HTML5 Document Parsing

```go
package main

import (
    "fmt"
    "github.com/khicago/markit"
)

func main() {
    html := `<!DOCTYPE html>
    <html>
        <head><title>My Page</title></head>
        <body>
            <h1>Welcome</h1>
            <p>Hello <strong>world</strong>!</p>
            <br>
            <img src="image.jpg" alt="Photo">
        </body>
    </html>`
    
    // Use HTML configuration for proper HTML5 parsing
    config := markit.DefaultHTMLConfig()
    parser := markit.NewParserWithConfig(html, config)
    doc, err := parser.Parse()
    if err != nil {
        panic(err)
    }
    
    // Find all paragraph elements
    paragraphs := findElementsByTag(doc.Root, "p")
    fmt.Printf("Found %d paragraphs\n", len(paragraphs))
}

func findElementsByTag(element *markit.Element, tagName string) []*markit.Element {
    var results []*markit.Element
    
    if element.TagName == tagName {
        results = append(results, element)
    }
    
    for _, child := range element.Children {
        if childEl, ok := child.(*markit.Element); ok {
            results = append(results, findElementsByTag(childEl, tagName)...)
        }
    }
    
    return results
}
```

### Custom Configuration

```go
func advancedConfiguration() {
    content := `<template>
        <component name="header">
            <slot>Default content</slot>
        </component>
    </template>`
    
    // Create custom configuration
    config := &markit.ParserConfig{
        CaseSensitive:      false,                     // HTML-style case insensitivity
        AllowSelfCloseTags: true,                      // Support <br/> style tags
        SkipComments:       true,                      // Skip comment nodes for performance
        TrimWhitespace:     true,                      // Normalize whitespace
        AttributeProcessor: markit.DefaultAttributeProcessor, // Custom attribute handling
    }
    
    parser := markit.NewParserWithConfig(content, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Process the parsed document
    processTemplate(doc)
}
```

### AST Traversal and Transformation

```go
// Custom visitor for extracting links
type LinkExtractor struct {
    Links []Link
}

type Link struct {
    URL  string
    Text string
}

func (le *LinkExtractor) VisitElement(element *markit.Element) error {
    if element.TagName == "a" {
        href, hasHref := element.Attributes["href"]
        if hasHref {
            link := Link{URL: href}
            
            // Extract text content
            markit.Walk(doc.Root, func(node markit.Node) bool {
                if element, ok := node.(*markit.Element); ok {
                    fmt.Printf("Element: <%s>\n", element.TagName)
                }
                return true // Continue walking
            })
            
            le.Links = append(le.Links, link)
        }
    }
    return nil
}

// Error handling with detailed information
func handleParsingErrors(content string) {
    parser := markit.NewParser(content)
    doc, err := parser.Parse()
    
    if err != nil {
        if parseErr, ok := err.(*markit.ParseError); ok {
            fmt.Printf("Parse error at line %d, column %d: %s\n", 
                parseErr.Line, parseErr.Column, parseErr.Message)
        } else {
            fmt.Printf("General error: %s\n", err.Error())
        }
        return
    }
    
    // Process successful parse result
    fmt.Printf("Successfully parsed document with %d elements\n", 
        countElements(doc.Root))
}
```

## 🔗 Links

- **GitHub**: [github.com/khicago/markit](https://github.com/khicago/markit)
- **Issues**: [Report bugs or request features](https://github.com/khicago/markit/issues)  
- **Discussions**: [Join the community discussion](https://github.com/khicago/markit/discussions)

## 📄 License

MarkIt is released under the MIT License. See [LICENSE](https://github.com/khicago/markit/blob/main/LICENSE) for details.

---

<div class="text-center">
<h3>Ready to get started?</h3>
<a href="https://github.com/khicago/markit" class="btn btn-outline">View on GitHub</a>
</div> 