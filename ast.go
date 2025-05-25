// Package markit provides an extensible markup language parser with pluggable
// tag bracket protocols for parsing various markup languages and custom formats.
//
// MarkIt is designed around the concept of configurable tag bracket protocols,
// allowing users to define custom markup syntaxes through plugins. Instead of
// being tied to specific markup languages like XML or HTML, MarkIt provides
// a flexible foundation where any tag-based syntax can be supported through
// protocol extensions.
//
// Core Design Philosophy:
//   - Tag Bracket Protocols: Define custom syntax like <a...b> where 'a' and 'b'
//     are configurable opening and closing bracket sequences
//   - Plugin Architecture: Load syntax plugins at parser initialization
//   - Language Agnostic: Support XML, HTML, and any custom markup through plugins
//   - Compile-time Configuration: Case sensitivity and other options configurable
//
// Key Features:
// - 🔧 Pluggable tag bracket protocols (<...>, <?...?>, <!...>, etc.)
// - 🎯 Universal syntax support (XML, HTML, custom formats)
// - ⚡ High-performance parsing with minimal overhead
// - 📍 Precise error reporting with position tracking
// - 🔀 Visitor pattern for flexible AST traversal
// - 🧩 Plugin system for extending syntax support
package markit

// Node 表示 AST 中的一个节点
type Node interface {
	// Type 返回节点类型
	Type() NodeType
	// Position 返回节点在源码中的位置
	Position() Position
	// String 返回节点的字符串表示
	String() string
}

// NodeType 表示节点类型
type NodeType int

const (
	NodeTypeDocument NodeType = iota
	NodeTypeElement
	NodeTypeText
	NodeTypeProcessingInstruction
	NodeTypeDoctype
	NodeTypeCDATA
	NodeTypeComment
)

// Document 表示文档根节点
type Document struct {
	Children []Node
	Pos      Position
}

func (d *Document) Type() NodeType     { return NodeTypeDocument }
func (d *Document) Position() Position { return d.Pos }
func (d *Document) String() string     { return "Document" }

// Element 表示元素节点
type Element struct {
	TagName    string
	Attributes map[string]string
	Children   []Node
	SelfClose  bool
	Pos        Position
}

func (e *Element) Type() NodeType     { return NodeTypeElement }
func (e *Element) Position() Position { return e.Pos }
func (e *Element) String() string     { return e.TagName }

// Text 表示文本节点
type Text struct {
	Content string
	Pos     Position
}

func (t *Text) Type() NodeType     { return NodeTypeText }
func (t *Text) Position() Position { return t.Pos }
func (t *Text) String() string     { return t.Content }

// ProcessingInstruction 表示处理指令节点
type ProcessingInstruction struct {
	Target  string
	Content string
	Pos     Position
}

func (pi *ProcessingInstruction) Type() NodeType     { return NodeTypeProcessingInstruction }
func (pi *ProcessingInstruction) Position() Position { return pi.Pos }
func (pi *ProcessingInstruction) String() string     { return pi.Target }

// Doctype 表示DOCTYPE声明节点
type Doctype struct {
	Content string
	Pos     Position
}

func (dt *Doctype) Type() NodeType     { return NodeTypeDoctype }
func (dt *Doctype) Position() Position { return dt.Pos }
func (dt *Doctype) String() string     { return dt.Content }

// CDATA 表示CDATA节点
type CDATA struct {
	Content string
	Pos     Position
}

func (cd *CDATA) Type() NodeType     { return NodeTypeCDATA }
func (cd *CDATA) Position() Position { return cd.Pos }
func (cd *CDATA) String() string     { return cd.Content }

// Comment 表示注释节点
type Comment struct {
	Content string
	Pos     Position
}

func (c *Comment) Type() NodeType     { return NodeTypeComment }
func (c *Comment) Position() Position { return c.Pos }
func (c *Comment) String() string     { return c.Content }

// AttributeProcessor 属性处理器接口
type AttributeProcessor interface {
	// ProcessAttribute 处理属性，返回处理后的键值对
	ProcessAttribute(key, value string) (string, interface{}, error)
	// IsBooleanAttribute 检查是否是布尔属性
	IsBooleanAttribute(key string) bool
}

// DefaultAttributeProcessor 默认属性处理器
type DefaultAttributeProcessor struct{}

func (p *DefaultAttributeProcessor) ProcessAttribute(key, value string) (string, interface{}, error) {
	// 如果值为空，认为是布尔属性
	if value == "" {
		return key, true, nil
	}
	return key, value, nil
}

func (p *DefaultAttributeProcessor) IsBooleanAttribute(key string) bool {
	// HTML5 标准布尔属性列表
	booleanAttrs := map[string]bool{
		"checked":   true,
		"disabled":  true,
		"selected":  true,
		"readonly":  true,
		"required":  true,
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
	}
	return booleanAttrs[key]
}
