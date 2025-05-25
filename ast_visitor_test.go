package markit

import (
	"strings"
	"testing"
)

// WalkTestVisitor 实现 Visitor 接口用于测试
type WalkTestVisitor struct {
	visitedNodes []Node
	visitedTypes []NodeType
	shouldStop   bool
	stopAt       string
}

func (v *WalkTestVisitor) VisitDocument(node *Document) error {
	v.visitedNodes = append(v.visitedNodes, node)
	v.visitedTypes = append(v.visitedTypes, node.Type())
	return nil
}

func (v *WalkTestVisitor) VisitElement(node *Element) error {
	v.visitedNodes = append(v.visitedNodes, node)
	v.visitedTypes = append(v.visitedTypes, node.Type())
	if v.shouldStop && node.TagName == v.stopAt {
		return &ParseError{Message: "stop"}
	}
	return nil
}

func (v *WalkTestVisitor) VisitText(node *Text) error {
	v.visitedNodes = append(v.visitedNodes, node)
	v.visitedTypes = append(v.visitedTypes, node.Type())
	return nil
}

func (v *WalkTestVisitor) VisitProcessingInstruction(node *ProcessingInstruction) error {
	v.visitedNodes = append(v.visitedNodes, node)
	v.visitedTypes = append(v.visitedTypes, node.Type())
	return nil
}

func (v *WalkTestVisitor) VisitDoctype(node *Doctype) error {
	v.visitedNodes = append(v.visitedNodes, node)
	v.visitedTypes = append(v.visitedTypes, node.Type())
	return nil
}

func (v *WalkTestVisitor) VisitCDATA(node *CDATA) error {
	v.visitedNodes = append(v.visitedNodes, node)
	v.visitedTypes = append(v.visitedTypes, node.Type())
	return nil
}

func (v *WalkTestVisitor) VisitComment(node *Comment) error {
	v.visitedNodes = append(v.visitedNodes, node)
	v.visitedTypes = append(v.visitedTypes, node.Type())
	return nil
}

// TestWalkTraversal 测试Walk函数的遍历功能
func TestWalkTraversal(t *testing.T) {
	input := `<root>
		<header>
			<title>Test</title>
		</header>
		<content>
			<p>Paragraph 1</p>
			<p>Paragraph 2</p>
		</content>
	</root>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	visitor := &WalkTestVisitor{}
	err = Walk(doc, visitor)
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}

	// 验证访问的节点数量
	expectedNodeCount := 8 // Document + Element(root) + Element(header) + Element(title) + Text + Element(content) + Element(p)*2 + Text*2
	if len(visitor.visitedNodes) < expectedNodeCount {
		t.Errorf("expected at least %d visited nodes, got %d", expectedNodeCount, len(visitor.visitedNodes))
	}

	// 验证第一个节点是Document
	if visitor.visitedTypes[0] != NodeTypeDocument {
		t.Errorf("expected first visited node to be Document, got %v", visitor.visitedTypes[0])
	}

	// 验证包含所有预期的节点类型
	hasDocument := false
	hasElement := false
	hasText := false

	for _, nodeType := range visitor.visitedTypes {
		switch nodeType {
		case NodeTypeDocument:
			hasDocument = true
		case NodeTypeElement:
			hasElement = true
		case NodeTypeText:
			hasText = true
		}
	}

	if !hasDocument {
		t.Error("expected to visit Document node")
	}
	if !hasElement {
		t.Error("expected to visit Element nodes")
	}
	if !hasText {
		t.Error("expected to visit Text nodes")
	}
}

// TestWalkEarlyTermination 测试Walk函数的提前终止功能
func TestWalkEarlyTermination(t *testing.T) {
	input := `<root>
		<first>Content 1</first>
		<second>Content 2</second>
		<third>Content 3</third>
	</root>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	visitor := &WalkTestVisitor{shouldStop: true, stopAt: "second"}
	err = Walk(doc, visitor)

	// 应该返回错误表示提前终止
	if err == nil {
		t.Error("expected error for early termination")
	}

	// 验证访问的元素
	var visitedElements []string
	for _, node := range visitor.visitedNodes {
		if element, ok := node.(*Element); ok {
			visitedElements = append(visitedElements, element.TagName)
		}
	}

	// 验证没有访问到"third"元素
	for _, tagName := range visitedElements {
		if tagName == "third" {
			t.Error("expected not to visit 'third' element due to early termination")
		}
	}
}

// CommentCollector 专门收集注释的访问者
type CommentCollector struct {
	comments []string
	elements []string
}

