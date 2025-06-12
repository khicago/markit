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

### Thread Safety Considerations

Configuration objects in MarkIt are not inherently thread-safe. When using MarkIt in concurrent environments, follow these guidelines:

```go
// AVOID: Sharing and modifying the same config in multiple goroutines
sharedConfig := markit.HTMLConfig()

go func() {
    // Unsafe - may cause race conditions
    sharedConfig.VoidElements["custom-tag"] = true
    parser1 := markit.NewParserWithConfig(content1, sharedConfig)
    // Use parser1...
}()

go func() {
    // Unsafe - concurrent modification
    sharedConfig.CaseSensitive = true
    parser2 := markit.NewParserWithConfig(content2, sharedConfig)
    // Use parser2...
}()
```

```go
// RECOMMENDED: Create separate configs for each goroutine
baseConfig := markit.HTMLConfig()

go func() {
    // Create a copy for this goroutine
    config1 := *baseConfig // Create a copy
    config1.VoidElements["custom-tag"] = true
    parser1 := markit.NewParserWithConfig(content1, &config1)
    // Use parser1...
}()

go func() {
    // Create a copy for this goroutine
    config2 := *baseConfig // Create a copy
    config2.CaseSensitive = true
    parser2 := markit.NewParserWithConfig(content2, &config2)
    // Use parser2...
}()
```

### Configuration Immutability

While the current API allows modifying configuration objects, it's best to treat them as immutable:

```go
// AVOID: Modifying configuration after parser creation
config := markit.HTMLConfig()
parser := markit.NewParserWithConfig(content, config)

// Don't do this - may lead to unexpected behavior
config.CaseSensitive = true
ast, _ := parser.Parse() // Parsing behavior depends on when the config is read
```

```go
// RECOMMENDED: Create a new configuration for each distinct use case
htmlConfig := markit.HTMLConfig()
parser1 := markit.NewParserWithConfig(content1, htmlConfig)

// For different settings, create a new config
customConfig := markit.HTMLConfig()
customConfig.CaseSensitive = true
customConfig.SkipComments = true
parser2 := markit.NewParserWithConfig(content2, customConfig)
```

### Future Immutable API

In upcoming releases, MarkIt will adopt an immutable configuration pattern:

```go
// Future API (planned for next release)
baseConfig := markit.HTMLConfig()

// Each method returns a new instance with the modified setting
htmlConfig := baseConfig.
    WithCaseSensitive(false).
    WithVoidElement("custom-element").
    WithSkipComments(true)

// Original config remains unchanged
// baseConfig != htmlConfig
```

### Performance Optimization

// ... existing code ...

## Configuration Examples

### Complete HTML5 Setup

```