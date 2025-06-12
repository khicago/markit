---
title: "Renderer - MarkIt Documentation"
description: "Complete guide to the MarkIt renderer system for converting AST back to markup"
keywords: ["renderer", "markup", "AST", "formatting", "validation"]
author: "MarkIt Team"
---

# Renderer

The MarkIt renderer system provides powerful capabilities for converting Abstract Syntax Trees (AST) back into formatted markup. It offers extensive customization options for output formatting, validation, and different rendering styles.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Renderer Types](#renderer-types)
- [Configuration Options](#configuration-options)
- [Rendering Methods](#rendering-methods)
- [Validation](#validation)
- [Advanced Usage](#advanced-usage)
- [Best Practices](#best-practices)

## Overview

The renderer system consists of several components:

- **Renderer**: Main rendering engine with customizable options
- **DebugRenderer**: Specialized renderer for AST debugging and visualization
- **RenderOptions**: Configuration structure for controlling rendering behavior
- **ValidationOptions**: Settings for document validation during rendering

## Quick Start

### Basic Rendering

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/khicago/markit"
)

func main() {
    // Parse a document
    parser := markit.NewParser(`<div><p>Hello World</p></div>`)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a renderer and render the document
    renderer := markit.NewRenderer()
    output := renderer.Render(doc)
    fmt.Println(output)
    // Output: <div><p>Hello World</p></div>
}
```

### Pretty-Printed Output

```go
// Create renderer with formatting options
opts := &markit.RenderOptions{
    Indent:      "  ",    // 2-space indentation
    CompactMode: false,   // Enable pretty-printing
    EscapeText:  true,    // Escape HTML entities
}

renderer := markit.NewRendererWithOptions(opts)
output, err := renderer.RenderToString(doc)
if err != nil {
    log.Fatal(err)
}

fmt.Println(output)
// Output:
// <div>
//   <p>Hello World</p>
// </div>
```

## Renderer Types

### Standard Renderer

The main renderer for production use:

```go
// Default renderer
renderer := markit.NewRenderer()

// Renderer with custom options
opts := &markit.RenderOptions{
    Indent:         "    ",  // 4-space indentation
    CompactMode:    false,   // Pretty-printed output
    SortAttributes: true,    // Sort attributes alphabetically
    EscapeText:     true,    // Escape special characters
}
renderer := markit.NewRendererWithOptions(opts)

// Renderer with configuration
config := &markit.Config{
    CaseSensitive: false,
    // ... other config options
}
renderer := markit.NewRendererWithConfig(config)
```

### Debug Renderer

Specialized renderer for AST visualization:

```go
// Create debug renderer
debugRenderer := markit.NewDebugRenderer()

// Render AST structure for debugging
debugOutput := debugRenderer.RenderToString(doc)
fmt.Println(debugOutput)
// Output shows the AST structure with indentation
```

## Configuration Options

### RenderOptions

```go
type RenderOptions struct {
    // Indentation string (e.g., "  ", "    ", "\t")
    Indent string
    
    // Whether to render in compact mode (no formatting)
    CompactMode bool
    
    // Whether to escape text content
    EscapeText bool
    
    // Whether to sort attributes alphabetically
    SortAttributes bool
    
    // Style for rendering empty elements
    EmptyElementStyle EmptyElementStyle
    
    // Whether to exclude XML/HTML declarations
    ExcludeDeclaration bool
    
    // Validation options
    Validation *ValidationOptions
}
```

### Empty Element Styles

```go
// Self-closing style: <img />
EmptyElementSelfClosing

// Paired tag style: <img></img>
EmptyElementPairedTag

// Void element style: <img> (HTML5 void elements)
EmptyElementVoid
```

### Validation Options

```go
type ValidationOptions struct {
    // Check for well-formed XML/HTML
    CheckWellFormed bool
    
    // Validate UTF-8 encoding
    CheckEncoding bool
    
    // Custom validation rules
    CustomRules []ValidationRule
}
```

## Rendering Methods

### Document Rendering

```go
// Render entire document
output := renderer.Render(doc)

// Render to string with error handling
output, err := renderer.RenderToString(doc)
if err != nil {
    log.Fatal(err)
}

// Render to writer (streaming)
var buf bytes.Buffer
err := renderer.RenderToWriter(doc, &buf)
if err != nil {
    log.Fatal(err)
}
```

### Element Rendering

```go
// Render specific element
element := doc.Children[0].(*markit.Element)
output := renderer.RenderElement(element)

// Render element to writer
var buf bytes.Buffer
err := renderer.RenderElementToWriter(element, &buf)
if err != nil {
    log.Fatal(err)
}
```

### Validation During Rendering

```go
// Render with validation
validationOpts := &markit.ValidationOptions{
    CheckWellFormed: true,
    CheckEncoding:   true,
}

renderer.SetValidation(validationOpts)
output, err := renderer.RenderWithValidation(doc)
if err != nil {
    if validationErr, ok := err.(*markit.ValidationError); ok {
        fmt.Printf("Validation error: %s\n", validationErr.Message)
    }
}
```

## Advanced Usage

### Custom Formatting

```go
// HTML5-style rendering
htmlRenderer := markit.NewRendererWithOptions(&markit.RenderOptions{
    Indent:            "  ",
    CompactMode:       false,
    EscapeText:        true,
    SortAttributes:    true,
    EmptyElementStyle: markit.EmptyElementVoid,
})

// XML-style rendering
xmlRenderer := markit.NewRendererWithOptions(&markit.RenderOptions{
    Indent:            "  ",
    CompactMode:       false,
    EscapeText:        true,
    SortAttributes:    true,
    EmptyElementStyle: markit.EmptyElementSelfClosing,
})
```

### Streaming Large Documents

```go
// For large documents, use streaming to avoid memory issues
file, err := os.Create("output.xml")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

renderer := markit.NewRenderer()
err = renderer.RenderToWriter(doc, file)
if err != nil {
    log.Fatal(err)
}
```

### Custom Validation

```go
// Define custom validation rules
customRule := func(node markit.Node) error {
    if elem, ok := node.(*markit.Element); ok {
        if elem.Name == "script" {
            return fmt.Errorf("script tags not allowed")
        }
    }
    return nil
}

validationOpts := &markit.ValidationOptions{
    CheckWellFormed: true,
    CustomRules:     []markit.ValidationRule{customRule},
}

renderer.SetValidation(validationOpts)
```

### AST Debugging

```go
// Use PrettyPrint for quick AST visualization
fmt.Println("AST Structure:")
fmt.Println(markit.PrettyPrint(doc))

// Use DebugRenderer for detailed formatting
debugRenderer := markit.NewDebugRenderer()
debugOutput := debugRenderer.RenderToString(doc)
fmt.Println("Detailed AST:")
fmt.Println(debugOutput)
```

## Best Practices

### Performance Optimization

1. **Reuse Renderers**: Create renderer instances once and reuse them
2. **Use Streaming**: For large documents, render directly to writers
3. **Disable Validation**: Skip validation in production if not needed
4. **Choose Appropriate Options**: Use compact mode for minimal output

```go
// Good: Reuse renderer
renderer := markit.NewRenderer()
for _, doc := range documents {
    output := renderer.Render(doc)
    // process output
}

// Good: Streaming for large documents
renderer.RenderToWriter(largeDoc, outputFile)
```

### Error Handling

```go
// Always handle rendering errors
output, err := renderer.RenderToString(doc)
if err != nil {
    if validationErr, ok := err.(*markit.ValidationError); ok {
        log.Printf("Validation failed: %s", validationErr.Message)
    } else {
        log.Printf("Rendering failed: %s", err.Error())
    }
    return
}
```

### Configuration Management

```go
// Define standard configurations
var (
    HTMLConfig = &markit.RenderOptions{
        Indent:            "  ",
        CompactMode:       false,
        EscapeText:        true,
        EmptyElementStyle: markit.EmptyElementVoid,
    }
    
    XMLConfig = &markit.RenderOptions{
        Indent:            "  ",
        CompactMode:       false,
        EscapeText:        true,
        EmptyElementStyle: markit.EmptyElementSelfClosing,
        SortAttributes:    true,
    }
    
    CompactConfig = &markit.RenderOptions{
        CompactMode: true,
        EscapeText:  true,
    }
)

// Use predefined configurations
htmlRenderer := markit.NewRendererWithOptions(HTMLConfig)
xmlRenderer := markit.NewRendererWithOptions(XMLConfig)
```

### Testing Rendered Output

```go
func TestRendering(t *testing.T) {
    parser := markit.NewParser(`<div><p>Test</p></div>`)
    doc, err := parser.Parse()
    require.NoError(t, err)
    
    renderer := markit.NewRendererWithOptions(&markit.RenderOptions{
        CompactMode: true,
    })
    
    output := renderer.Render(doc)
    assert.Equal(t, `<div><p>Test</p></div>`, output)
}
```

## Integration Examples

### Web Server Integration

```go
func renderHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request content
    parser := markit.NewParser(requestContent)
    doc, err := parser.Parse()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Render with validation
    renderer := markit.NewRendererWithOptions(&markit.RenderOptions{
        EscapeText: true,
        Validation: &markit.ValidationOptions{
            CheckWellFormed: true,
        },
    })
    
    err = renderer.RenderToWriter(doc, w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
```

### Template Processing

```go
func processTemplate(templateContent string, data interface{}) (string, error) {
    // Parse template
    parser := markit.NewParser(templateContent)
    doc, err := parser.Parse()
    if err != nil {
        return "", err
    }
    
    // Process template with data
    processedDoc := processTemplateData(doc, data)
    
    // Render final output
    renderer := markit.NewRendererWithOptions(&markit.RenderOptions{
        Indent:      "  ",
        CompactMode: false,
        EscapeText:  true,
    })
    
    return renderer.RenderToString(processedDoc)
}
```

The MarkIt renderer provides a comprehensive solution for converting AST back to markup with full control over formatting, validation, and output style. Whether you need compact output for production or pretty-printed markup for development, the renderer system has you covered. 