func (c *CommentCollector) VisitDocument(node *Document) error { return nil }
func (c *CommentCollector) VisitElement(node *Element) error {
	if node.TagName != "root" {
		c.elements = append(c.elements, node.TagName)
	}
	return nil
}
func (c *CommentCollector) VisitText(node *Text) error                                   { return nil }
func (c *CommentCollector) VisitProcessingInstruction(node *ProcessingInstruction) error { return nil }
func (c *CommentCollector) VisitDoctype(node *Doctype) error                             { return nil }
func (c *CommentCollector) VisitCDATA(node *CDATA) error                                 { return nil }
func (c *CommentCollector) VisitComment(node *Comment) error {
	c.comments = append(c.comments, strings.TrimSpace(node.Content))
	return nil
}

// TestWalkWithComments 测试包含注释的遍历
func TestWalkWithComments(t *testing.T) {
	input := `<root>
		<!-- This is a comment -->
		<element>Content</element>
		<!-- Another comment -->
	</root>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	collector := &CommentCollector{}
	err = Walk(doc, collector)
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}

	// 验证注释内容
	expectedComments := []string{"This is a comment", "Another comment"}
	if len(collector.comments) != len(expectedComments) {
		t.Errorf("expected %d comments, got %d", len(expectedComments), len(collector.comments))
	}

	for i, expected := range expectedComments {
		if i < len(collector.comments) && collector.comments[i] != expected {
			t.Errorf("expected comment %d to be %q, got %q", i, expected, collector.comments[i])
		}
	}

	// 验证元素
	if len(collector.elements) != 1 || collector.elements[0] != "element" {
		t.Errorf("expected one 'element', got %v", collector.elements)
	}
}

// ElementOrderCollector 收集元素遍历顺序的访问者
type ElementOrderCollector struct {
	elementOrder []string
}

func (e *ElementOrderCollector) VisitDocument(node *Document) error { return nil }
func (e *ElementOrderCollector) VisitElement(node *Element) error {
	e.elementOrder = append(e.elementOrder, node.TagName)
	return nil
}
func (e *ElementOrderCollector) VisitText(node *Text) error { return nil }
func (e *ElementOrderCollector) VisitProcessingInstruction(node *ProcessingInstruction) error {
	return nil
}
func (e *ElementOrderCollector) VisitDoctype(node *Doctype) error { return nil }
func (e *ElementOrderCollector) VisitCDATA(node *CDATA) error     { return nil }
func (e *ElementOrderCollector) VisitComment(node *Comment) error { return nil }

// TestWalkDepthFirstOrder 测试深度优先遍历顺序
func TestWalkDepthFirstOrder(t *testing.T) {
	input := `<a>
		<b>
			<c>Text C</c>
		</b>
		<d>Text D</d>
	</a>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	collector := &ElementOrderCollector{}
	err = Walk(doc, collector)
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}

	// 验证深度优先遍历顺序：a -> b -> c -> d
	expectedOrder := []string{"a", "b", "c", "d"}
	if len(collector.elementOrder) != len(expectedOrder) {
		t.Errorf("expected %d elements, got %d", len(expectedOrder), len(collector.elementOrder))
	}

	for i, expected := range expectedOrder {
		if i < len(collector.elementOrder) && collector.elementOrder[i] != expected {
			t.Errorf("expected element %d to be %q, got %q", i, expected, collector.elementOrder[i])
		}
	}
}

// SelfClosingElementCollector 收集自闭合元素的访问者
type SelfClosingElementCollector struct {
	selfClosingElements []string
	regularElements     []string
}

func (s *SelfClosingElementCollector) VisitDocument(node *Document) error { return nil }
func (s *SelfClosingElementCollector) VisitElement(node *Element) error {
	if node.TagName != "container" {
		if node.SelfClose {
			s.selfClosingElements = append(s.selfClosingElements, node.TagName)
		} else {
			s.regularElements = append(s.regularElements, node.TagName)
		}
	}
	return nil
}
func (s *SelfClosingElementCollector) VisitText(node *Text) error { return nil }
func (s *SelfClosingElementCollector) VisitProcessingInstruction(node *ProcessingInstruction) error {
	return nil
}
func (s *SelfClosingElementCollector) VisitDoctype(node *Doctype) error { return nil }
func (s *SelfClosingElementCollector) VisitCDATA(node *CDATA) error     { return nil }
func (s *SelfClosingElementCollector) VisitComment(node *Comment) error { return nil }

