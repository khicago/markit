package markit

import (
	"testing"
)

// TestASTNodeTypes 测试AST节点类型
func TestASTNodeTypes(t *testing.T) {
	// 创建测试文档
	doc := &Document{
		Children: []Node{
			&Element{
				TagName: "root",
				Attributes: map[string]string{
					"id": "main",
				},
				Children: []Node{
					&Text{Content: "Hello World"},
					&Element{
						TagName:   "img",
						SelfClose: true,
						Attributes: map[string]string{
							"src": "image.png",
							"alt": "test image",
						},
					},
				},
			},
		},
	}

	// 测试文档节点
	if len(doc.Children) != 1 {
		t.Errorf("expected 1 child in document, got %d", len(doc.Children))
	}

	if doc.Type() != NodeTypeDocument {
		t.Errorf("expected NodeTypeDocument, got %v", doc.Type())
	}

	if doc.String() != "Document" {
		t.Errorf("expected 'Document', got %q", doc.String())
	}

	// 测试元素节点
	rootElement, ok := doc.Children[0].(*Element)
	if !ok {
		t.Fatalf("expected Element, got %T", doc.Children[0])
	}

	if rootElement.Type() != NodeTypeElement {
		t.Errorf("expected NodeTypeElement, got %v", rootElement.Type())
	}

	if rootElement.TagName != "root" {
		t.Errorf("expected tag name 'root', got %q", rootElement.TagName)
	}

	if rootElement.String() != "root" {
		t.Errorf("expected 'root', got %q", rootElement.String())
	}

	if rootElement.Attributes["id"] != "main" {
		t.Errorf("expected id attribute 'main', got %q", rootElement.Attributes["id"])
	}

	if len(rootElement.Children) != 2 {
		t.Errorf("expected 2 children in root element, got %d", len(rootElement.Children))
	}

	// 测试文本节点
	textNode, ok := rootElement.Children[0].(*Text)
	if !ok {
		t.Fatalf("expected Text node, got %T", rootElement.Children[0])
	}

	if textNode.Type() != NodeTypeText {
		t.Errorf("expected NodeTypeText, got %v", textNode.Type())
	}

	if textNode.Content != "Hello World" {
		t.Errorf("expected text content 'Hello World', got %q", textNode.Content)
	}

	if textNode.String() != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", textNode.String())
	}

	// 测试自闭合元素
	imgElement, ok := rootElement.Children[1].(*Element)
	if !ok {
		t.Fatalf("expected Element, got %T", rootElement.Children[1])
	}

	if imgElement.TagName != "img" {
		t.Errorf("expected tag name 'img', got %q", imgElement.TagName)
	}

	if !imgElement.SelfClose {
		t.Error("expected self-close element")
	}

	if imgElement.Attributes["src"] != "image.png" {
		t.Errorf("expected src attribute 'image.png', got %q", imgElement.Attributes["src"])
	}

	if imgElement.Attributes["alt"] != "test image" {
		t.Errorf("expected alt attribute 'test image', got %q", imgElement.Attributes["alt"])
	}
}

