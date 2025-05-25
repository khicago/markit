package markit

import (
	"fmt"
	"testing"
)

// TestParseNodeErrorHandling 测试parseNode函数的错误处理分支
func TestParseNodeErrorHandling(t *testing.T) {
	t.Run("Unknown token type", func(t *testing.T) {
		// 创建一个包含未知token类型的测试场景
		// 这需要通过模拟或者特殊输入来触发
		parser := NewParser("<div>content</div>")

		// 手动设置一个未知的token类型来测试default分支
		parser.current = Token{
			Type:     TokenType(999), // 未知的token类型
			Value:    "unknown",
			Position: Position{Line: 1, Column: 1},
		}

		_, err := parser.parseNode()
		if err == nil {
			t.Error("expected error for unknown token type, got nil")
		}

		parseErr, ok := err.(*ParseError)
		if !ok {
			t.Errorf("expected ParseError, got %T", err)
		}

		expectedMsg := "unexpected token UNKNOWN(999)"
		if parseErr.Message != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, parseErr.Message)
		}
	})

	t.Run("Error token handling", func(t *testing.T) {
		// 创建一个会产生错误token的输入
		parser := NewParser("<123invalid>")

		_, err := parser.Parse()
		if err == nil {
			t.Error("expected error for invalid tag name, got nil")
		}

		parseErr, ok := err.(*ParseError)
		if !ok {
			t.Errorf("expected ParseError, got %T", err)
		}

		if parseErr.Message != "invalid tag name" {
			t.Errorf("expected 'invalid tag name', got %q", parseErr.Message)
		}
	})

	t.Run("Skip comments configuration", func(t *testing.T) {
		config := DefaultConfig()
		config.SkipComments = true

		input := "<!-- comment --><div>content</div>"
		parser := NewParserWithConfig(input, config)

		doc, err := parser.Parse()
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		// 应该只有一个元素，注释被跳过了
		if len(doc.Children) != 1 {
			t.Errorf("expected 1 child (comment skipped), got %d", len(doc.Children))
		}

		element, ok := doc.Children[0].(*Element)
		if !ok {
			t.Errorf("expected Element, got %T", doc.Children[0])
		}

		if element.TagName != "div" {
			t.Errorf("expected tag name 'div', got %q", element.TagName)
		}
	})
}

// TestSpecialNodeTypeParsing 测试特殊节点类型的解析
func TestSpecialNodeTypeParsing(t *testing.T) {
	t.Run("ProcessingInstruction parsing", func(t *testing.T) {
		testProcessingInstructionParsing(t)
	})

	t.Run("Doctype parsing", func(t *testing.T) {
		testDoctypeParsing(t)
	})

	t.Run("CDATA parsing", func(t *testing.T) {
		testCDATAParsing(t)
	})

	t.Run("Wrong token type for ProcessingInstruction", func(t *testing.T) {
		testWrongTokenTypeForProcessingInstruction(t)
	})

	t.Run("Wrong token type for Doctype", func(t *testing.T) {
		testWrongTokenTypeForDoctype(t)
	})

	t.Run("Wrong token type for CDATA", func(t *testing.T) {
		testWrongTokenTypeForCDATA(t)
	})
}

// testProcessingInstructionParsing 测试ProcessingInstruction解析
func testProcessingInstructionParsing(t *testing.T) {
	parser := NewParser("")
	parser.current = Token{
		Type:     TokenProcessingInstruction,
		Value:    "xml version=\"1.0\"",
		Position: Position{Line: 1, Column: 1},
	}
	parser.peek = Token{Type: TokenEOF}

	node, err := parser.parseProcessingInstruction()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	pi, ok := node.(*ProcessingInstruction)
	if !ok {
		t.Fatalf("expected ProcessingInstruction, got %T", node)
	}

	if pi.Target != "xml version=\"1.0\"" {
		t.Errorf("expected target 'xml version=\"1.0\"', got %q", pi.Target)
	}

	if pi.Content != "xml version=\"1.0\"" {
		t.Errorf("expected content 'xml version=\"1.0\"', got %q", pi.Content)
	}
}

