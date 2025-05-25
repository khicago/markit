package markit

import (
	"errors"
	"testing"
)

// TestParseNodeSpecialTokenTypes 测试parseNode函数中的特殊token类型
func TestParseNodeSpecialTokenTypes(t *testing.T) {
	t.Run("ProcessingInstruction token", func(t *testing.T) {
		// 创建一个包含处理指令token的lexer
		lexer := &Lexer{}
		parser := &Parser{
			lexer: lexer,
			current: Token{
				Type:     TokenProcessingInstruction,
				Value:    "xml version=\"1.0\"",
				Position: Position{Line: 1, Column: 1},
			},
			config: &ParserConfig{
				CaseSensitive:      true,
				AllowSelfCloseTags: true,
				SkipComments:       false,
			},
		}

		node, err := parser.parseNode()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		pi, ok := node.(*ProcessingInstruction)
		if !ok {
			t.Fatalf("expected ProcessingInstruction, got %T", node)
		}

		if pi.Content != "xml version=\"1.0\"" {
			t.Errorf("expected content 'xml version=\"1.0\"', got %q", pi.Content)
		}
	})

	t.Run("Doctype token", func(t *testing.T) {
		lexer := &Lexer{}
		parser := &Parser{
			lexer: lexer,
			current: Token{
				Type:     TokenDoctype,
				Value:    "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\"",
				Position: Position{Line: 1, Column: 1},
			},
			config: &ParserConfig{
				CaseSensitive:      true,
				AllowSelfCloseTags: true,
				SkipComments:       false,
			},
		}

		node, err := parser.parseNode()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		doctype, ok := node.(*Doctype)
		if !ok {
			t.Fatalf("expected Doctype, got %T", node)
		}

		if doctype.Content != "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\"" {
			t.Errorf("expected content 'html PUBLIC \"-//W3C//DTD HTML 4.01//EN\"', got %q", doctype.Content)
		}
	})

	t.Run("CDATA token", func(t *testing.T) {
		lexer := &Lexer{}
		parser := &Parser{
			lexer: lexer,
			current: Token{
				Type:     TokenCDATA,
				Value:    "function() { return 'test'; }",
				Position: Position{Line: 1, Column: 1},
			},
			config: &ParserConfig{
				CaseSensitive:      true,
				AllowSelfCloseTags: true,
				SkipComments:       false,
			},
		}

		node, err := parser.parseNode()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cdata, ok := node.(*CDATA)
		if !ok {
			t.Fatalf("expected CDATA, got %T", node)
		}

		if cdata.Content != "function() { return 'test'; }" {
			t.Errorf("expected content 'function() { return 'test'; }', got %q", cdata.Content)
		}
	})
}

// TestParseElementErrorCondition 测试parseElement中的错误条件
func TestParseElementErrorCondition(t *testing.T) {
	t.Run("Non-open tag token", func(t *testing.T) {
		lexer := &Lexer{}
		parser := &Parser{
			lexer: lexer,
			current: Token{
				Type:     TokenText,
				Value:    "not a tag",
				Position: Position{Line: 1, Column: 1},
			},
			config: &ParserConfig{
				CaseSensitive:      true,
				AllowSelfCloseTags: true,
				SkipComments:       false,
			},
		}

		_, err := parser.parseElement()
		if err == nil {
			t.Fatal("expected error for non-open tag token")
		}

		parseErr, ok := err.(*ParseError)
		if !ok {
			t.Fatalf("expected ParseError, got %T", err)
		}

		expectedMsg := "expected open tag, got TEXT"
		if parseErr.Message != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, parseErr.Message)
		}
	})
}

// TestWalkVisitorDocumentError 测试Walk函数中VisitDocument的错误处理
func TestWalkVisitorDocumentError(t *testing.T) {
	// 创建一个会在VisitDocument时返回错误的visitor
	errorVisitor := &errorDocumentVisitor{
		shouldError: true,
	}

	// 创建一个简单的文档
	doc := &Document{
		Children: []Node{
			&Text{Content: "test", Pos: Position{Line: 1, Column: 1}},
		},
	}

	// 调用Walk函数
	err := Walk(doc, errorVisitor)
	if err == nil {
		t.Fatal("expected error from VisitDocument")
	}

	expectedMsg := "document visit error"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

// errorDocumentVisitor 是一个会在VisitDocument时返回错误的visitor
type errorDocumentVisitor struct {
	shouldError bool
}

func (v *errorDocumentVisitor) VisitDocument(doc *Document) error {
	if v.shouldError {
		return errors.New("document visit error")
	}
	return nil
}

func (v *errorDocumentVisitor) VisitElement(elem *Element) error {
	return nil
}

func (v *errorDocumentVisitor) VisitText(text *Text) error {
	return nil
}

func (v *errorDocumentVisitor) VisitProcessingInstruction(pi *ProcessingInstruction) error {
	return nil
}

func (v *errorDocumentVisitor) VisitDoctype(doctype *Doctype) error {
	return nil
}

func (v *errorDocumentVisitor) VisitCDATA(cdata *CDATA) error {
	return nil
}

func (v *errorDocumentVisitor) VisitComment(comment *Comment) error {
	return nil
}