// TestWalkWithSelfClosingElements 测试包含自闭合元素的遍历
func TestWalkWithSelfClosingElements(t *testing.T) {
	input := `<container>
		<img src="test.jpg" />
		<p>Text content</p>
		<br />
	</container>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	collector := &SelfClosingElementCollector{}
	err = Walk(doc, collector)
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}

	// 验证自闭合元素
	expectedSelfClosing := []string{"img", "br"}
	if len(collector.selfClosingElements) != len(expectedSelfClosing) {
		t.Errorf("expected %d self-closing elements, got %d", len(expectedSelfClosing), len(collector.selfClosingElements))
	}

	for i, expected := range expectedSelfClosing {
		if i < len(collector.selfClosingElements) && collector.selfClosingElements[i] != expected {
			t.Errorf("expected self-closing element %d to be %q, got %q", i, expected, collector.selfClosingElements[i])
		}
	}

	// 验证常规元素
	expectedRegular := []string{"p"}
	if len(collector.regularElements) != len(expectedRegular) {
		t.Errorf("expected %d regular elements, got %d", len(expectedRegular), len(collector.regularElements))
	}

	for i, expected := range expectedRegular {
		if i < len(collector.regularElements) && collector.regularElements[i] != expected {
			t.Errorf("expected regular element %d to be %q, got %q", i, expected, collector.regularElements[i])
		}
	}
}

// CountingVisitor 计数访问者
type CountingVisitor struct {
	count int
}

func (c *CountingVisitor) VisitDocument(node *Document) error { c.count++; return nil }
func (c *CountingVisitor) VisitElement(node *Element) error   { c.count++; return nil }
func (c *CountingVisitor) VisitText(node *Text) error         { c.count++; return nil }
func (c *CountingVisitor) VisitProcessingInstruction(node *ProcessingInstruction) error {
	c.count++
	return nil
}
func (c *CountingVisitor) VisitDoctype(node *Doctype) error { c.count++; return nil }
func (c *CountingVisitor) VisitCDATA(node *CDATA) error     { c.count++; return nil }
func (c *CountingVisitor) VisitComment(node *Comment) error { c.count++; return nil }

// TestWalkEmptyDocument 测试空文档的遍历
func TestWalkEmptyDocument(t *testing.T) {
	doc := &Document{
		Children: []Node{},
	}

	visitor := &CountingVisitor{}
	err := Walk(doc, visitor)
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}

	// 应该只访问Document节点本身
	if visitor.count != 1 {
		t.Errorf("expected to visit 1 node (Document), got %d", visitor.count)
	}
}

// NodeTypeCollector 收集节点类型的访问者
type NodeTypeCollector struct {
	nodeTypes []NodeType
}

func (n *NodeTypeCollector) VisitDocument(node *Document) error {
	n.nodeTypes = append(n.nodeTypes, node.Type())
	return nil
}
func (n *NodeTypeCollector) VisitElement(node *Element) error {
	n.nodeTypes = append(n.nodeTypes, node.Type())
	return nil
}
func (n *NodeTypeCollector) VisitText(node *Text) error {
	n.nodeTypes = append(n.nodeTypes, node.Type())
	return nil
}
func (n *NodeTypeCollector) VisitProcessingInstruction(node *ProcessingInstruction) error {
	n.nodeTypes = append(n.nodeTypes, node.Type())
	return nil
}
func (n *NodeTypeCollector) VisitDoctype(node *Doctype) error {
	n.nodeTypes = append(n.nodeTypes, node.Type())
	return nil
}
func (n *NodeTypeCollector) VisitCDATA(node *CDATA) error {
	n.nodeTypes = append(n.nodeTypes, node.Type())
	return nil
}
func (n *NodeTypeCollector) VisitComment(node *Comment) error {
	n.nodeTypes = append(n.nodeTypes, node.Type())
	return nil
}

// TestWalkSingleTextNode 测试只包含文本节点的文档
func TestWalkSingleTextNode(t *testing.T) {
	input := "Just plain text"

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	collector := &NodeTypeCollector{}
	err = Walk(doc, collector)
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}

	// 应该访问Document和Text节点
	expectedTypes := []NodeType{NodeTypeDocument, NodeTypeText}
	if len(collector.nodeTypes) != len(expectedTypes) {
		t.Errorf("expected %d nodes, got %d", len(expectedTypes), len(collector.nodeTypes))
	}

	for i, expected := range expectedTypes {
		if i < len(collector.nodeTypes) && collector.nodeTypes[i] != expected {
			t.Errorf("expected node type %d to be %v, got %v", i, expected, collector.nodeTypes[i])
		}
	}
}
