---
layout: default
title: "MarkIt - Fast XML/HTML Parser for Go"
description: "A high-performance, memory-efficient XML and HTML parser for Go with a simple API and extensive configuration options."
keywords: "Go, XML, HTML, parser, AST, performance, memory-efficient, llm, ai"
author: "Khicago Team"
---

# MarkIt Parser

<div class="text-center mb-2">
    <p style="font-size: 1.2rem; color: #666;">A fast, memory-efficient XML and HTML parser for Go</p>
    <div style="margin: 2rem 0;">
        <span class="badge badge-success">Go 1.19+</span>
        <span class="badge badge-info">Zero Dependencies</span>
        <span class="badge badge-warning">High Performance</span>
    </div>
</div>

## 🚀 Features

- **High Performance**: Optimized for speed and memory efficiency
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **Flexible Configuration**: Support for both XML and HTML parsing modes
- **Rich AST**: Complete Abstract Syntax Tree with full node information
- **Easy to Use**: Simple, intuitive API that's easy to learn and use
- **Extensible**: Custom attribute processors and protocol matchers
- **Well Tested**: Comprehensive test suite with 99.5% code coverage
- **Production Ready**: Used in production environments

## 📦 Installation

```bash
go get github.com/khicago/xmlite
```

## 🏃‍♂️ Quick Start

### Basic XML Parsing

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/khicago/xmlite"
)

func main() {
    xml := `
    <bookstore>
        <book id="1" category="fiction">
            <title>The Great Gatsby</title>
            <author>F. Scott Fitzgerald</author>
            <price currency="USD">12.99</price>
        </book>
    </bookstore>`
    
    // Parse with default XML configuration
    parser := xmlite.NewParser(xml)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Access the root element
    root := doc.Root
    fmt.Printf("Root element: %s\n", root.Name)
    
    // Find book elements
    books := root.FindElements("book")
    for _, book := range books {
        title := book.FindElement("title")
        author := book.FindElement("author")
        price := book.FindElement("price")
        
        fmt.Printf("Book: %s by %s - %s %s\n",
            title.Text(),
            author.Text(),
            price.GetAttribute("currency"),
            price.Text())
    }
}
```

### HTML5 Parsing

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/khicago/xmlite"
)

func main() {
    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>My Page</title>
        <meta charset="utf-8">
    </head>
    <body>
        <h1>Welcome</h1>
        <p>This is a <strong>sample</strong> page.</p>
        <img src="image.jpg" alt="Sample Image">
    </body>
    </html>`
    
    // Parse with HTML5 configuration
    config := xmlite.DefaultHTMLConfig()
    parser := xmlite.NewParserWithConfig(html, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Find all images
    images := doc.Root.FindElementsRecursive("img")
    for _, img := range images {
        src := img.GetAttribute("src")
        alt := img.GetAttribute("alt")
        fmt.Printf("Image: %s (%s)\n", src, alt)
    }
}
```

## 🎯 Use Cases

MarkIt is perfect for a wide range of applications:

### Web Scraping
Extract data from HTML pages with ease:
```go
// Parse product information
products := doc.Root.FindElementsByClass("product")
for _, product := range products {
    name := product.FindElement("h2").Text()
    price := product.FindElementByClass("price").Text()
    // Process product data...
}
```

### Configuration Files
Parse XML configuration files:
```go
config := doc.Root.FindElement("database")
host := config.GetAttribute("host")
port := config.GetAttribute("port")
```

### API Response Processing
Handle XML API responses:
```go
users := doc.Root.FindElements("user")
for _, user := range users {
    id := user.GetAttribute("id")
    name := user.FindElement("name").Text()
    email := user.FindElement("email").Text()
}
```

### Template Processing
Build template engines and content processors:
```go
// Process template variables
variables := doc.Root.FindElementsRecursive("var")
for _, v := range variables {
    name := v.GetAttribute("name")
    value := v.Text()
    // Replace template variables...
}
```

## 🔧 Configuration Options

MarkIt provides flexible configuration options for different parsing needs:

| Configuration | Best For | Features |
|---------------|----------|----------|
| `DefaultXMLConfig()` | XML documents | Strict parsing, case-sensitive |
| `DefaultHTMLConfig()` | HTML documents | Void elements, case-insensitive |
| Custom Config | Specific needs | Full control over parsing behavior |

### Custom Configuration Example

```go
config := &xmlite.ParserConfig{
    CaseSensitive:      false,
    AllowSelfCloseTags: true,
    SkipComments:       false,
    VoidElements:       []string{"br", "hr", "img", "input"},
    AttributeProcessor: xmlite.DefaultAttributeProcessor,
}

parser := xmlite.NewParserWithConfig(content, config)
```

## 📊 Performance

MarkIt is designed for high performance:

- **Memory Efficient**: Minimal memory allocation during parsing
- **Fast Parsing**: Optimized tokenization and AST construction
- **Scalable**: Handles large documents efficiently
- **Concurrent Safe**: Safe for use in concurrent applications

### Benchmark Results

```
BenchmarkParseXML-8     	   10000	    120543 ns/op	   45234 B/op	     892 allocs/op
BenchmarkParseHTML-8    	    8000	    145678 ns/op	   52341 B/op	    1023 allocs/op
```

## 🛠️ Advanced Features

### Custom Attribute Processing
```go
customProcessor := func(key, value string) (string, string) {
    // Custom attribute processing logic
    return strings.ToLower(key), strings.TrimSpace(value)
}

config.AttributeProcessor = customProcessor
```

### AST Traversal
```go
// Walk the entire AST
xmlite.Walk(doc.Root, func(node xmlite.Node) bool {
    if element, ok := node.(*xmlite.Element); ok {
        fmt.Printf("Element: %s\n", element.Name)
    }
    return true // Continue traversal
})
```

### Error Handling
```go
doc, err := parser.Parse()
if err != nil {
    if parseErr, ok := err.(*xmlite.ParseError); ok {
        fmt.Printf("Parse error at line %d, column %d: %s\n",
            parseErr.Line, parseErr.Column, parseErr.Message)
    }
}
```

## 📚 Documentation

- **[Getting Started](getting-started)** - Installation and basic usage
- **[Configuration](configuration)** - Detailed configuration options
- **[API Reference](api-reference)** - Complete API documentation
- **[Examples](examples)** - Practical usage examples
- **[FAQ](faq)** - Frequently asked questions
- **[Contributing](contributing)** - How to contribute to the project

## 🤝 Community

- **GitHub**: [github.com/khicago/xmlite](https://github.com/khicago/xmlite)
- **Issues**: [Report bugs or request features](https://github.com/khicago/xmlite/issues)
- **Discussions**: [Join the community discussion](https://github.com/khicago/xmlite/discussions)

## 📄 License

MarkIt is released under the MIT License. See [LICENSE](https://github.com/khicago/xmlite/blob/main/LICENSE) for details.

---

<div class="text-center mt-2">
    <a href="getting-started" class="btn">Get Started</a>
    <a href="https://github.com/khicago/xmlite" class="btn btn-outline">View on GitHub</a>
</div> 