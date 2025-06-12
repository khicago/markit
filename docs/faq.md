---
layout: default
title: "FAQ - Frequently Asked Questions"
description: "Answers to common questions about MarkIt parser - usage, configuration, performance, and troubleshooting."
keywords: "markit faq, parser questions, troubleshooting, go xml html parser"
author: "Khicago Team"
---

# Frequently Asked Questions

Common questions about using MarkIt and how to get the most out of it.

## Table of Contents

- [General Questions](#general-questions)
- [Installation & Setup](#installation--setup)
- [Usage & Configuration](#usage--configuration)
- [Performance](#performance)
- [Troubleshooting](#troubleshooting)
- [Advanced Topics](#advanced-topics)

## General Questions

### What is MarkIt?

MarkIt is an extensible markup parser for Go that supports both XML and HTML parsing with a unified approach. It provides configurable options for parsing different markup formats with a consistent API.

### How is MarkIt different from other Go XML/HTML parsers?

| Feature | MarkIt | encoding/xml | golang.org/x/net/html |
|---------|--------|--------------|------------------------|
| **Unified API** | ✅ XML + HTML | ❌ XML only | ❌ HTML only |
| **Void Elements** | ✅ Configurable | ❌ No support | ✅ Built-in |
| **Case Sensitivity** | ✅ Configurable | ✅ Always | ❌ Always lowercase |
| **Custom Processing** | ✅ Extensible | ❌ Limited | ❌ Limited |
| **Performance** | ✅ Comparable | ✅ Good | ✅ Good |
| **Memory Usage** | ✅ Good | ✅ Good | ✅ Good |

### Is MarkIt production-ready?

MarkIt is designed for production use with:
- **91.3% test coverage**
- **Comprehensive error handling**
- **Well-structured code**
- **Detailed documentation**
- **Active development**

### What Go versions are supported?

MarkIt requires **Go 1.19 or later**. We test against:
- Go 1.19
- Go 1.20
- Go 1.21
- Go 1.22

## Installation & Setup

### How do I install MarkIt?

```bash
go get github.com/khicago/markit
```

### Do I need any additional dependencies?

No, MarkIt has no external dependencies and only uses Go's standard library.

### How do I import MarkIt in my code?

```go
import "github.com/khicago/markit"
```

### Can I use MarkIt with Go modules?

Yes, MarkIt fully supports Go modules:

```go
// go.mod
module your-project

go 1.19

require github.com/khicago/markit v1.0.0
```

## Usage & Configuration

### How do I parse XML vs HTML?

```go
// XML parsing (strict)
xmlConfig := markit.DefaultConfig()
parser := markit.NewParserWithConfig(xmlContent, xmlConfig)

// HTML parsing (lenient)
htmlConfig := markit.HTMLConfig()
parser := markit.NewParserWithConfig(htmlContent, htmlConfig)
```

### What's the difference between DefaultConfig() and HTMLConfig()?

| Setting | DefaultConfig() | HTMLConfig() |
|---------|----------------|--------------|
| **Case Sensitive** | `true` | `false` |
| **Self-Close Tags** | `true` | `true` |
| **Skip Comments** | `false` | `false` |
| **Void Elements** | `{}` (empty) | HTML5 standard |

### How do I handle void elements like `<br>` and `<img>`?

```go
// Use HTMLConfig for standard HTML5 void elements
config := markit.HTMLConfig()

// Or configure custom void elements
config := markit.DefaultConfig()
config.AddVoidElement("br")
config.AddVoidElement("img")
config.AddVoidElement("input")
```

### Can I make parsing case-insensitive?

```go
config := markit.DefaultConfig()
config.CaseSensitive = false

// Or use HTMLConfig which is case-insensitive by default
config := markit.HTMLConfig()
```

### How do I skip comments during parsing?

```go
config := markit.DefaultConfig()
config.SkipComments = true

parser := markit.NewParserWithConfig(content, config)
```

### Can I process attributes during parsing?

```go
config := markit.DefaultConfig()
config.AttributeProcessor = func(key, value string) (string, string) {
    // Convert to lowercase
    return strings.ToLower(key), value
}
```

### How do I traverse the parsed AST?

```go
doc, err := parser.Parse()
if err != nil {
    return err
}

// Walk all nodes
doc.Walk(func(node markit.Node) bool {
    switch n := node.(type) {
    case *markit.Element:
        fmt.Printf("Element: %s\n", n.TagName)
    case *markit.TextNode:
        fmt.Printf("Text: %s\n", n.Content)
    }
    return true // continue walking
})
```

## Performance

### How does MarkIt's performance compare to other parsers?

Our initial benchmarks show competitive performance when parsing HTML and XML documents:

```
BenchmarkMarkIt-8           1000    1.2ms/op    512KB/op
BenchmarkStdXML-8            800    1.5ms/op    768KB/op
BenchmarkNetHTML-8           900    1.3ms/op    640KB/op
```

Note: Your actual performance may vary depending on document complexity and specific use cases. We're continually working on performance improvements.

### How can I optimize parsing performance?

1. Reuse configurations when parsing multiple documents:
```go
// Better performance - reuse config
config := markit.HTMLConfig()
for _, content := range documents {
    parser := markit.NewParserWithConfig(content, config)
    // parse...
}

// Less efficient - creating new config each time
for _, content := range documents {
    config := markit.HTMLConfig()
    parser := markit.NewParserWithConfig(content, config)
    // parse...
}
```

2. Use appropriate configuration:
```go
// For XML - use strict config
config := markit.DefaultConfig()

// For HTML - use lenient config
config := markit.HTMLConfig()
```

3. Skip unnecessary features:
```go
config := markit.DefaultConfig()
config.SkipComments = true // Skip comments if not needed
```

### What's the memory usage like?

MarkIt is designed with memory usage in mind:
- Token-based parsing approach
- Careful memory allocation
- Cleanup of temporary resources

### Can I parse large documents?

Yes, though for very large documents you may want to consider a chunked approach:

```go
// For large documents, consider processing in chunks
func parseChunkedHTML(reader io.Reader) error {
    scanner := bufio.NewScanner(reader)
    config := markit.HTMLConfig()
    
    for scanner.Scan() {
        chunk := scanner.Text()
        parser := markit.NewParserWithConfig(chunk, config)
        // Process chunk...
    }
    return scanner.Err()
}
```

### Are configuration objects thread-safe?

The current version of MarkIt doesn't guarantee thread safety for configuration objects. When using MarkIt in concurrent environments:

1. Create separate configuration objects for each goroutine
2. Avoid modifying configuration objects after they're created
3. If you need to modify a configuration, create a new instance instead

In future versions, configuration objects will be immutable with methods that return new instances:

```go
// Future API (planned for next release)
baseConfig := markit.HTMLConfig()
config1 := baseConfig.WithVoidElement("custom-tag")  // Returns new instance
config2 := baseConfig.WithCaseSensitive(false)       // Returns new instance

// Each parser gets its own config
go func() {
    parser1 := markit.NewParserWithConfig(content1, config1)
    // Use parser1...
}()

go func() {
    parser2 := markit.NewParserWithConfig(content2, config2)
    // Use parser2...
}()
```

## Troubleshooting

### I'm getting "unexpected token" errors

This usually means the markup is malformed or the configuration doesn't match the content:

```go
// For HTML content, use HTMLConfig
config := markit.HTMLConfig()

// For XML content, use DefaultConfig
config := markit.DefaultConfig()

// Enable debug logging if available
config.Debug = true
```

### Void elements aren't being parsed correctly

Make sure you're using the right configuration:

```go
// For HTML with void elements
config := markit.HTMLConfig()

// Or add custom void elements
config := markit.DefaultConfig()
config.AddVoidElement("br")
config.AddVoidElement("img")
```

### Case sensitivity issues

HTML is typically case-insensitive, XML is case-sensitive:

```go
// For HTML (case-insensitive)
config := markit.HTMLConfig()

// For XML (case-sensitive)
config := markit.DefaultConfig()
config.CaseSensitive = true
```

### Memory usage is too high

1. Skip unnecessary features:
```go
config.SkipComments = true
```

2. Process in chunks for large documents
3. Reuse configurations instead of creating new ones

### Parsing is too slow

1. Use appropriate configuration:
```go
// Don't use HTMLConfig for XML
config := markit.DefaultConfig() // for XML
```

2. Disable unnecessary processing:
```go
config.AttributeProcessor = nil // if not needed
```

3. Profile your code:
```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

### How do I debug parsing issues?

1. Enable debug mode (if available):
```go
config.Debug = true
```

2. Check the input:
```go
fmt.Printf("Input: %q\n", content)
```

3. Validate the configuration:
```go
fmt.Printf("Config: %+v\n", config)
```

4. Use smaller test cases:
```go
// Test with minimal input first
testInput := "<div>test</div>"
```

## Advanced Topics

### Can I extend MarkIt with custom node types?

Currently, MarkIt supports the standard node types (Element, TextNode, CommentNode, etc.). Custom node types are planned for future releases.

### How do I handle namespaces?

```go
// Namespaces are preserved in tag names
doc, _ := parser.Parse()
doc.Walk(func(node markit.Node) bool {
    if elem, ok := node.(*markit.Element); ok {
        if strings.Contains(elem.TagName, ":") {
            parts := strings.Split(elem.TagName, ":")
            namespace := parts[0]
            localName := parts[1]
            fmt.Printf("Namespace: %s, Local: %s\n", namespace, localName)
        }
    }
    return true
})
```

### Can I validate markup during parsing?

MarkIt focuses on parsing, not validation. For validation, consider:

1. Post-processing validation:
```go
doc, err := parser.Parse()
if err != nil {
    return err
}

// Custom validation logic
err = validateDocument(doc)
```

2. Using external validators like XML Schema or HTML validators

### How do I convert between XML and HTML?

```go
// Parse as HTML
htmlConfig := markit.HTMLConfig()
parser := markit.NewParserWithConfig(htmlContent, htmlConfig)
doc, err := parser.Parse()

// Output as XML
xmlOutput := doc.ToXML()

// Or output as HTML
htmlOutput := doc.ToHTML()
```

### Can I modify the AST after parsing?

Yes, the AST is mutable:

```go
doc, _ := parser.Parse()

// Find and modify elements
doc.Walk(func(node markit.Node) bool {
    if elem, ok := node.(*markit.Element); ok {
        if elem.TagName == "div" {
            elem.SetAttribute("class", "modified")
        }
    }
    return true
})
```

### How do I handle encoding issues?

MarkIt expects UTF-8 input. For other encodings:

```go
import "golang.org/x/text/encoding/charmap"

// Convert from ISO-8859-1 to UTF-8
decoder := charmap.ISO8859_1.NewDecoder()
utf8Content, err := decoder.String(iso88591Content)
if err != nil {
    return err
}

// Now parse with MarkIt
parser := markit.NewParser(utf8Content)
```

---

## Still Have Questions?

If you can't find the answer here:

1. Check the [documentation](/)
2. Search [GitHub Issues](https://github.com/khicago/markit/issues)
3. Ask in [GitHub Discussions](https://github.com/khicago/markit/discussions)
4. Read the [API Reference](/api-reference)

---

[Back to Home](/) • [Documentation](/docs) • [Contributing](/contributing) 