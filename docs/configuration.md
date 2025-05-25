---
layout: default
title: "Configuration - Parser Settings"
description: "Complete guide to MarkIt parser configuration options - case sensitivity, void elements, comments, and custom settings."
keywords: "markit configuration, parser settings, go xml parser, html parser configuration"
author: "Khicago Team"
---

# Configuration Guide

> **Master MarkIt's flexible configuration system**

MarkIt's power lies in its configurability. This guide covers all configuration options, from basic settings to advanced customizations that let you parse any markup format.

## 📋 Table of Contents

- [Configuration Overview](#configuration-overview)
- [Pre-built Configurations](#pre-built-configurations)
- [Configuration Options](#configuration-options)
- [Custom Configurations](#custom-configurations)
- [Advanced Settings](#advanced-settings)
- [Performance Tuning](#performance-tuning)
- [Best Practices](#best-practices)

## Configuration Overview

MarkIt uses the `ParserConfig` struct to control parsing behavior. You can use pre-built configurations or create custom ones for specific needs.

### Basic Configuration Structure

```go
type ParserConfig struct {
    // Case sensitivity for tag and attribute names
    CaseSensitive bool
    
    // Allow XML-style self-closing tags (<tag />)
    AllowSelfCloseTags bool
    
    // Skip comment nodes in the AST
    SkipComments bool
    
    // HTML5 void elements that don't need closing tags
    VoidElements map[string]bool
    
    // Custom attribute processor
    AttributeProcessor AttributeProcessor
    
    // Protocol matcher for custom markup languages
    ProtocolMatcher *ProtocolMatcher
}
```

## Pre-built Configurations

### Default Configuration

```go
config := markit.DefaultConfig()

// Settings:
// - CaseSensitive: true
// - AllowSelfCloseTags: true
// - SkipComments: false
// - VoidElements: empty
// - AttributeProcessor: DefaultAttributeProcessor
// - ProtocolMatcher: standard XML protocols
```

**Best for:**
- XML documents
- Configuration files
- Structured data formats
- When you need strict parsing

### HTML Configuration

```go
config := markit.HTMLConfig()

// Settings:
// - CaseSensitive: false
// - AllowSelfCloseTags: true
// - SkipComments: false
// - VoidElements: all HTML5 void elements
// - AttributeProcessor: DefaultAttributeProcessor
// - ProtocolMatcher: standard XML protocols
```

**Best for:**
- HTML5 documents
- Web scraping
- Template processing
- When you need HTML5 compatibility

## Configuration Options

### Case Sensitivity

Controls whether tag and attribute names are case-sensitive.

```go
// Case-sensitive (XML style)
config := markit.DefaultConfig()
config.CaseSensitive = true

// <Tag> and <tag> are different elements
parser := markit.NewParserWithConfig(`<Tag><tag></tag></Tag>`, config)
ast, _ := parser.Parse()
```

```go
// Case-insensitive (HTML style)
config := markit.DefaultConfig()
config.CaseSensitive = false

// <Tag> and <tag> are the same element
parser := markit.NewParserWithConfig(`<Tag></tag>`, config)
ast, _ := parser.Parse() // Works fine
```

**Impact on Performance:**
- `CaseSensitive = true`: Faster string comparisons
- `CaseSensitive = false`: Requires string normalization

### Self-Closing Tags

Controls whether XML-style self-closing tags are allowed.

```go
// Allow self-closing tags
config := markit.DefaultConfig()
config.AllowSelfCloseTags = true

content := `<root>
    <item />
    <another-item />
</root>`

parser := markit.NewParserWithConfig(content, config)
ast, _ := parser.Parse() // Success
```

```go
// Disallow self-closing tags
config := markit.DefaultConfig()
config.AllowSelfCloseTags = false

content := `<root><item /></root>`
parser := markit.NewParserWithConfig(content, config)
_, err := parser.Parse() // Error: self-closing tags not allowed
```

### Comment Handling

Controls whether comments are included in the AST or skipped.

```go
// Include comments in AST
config := markit.DefaultConfig()
config.SkipComments = false

content := `<root>
    <!-- This comment will be in the AST -->
    <item>Content</item>
</root>`

parser := markit.NewParserWithConfig(content, config)
ast, _ := parser.Parse()
// AST includes comment nodes
```

```go
// Skip comments
config := markit.DefaultConfig()
config.SkipComments = true

content := `<root>
    <!-- This comment will be ignored -->
    <item>Content</item>
</root>`

parser := markit.NewParserWithConfig(content, config)
ast, _ := parser.Parse()
// AST does not include comment nodes
```

### Void Elements

Configure which elements are treated as void (self-closing without explicit syntax).

```go
// HTML5 void elements
config := markit.HTMLConfig()
// Includes: br, hr, img, input, meta, link, etc.

// Check if element is void
isVoid := config.IsVoidElement("br") // true
isVoid = config.IsVoidElement("div") // false
```

```go
// Custom void elements
config := markit.DefaultConfig()
config.SetVoidElements([]string{"icon", "spacer", "divider"})

// Add individual void element
config.AddVoidElement("separator")

// Remove void element
config.RemoveVoidElement("divider")

// Check support
supported := config.IsVoidElement("icon") // true
```

## Custom Configurations

### Creating Custom Configurations

```go
// Start with a base configuration
config := markit.DefaultConfig()

// Customize for your needs
config.CaseSensitive = false        // HTML-style
config.SkipComments = true          // Performance optimization
config.SetVoidElements([]string{    // Custom void elements
    "ui-icon",
    "ui-spacer",
    "ui-divider",
})

// Use the custom configuration
parser := markit.NewParserWithConfig(content, config)
```

### Configuration for Different Use Cases

#### Web Scraping Configuration

```go
func WebScrapingConfig() *markit.ParserConfig {
    config := markit.HTMLConfig()
    config.SkipComments = true  // Comments not needed for scraping
    return config
}

// Usage
config := WebScrapingConfig()
parser := markit.NewParserWithConfig(htmlContent, config)
```

#### Template Engine Configuration

```go
func TemplateConfig() *markit.ParserConfig {
    config := markit.DefaultConfig()
    config.CaseSensitive = false
    config.SkipComments = false  // Comments might contain template logic
    
    // Add custom template elements as void
    config.SetVoidElements([]string{
        "include",
        "import",
        "placeholder",
    })
    
    return config
}
```

#### Configuration File Parser

```go
func ConfigFileConfig() *markit.ParserConfig {
    config := markit.DefaultConfig()
    config.CaseSensitive = true   // Strict parsing
    config.SkipComments = false   // Comments are documentation
    // No void elements for config files
    return config
}
```

#### API Response Parser

```go
func APIResponseConfig() *markit.ParserConfig {
    config := markit.DefaultConfig()
    config.SkipComments = true    // APIs don't usually have comments
    config.CaseSensitive = true   // Strict field names
    return config
}
```

## Advanced Settings

### Custom Attribute Processor

Create custom attribute processing logic:

```go
type CustomAttributeProcessor struct{}

func (cap *CustomAttributeProcessor) ProcessAttributes(attrs map[string]string) map[string]string {
    processed := make(map[string]string)
    
    for key, value := range attrs {
        // Custom processing logic
        switch key {
        case "data-json":
            // Parse JSON attributes
            processed[key] = processJSONAttribute(value)
        case "style":
            // Normalize CSS styles
            processed[key] = normalizeCSS(value)
        default:
            processed[key] = value
        }
    }
    
    return processed
}

func (cap *CustomAttributeProcessor) ProcessAttribute(key, value string) (string, string) {
    // Individual attribute processing
    return key, value
}

func (cap *CustomAttributeProcessor) IsBooleanAttribute(key string) bool {
    // Define custom boolean attributes
    booleanAttrs := map[string]bool{
        "checked":   true,
        "disabled":  true,
        "readonly":  true,
        "required":  true,
        "selected":  true,
        "autofocus": true,
        "autoplay":  true,
        "controls":  true,
        "defer":     true,
        "hidden":    true,
        "loop":      true,
        "multiple":  true,
        "muted":     true,
        "open":      true,
        "reversed":  true,
        "scoped":    true,
        // Add custom boolean attributes
        "my-flag":   true,
        "enabled":   true,
    }
    
    return booleanAttrs[key]
}

// Usage
config := markit.DefaultConfig()
config.AttributeProcessor = &CustomAttributeProcessor{}
```

### Protocol Matcher Customization

For advanced users who want to support custom markup languages:

```go
// This is an advanced feature - most users won't need this
// See the Custom Protocols documentation for details

config := markit.DefaultConfig()
// Custom protocol matcher setup would go here
// This allows parsing of custom markup like {{...}} or [...] syntax
```

## Performance Tuning

### Memory Optimization

```go
// Minimize memory usage
config := markit.DefaultConfig()
config.SkipComments = true  // Reduces AST node count

// Reuse configurations
var globalConfig = markit.HTMLConfig()

func parseMultipleDocuments(docs []string) {
    for _, doc := range docs {
        // Reuse config instead of creating new ones
        parser := markit.NewParserWithConfig(doc, globalConfig)
        ast, _ := parser.Parse()
        processDocument(ast)
    }
}
```

### Speed Optimization

```go
// Optimize for speed
config := markit.DefaultConfig()
config.CaseSensitive = true    // Faster string comparisons
config.SkipComments = true     // Skip comment processing

// For HTML parsing where case doesn't matter much
config.CaseSensitive = false   // Slower but more flexible
```

### Batch Processing Configuration

```go
// Configuration for processing many documents
func BatchProcessingConfig() *markit.ParserConfig {
    config := markit.DefaultConfig()
    config.SkipComments = true     // Skip unnecessary nodes
    config.CaseSensitive = true    // Faster comparisons
    return config
}

// Pre-allocate and reuse
var batchConfig = BatchProcessingConfig()

func processBatch(documents []string) {
    for _, doc := range documents {
        parser := markit.NewParserWithConfig(doc, batchConfig)
        ast, _ := parser.Parse()
        // Process...
    }
}
```

## Best Practices

### 1. Choose the Right Base Configuration

```go
// For HTML documents
config := markit.HTMLConfig()

// For XML/structured data
config := markit.DefaultConfig()

// For custom markup
config := markit.DefaultConfig()
// Then customize as needed
```

### 2. Reuse Configurations

```go
// ❌ Don't create new configs repeatedly
func badExample(docs []string) {
    for _, doc := range docs {
        config := markit.HTMLConfig() // Wasteful
        parser := markit.NewParserWithConfig(doc, config)
        // ...
    }
}

// ✅ Reuse configurations
var sharedConfig = markit.HTMLConfig()

func goodExample(docs []string) {
    for _, doc := range docs {
        parser := markit.NewParserWithConfig(doc, sharedConfig)
        // ...
    }
}
```

### 3. Configure for Your Use Case

```go
// Web scraping - skip comments, case insensitive
webConfig := markit.HTMLConfig()
webConfig.SkipComments = true

// API parsing - strict, fast
apiConfig := markit.DefaultConfig()
apiConfig.CaseSensitive = true
apiConfig.SkipComments = true

// Template processing - flexible, preserve comments
templateConfig := markit.DefaultConfig()
templateConfig.CaseSensitive = false
templateConfig.SkipComments = false
```

### 4. Validate Configuration

```go
func validateConfig(config *markit.ParserConfig) error {
    if config == nil {
        return fmt.Errorf("configuration cannot be nil")
    }
    
    if config.VoidElements == nil {
        config.VoidElements = make(map[string]bool)
    }
    
    if config.AttributeProcessor == nil {
        config.AttributeProcessor = &markit.DefaultAttributeProcessor{}
    }
    
    return nil
}
```

### 5. Document Your Configuration

```go
// DocumentConfig creates a configuration optimized for
// parsing documentation files with custom elements
func DocumentConfig() *markit.ParserConfig {
    config := markit.DefaultConfig()
    config.CaseSensitive = false  // Flexible tag names
    config.SkipComments = false   // Preserve documentation comments
    
    // Add documentation-specific void elements
    config.SetVoidElements([]string{
        "toc",        // Table of contents placeholder
        "pagebreak",  // Page break marker
        "include",    // File inclusion marker
    })
    
    return config
}
```

## Configuration Examples

### Complete HTML5 Setup

```go
func setupHTML5Parser(content string) (*markit.Document, error) {
    config := markit.HTMLConfig()
    
    // Optional customizations
    config.SkipComments = false  // Keep comments for debugging
    
    // Add custom void elements for web components
    config.AddVoidElement("my-icon")
    config.AddVoidElement("my-spacer")
    
    parser := markit.NewParserWithConfig(content, config)
    return parser.Parse()
}
```

### Strict XML Parser

```go
func setupXMLParser(content string) (*markit.Document, error) {
    config := markit.DefaultConfig()
    
    // Strict XML settings
    config.CaseSensitive = true
    config.AllowSelfCloseTags = true
    config.SkipComments = false
    
    // No void elements in strict XML
    config.VoidElements = make(map[string]bool)
    
    parser := markit.NewParserWithConfig(content, config)
    return parser.Parse()
}
```

### Performance-Optimized Parser

```go
func setupFastParser(content string) (*markit.Document, error) {
    config := markit.DefaultConfig()
    
    // Optimize for speed
    config.CaseSensitive = true  // Faster comparisons
    config.SkipComments = true   // Skip comment processing
    
    parser := markit.NewParserWithConfig(content, config)
    return parser.Parse()
}
```

---

## Next Steps

- [🧩 Void Elements](void-elements) - Learn about HTML5 void elements
- [🌳 AST Traversal](ast-traversal) - Master AST walking and transformation
- [🔧 Custom Protocols](custom-protocols) - Create custom markup languages
- [💡 Examples](examples) - See real-world configuration examples

---

<div align="center">

**[🏠 Back to Home](/)** • **[📋 Report Issues](https://github.com/khicago/markit/issues)** • **[💬 Discussions](https://github.com/khicago/markit/discussions)**

</div> 