// testDoctypeParsing 测试Doctype解析
func testDoctypeParsing(t *testing.T) {
	parser := NewParser("")
	parser.current = Token{
		Type:     TokenDoctype,
		Value:    "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\"",
		Position: Position{Line: 1, Column: 1},
	}
	parser.peek = Token{Type: TokenEOF}

	node, err := parser.parseDoctype()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	doctype, ok := node.(*Doctype)
	if !ok {
		t.Fatalf("expected Doctype, got %T", node)
	}

	expected := "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\""
	if doctype.Content != expected {
		t.Errorf("expected content %q, got %q", expected, doctype.Content)
	}
}

// testCDATAParsing 测试CDATA解析
func testCDATAParsing(t *testing.T) {
	parser := NewParser("")
	parser.current = Token{
		Type:     TokenCDATA,
		Value:    "function() { return 'test'; }",
		Position: Position{Line: 1, Column: 1},
	}
	parser.peek = Token{Type: TokenEOF}

	node, err := parser.parseCDATA()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	cdata, ok := node.(*CDATA)
	if !ok {
		t.Fatalf("expected CDATA, got %T", node)
	}

	expected := "function() { return 'test'; }"
	if cdata.Content != expected {
		t.Errorf("expected content %q, got %q", expected, cdata.Content)
	}
}

// testWrongTokenTypeForProcessingInstruction 测试ProcessingInstruction错误token类型
func testWrongTokenTypeForProcessingInstruction(t *testing.T) {
	parser := NewParser("")
	parser.current = Token{
		Type:     TokenText, // 错误的token类型
		Value:    "text",
		Position: Position{Line: 1, Column: 1},
	}

	_, err := parser.parseProcessingInstruction()
	if err == nil {
		t.Error("expected error for wrong token type, got nil")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %T", err)
	}

	expectedMsg := "expected processing instruction token, got TEXT"
	if parseErr.Message != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, parseErr.Message)
	}
}

// testWrongTokenTypeForDoctype 测试Doctype错误token类型
func testWrongTokenTypeForDoctype(t *testing.T) {
	parser := NewParser("")
	parser.current = Token{
		Type:     TokenText,
		Value:    "text",
		Position: Position{Line: 1, Column: 1},
	}

	_, err := parser.parseDoctype()
	if err == nil {
		t.Error("expected error for wrong token type, got nil")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %T", err)
	}

	expectedMsg := "expected doctype token, got TEXT"
	if parseErr.Message != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, parseErr.Message)
	}
}

// testWrongTokenTypeForCDATA 测试CDATA错误token类型
func testWrongTokenTypeForCDATA(t *testing.T) {
	parser := NewParser("")
	parser.current = Token{
		Type:     TokenText,
		Value:    "text",
		Position: Position{Line: 1, Column: 1},
	}

	_, err := parser.parseCDATA()
	if err == nil {
		t.Error("expected error for wrong token type, got nil")
	}

	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %T", err)
	}

	expectedMsg := "expected CDATA token, got TEXT"
	if parseErr.Message != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, parseErr.Message)
	}
}

// TestParseTextErrorHandling 测试parseText函数的错误处理
func TestParseTextErrorHandling(t *testing.T) {
	t.Run("Wrong token type for text", func(t *testing.T) {
		parser := NewParser("")
		parser.current = Token{
			Type:     TokenOpenTag,
			Value:    "div",
			Position: Position{Line: 1, Column: 1},
		}

		_, err := parser.parseText()
		if err == nil {
			t.Error("expected error for wrong token type, got nil")
		}

		parseErr, ok := err.(*ParseError)
		if !ok {
			t.Errorf("expected ParseError, got %T", err)
		}

		expectedMsg := "expected text token, got OPEN_TAG"
		if parseErr.Message != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, parseErr.Message)
		}
	})
}

