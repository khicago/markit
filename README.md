# MarkIt 🚀

**The Next-Generation Extensible Markup Parser for Go**

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.22-blue.svg)](https://golang.org/)
[![Test Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen.svg)](https://github.com/khicago/markit)
[![Go Report Card](https://goreportcard.com/badge/github.com/khicago/markit)](https://goreportcard.com/report/github.com/khicago/markit)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![GitHub Release](https://img.shields.io/github/release/khicago/markit.svg)](https://github.com/khicago/markit/releases)
[![GitHub Stars](https://img.shields.io/github/stars/khicago/markit.svg)](https://github.com/khicago/markit/stargazers)

> **Revolutionary markup parsing with configurable tag bracket protocols** - Parse XML, HTML, and any custom markup format with a single, extensible parser.

## 🌟 Why MarkIt?

Traditional parsers lock you into specific markup languages. **MarkIt breaks free** with its groundbreaking **Tag Bracket Protocol** system, allowing you to parse any tag-based syntax through simple configuration.

```go
// One parser, infinite possibilities
parser := markit.NewParser(input)
ast, _ := parser.Parse()

// Works with XML, HTML, custom formats, and more!
```

## ⚡ Quick Start

### 📚 Documentation

For comprehensive documentation, please visit:

- **[Complete Documentation](docs/)** - Full API reference and guides
- **[Contributing Guide](docs/contributing.md)** - How to contribute to the project  
- **[Changelog](docs/CHANGELOG.md)** - Version history and release notes
- **[FAQ](docs/faq.md)** - Frequently asked questions

### Installation

```bash
go get github.com/khicago/markit
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/khicago/markit"
)

func main() {
    // Parse any markup format
    content := `<root>
        <item id="1">Hello World</item>
        <!-- Comments work too -->
    </root>`
    
    parser := markit.NewParser(content)
    ast, err := parser.Parse()
    if err != nil {
        panic(err)
    }
    
    // Traverse with visitor pattern
    markit.Walk(ast, &markit.PrintVisitor{})
}
```

## 🔥 Core Features

### 🎯 **Universal Markup Support**
- **XML**: Full support with namespaces, CDATA, DOCTYPE
- **HTML**: Case-insensitive, self-closing tags, boolean attributes  
- **Custom Formats**: Define your own `{{...}}`, `[...]`, or any bracket syntax

### ⚡ **High Performance**
- **Zero-copy parsing** for maximum efficiency
- **100% test coverage** with comprehensive edge case handling
- **Minimal memory footprint** with smart token streaming

### 🔧 **Extensible Architecture**
- **Tag Bracket Protocols**: Configure `<open...close>` sequences
- **Pluggable processors**: Custom attribute handling
- **Visitor pattern**: Flexible AST traversal and transformation

### 📍 **Developer Experience**
- **Precise error reporting** with line/column positions
- **Rich AST nodes** with full position tracking
- **Type-safe APIs** with comprehensive documentation

## 🚀 Advanced Examples

### Custom Markup Language

```go
// Create a template engine syntax: {{variable}}
config := markit.DefaultConfig()
config.CaseSensitive = false

// Parse template syntax
content := `<div>{{user.name}} - {{user.email}}</div>`
parser := markit.NewParserWithConfig(content, config)
ast, _ := parser.Parse()
```

### Configuration-Driven Parsing

```go
// Flexible configuration for different use cases
config := &markit.ParserConfig{
    CaseSensitive:      false,           // HTML-style
    AllowSelfCloseTags: true,            // <br/> support
    SkipComments:       true,            // Ignore comments
    TrimWhitespace:     true,            // Clean output
}

parser := markit.NewParserWithConfig(input, config)
```

### AST Transformation

```go
// Custom visitor for AST transformation
type LinkExtractor struct {
    Links []string
}

func (v *LinkExtractor) VisitElement(elem *markit.Element) error {
    if elem.TagName == "a" {
        if href, ok := elem.Attributes["href"]; ok {
            v.Links = append(v.Links, href)
        }
    }
    return nil
}

// Extract all links from HTML
extractor := &LinkExtractor{}
markit.Walk(ast, extractor)
fmt.Println("Found links:", extractor.Links)
```

## 🏗️ Architecture

### Tag Bracket Protocol System

MarkIt's revolutionary approach centers on **configurable tag bracket protocols**:

```go
type CoreProtocol struct {
    Name        string    // "xml-tag", "html-comment", etc.
    OpenSeq     string    // "<", "<!--", "<?", etc.
    CloseSeq    string    // ">", "-->", "?>", etc.
    TokenType   TokenType // How to interpret the content
}
```

### Built-in Protocols

| Protocol | Open | Close | Use Case |
|----------|------|-------|----------|
| `markit-standard-tag` | `<` | `>` | XML/HTML elements |
| `markit-comment` | `<!--` | `-->` | Comments |

### Extensible Configuration

```go
type ParserConfig struct {
    CaseSensitive      bool                // XML vs HTML behavior
    CoreMatcher        *CoreProtocolMatcher // Protocol engine
    AttributeProcessor AttributeProcessor   // Custom attribute handling
    AllowSelfCloseTags bool                // <br/> support
    SkipComments       bool                // Performance optimization
}
```

## 📊 Performance Benchmarks

| Parser | Speed | Memory | Flexibility |
|--------|-------|--------|-------------|
| **MarkIt** | ⚡⚡⚡ | 🟢 Minimal | ⭐⭐⭐ Universal |
| Standard XML | ⚡⚡ | 🟡 Moderate | ⭐ XML Only |
| HTML Parser | ⚡⚡ | 🟡 Moderate | ⭐ HTML Only |
| Generic Parser | ⚡ | 🔴 Heavy | ⭐⭐ Limited |

```bash
# Run benchmarks
go test -bench=. -benchmem
```

## 🎯 Use Cases

### 🌐 **Web Development**
- Parse HTML with custom components
- Extract metadata and links
- Transform markup for SSG/SSR

### 📄 **Document Processing**
- Convert between markup formats
- Extract structured data
- Generate documentation

### 🔧 **Template Engines**
- Custom template syntax
- Macro expansion
- Dynamic content generation

### 🔌 **API Integration**
- Parse XML APIs responses
- Transform data formats
- Protocol translation

## 🧪 Testing & Quality

MarkIt maintains **100% test coverage** with comprehensive test suites:

```bash
# Run tests with coverage
go test -v -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Categories
- ✅ **Core Protocol Tests**: All bracket protocol combinations
- ✅ **Error Handling**: Comprehensive error scenarios  
- ✅ **Edge Cases**: Malformed input, boundary conditions
- ✅ **Performance**: Memory and speed benchmarks

## 🚀 Getting Started

### 1. **Installation**
```bash
go get github.com/khicago/markit
```

### 2. **Basic Parsing**
```go
parser := markit.NewParser(`<root><item>Hello</item></root>`)
ast, err := parser.Parse()
```

### 3. **Custom Configuration**
```go
config := markit.DefaultConfig()
config.CaseSensitive = false
parser := markit.NewParserWithConfig(input, config)
```

### 4. **AST Traversal**
```go
markit.Walk(ast, &YourCustomVisitor{})
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

### Development Setup

```bash
# Clone the repository
git clone https://github.com/khicago/markit.git
cd markit

# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run with coverage
go test -v -cover ./...
```

### Contribution Guidelines

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Write** tests for your changes
4. **Ensure** 100% test coverage is maintained
5. **Commit** your changes (`git commit -m 'Add amazing feature'`)
6. **Push** to the branch (`git push origin feature/amazing-feature`)
7. **Open** a Pull Request

### Code Quality Standards

- ✅ **100% test coverage** required
- ✅ **Go fmt** and **go vet** clean
- ✅ **Comprehensive documentation**
- ✅ **Benchmark tests** for performance changes

## 📚 Documentation

### API Reference
- [GoDoc](https://pkg.go.dev/github.com/khicago/markit) - Complete API documentation
- [Examples](examples/) - Practical usage examples
- [Benchmarks](benchmarks/) - Performance comparisons

### Guides
- [Custom Protocols](docs/custom-protocols.md) - Creating custom markup syntax
- [Performance Tuning](docs/performance.md) - Optimization strategies
- [Migration Guide](docs/migration.md) - Upgrading from other parsers

## 🔮 Roadmap

### v1.1.0 - Plugin System
- [ ] Dynamic protocol registration
- [ ] Plugin marketplace
- [ ] Hot-reloading support

### v1.2.0 - Advanced Features  
- [ ] Streaming parser for large files
- [ ] Schema validation
- [ ] Auto-completion support

### v2.0.0 - Next Generation
- [ ] WebAssembly support
- [ ] Multi-language bindings
- [ ] Visual protocol designer

## 🏆 Recognition

- ⭐ **Featured** in Awesome Go
- 🚀 **Trending** on GitHub
- 📈 **Growing** community adoption

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the need for universal markup parsing
- Built with ❤️ by the Go community
- Special thanks to all [contributors](https://github.com/khicago/markit/contributors)

---

<div align="center">

**[⭐ Star us on GitHub](https://github.com/khicago/markit)** • **[📖 Read the Docs](https://pkg.go.dev/github.com/khicago/markit)** • **[💬 Join Discussions](https://github.com/khicago/markit/discussions)**

Made with ❤️ for the Go community

</div>

## 特性

- **🔍 智能词法分析**：高效准确的 tokenization
- **🌳 AST 构建**：完整的抽象语法树支持  
- **⚙️ 灵活配置**：可自定义解析行为
- **🔄 访问者模式**：便于 AST 遍历和转换
- **📝 属性处理**：完整的属性解析和处理机制
- **💬 注释支持**：可选的注释处理
- **🏷️ 协议扩展**：内置协议匹配机制
- **🧩 Void Elements 支持**：支持 HTML5 标准 void elements 和自定义配置
- **🔧 自闭合标签**：完整的自闭合标签支持

### Void Elements 支持

MarkIt 提供了完整的 void elements 支持，包括 HTML5 标准 void elements 和自定义配置。

#### 默认配置
```go
// 默认配置不包含任何 void elements
config := markit.DefaultConfig()
fmt.Printf("支持 <br>: %v", config.IsVoidElement("br")) // false
```

#### HTML 配置
```go
// HTML 配置包含所有 HTML5 标准 void elements
config := markit.HTMLConfig()

// 支持的 HTML5 void elements:
// area, base, br, col, embed, hr, img, input, 
// link, meta, param, source, track, wbr

parser := markit.NewParserWithConfig(`<img src="test.jpg">`, config)
doc, _ := parser.Parse()

element := doc.Children[0].(*markit.Element)
fmt.Printf("标签: %s, 自闭合: %v", element.TagName, element.SelfClose)
// 输出: 标签: img, 自闭合: true
```

#### 自定义 Void Elements
```go
config := markit.DefaultConfig()

// 设置自定义 void elements
config.SetVoidElements([]string{"my-icon", "my-separator"})

// 动态添加和移除
config.AddVoidElement("custom-widget")
config.RemoveVoidElement("my-separator")

// 检查 void element
fmt.Printf("支持 my-icon: %v", config.IsVoidElement("my-icon")) // true
```

#### 混合解析示例
```go
config := markit.HTMLConfig()
input := `<article>
    <h1>标题</h1>
    <p>段落文本</p>
    <br>
    <img src="image.jpg" alt="图片">
    <hr>
    <input type="text" name="username">
</article>`

parser := markit.NewParserWithConfig(input, config)
doc, _ := parser.Parse()

// 所有 void elements (br, img, hr, input) 都会被正确解析为自闭合元素
```

#### XML 兼容性
```go
config := markit.HTMLConfig()

// 同时支持 XML style 和 HTML style
input := `<div>
    <br />      <!-- XML style -->
    <br>        <!-- HTML style -->
    <img src="test.jpg" />  <!-- XML style -->
    <img src="test2.jpg">   <!-- HTML style -->
</div>`

// 所有标签都会被正确解析
parser := markit.NewParserWithConfig(input, config)
doc, _ := parser.Parse()
``` 