package markit

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode/utf8"
)

// RenderOptions 渲染选项配置
type RenderOptions struct {
	// Indent 缩进字符串，如 "  " 表示两个空格，"\t" 表示制表符
	Indent string
	// EscapeText 是否转义文本内容（默认：true）
	EscapeText bool
	// PreserveSpace 是否保留空白字符
	PreserveSpace bool
	// CompactMode 小元素的单行输出模式
	CompactMode bool
	// SortAttributes 是否按字母顺序排序属性
	SortAttributes bool
	// EmptyElementStyle 空元素的样式
	EmptyElementStyle EmptyElementStyle
	// IncludeDeclaration 是否包含声明行（如 <?xml...?>, <!DOCTYPE...> 等）
	IncludeDeclaration bool
}

// EmptyElementStyle 空元素样式枚举
type EmptyElementStyle int

const (
	// SelfClosingStyle 自闭合风格 <tag/>
	SelfClosingStyle EmptyElementStyle = iota
	// PairedTagStyle 配对标签风格 <tag></tag>
	PairedTagStyle
	// VoidElementStyle 基于配置的 void 元素样式
	VoidElementStyle
)

// ValidationOptions 验证选项
type ValidationOptions struct {
	// CheckWellFormed 验证格式良好性
	CheckWellFormed bool
	// CheckEncoding 验证字符编码
	CheckEncoding bool
	// CheckNesting 检查元素嵌套规则
	CheckNesting bool
}

// ValidationError 验证错误
type ValidationError struct {
	Message  string
	Position Position
	NodeType NodeType
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error at line %d, column %d: %s",
		e.Position.Line, e.Position.Column, e.Message)
}

// Renderer 通用标记语言渲染器
type Renderer struct {
	options    *RenderOptions
	config     *ParserConfig
	validation *ValidationOptions
}

// NewRenderer 创建默认渲染器
func NewRenderer() *Renderer {
	return &Renderer{
		options: &RenderOptions{
			Indent:             "  ",
			EscapeText:         true,
			PreserveSpace:      false,
			CompactMode:        false,
			SortAttributes:     false,
			EmptyElementStyle:  SelfClosingStyle,
			IncludeDeclaration: true,
		},
	}
}

// NewRendererWithOptions 创建带选项的渲染器
func NewRendererWithOptions(opts *RenderOptions) *Renderer {
	if opts == nil {
		return NewRenderer()
	}

	// 创建选项副本以避免外部修改
	options := *opts
	return &Renderer{
		options: &options,
	}
}

// NewRendererWithConfig 创建带配置的渲染器
func NewRendererWithConfig(config *ParserConfig, opts *RenderOptions) *Renderer {
	renderer := NewRendererWithOptions(opts)
	renderer.config = config
	return renderer
}

// SetOptions 设置渲染选项
func (r *Renderer) SetOptions(opts *RenderOptions) {
	if opts != nil {
		// 创建副本以避免外部修改
		options := *opts
		r.options = &options
	}
}

// SetConfig 设置解析器配置
func (r *Renderer) SetConfig(config *ParserConfig) {
	r.config = config
}

// SetValidation 设置验证选项
func (r *Renderer) SetValidation(validation *ValidationOptions) {
	r.validation = validation
}

// Render 渲染文档为字符串（向后兼容）
func (r *Renderer) Render(doc *Document) string {
	result, _ := r.RenderToString(doc)
	return result
}

// RenderToString 渲染文档为字符串
func (r *Renderer) RenderToString(doc *Document) (string, error) {
	if doc == nil {
		return "", fmt.Errorf("document is nil")
	}

	// 执行验证
	if r.validation != nil {
		if err := r.validateDocument(doc); err != nil {
			return "", err
		}
	}

	var sb strings.Builder
	if err := r.RenderToWriter(doc, &sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// RenderToWriter 渲染文档到 Writer
func (r *Renderer) RenderToWriter(doc *Document, w io.Writer) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}
	if w == nil {
		return fmt.Errorf("writer is nil")
	}

	// 执行验证
	if r.validation != nil {
		if err := r.validateDocument(doc); err != nil {
			return err
		}
	}

	// 渲染文档节点
	for _, child := range doc.Children {
		if err := r.renderNode(child, w, 0); err != nil {
			return err
		}
	}

	return nil
}