// TestParseSelfCloseElementErrorHandling 测试parseSelfCloseElement函数的错误处理
func TestParseSelfCloseElementErrorHandling(t *testing.T) {
	t.Run("Wrong token type for self-close element", func(t *testing.T) {
		parser := NewParser("")
		parser.current = Token{
			Type:     TokenOpenTag,
			Value:    "div",
			Position: Position{Line: 1, Column: 1},
		}

		_, err := parser.parseSelfCloseElement()
		if err == nil {
			t.Error("expected error for wrong token type, got nil")
		}

		parseErr, ok := err.(*ParseError)
		if !ok {
			t.Errorf("expected ParseError, got %T", err)
		}

		expectedMsg := "expected self-close tag, got OPEN_TAG"
		if parseErr.Message != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, parseErr.Message)
		}
	})
}

// TestParseCommentErrorHandling 测试parseComment函数的错误处理
func TestParseCommentErrorHandling(t *testing.T) {
	t.Run("Wrong token type for comment", func(t *testing.T) {
		parser := NewParser("")
		parser.current = Token{
			Type:     TokenText,
			Value:    "text",
			Position: Position{Line: 1, Column: 1},
		}

		_, err := parser.parseComment()
		if err == nil {
			t.Error("expected error for wrong token type, got nil")
		}

		parseErr, ok := err.(*ParseError)
		if !ok {
			t.Errorf("expected ParseError, got %T", err)
		}

		expectedMsg := "expected comment token, got TEXT"
		if parseErr.Message != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, parseErr.Message)
		}
	})
}

// TestWalkErrorBranches 测试Walk函数的错误分支
func TestWalkErrorBranches(t *testing.T) {
	t.Run("Error in VisitElement", func(t *testing.T) {
		element := &Element{
			TagName: "div",
			Children: []Node{
				&Text{Content: "test"},
			},
		}

		visitor := &ErrorOnElementVisitor{}
		err := Walk(element, visitor)

		if err == nil {
			t.Error("expected error from Walk, got nil")
		}

		if err.Error() != "error visiting element" {
			t.Errorf("expected 'error visiting element', got %q", err.Error())
		}
	})

	t.Run("Error in child traversal", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "div",
					Children: []Node{
						&Text{Content: "test"},
					},
				},
			},
		}

		visitor := &ErrorOnTextVisitor{}
		err := Walk(doc, visitor)

		if err == nil {
			t.Error("expected error from Walk, got nil")
		}

		if err.Error() != "error visiting text" {
			t.Errorf("expected 'error visiting text', got %q", err.Error())
		}
	})
}

// ErrorOnElementVisitor 在访问Element时返回错误的访问者
type ErrorOnElementVisitor struct{}

func (v *ErrorOnElementVisitor) VisitDocument(doc *Document) error { return nil }
func (v *ErrorOnElementVisitor) VisitElement(elem *Element) error {
	return fmt.Errorf("error visiting element")
}
func (v *ErrorOnElementVisitor) VisitText(text *Text) error { return nil }
func (v *ErrorOnElementVisitor) VisitProcessingInstruction(pi *ProcessingInstruction) error {
	return nil
}
func (v *ErrorOnElementVisitor) VisitDoctype(doctype *Doctype) error { return nil }
func (v *ErrorOnElementVisitor) VisitCDATA(cdata *CDATA) error       { return nil }
func (v *ErrorOnElementVisitor) VisitComment(comment *Comment) error { return nil }

// ErrorOnTextVisitor 在访问Text时返回错误的访问者
type ErrorOnTextVisitor struct{}

func (v *ErrorOnTextVisitor) VisitDocument(doc *Document) error { return nil }
func (v *ErrorOnTextVisitor) VisitElement(elem *Element) error  { return nil }
func (v *ErrorOnTextVisitor) VisitText(text *Text) error {
	return fmt.Errorf("error visiting text")
}
func (v *ErrorOnTextVisitor) VisitProcessingInstruction(pi *ProcessingInstruction) error { return nil }
func (v *ErrorOnTextVisitor) VisitDoctype(doctype *Doctype) error                        { return nil }
func (v *ErrorOnTextVisitor) VisitCDATA(cdata *CDATA) error                              { return nil }
func (v *ErrorOnTextVisitor) VisitComment(comment *Comment) error                        { return nil }
