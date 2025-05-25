---
layout: default
title: "Getting Started - MarkIt Parser"
description: "Quick start guide for MarkIt - installation, basic usage, configuration, and your first parsing examples."
keywords: "markit getting started, go xml parser, html parser, installation guide, llm, ai"
author: "Khicago Team"
---

# Getting Started

> **Your journey to powerful markup parsing starts here**

This guide will get you up and running with MarkIt in minutes. From installation to your first parsed document, we'll cover everything you need to know.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Basic Usage](#basic-usage)
- [Configuration Basics](#configuration-basics)
- [Common Use Cases](#common-use-cases)
- [Error Handling](#error-handling)
- [Next Steps](#next-steps)

## Prerequisites

- **Go 1.22 or later** - MarkIt uses modern Go features
- **Basic Go knowledge** - Understanding of packages, interfaces, and error handling
- **Text editor or IDE** - VS Code, GoLand, or your favorite editor

### Verify Go Installation

```bash
$ go version
go version go1.22.0 linux/amd64  # Should be 1.22+
```

## Installation

### Using go get (Recommended)

```bash
go get github.com/khicago/markit
```

### Using go.mod

Add to your `go.mod` file:

```go
module your-project

go 1.22

require github.com/khicago/markit latest
```

Then run:

```bash
go mod download
```

### Verify Installation

Create a simple test file:

```go
package main

import (
    "fmt"
    "github.com/khicago/markit"
)

func main() {
    parser := markit.NewParser("<test>Hello MarkIt!</test>")
    ast, err := parser.Parse()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Parsed %d nodes\n", len(ast.Children))
}
```

Run it:

```bash
go run main.go
# Output: Parsed 1 nodes
```

## Basic Usage

### Your First Parse

```go
package main

import (
    "fmt"
    "github.com/khicago/markit"
)

func main() {
    // Simple XML content
    content := `<bookshelf>
        <book id="1" author="Jane Doe">
            <title>Go Programming</title>
            <year>2024</year>
            <!-- This is a comment -->
        </book>
        <book id="2" author="John Smith">
            <title>Advanced Parsing</title>
            <year>2023</year>
        </book>
    </bookshelf>`
    
    // Create parser and parse
    parser := markit.NewParser(content)
    ast, err := parser.Parse()
    if err != nil {
        panic(err)
    }
    
    // Access the root document
    fmt.Printf("Document has %d children\n", len(ast.Children))
    
    // Get the bookshelf element
    if len(ast.Children) > 0 {
        if bookshelf, ok := ast.Children[0].(*markit.Element); ok {
            fmt.Printf("Root element: <%s>\n", bookshelf.TagName)
            fmt.Printf("Number of books: %d\n", len(bookshelf.Children))
        }
    }
}
```

### Working with Attributes

```go
func exploreAttributes() {
    content := `<user id="123" name="Alice" active="true" age="25" />`
    
    parser := markit.NewParser(content)
    ast, _ := parser.Parse()
    
    if element, ok := ast.Children[0].(*markit.Element); ok {
        fmt.Printf("Element: <%s>\n", element.TagName)
        fmt.Printf("Self-closing: %v\n", element.SelfClose)
        
        // Access attributes
        for key, value := range element.Attributes {
            fmt.Printf("  %s = \"%s\"\n", key, value)
        }
        
        // Get specific attribute
        if id, exists := element.Attributes["id"]; exists {
            fmt.Printf("User ID: %s\n", id)
        }
    }
}
```

### Traversing the AST

```go
func walkExample() {
    content := `<article>
        <h1>Title</h1>
        <p>First paragraph</p>
        <p>Second paragraph</p>
        <!-- End of content -->
    </article>`
    
    parser := markit.NewParser(content)
    ast, _ := parser.Parse()
    
    // Use the built-in print visitor
    fmt.Println("=== AST Structure ===")
    markit.Walk(ast, &markit.PrintVisitor{})
    
    // Or create a custom visitor
    collector := &TextCollector{}
    markit.Walk(ast, collector)
    
    fmt.Println("\n=== Extracted Text ===")
    for _, text := range collector.Texts {
        fmt.Printf("Text: %q\n", text)
    }
}

// Custom visitor to collect all text content
type TextCollector struct {
    Texts []string
}

func (tc *TextCollector) VisitDocument(doc *markit.Document) error {
    return nil // Continue traversal
}

func (tc *TextCollector) VisitElement(elem *markit.Element) error {
    return nil // Continue traversal
}

func (tc *TextCollector) VisitText(text *markit.Text) error {
    tc.Texts = append(tc.Texts, text.Content)
    return nil
}

func (tc *TextCollector) VisitComment(comment *markit.Comment) error {
    return nil // Skip comments
}
```

## Configuration Basics

MarkIt is highly configurable. Here are the most common configuration patterns:

### Default Configuration

```go
// Default: XML-like behavior
config := markit.DefaultConfig()
parser := markit.NewParserWithConfig(content, config)

// Configuration properties:
// - CaseSensitive: true
// - AllowSelfCloseTags: true  
// - SkipComments: false
// - VoidElements: none
```

### HTML Configuration

```go
// HTML5: Case-insensitive with void elements
config := markit.HTMLConfig()
parser := markit.NewParserWithConfig(htmlContent, config)

// Includes all HTML5 void elements:
// br, hr, img, input, meta, link, etc.
```

### Custom Configuration

```go
config := &markit.ParserConfig{
    CaseSensitive:      false,  // HTML-style
    AllowSelfCloseTags: true,   // Allow <br/>
    SkipComments:       false,  // Include comments in AST
    VoidElements:       make(map[string]bool),
}

// Add custom void elements
config.SetVoidElements([]string{"my-component", "icon", "spacer"})

parser := markit.NewParserWithConfig(content, config)
```

## Common Use Cases

### 1. Parsing HTML Documents

```go
func parseHTML() {
    html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Sample Page</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <h1>Welcome</h1>
    <p>This is a <strong>sample</strong> page.</p>
    <img src="image.jpg" alt="Sample Image">
    <hr>
</body>
</html>`
    
    config := markit.HTMLConfig()
    parser := markit.NewParserWithConfig(html, config)
    ast, err := parser.Parse()
    
    if err != nil {
        fmt.Printf("Parse error: %v\n", err)
        return
    }
    
    // Extract page title
    titleExtractor := &TitleExtractor{}
    markit.Walk(ast, titleExtractor)
    fmt.Printf("Page title: %s\n", titleExtractor.Title)
}

type TitleExtractor struct {
    Title string
    inTitle bool
}

func (te *TitleExtractor) VisitElement(elem *markit.Element) error {
    if elem.TagName == "title" {
        te.inTitle = true
    }
    return nil
}

func (te *TitleExtractor) VisitText(text *markit.Text) error {
    if te.inTitle {
        te.Title = text.Content
        te.inTitle = false
    }
    return nil
}

func (te *TitleExtractor) VisitDocument(doc *markit.Document) error { return nil }
func (te *TitleExtractor) VisitComment(comment *markit.Comment) error { return nil }
```

### 2. Configuration Files

```go
func parseConfig() {
    configXML := `<configuration>
        <database>
            <host>localhost</host>
            <port>5432</port>
            <name>myapp</name>
        </database>
        <logging level="info" file="/var/log/app.log" />
        <features>
            <feature name="auth" enabled="true" />
            <feature name="cache" enabled="false" />
        </features>
    </configuration>`
    
    parser := markit.NewParser(configXML)
    ast, _ := parser.Parse()
    
    config := extractConfig(ast)
    fmt.Printf("Database: %s:%s/%s\n", 
        config.Database.Host, 
        config.Database.Port, 
        config.Database.Name)
}

type AppConfig struct {
    Database struct {
        Host string
        Port string
        Name string
    }
    Features map[string]bool
}

func extractConfig(ast *markit.Document) AppConfig {
    var config AppConfig
    config.Features = make(map[string]bool)
    
    // Implementation left as exercise
    // Traverse AST and populate config struct
    
    return config
}
```

### 3. Template Processing

```go
func processTemplate() {
    template := `<page>
        <header>
            <title>{{.Title}}</title>
            <meta name="author" content="{{.Author}}">
        </header>
        <content>
            {{range .Articles}}
            <article id="{{.ID}}">
                <h2>{{.Title}}</h2>
                <p>{{.Content}}</p>
            </article>
            {{end}}
        </content>
    </page>`
    
    parser := markit.NewParser(template)
    ast, _ := parser.Parse()
    
    // Process template variables
    processor := &TemplateProcessor{
        Variables: map[string]string{
            "Title":  "My Blog",
            "Author": "Jane Doe",
        },
    }
    
    markit.Walk(ast, processor)
}

type TemplateProcessor struct {
    Variables map[string]string
}

func (tp *TemplateProcessor) VisitText(text *markit.Text) error {
    // Process template variables in text content
    // Implementation depends on your template syntax
    return nil
}

func (tp *TemplateProcessor) VisitDocument(doc *markit.Document) error { return nil }
func (tp *TemplateProcessor) VisitElement(elem *markit.Element) error { return nil }
func (tp *TemplateProcessor) VisitComment(comment *markit.Comment) error { return nil }
```

## Error Handling

### Parse Errors

```go
func handleErrors() {
    // Malformed XML
    badContent := `<root>
        <unclosed>Content
        <mismatched>Content</wrong>
    </root>`
    
    parser := markit.NewParser(badContent)
    ast, err := parser.Parse()
    
    if err != nil {
        fmt.Printf("Parse failed: %v\n", err)
        
        // Errors include position information
        if parseErr, ok := err.(*markit.ParseError); ok {
            fmt.Printf("Error at line %d, column %d\n", 
                parseErr.Position.Line, 
                parseErr.Position.Column)
        }
        return
    }
    
    // Continue with successful parse
    fmt.Printf("Parsed successfully: %d nodes\n", len(ast.Children))
}
```

### Validation

```go
func validateDocument() {
    content := `<book>
        <title></title>  <!-- Empty title -->
        <author>Jane Doe</author>
        <isbn>123</isbn>  <!-- Invalid ISBN -->
    </book>`
    
    parser := markit.NewParser(content)
    ast, _ := parser.Parse()
    
    validator := &BookValidator{}
    err := markit.Walk(ast, validator)
    
    if err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Document is valid")
    }
}

type BookValidator struct{}

func (bv *BookValidator) VisitElement(elem *markit.Element) error {
    switch elem.TagName {
    case "title":
        // Check if title has content
        if len(elem.Children) == 0 {
            return fmt.Errorf("title cannot be empty")
        }
    case "isbn":
        // Validate ISBN format
        // Implementation details...
    }
    return nil
}

func (bv *BookValidator) VisitDocument(doc *markit.Document) error { return nil }
func (bv *BookValidator) VisitText(text *markit.Text) error { return nil }
func (bv *BookValidator) VisitComment(comment *markit.Comment) error { return nil }
```

## Performance Tips

### Reuse Configurations

```go
// ❌ Don't create new configs repeatedly
func badParsing(documents []string) {
    for _, doc := range documents {
        config := markit.HTMLConfig() // Wasteful
        parser := markit.NewParserWithConfig(doc, config)
        // ...
    }
}

// ✅ Reuse configuration
var sharedConfig = markit.HTMLConfig()

func goodParsing(documents []string) {
    for _, doc := range documents {
        parser := markit.NewParserWithConfig(doc, sharedConfig)
        // ...
    }
}
```

### Skip Unnecessary Features

```go
// For performance-critical parsing
config := markit.DefaultConfig()
config.SkipComments = true  // Skip comment nodes

// For HTML without case sensitivity needs
config.CaseSensitive = true  // Faster string comparisons
```

## Next Steps

Now that you've learned the basics, explore these advanced topics:

### 🧩 **Void Elements**
Learn about HTML5 void elements like `<br>`, `<img>`, and custom void elements.
**[Read the Void Elements Guide →](void-elements)**

### ⚙️ **Configuration**  
Deep dive into all configuration options and customization possibilities.
**[Explore Configuration →](configuration)**

### 🌳 **AST Traversal**
Master the visitor pattern and AST transformation techniques.
**[Learn AST Traversal →](ast-traversal)**

### 🔧 **Custom Protocols**
Create your own markup languages with custom tag bracket protocols.
**[Build Custom Protocols →](custom-protocols)**

### 💡 **Examples**
See real-world examples and advanced usage patterns.
**[Browse Examples →](examples)**

## Common Questions

### Q: Can MarkIt parse HTML with missing closing tags?

A: MarkIt is designed for well-formed markup. For HTML with missing tags, consider using a specialized HTML parser first, then MarkIt for processing.

### Q: How does MarkIt compare to encoding/xml?

A: MarkIt is more flexible and extensible, supporting custom protocols and advanced configuration. Use encoding/xml for simple XML needs, MarkIt for complex or custom markup.

### Q: Can I parse multiple documents concurrently?

A: Yes! MarkIt is thread-safe. You can parse different documents in parallel using separate parser instances.

```go
func concurrentParsing(documents []string) {
    var wg sync.WaitGroup
    results := make(chan *markit.Document, len(documents))
    
    for _, doc := range documents {
        wg.Add(1)
        go func(content string) {
            defer wg.Done()
            parser := markit.NewParser(content)
            ast, _ := parser.Parse()
            results <- ast
        }(doc)
    }
    
    wg.Wait()
    close(results)
}
```

---

## Get Help

- **📋 Issues**: [Report problems](https://github.com/khicago/markit/issues)
- **💬 Discussions**: [Ask questions](https://github.com/khicago/markit/discussions)  
- **📚 API Docs**: [Complete reference](https://pkg.go.dev/github.com/khicago/markit)

---

<div align="center">

**[🏠 Back to Home](/)** • **[🧩 Void Elements →](void-elements)** • **[⚙️ Configuration →](configuration)**

</div> 