// RenderElement 渲染单个元素为字符串
func (r *Renderer) RenderElement(elem *Element) (string, error) {
	if elem == nil {
		return "", fmt.Errorf("element is nil")
	}

	var sb strings.Builder
	if err := r.RenderElementToWriter(elem, &sb); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// RenderElementToWriter 渲染单个元素到 Writer
func (r *Renderer) RenderElementToWriter(elem *Element, w io.Writer) error {
	if elem == nil {
		return fmt.Errorf("element is nil")
	}
	if w == nil {
		return fmt.Errorf("writer is nil")
	}

	return r.renderNode(elem, w, 0)
}

// RenderWithValidation 带验证的渲染
func (r *Renderer) RenderWithValidation(doc *Document, opts *ValidationOptions) (string, error) {
	if doc == nil {
		return "", fmt.Errorf("document is nil")
	}

	// 临时设置验证选项
	oldValidation := r.validation
	r.validation = opts
	defer func() {
		r.validation = oldValidation
	}()

	return r.RenderToString(doc)
}

// renderNode 渲染单个节点
func (r *Renderer) renderNode(node Node, w io.Writer, depth int) error {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *Document:
		return r.renderDocument(n, w, depth)
	case *Element:
		return r.renderElement(n, w, depth)
	case *Text:
		return r.renderText(n, w, depth)
	case *Comment:
		return r.renderComment(n, w, depth)
	case *ProcessingInstruction:
		return r.renderProcessingInstruction(n, w, depth)
	case *Doctype:
		return r.renderDoctype(n, w, depth)
	case *CDATA:
		return r.renderCDATA(n, w, depth)
	default:
		return fmt.Errorf("unknown node type: %T", node)
	}
}

// renderDocument 渲染文档节点
func (r *Renderer) renderDocument(doc *Document, w io.Writer, depth int) error {
	for _, child := range doc.Children {
		if err := r.renderNode(child, w, depth); err != nil {
			return err
		}
	}
	return nil
}

