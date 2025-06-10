package markit

import (
	"fmt"
	"sort"
	"strings"
)

// Parser 语法分析器
type Parser struct {
	lexer     *Lexer
	current   Token
	peek      Token
	processor AttributeProcessor
	config    *ParserConfig
}

// NewParser 创建新的语法分析器（使用默认配置）
func NewParser(input string) *Parser {
	return NewParserWithConfig(input, DefaultConfig())
}

// NewParserWithConfig 创建带配置的语法分析器
func NewParserWithConfig(input string, config *ParserConfig) *Parser {
	lexer := NewLexerWithConfig(input, config)
	p := &Parser{
		lexer:     lexer,
		processor: config.AttributeProcessor,
		config:    config,
	}

	// 读取前两个 token，跳过注释
	p.nextToken()
	p.nextToken()

	// 如果配置要求跳过注释，则跳过它们
	if p.config.SkipComments {
		for p.current.Type == TokenComment {
			p.nextToken()
		}
	}

	return p
}

// SetAttributeProcessor 设置属性处理器
func (p *Parser) SetAttributeProcessor(processor AttributeProcessor) {
	p.processor = processor
	p.config.AttributeProcessor = processor
}

// GetConfig 获取解析器配置
func (p *Parser) GetConfig() *ParserConfig {
	return p.config
}

// SetConfig 设置解析器配置
func (p *Parser) SetConfig(config *ParserConfig) {
	p.config = config
	p.processor = config.AttributeProcessor
	// 更新lexer配置
	p.lexer.SetConfig(config)
}

// Parse 解析输入并返回 AST
func (p *Parser) Parse() (*Document, error) {
	doc := &Document{
		Children: []Node{},
		Pos:      p.current.Position,
	}

	for p.current.Type != TokenEOF {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			doc.Children = append(doc.Children, node)
		}
	}

	return doc, nil
}

// parseNode 解析一个节点
func (p *Parser) parseNode() (Node, error) {
	// 如果配置要求跳过注释，则跳过注释token
	if p.config.SkipComments && p.current.Type == TokenComment {
		p.nextToken()
		return p.parseNode() // 递归解析下一个节点
	}

	switch p.current.Type {
	case TokenText:
		return p.parseText()
	case TokenOpenTag:
		return p.parseElement()
	case TokenSelfCloseTag:
		return p.parseSelfCloseElement()
	case TokenProcessingInstruction:
		return p.parseProcessingInstruction()
	case TokenDoctype:
		return p.parseDoctype()
	case TokenCDATA:
		return p.parseCDATA()
	case TokenComment:
		return p.parseComment()
	case TokenError:
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  p.current.Value,
		}
	case TokenEOF:
		return nil, nil
	default:
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("unexpected token %s", p.current.Type),
		}
	}
}

// parseText 解析文本节点
func (p *Parser) parseText() (Node, error) {
	if p.current.Type != TokenText {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("expected text token, got %s", p.current.Type),
		}
	}

	text := &Text{
		Content: p.current.Value,
		Pos:     p.current.Position,
	}

	p.nextToken()
	return text, nil
}

// parseElement 解析元素节点
func (p *Parser) parseElement() (Node, error) {
	if p.current.Type != TokenOpenTag {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("expected open tag, got %s", p.current.Type),
		}
	}

	element := &Element{
		TagName:    p.current.Value,
		Attributes: p.current.Attributes,
		Children:   []Node{},
		SelfClose:  false,
		Pos:        p.current.Position,
	}

	tagName := p.current.Value
	p.nextToken()

	// 检查是否是 void element
	if p.config != nil && p.config.IsVoidElement(tagName) {
		// void element 不需要结束标签，直接返回自闭合元素
		element.SelfClose = true
		return element, nil
	}

	// 解析子节点
	for p.current.Type != TokenCloseTag && p.current.Type != TokenEOF {
		child, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if child != nil {
			element.Children = append(element.Children, child)
		}
	}

	// 检查结束标签
	if p.current.Type != TokenCloseTag {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("expected close tag for <%s>, got %s", tagName, p.current.Type),
		}
	}

	if p.current.Value != tagName {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("mismatched tags: expected </%s>, got </%s>", tagName, p.current.Value),
		}
	}

	p.nextToken()
	return element, nil
}

// parseSelfCloseElement 解析自闭合元素
func (p *Parser) parseSelfCloseElement() (Node, error) {
	if p.current.Type != TokenSelfCloseTag {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("expected self-close tag, got %s", p.current.Type),
		}
	}

	element := &Element{
		TagName:    p.current.Value,
		Attributes: p.current.Attributes,
		Children:   []Node{},
		SelfClose:  true,
		Pos:        p.current.Position,
	}

	p.nextToken()
	return element, nil
}

// parseProcessingInstruction 解析处理指令
func (p *Parser) parseProcessingInstruction() (Node, error) {
	if p.current.Type != TokenProcessingInstruction {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("expected processing instruction token, got %s", p.current.Type),
		}
	}

	pi := &ProcessingInstruction{
		Target:  p.current.Value,
		Content: p.current.Value,
		Pos:     p.current.Position,
	}

	p.nextToken()
	return pi, nil
}

// parseDoctype 解析DOCTYPE声明
func (p *Parser) parseDoctype() (Node, error) {
	if p.current.Type != TokenDoctype {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("expected doctype token, got %s", p.current.Type),
		}
	}

	doctype := &Doctype{
		Content: p.current.Value,
		Pos:     p.current.Position,
	}

	p.nextToken()
	return doctype, nil
}