// TestAllNodeTypes 测试所有节点类型的基本方法
func TestAllNodeTypes(t *testing.T) {
	t.Run("Document methods", func(t *testing.T) {
		doc := &Document{
			Pos:      Position{Line: 1, Column: 1, Offset: 0},
			Children: []Node{},
		}

		if doc.Type() != NodeTypeDocument {
			t.Errorf("Expected NodeTypeDocument, got %v", doc.Type())
		}

		pos := doc.Position()
		if pos.Line != 1 || pos.Column != 1 || pos.Offset != 0 {
			t.Errorf("Expected position 1:1:0, got %d:%d:%d", pos.Line, pos.Column, pos.Offset)
		}

		if doc.String() != "Document" {
			t.Errorf("Expected 'Document', got '%s'", doc.String())
		}
	})

	t.Run("Comment methods", func(t *testing.T) {
		comment := &Comment{
			Content: "This is a comment with some text",
			Pos:     Position{Line: 3, Column: 5, Offset: 75},
		}

		if comment.Type() != NodeTypeComment {
			t.Errorf("Expected NodeTypeComment, got %v", comment.Type())
		}

		pos := comment.Position()
		if pos.Line != 3 || pos.Column != 5 || pos.Offset != 75 {
			t.Errorf("Expected position 3:5:75, got %d:%d:%d", pos.Line, pos.Column, pos.Offset)
		}

		if comment.String() != "This is a comment with some text" {
			t.Errorf("Expected comment content, got '%s'", comment.String())
		}
	})

	t.Run("Text Position method", func(t *testing.T) {
		text := &Text{
			Content: "Some text content",
			Pos:     Position{Line: 4, Column: 8, Offset: 95},
		}

		pos := text.Position()
		if pos.Line != 4 || pos.Column != 8 || pos.Offset != 95 {
			t.Errorf("Expected position 4:8:95, got %d:%d:%d", pos.Line, pos.Column, pos.Offset)
		}
	})

	t.Run("Element Position method", func(t *testing.T) {
		element := &Element{
			TagName: "div",
			Pos:     Position{Line: 2, Column: 3, Offset: 25},
		}

		pos := element.Position()
		if pos.Line != 2 || pos.Column != 3 || pos.Offset != 25 {
			t.Errorf("Expected position 2:3:25, got %d:%d:%d", pos.Line, pos.Column, pos.Offset)
		}
	})
}

// TestProcessingInstructionNode 测试ProcessingInstruction节点
func TestProcessingInstructionNode(t *testing.T) {
	pi := &ProcessingInstruction{
		Target:  "xml",
		Content: "version=\"1.0\" encoding=\"UTF-8\"",
		Pos:     Position{Line: 1, Column: 1, Offset: 0},
	}

	if pi.Type() != NodeTypeProcessingInstruction {
		t.Errorf("Expected NodeTypeProcessingInstruction, got %v", pi.Type())
	}

	pos := pi.Position()
	if pos.Line != 1 || pos.Column != 1 || pos.Offset != 0 {
		t.Errorf("Expected position 1:1:0, got %d:%d:%d", pos.Line, pos.Column, pos.Offset)
	}

	if pi.String() != "xml" {
		t.Errorf("Expected PI target, got '%s'", pi.String())
	}
}

// TestDoctypeNode 测试Doctype节点
func TestDoctypeNode(t *testing.T) {
	doctype := &Doctype{
		Content: "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\"",
		Pos:     Position{Line: 2, Column: 1, Offset: 50},
	}

	if doctype.Type() != NodeTypeDoctype {
		t.Errorf("Expected NodeTypeDoctype, got %v", doctype.Type())
	}

	pos := doctype.Position()
	if pos.Line != 2 || pos.Column != 1 || pos.Offset != 50 {
		t.Errorf("Expected position 2:1:50, got %d:%d:%d", pos.Line, pos.Column, pos.Offset)
	}

	if doctype.String() != "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\"" {
		t.Errorf("Expected DOCTYPE content, got '%s'", doctype.String())
	}
}

// TestCDATANode 测试CDATA节点
func TestCDATANode(t *testing.T) {
	cdata := &CDATA{
		Content: "function() { return 'Hello World'; }",
		Pos:     Position{Line: 5, Column: 10, Offset: 120},
	}

	if cdata.Type() != NodeTypeCDATA {
		t.Errorf("Expected NodeTypeCDATA, got %v", cdata.Type())
	}

	pos := cdata.Position()
	if pos.Line != 5 || pos.Column != 10 || pos.Offset != 120 {
		t.Errorf("Expected position 5:10:120, got %d:%d:%d", pos.Line, pos.Column, pos.Offset)
	}

	if cdata.String() != "function() { return 'Hello World'; }" {
		t.Errorf("Expected CDATA content, got '%s'", cdata.String())
	}
}