// renderElement 渲染元素节点
func (r *Renderer) renderElement(elem *Element, w io.Writer, depth int) error {
	indent := strings.Repeat(r.options.Indent, depth)

	// 如果不是紧凑模式且不是顶层元素，添加缩进
	if !r.options.CompactMode && depth > 0 {
		if _, err := w.Write([]byte(indent)); err != nil {
			return err
		}
	}

	// 开始标签
	if _, err := w.Write([]byte("<")); err != nil {
		return err
	}
	if _, err := w.Write([]byte(elem.TagName)); err != nil {
		return err
	}

	// 渲染属性
	if err := r.renderAttributes(elem, w); err != nil {
		return err
	}

	// 处理自闭合元素
	if elem.SelfClose {
		switch r.options.EmptyElementStyle {
		case SelfClosingStyle:
			if _, err := w.Write([]byte(" />")); err != nil {
				return err
			}
		case PairedTagStyle:
			if _, err := w.Write([]byte("></")); err != nil {
				return err
			}
			if _, err := w.Write([]byte(elem.TagName)); err != nil {
				return err
			}
			if _, err := w.Write([]byte(">")); err != nil {
				return err
			}
		case VoidElementStyle:
			if r.config != nil && r.config.IsVoidElement(elem.TagName) {
				if _, err := w.Write([]byte(">")); err != nil {
					return err
				}
			} else {
				if _, err := w.Write([]byte(" />")); err != nil {
					return err
				}
			}
		default:
			if _, err := w.Write([]byte(" />")); err != nil {
				return err
			}
		}
		// 自闭合元素后换行
		if !r.options.CompactMode {
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		}
		return nil
	}

	if _, err := w.Write([]byte(">")); err != nil {
		return err
	}

	// 渲染子节点
	if len(elem.Children) > 0 {
		// 检查是否有非文本子节点
		hasNonTextChild := false
		for _, child := range elem.Children {
			if _, ok := child.(*Text); !ok {
				hasNonTextChild = true
				break
			}
		}

		// 检查是否只有一个文本子节点
		isSingleTextChild := len(elem.Children) == 1
		if textChild, ok := elem.Children[0].(*Text); ok && isSingleTextChild {
			// 单个文本子节点的情况
			// 对于单行简单文本，添加换行和缩进
			if !r.options.CompactMode && !strings.ContainsAny(textChild.Content, "\n\r") {
				if _, err := w.Write([]byte("\n")); err != nil {
					return err
				}
				if _, err := w.Write([]byte(strings.Repeat(r.options.Indent, depth+1))); err != nil {
					return err
				}
			}
			if err := r.renderText(textChild, w, depth+1); err != nil {
				return err
			}
			// 单个文本子节点后也需要换行和缩进
			if !r.options.CompactMode && !strings.ContainsAny(textChild.Content, "\n\r") {
				if _, err := w.Write([]byte("\n")); err != nil {
					return err
				}
				if _, err := w.Write([]byte(indent)); err != nil {
					return err
				}
			}
		} else {
			// 多个子节点或包含非文本节点的情况
			if !r.options.CompactMode {
				if _, err := w.Write([]byte("\n")); err != nil {
					return err
				}
			}

			for _, child := range elem.Children {
				if err := r.renderNode(child, w, depth+1); err != nil {
					return err
				}
			}

			// 结束标签前的缩进（只有在有非文本子节点时）
			if !r.options.CompactMode && hasNonTextChild {
				if _, err := w.Write([]byte(indent)); err != nil {
					return err
				}
			}
		}
	}

	// 结束标签
	if _, err := w.Write([]byte("</")); err != nil {
		return err
	}
	if _, err := w.Write([]byte(elem.TagName)); err != nil {
		return err
	}
	if _, err := w.Write([]byte(">")); err != nil {
		return err
	}

	// 元素后换行
	if !r.options.CompactMode {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

// renderAttributes 渲染属性
func (r *Renderer) renderAttributes(elem *Element, w io.Writer) error {
	if elem.Attributes == nil || len(elem.Attributes) == 0 {
		return nil
	}

	// 获取属性键并排序（如果需要）
	keys := make([]string, 0, len(elem.Attributes))
	for key := range elem.Attributes {
		keys = append(keys, key)
	}

	if r.options.SortAttributes {
		sort.Strings(keys)
	}

	// 渲染属性
	for _, key := range keys {
		value := elem.Attributes[key]
		if _, err := w.Write([]byte(" ")); err != nil {
			return err
		}
		if _, err := w.Write([]byte(key)); err != nil {
			return err
		}

		if value != "" {
			escapedValue := value
			if r.options.EscapeText {
				escapedValue = escapeText(value)
			}
			if _, err := w.Write([]byte(`="`)); err != nil {
				return err
			}
			if _, err := w.Write([]byte(escapedValue)); err != nil {
				return err
			}
			if _, err := w.Write([]byte(`"`)); err != nil {
				return err
			}
		}
	}

	return nil
}

// renderText 渲染文本节点
func (r *Renderer) renderText(text *Text, w io.Writer, depth int) error {
	content := text.Content
	if r.options.EscapeText {
		content = escapeText(content)
	}

	// 如果不是紧凑模式，并且文本包含换行或者是多行文本，需要处理缩进
	if !r.options.CompactMode && strings.ContainsAny(content, "\n\r\t") {
		// 对于包含换行的文本，保持原有格式但添加适当的缩进
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if i > 0 {
				if _, err := w.Write([]byte("\n")); err != nil {
					return err
				}
				if strings.TrimSpace(line) != "" { // 只对非空行添加缩进
					if _, err := w.Write([]byte(strings.Repeat(r.options.Indent, depth))); err != nil {
						return err
					}
				}
			}
			if _, err := w.Write([]byte(line)); err != nil {
				return err
			}
		}
	} else {
		// 简单文本直接输出
		if _, err := w.Write([]byte(content)); err != nil {
			return err
		}
	}
	return nil
}

// renderComment 渲染注释节点
func (r *Renderer) renderComment(comment *Comment, w io.Writer, depth int) error {
	if !r.options.CompactMode && depth > 0 {
		if err := r.writeIndent(w, depth); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("<!--" + comment.Content + "-->")); err != nil {
		return err
	}

	if !r.options.CompactMode {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

// renderProcessingInstruction 渲染处理指令节点
func (r *Renderer) renderProcessingInstruction(pi *ProcessingInstruction, w io.Writer, depth int) error {
	// 如果不包含声明，跳过处理指令
	if !r.options.IncludeDeclaration {
		return nil
	}

	if !r.options.CompactMode && depth > 0 {
		if err := r.writeIndent(w, depth); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("<?" + pi.Target)); err != nil {
		return err
	}

	if pi.Content != "" {
		if _, err := w.Write([]byte(" " + pi.Content)); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("?>")); err != nil {
		return err
	}

	if !r.options.CompactMode {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

// renderDoctype 渲染 DOCTYPE 节点
func (r *Renderer) renderDoctype(doctype *Doctype, w io.Writer, depth int) error {
	// 如果不包含声明，跳过 DOCTYPE
	if !r.options.IncludeDeclaration {
		return nil
	}

	if !r.options.CompactMode && depth > 0 {
		if err := r.writeIndent(w, depth); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("<!DOCTYPE " + doctype.Content + ">")); err != nil {
		return err
	}

	if !r.options.CompactMode {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

// renderCDATA 渲染 CDATA 节点
func (r *Renderer) renderCDATA(cdata *CDATA, w io.Writer, depth int) error {
	if !r.options.CompactMode && depth > 0 {
		if err := r.writeIndent(w, depth); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte("<![CDATA[" + cdata.Content + "]]>")); err != nil {
		return err
	}

	if !r.options.CompactMode {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

// writeIndent 写入缩进
func (r *Renderer) writeIndent(w io.Writer, depth int) error {
	for i := 0; i < depth; i++ {
		if _, err := w.Write([]byte(r.options.Indent)); err != nil {
			return err
		}
	}
	return nil
}

// isSmallElement 判断是否为小元素（适合紧凑模式）
func (r *Renderer) isSmallElement(elem *Element) bool {
	if len(elem.Children) == 0 {
		return true
	}

	if len(elem.Children) == 1 {
		if text, ok := elem.Children[0].(*Text); ok {
			return len(strings.TrimSpace(text.Content)) < 50
		}
	}

	return false
}

// isOnlyTextChildren 判断是否只有文本子节点
func (r *Renderer) isOnlyTextChildren(elem *Element) bool {
	for _, child := range elem.Children {
		if _, ok := child.(*Text); !ok {
			return false
		}
	}
	return true
}

// validateDocument 验证文档
func (r *Renderer) validateDocument(doc *Document) error {
	if r.validation == nil {
		return nil
	}

	var errors []error

	// 遍历文档检查各种验证规则
	for _, child := range doc.Children {
		if err := r.validateNode(child); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors[0] // 返回第一个错误
	}

	return nil
}

// validateNode 验证单个节点
func (r *Renderer) validateNode(node Node) error {
	if r.validation == nil {
		return nil
	}

	switch n := node.(type) {
	case *Element:
		return r.validateElement(n)
	case *Text:
		return r.validateText(n)
	default:
		return nil
	}
}

// validateElement 验证元素节点
func (r *Renderer) validateElement(elem *Element) error {
	if r.validation.CheckWellFormed {
		// 检查标签名是否有效
		if !isValidTagName(elem.TagName) {
			return &ValidationError{
				Message:  fmt.Sprintf("invalid tag name: %s", elem.TagName),
				Position: elem.Position(),
				NodeType: NodeTypeElement,
			}
		}

		// 检查属性名是否有效
		for attrName := range elem.Attributes {
			if !isValidAttributeName(attrName) {
				return &ValidationError{
					Message:  fmt.Sprintf("invalid attribute name: %s", attrName),
					Position: elem.Position(),
					NodeType: NodeTypeElement,
				}
			}
		}
	}

	// 递归验证子节点
	for _, child := range elem.Children {
		if err := r.validateNode(child); err != nil {
			return err
		}
	}

	return nil
}

// validateText 验证文本节点
func (r *Renderer) validateText(text *Text) error {
	if r.validation == nil || !r.validation.CheckEncoding {
		return nil
	}

	// 检查 UTF-8 编码是否有效
	if !utf8.ValidString(text.Content) {
		return &ValidationError{
			Message:  "invalid UTF-8 encoding in text content",
			Position: text.Position(),
			NodeType: NodeTypeText,
		}
	}

	return nil
}

// isValidTagName 检查标签名是否有效
func isValidTagName(name string) bool {
	if name == "" {
		return false
	}

	// 通用标签名规则：以字母或下划线开头，后续可包含字母、数字、连字符、下划线、点
	for i, r := range name {
		if i == 0 {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_') {
				return false
			}
		} else {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.') {
				return false
			}
		}
	}

	return true
}

// isValidAttributeName 检查属性名是否有效
func isValidAttributeName(name string) bool {
	return isValidTagName(name) // 使用相同的规则
}

// escapeText 转义文本内容
func escapeText(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
