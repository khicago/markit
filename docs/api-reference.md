---
layout: default
title: "API Reference - Complete Documentation"
description: "Complete API reference for MarkIt parser - all types, methods, and interfaces with examples."
keywords: "markit api reference, go xml parser api, html parser documentation"
author: "MarkIt Team"
---

# API Reference

> **Complete reference for all MarkIt types and methods**

This document provides comprehensive API documentation for the MarkIt parser library.

## 📋 Table of Contents

- [Core Types](#core-types)
- [Parser](#parser)
- [Configuration](#configuration)
- [AST Nodes](#ast-nodes)
- [Traversal](#traversal)
- [Attributes](#attributes)
- [Error Handling](#error-handling)
- [Utilities](#utilities)

## Core Types

### Parser

The main parser type for processing markup documents.

```go
type Parser struct {
    // Private fields
}
```

#### Constructor Functions

##### `NewParser(content string) *Parser`

Creates a new parser with default configuration.

```go
parser := markit.NewParser(`<root><item>content</item></root>`)
```

**Parameters:**
- `content` (string): The markup content to parse

**Returns:**
- `*Parser`: A new parser instance

##### `NewParserWithConfig(content string, config *ParserConfig) *Parser`

Creates a new parser with custom configuration.

```go
config := markit.HTMLConfig()
parser := markit.NewParserWithConfig(content, config)
```

**Parameters:**
- `content` (string): The markup content to parse
- `config` (*ParserConfig): Parser configuration

**Returns:**
- `*Parser`: A new parser instance

#### Methods

##### `Parse() (*Document, error)`

Parses the content and returns the AST.

```go
ast, err := parser.Parse()
if err != nil {
    log.Fatal(err)
}
```

**Returns:**
- `*Document`: The parsed document AST
- `error`: Parse error if any

##### `ParseWithOptions(options ParseOptions) (*Document, error)`

Parses with additional options.

```go
options := markit.ParseOptions{
    StrictMode: true,
    MaxDepth:   100,
}
ast, err := parser.ParseWithOptions(options)
```

**Parameters:**
- `options` (ParseOptions): Additional parsing options

**Returns:**
- `*Document`: The parsed document AST
- `error`: Parse error if any

### Document

Represents the root of the parsed document.

```go
type Document struct {
    Root     *Element
    Metadata map[string]interface{}
}
```

#### Fields

- `Root` (*Element): The root element of the document
- `Metadata` (map[string]interface{}): Document metadata

#### Methods

##### `String() string`

Returns a string representation of the document.

```go
fmt.Println(doc.String())
```

##### `ToXML() string`

Converts the document back to XML format.

```go
xmlContent := doc.ToXML()
```

##### `ToHTML() string`

Converts the document to HTML format.

```go
htmlContent := doc.ToHTML()
```

## Configuration

### ParserConfig

Configuration for parser behavior.

```go
type ParserConfig struct {
    CaseSensitive      bool
    AllowSelfCloseTags bool
    SkipComments       bool
    VoidElements       map[string]bool
    AttributeProcessor AttributeProcessor
    ProtocolMatcher    *ProtocolMatcher
}
```

#### Constructor Functions

##### `DefaultConfig() *ParserConfig`

Returns the default parser configuration.

```go
config := markit.DefaultConfig()
```

**Returns:**
- `*ParserConfig`: Default configuration

##### `HTMLConfig() *ParserConfig`

Returns HTML5-optimized configuration.

```go
config := markit.HTMLConfig()
```

**Returns:**
- `*ParserConfig`: HTML5 configuration

#### Methods

##### `IsVoidElement(tagName string) bool`

Checks if a tag is a void element.

```go
isVoid := config.IsVoidElement("br") // true for HTML config
```

**Parameters:**
- `tagName` (string): Tag name to check

**Returns:**
- `bool`: True if the tag is a void element

##### `AddVoidElement(tagName string)`

Adds a void element.

```go
config.AddVoidElement("my-icon")
```

**Parameters:**
- `tagName` (string): Tag name to add as void element

##### `RemoveVoidElement(tagName string)`

Removes a void element.

```go
config.RemoveVoidElement("br")
```

**Parameters:**
- `tagName` (string): Tag name to remove from void elements

##### `SetVoidElements(elements []string)`

Sets the complete list of void elements.

```go
config.SetVoidElements([]string{"br", "hr", "img"})
```

**Parameters:**
- `elements` ([]string): List of void element tag names

##### `NormalizeCase(tagName string) string`

Normalizes tag name case based on configuration.

```go
normalized := config.NormalizeCase("DIV") // "div" if case insensitive
```

**Parameters:**
- `tagName` (string): Tag name to normalize

**Returns:**
- `string`: Normalized tag name

### ParseOptions

Additional options for parsing.

```go
type ParseOptions struct {
    StrictMode bool
    MaxDepth   int
    Timeout    time.Duration
}
```

#### Fields

- `StrictMode` (bool): Enable strict parsing mode
- `MaxDepth` (int): Maximum nesting depth (0 = unlimited)
- `Timeout` (time.Duration): Parse timeout (0 = no timeout)

## AST Nodes

### Node

Base interface for all AST nodes.

```go
type Node interface {
    Type() NodeType
    String() string
    Parent() Node
    SetParent(Node)
}
```

#### Methods

##### `Type() NodeType`

Returns the node type.

```go
nodeType := node.Type()
```

**Returns:**
- `NodeType`: The type of the node

##### `String() string`

Returns string representation.

```go
content := node.String()
```

**Returns:**
- `string`: String representation of the node

##### `Parent() Node`

Returns the parent node.

```go
parent := node.Parent()
```

**Returns:**
- `Node`: Parent node (nil if root)

##### `SetParent(Node)`

Sets the parent node.

```go
node.SetParent(parentNode)
```

**Parameters:**
- `Node`: The parent node

### Element

Represents an XML/HTML element.

```go
type Element struct {
    TagName    string
    Attributes map[string]string
    Children   []Node
    SelfClosed bool
    parent     Node
}
```

#### Fields

- `TagName` (string): The element's tag name
- `Attributes` (map[string]string): Element attributes
- `Children` ([]Node): Child nodes
- `SelfClosed` (bool): Whether the element is self-closed

#### Methods

##### `AddChild(child Node)`

Adds a child node.

```go
element.AddChild(textNode)
```

**Parameters:**
- `child` (Node): Child node to add

##### `RemoveChild(child Node) bool`

Removes a child node.

```go
removed := element.RemoveChild(childNode)
```

**Parameters:**
- `child` (Node): Child node to remove

**Returns:**
- `bool`: True if child was found and removed

##### `GetAttribute(name string) (string, bool)`

Gets an attribute value.

```go
value, exists := element.GetAttribute("id")
```

**Parameters:**
- `name` (string): Attribute name

**Returns:**
- `string`: Attribute value
- `bool`: True if attribute exists

##### `SetAttribute(name, value string)`

Sets an attribute value.

```go
element.SetAttribute("class", "highlight")
```

**Parameters:**
- `name` (string): Attribute name
- `value` (string): Attribute value

##### `RemoveAttribute(name string) bool`

Removes an attribute.

```go
removed := element.RemoveAttribute("style")
```

**Parameters:**
- `name` (string): Attribute name

**Returns:**
- `bool`: True if attribute was found and removed

##### `HasAttribute(name string) bool`

Checks if an attribute exists.

```go
hasClass := element.HasAttribute("class")
```

**Parameters:**
- `name` (string): Attribute name

**Returns:**
- `bool`: True if attribute exists

##### `FindChildByTag(tagName string) *Element`

Finds first child element by tag name.

```go
child := element.FindChildByTag("div")
```

**Parameters:**
- `tagName` (string): Tag name to search for

**Returns:**
- `*Element`: First matching child element (nil if not found)

##### `FindChildrenByTag(tagName string) []*Element`

Finds all child elements by tag name.

```go
children := element.FindChildrenByTag("li")
```

**Parameters:**
- `tagName` (string): Tag name to search for

**Returns:**
- `[]*Element`: All matching child elements

##### `FindDescendantByTag(tagName string) *Element`

Finds first descendant element by tag name (recursive).

```go
descendant := element.FindDescendantByTag("span")
```

**Parameters:**
- `tagName` (string): Tag name to search for

**Returns:**
- `*Element`: First matching descendant element (nil if not found)

##### `FindDescendantsByTag(tagName string) []*Element`

Finds all descendant elements by tag name (recursive).

```go
descendants := element.FindDescendantsByTag("a")
```

**Parameters:**
- `tagName` (string): Tag name to search for

**Returns:**
- `[]*Element`: All matching descendant elements

### TextNode

Represents text content.

```go
type TextNode struct {
    Content string
    parent  Node
}
```

#### Fields

- `Content` (string): The text content

#### Methods

##### `Trim() string`

Returns trimmed text content.

```go
trimmed := textNode.Trim()
```

**Returns:**
- `string`: Trimmed text content

##### `IsWhitespace() bool`

Checks if the text is only whitespace.

```go
isWhitespace := textNode.IsWhitespace()
```

**Returns:**
- `bool`: True if content is only whitespace

### CommentNode

Represents an XML/HTML comment.

```go
type CommentNode struct {
    Content string
    parent  Node
}
```

#### Fields

- `Content` (string): The comment content

### CDATANode

Represents a CDATA section.

```go
type CDATANode struct {
    Content string
    parent  Node
}
```

#### Fields

- `Content` (string): The CDATA content

### ProcessingInstructionNode

Represents a processing instruction.

```go
type ProcessingInstructionNode struct {
    Target string
    Data   string
    parent Node
}
```

#### Fields

- `Target` (string): The PI target
- `Data` (string): The PI data

### DoctypeNode

Represents a DOCTYPE declaration.

```go
type DoctypeNode struct {
    Name     string
    PublicID string
    SystemID string
    parent   Node
}
```

#### Fields

- `Name` (string): DOCTYPE name
- `PublicID` (string): Public identifier
- `SystemID` (string): System identifier

### NodeType

Enumeration of node types.

```go
type NodeType int

const (
    ElementNodeType NodeType = iota
    TextNodeType
    CommentNodeType
    CDATANodeType
    ProcessingInstructionNodeType
    DoctypeNodeType
)
```

#### Constants

- `ElementNodeType`: Element node
- `TextNodeType`: Text node
- `CommentNodeType`: Comment node
- `CDATANodeType`: CDATA node
- `ProcessingInstructionNodeType`: Processing instruction node
- `DoctypeNodeType`: DOCTYPE node

## Traversal

### Walker

Interface for AST traversal.

```go
type Walker interface {
    Walk(node Node) error
}
```

#### Methods

##### `Walk(node Node) error`

Walks through the AST starting from the given node.

```go
err := walker.Walk(document.Root)
```

**Parameters:**
- `node` (Node): Starting node for traversal

**Returns:**
- `error`: Error if traversal fails

### Visitor

Interface for visiting nodes during traversal.

```go
type Visitor interface {
    VisitEnter(node Node) error
    VisitLeave(node Node) error
}
```

#### Methods

##### `VisitEnter(node Node) error`

Called when entering a node.

```go
err := visitor.VisitEnter(node)
```

**Parameters:**
- `node` (Node): The node being entered

**Returns:**
- `error`: Error to stop traversal

##### `VisitLeave(node Node) error`

Called when leaving a node.

```go
err := visitor.VisitLeave(node)
```

**Parameters:**
- `node` (Node): The node being left

**Returns:**
- `error`: Error to stop traversal

### Built-in Walkers

#### `NewDepthFirstWalker(visitor Visitor) Walker`

Creates a depth-first walker.

```go
walker := markit.NewDepthFirstWalker(myVisitor)
```

**Parameters:**
- `visitor` (Visitor): Visitor to use during traversal

**Returns:**
- `Walker`: Depth-first walker

#### `NewBreadthFirstWalker(visitor Visitor) Walker`

Creates a breadth-first walker.

```go
walker := markit.NewBreadthFirstWalker(myVisitor)
```

**Parameters:**
- `visitor` (Visitor): Visitor to use during traversal

**Returns:**
- `Walker`: Breadth-first walker

### Utility Functions

#### `WalkDepthFirst(node Node, visitor Visitor) error`

Convenience function for depth-first traversal.

```go
err := markit.WalkDepthFirst(document.Root, myVisitor)
```

**Parameters:**
- `node` (Node): Starting node
- `visitor` (Visitor): Visitor to use

**Returns:**
- `error`: Traversal error if any

#### `WalkBreadthFirst(node Node, visitor Visitor) error`

Convenience function for breadth-first traversal.

```go
err := markit.WalkBreadthFirst(document.Root, myVisitor)
```

**Parameters:**
- `node` (Node): Starting node
- `visitor` (Visitor): Visitor to use

**Returns:**
- `error`: Traversal error if any

## Attributes

### AttributeProcessor

Interface for custom attribute processing.

```go
type AttributeProcessor interface {
    ProcessAttributes(attrs map[string]string) map[string]string
    ProcessAttribute(key, value string) (string, string)
    IsBooleanAttribute(key string) bool
}
```

#### Methods

##### `ProcessAttributes(attrs map[string]string) map[string]string`

Processes all attributes at once.

```go
processed := processor.ProcessAttributes(attributes)
```

**Parameters:**
- `attrs` (map[string]string): Raw attributes

**Returns:**
- `map[string]string`: Processed attributes

##### `ProcessAttribute(key, value string) (string, string)`

Processes a single attribute.

```go
newKey, newValue := processor.ProcessAttribute("data-value", "123")
```

**Parameters:**
- `key` (string): Attribute name
- `value` (string): Attribute value

**Returns:**
- `string`: Processed attribute name
- `string`: Processed attribute value

##### `IsBooleanAttribute(key string) bool`

Checks if an attribute is boolean.

```go
isBool := processor.IsBooleanAttribute("checked")
```

**Parameters:**
- `key` (string): Attribute name

**Returns:**
- `bool`: True if attribute is boolean

### DefaultAttributeProcessor

Default implementation of AttributeProcessor.

```go
type DefaultAttributeProcessor struct{}
```

#### Constructor

##### `NewDefaultAttributeProcessor() *DefaultAttributeProcessor`

Creates a new default attribute processor.

```go
processor := markit.NewDefaultAttributeProcessor()
```

**Returns:**
- `*DefaultAttributeProcessor`: New processor instance

## Error Handling

### ParseError

Error type for parsing errors.

```go
type ParseError struct {
    Message  string
    Line     int
    Column   int
    Position int
}
```

#### Fields

- `Message` (string): Error message
- `Line` (int): Line number where error occurred
- `Column` (int): Column number where error occurred
- `Position` (int): Character position where error occurred

#### Methods

##### `Error() string`

Returns the error message.

```go
fmt.Println(err.Error())
```

**Returns:**
- `string`: Formatted error message

##### `String() string`

Returns detailed error information.

```go
fmt.Println(parseErr.String())
```

**Returns:**
- `string`: Detailed error information

### Error Types

#### `ErrInvalidSyntax`

Returned when markup syntax is invalid.

```go
var ErrInvalidSyntax = errors.New("invalid markup syntax")
```

#### `ErrUnexpectedEOF`

Returned when input ends unexpectedly.

```go
var ErrUnexpectedEOF = errors.New("unexpected end of input")
```

#### `ErrInvalidCharacter`

Returned when an invalid character is encountered.

```go
var ErrInvalidCharacter = errors.New("invalid character")
```

#### `ErrMaxDepthExceeded`

Returned when maximum nesting depth is exceeded.

```go
var ErrMaxDepthExceeded = errors.New("maximum nesting depth exceeded")
```

## Utilities

### String Functions

#### `EscapeXML(s string) string`

Escapes XML special characters.

```go
escaped := markit.EscapeXML(`<tag attr="value">content</tag>`)
```

**Parameters:**
- `s` (string): String to escape

**Returns:**
- `string`: Escaped string

#### `UnescapeXML(s string) string`

Unescapes XML entities.

```go
unescaped := markit.UnescapeXML("&lt;tag&gt;")
```

**Parameters:**
- `s` (string): String to unescape

**Returns:**
- `string`: Unescaped string

#### `EscapeHTML(s string) string`

Escapes HTML special characters.

```go
escaped := markit.EscapeHTML(`<script>alert("xss")</script>`)
```

**Parameters:**
- `s` (string): String to escape

**Returns:**
- `string`: Escaped string

#### `UnescapeHTML(s string) string`

Unescapes HTML entities.

```go
unescaped := markit.UnescapeHTML("&amp;nbsp;")
```

**Parameters:**
- `s` (string): String to unescape

**Returns:**
- `string`: Unescaped string

### Validation Functions

#### `IsValidTagName(name string) bool`

Validates if a string is a valid tag name.

```go
valid := markit.IsValidTagName("my-element") // true
```

**Parameters:**
- `name` (string): Tag name to validate

**Returns:**
- `bool`: True if valid tag name

#### `IsValidAttributeName(name string) bool`

Validates if a string is a valid attribute name.

```go
valid := markit.IsValidAttributeName("data-value") // true
```

**Parameters:**
- `name` (string): Attribute name to validate

**Returns:**
- `bool`: True if valid attribute name

### Conversion Functions

#### `NodeToMap(node Node) map[string]interface{}`

Converts a node to a map representation.

```go
nodeMap := markit.NodeToMap(element)
```

**Parameters:**
- `node` (Node): Node to convert

**Returns:**
- `map[string]interface{}`: Map representation

#### `MapToNode(data map[string]interface{}) (Node, error)`

Converts a map to a node.

```go
node, err := markit.MapToNode(nodeMap)
```

**Parameters:**
- `data` (map[string]interface{}): Map data

**Returns:**
- `Node`: Converted node
- `error`: Conversion error if any

#### `NodeToJSON(node Node) ([]byte, error)`

Converts a node to JSON.

```go
jsonData, err := markit.NodeToJSON(element)
```

**Parameters:**
- `node` (Node): Node to convert

**Returns:**
- `[]byte`: JSON data
- `error`: Conversion error if any

#### `JSONToNode(data []byte) (Node, error)`

Converts JSON to a node.

```go
node, err := markit.JSONToNode(jsonData)
```

**Parameters:**
- `data` ([]byte): JSON data

**Returns:**
- `Node`: Converted node
- `error`: Conversion error if any

---

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/khicago/markit"
)

func main() {
    // Parse HTML
    parser := markit.NewParserWithConfig(`
        <html>
            <head><title>Example</title></head>
            <body>
                <h1>Hello World</h1>
                <p>This is a paragraph.</p>
            </body>
        </html>
    `, markit.HTMLConfig())
    
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Find title
    title := doc.Root.FindDescendantByTag("title")
    if title != nil && len(title.Children) > 0 {
        if textNode, ok := title.Children[0].(*markit.TextNode); ok {
            fmt.Println("Title:", textNode.Content)
        }
    }
}
```

### Custom Visitor

```go
type ElementCounter struct {
    Count map[string]int
}

func (ec *ElementCounter) VisitEnter(node markit.Node) error {
    if element, ok := node.(*markit.Element); ok {
        ec.Count[element.TagName]++
    }
    return nil
}

func (ec *ElementCounter) VisitLeave(node markit.Node) error {
    return nil
}

// Usage
counter := &ElementCounter{Count: make(map[string]int)}
err := markit.WalkDepthFirst(doc.Root, counter)
```

---

<div align="center">

**[🏠 Back to Home](/)** • **[📋 Report Issues](https://github.com/khicago/markit/issues)** • **[💬 Discussions](https://github.com/khicago/markit/discussions)**

</div> 