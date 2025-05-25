---
layout: default
title: "Void Elements - HTML5 Self-Closing Tags"
description: "Complete guide to MarkIt's void elements support - HTML5 standard void elements, custom configurations, and XML compatibility."
keywords: "void elements, html5, self-closing tags, xml, go parser, markit, llm, ai"
author: "Khicago Team"
---

# Void Elements Support

> **Complete support for HTML5 void elements and custom configurations**

MarkIt provides comprehensive support for void elements - HTML tags that don't have closing tags and can't contain content. This feature enables seamless parsing of HTML5 documents while maintaining full XML compatibility.

## 📋 Table of Contents

- [What are Void Elements?](#what-are-void-elements)
- [HTML5 Standard Support](#html5-standard-support)
- [Configuration Options](#configuration-options)
- [Code Examples](#code-examples)
- [Best Practices](#best-practices)
- [XML Compatibility](#xml-compatibility)
- [Performance Considerations](#performance-considerations)
- [Troubleshooting](#troubleshooting)

## What are Void Elements?

Void elements are HTML tags that cannot have content and do not require closing tags. According to the HTML5 specification, these elements are self-closing by nature.

### HTML5 vs XML Syntax

```html
<!-- HTML5 Style (void elements) -->
<br>
<img src="photo.jpg" alt="Photo">
<input type="text" name="username">

<!-- XML Style (self-closing) -->
<br />
<img src="photo.jpg" alt="Photo" />
<input type="text" name="username" />
```

**Key Differences:**
- **HTML5 Style**: No closing slash, no closing tag
- **XML Style**: Closing slash, explicit self-closing syntax
- **MarkIt**: Supports both styles seamlessly

## HTML5 Standard Support

MarkIt includes built-in support for all HTML5 standard void elements:

| Element | Description | Common Use Case |
|---------|-------------|-----------------|
| `area` | Image map area | Interactive images |
| `base` | Document base URL | HTML head section |
| `br` | Line break | Text formatting |
| `col` | Table column | Table layout |
| `embed` | External content | Media embedding |
| `hr` | Horizontal rule | Section dividers |
| `img` | Image | Media content |
| `input` | Form input | User input |
| `link` | External resource | CSS/favicon links |
| `meta` | Document metadata | SEO and meta info |
| `param` | Object parameter | Deprecated elements |
| `source` | Media source | Audio/video sources |
| `track` | Text track | Video subtitles |
| `wbr` | Line break opportunity | Text wrapping |

### Enabling HTML5 Support

```go
package main

import "github.com/khicago/markit"

func main() {
    // Use HTML configuration for built-in void elements
    config := markit.HTMLConfig()
    
    // Verify void element support
    fmt.Printf("Supports <br>: %v\n", config.IsVoidElement("br"))     // true
    fmt.Printf("Supports <img>: %v\n", config.IsVoidElement("img"))   // true
    fmt.Printf("Supports <div>: %v\n", config.IsVoidElement("div"))   // false
}
```

## Configuration Options

### Default Configuration

```go
// Default configuration - no void elements
config := markit.DefaultConfig()
fmt.Printf("Supports void elements: %v", config.IsVoidElement("br")) // false

// This will fail because <br> expects a closing tag
parser := markit.NewParserWithConfig("<br>", config)
_, err := parser.Parse() // Error: expected closing tag
```

### HTML Configuration

```go
// HTML configuration with all standard void elements
config := markit.HTMLConfig()

// Parse HTML5 void elements
content := `<article>
    <h1>Article Title</h1>
    <p>First paragraph</p>
    <br>
    <img src="image.jpg" alt="Description">
    <hr>
    <p>Second paragraph</p>
</article>`

parser := markit.NewParserWithConfig(content, config)
ast, err := parser.Parse()
if err != nil {
    panic(err)
}
```

### Custom Void Elements

```go
// Start with default configuration
config := markit.DefaultConfig()

// Add custom void elements
config.SetVoidElements([]string{"my-icon", "my-separator", "custom-widget"})

// Or add individually
config.AddVoidElement("special-tag")

// Check if element is void
if config.IsVoidElement("my-icon") {
    fmt.Println("my-icon is configured as void element")
}

// Remove void element
config.RemoveVoidElement("special-tag")
```

### Dynamic Management

```go
config := markit.HTMLConfig()

// Runtime modification
config.AddVoidElement("custom-void")
config.RemoveVoidElement("wbr") // Remove standard void element

// Bulk operations
customElements := []string{"icon", "spacer", "divider"}
config.SetVoidElements(customElements) // Replaces all existing
```

## Code Examples

### Basic HTML5 Parsing

```go
package main

import (
    "fmt"
    "github.com/khicago/markit"
)

func parseHTML5() {
    config := markit.HTMLConfig()
    
    html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Test Page</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <img src="logo.png" alt="Company Logo">
    <hr>
    <form>
        <input type="text" name="username" placeholder="Username">
        <input type="password" name="password" placeholder="Password">
        <input type="submit" value="Login">
    </form>
</body>
</html>`

    parser := markit.NewParserWithConfig(html, config)
    ast, err := parser.Parse()
    if err != nil {
        panic(err)
    }
    
    // All void elements are correctly parsed as self-closing
    markit.Walk(ast, &markit.PrintVisitor{})
}
```

### Custom Void Elements for Components

```go
func parseCustomComponents() {
    config := markit.DefaultConfig()
    
    // Define custom void elements for UI components
    config.SetVoidElements([]string{
        "ui-icon",
        "ui-spacer", 
        "ui-divider",
        "ui-spinner",
        "ui-avatar",
    })
    
    component := `<ui-card>
        <ui-avatar user="john" size="large">
        <ui-spacer height="20">
        <h2>User Profile</h2>
        <ui-divider>
        <p>User information content</p>
        <ui-spinner loading="true">
    </ui-card>`
    
    parser := markit.NewParserWithConfig(component, config)
    ast, _ := parser.Parse()
    
    // Custom void elements are parsed correctly
    printElementInfo(ast)
}

func printElementInfo(node markit.Node) {
    switch n := node.(type) {
    case *markit.Element:
        selfClose := ""
        if n.SelfClose {
            selfClose = " (void)"
        }
        fmt.Printf("Element: <%s>%s\n", n.TagName, selfClose)
        
        for _, child := range n.Children {
            printElementInfo(child)
        }
    case *markit.Document:
        for _, child := range n.Children {
            printElementInfo(child)
        }
    }
}
```

### Mixed Content with Attributes

```go
func parseComplexHTML() {
    config := markit.HTMLConfig()
    
    content := `<article class="blog-post" id="post-123">
        <header>
            <h1>Advanced Web Development</h1>
            <meta property="article:author" content="Jane Doe">
            <meta property="article:published_time" content="2024-01-15">
        </header>
        
        <main>
            <p>Introduction paragraph with <br> line breaks.</p>
            
            <img src="diagram.svg" 
                 alt="Architecture Diagram" 
                 width="800" 
                 height="600" 
                 loading="lazy">
                 
            <hr class="section-divider">
            
            <form action="/subscribe" method="post">
                <input type="email" 
                       name="email" 
                       placeholder="Enter your email" 
                       required>
                <input type="hidden" name="source" value="blog">
                <input type="submit" value="Subscribe" class="btn-primary">
            </form>
        </main>
    </article>`
    
    parser := markit.NewParserWithConfig(content, config)
    ast, _ := parser.Parse()
    
    // Extract all void elements with attributes
    extractor := &VoidElementExtractor{}
    markit.Walk(ast, extractor)
    
    for _, elem := range extractor.VoidElements {
        fmt.Printf("Void element: <%s>\n", elem.TagName)
        for key, value := range elem.Attributes {
            fmt.Printf("  %s=\"%s\"\n", key, value)
        }
    }
}

type VoidElementExtractor struct {
    VoidElements []*markit.Element
}

func (v *VoidElementExtractor) VisitElement(elem *markit.Element) error {
    if elem.SelfClose {
        v.VoidElements = append(v.VoidElements, elem)
    }
    return nil
}

func (v *VoidElementExtractor) VisitDocument(doc *markit.Document) error { return nil }
func (v *VoidElementExtractor) VisitText(text *markit.Text) error { return nil }
func (v *VoidElementExtractor) VisitComment(comment *markit.Comment) error { return nil }
```

## Best Practices

### 1. Choose the Right Configuration

```go
// For HTML5 documents
config := markit.HTMLConfig()

// For XML documents  
config := markit.DefaultConfig()
config.AllowSelfCloseTags = true

// For custom markup languages
config := markit.DefaultConfig()
config.SetVoidElements([]string{"custom-void-1", "custom-void-2"})
```

### 2. Handle Case Sensitivity

```go
config := markit.HTMLConfig()
config.CaseSensitive = false // HTML is case-insensitive

// Both will be recognized as void elements
parser1 := markit.NewParserWithConfig("<BR>", config)
parser2 := markit.NewParserWithConfig("<br>", config)
```

### 3. Validate Void Element Usage

```go
func validateVoidElement(config *markit.ParserConfig, tagName string, hasContent bool) error {
    if config.IsVoidElement(tagName) && hasContent {
        return fmt.Errorf("void element <%s> cannot have content", tagName)
    }
    return nil
}
```

### 4. Performance Optimization

```go
// Pre-configure for better performance
var htmlConfig = markit.HTMLConfig()

func parseMultipleDocuments(documents []string) {
    for _, doc := range documents {
        // Reuse configuration instead of creating new ones
        parser := markit.NewParserWithConfig(doc, htmlConfig)
        ast, _ := parser.Parse()
        processAST(ast)
    }
}
```

## XML Compatibility

MarkIt maintains full compatibility with XML-style self-closing tags:

```go
config := markit.HTMLConfig()

// Mixed XML and HTML styles work together
mixed := `<document>
    <br />          <!-- XML style -->
    <br>            <!-- HTML style -->
    <img src="1.jpg" />  <!-- XML style -->
    <img src="2.jpg">    <!-- HTML style -->
</document>`

parser := markit.NewParserWithConfig(mixed, config)
ast, _ := parser.Parse()

// Both styles are parsed as self-closing elements
```

### XML-First Approach

```go
config := markit.DefaultConfig()
config.AllowSelfCloseTags = true
// Don't set void elements for strict XML parsing

// Only XML-style self-closing tags work
xmlContent := `<root>
    <element />
    <another-element />
</root>`
```

## Performance Considerations

### Configuration Reuse

```go
// ❌ Inefficient - creates new config each time
func parseDocuments(docs []string) {
    for _, doc := range docs {
        config := markit.HTMLConfig() // New config each iteration
        parser := markit.NewParserWithConfig(doc, config)
        // ... parse
    }
}

// ✅ Efficient - reuse configuration
var globalConfig = markit.HTMLConfig()

func parseDocuments(docs []string) {
    for _, doc := range docs {
        parser := markit.NewParserWithConfig(doc, globalConfig)
        // ... parse
    }
}
```

### Memory Usage

```go
// Void element detection is O(1) with map lookup
config := markit.HTMLConfig()

// Very fast - no string comparison loops
isVoid := config.IsVoidElement("img") // O(1) operation
```

## Troubleshooting

### Common Issues

**1. Void Element Not Recognized**
```go
// Problem: Using default config for HTML
config := markit.DefaultConfig()
parser := markit.NewParserWithConfig("<br>", config)
// Error: expected closing tag

// Solution: Use HTML config
config := markit.HTMLConfig()
parser := markit.NewParserWithConfig("<br>", config)
// Success: parsed as void element
```

**2. Case Sensitivity Issues**
```go
// Problem: Case-sensitive config with uppercase HTML
config := markit.HTMLConfig()
config.CaseSensitive = true
parser := markit.NewParserWithConfig("<BR>", config)
// May not recognize as void element

// Solution: Ensure consistent casing or use case-insensitive mode
config.CaseSensitive = false
```

**3. Custom Void Elements Not Working**
```go
// Problem: Adding to wrong configuration
config := markit.DefaultConfig()
config.AddVoidElement("my-element")
// But using different config instance for parsing

// Solution: Use same config instance
parser := markit.NewParserWithConfig(content, config)
```

### Debugging Tips

```go
func debugVoidElements(config *markit.ParserConfig, tagName string) {
    fmt.Printf("Is '%s' a void element? %v\n", 
        tagName, config.IsVoidElement(tagName))
    
    normalized := config.NormalizeCase(tagName)
    fmt.Printf("Normalized form: '%s'\n", normalized)
    fmt.Printf("Case sensitive: %v\n", config.CaseSensitive)
}

// Usage
config := markit.HTMLConfig()
debugVoidElements(config, "BR")
debugVoidElements(config, "br")
debugVoidElements(config, "custom-element")
```

## API Reference

### Configuration Methods

```go
// Check if element is void
IsVoidElement(tagName string) bool

// Add single void element
AddVoidElement(tagName string)

// Remove void element
RemoveVoidElement(tagName string)

// Set complete void elements list
SetVoidElements(elements []string)

// Normalize case based on configuration
NormalizeCase(s string) string
```

### Pre-built Configurations

```go
// Default configuration (no void elements)
config := markit.DefaultConfig()

// HTML5 configuration (all standard void elements)
config := markit.HTMLConfig()
```

---

## Next Steps

- [⚙️ Configuration Guide](configuration) - Explore all parser options
- [🌳 AST Traversal](ast-traversal) - Learn to walk and transform syntax trees
- [💡 Examples](examples) - See real-world usage examples
- [📚 API Reference](api-reference) - Complete API documentation

---

<div align="center">

**[🏠 Back to Home](/)** • **[📋 Report Issues](https://github.com/khicago/markit/issues)** • **[💬 Discussions](https://github.com/khicago/markit/discussions)**

</div> 