// parseCDATA 解析CDATA节点
func (p *Parser) parseCDATA() (Node, error) {
	if p.current.Type != TokenCDATA {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("expected CDATA token, got %s", p.current.Type),
		}
	}

	cdata := &CDATA{
		Content: p.current.Value,
		Pos:     p.current.Position,
	}

	p.nextToken()
	return cdata, nil
}

// parseComment 解析注释节点
func (p *Parser) parseComment() (Node, error) {
	if p.current.Type != TokenComment {
		return nil, &ParseError{
			Position: p.current.Position,
			Message:  fmt.Sprintf("expected comment token, got %s", p.current.Type),
		}
	}

	comment := &Comment{
		Content: p.current.Value,
		Pos:     p.current.Position,
	}

	p.nextToken()
	return comment, nil
}

// nextToken 移动到下一个 token
func (p *Parser) nextToken() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()

	// 不在这里跳过注释，让parseNode处理
}

// ParseError 解析错误
type ParseError struct {
	Position Position
	Message  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %s: %s", e.Position, e.Message)
}

// Visitor 访问者接口，用于遍历 AST
type Visitor interface {
	VisitDocument(*Document) error
	VisitElement(*Element) error
	VisitText(*Text) error
	VisitProcessingInstruction(*ProcessingInstruction) error
	VisitDoctype(*Doctype) error
	VisitCDATA(*CDATA) error
	VisitComment(*Comment) error
}

// Walk 遍历 AST
func Walk(node Node, visitor Visitor) error {
	switch n := node.(type) {
	case *Document:
		if err := visitor.VisitDocument(n); err != nil {
			return err
		}
		for _, child := range n.Children {
			if err := Walk(child, visitor); err != nil {
				return err
			}
		}
	case *Element:
		if err := visitor.VisitElement(n); err != nil {
			return err
		}
		for _, child := range n.Children {
			if err := Walk(child, visitor); err != nil {
				return err
			}
		}
	case *Text:
		return visitor.VisitText(n)
	case *ProcessingInstruction:
		return visitor.VisitProcessingInstruction(n)
	case *Doctype:
		return visitor.VisitDoctype(n)
	case *CDATA:
		return visitor.VisitCDATA(n)
	case *Comment:
		return visitor.VisitComment(n)
	}
	return nil
}

// PrettyPrint 美化打印 AST
func PrettyPrint(node Node) string {
	debugRenderer := NewDebugRenderer()
	return debugRenderer.RenderDebug(node)
}

// DebugRenderer 调试渲染器，专门用于AST结构展示
type DebugRenderer struct {
	*Renderer
}

// NewDebugRenderer 创建调试渲染器
func NewDebugRenderer() *DebugRenderer {
	opts := &RenderOptions{
		Indent:         "  ",
		EscapeText:     false, // 调试时不转义，显示原始内容
		CompactMode:    false,
		SortAttributes: true, // 调试时排序属性，保证输出一致性
	}
	
	return &DebugRenderer{
		Renderer: NewRendererWithOptions(opts),
	}
}

// RenderDebug 渲染调试信息
func (dr *DebugRenderer) RenderDebug(node Node) string {
	var sb strings.Builder
	dr.renderDebugNode(node, &sb, 0)
	return sb.String()
}

// renderDebugNode 渲染调试节点，复用Renderer的基础设施
func (dr *DebugRenderer) renderDebugNode(node Node, sb *strings.Builder, depth int) {
	if node == nil {
		return
	}

	indentStr := strings.Repeat(dr.options.Indent, depth)

	switch n := node.(type) {
	case *Document:
		sb.WriteString(fmt.Sprintf("%sDocument\n", indentStr))
		for _, child := range n.Children {
			dr.renderDebugNode(child, sb, depth+1)
		}
	case *Element:
		sb.WriteString(fmt.Sprintf("%s<%s", indentStr, n.TagName))
		
		// 复用Renderer的属性处理逻辑
		if len(n.Attributes) > 0 {
			// 获取排序后的属性键
			keys := make([]string, 0, len(n.Attributes))
			for key := range n.Attributes {
				keys = append(keys, key)
			}
			if dr.options.SortAttributes {
				sort.Strings(keys)
			}
			
			for _, key := range keys {
				value := n.Attributes[key]
				if value == "" {
					sb.WriteString(fmt.Sprintf(" %s", key))
				} else {
					sb.WriteString(fmt.Sprintf(" %s=%q", key, value))
				}
			}
		}
		
		if n.SelfClose {
			sb.WriteString(" />\n")
		} else {
			sb.WriteString(">\n")
			for _, child := range n.Children {
				dr.renderDebugNode(child, sb, depth+1)
			}
			sb.WriteString(fmt.Sprintf("%s</%s>\n", indentStr, n.TagName))
		}
	case *Text:
		sb.WriteString(fmt.Sprintf("%sText: %q\n", indentStr, n.Content))
	case *ProcessingInstruction:
		sb.WriteString(fmt.Sprintf("%sPI: %q\n", indentStr, n.Content))
	case *Doctype:
		sb.WriteString(fmt.Sprintf("%sDoctype: %q\n", indentStr, n.Content))
	case *CDATA:
		sb.WriteString(fmt.Sprintf("%sCDATA: %q\n", indentStr, n.Content))
	case *Comment:
		sb.WriteString(fmt.Sprintf("%sComment: %q\n", indentStr, n.Content))
	}
}
