package markit

import (
	"strings"
	"testing"
)

func TestCoreProtocols(t *testing.T) {
	protocols := GetCoreProtocols()

	if len(protocols) != 2 {
		t.Errorf("Expected 2 core protocols, got %d", len(protocols))
	}

	// 检查标准标签协议
	found := false
	for _, p := range protocols {
		if p.Name == "markit-standard-tag" {
			found = true
			if p.OpenSeq != "<" || p.CloseSeq != ">" {
				t.Errorf("Standard tag protocol has wrong sequences: open=%s, close=%s", p.OpenSeq, p.CloseSeq)
			}
		}
	}
	if !found {
		t.Error("Standard tag protocol not found")
	}

	// 检查注释协议
	found = false
	for _, p := range protocols {
		if p.Name == "markit-comment" {
			found = true
			if p.OpenSeq != "<!--" || p.CloseSeq != "-->" {
				t.Errorf("Comment protocol has wrong sequences: open=%s, close=%s", p.OpenSeq, p.CloseSeq)
			}
		}
	}
	if !found {
		t.Error("Comment protocol not found")
	}
}

func TestCoreProtocolMatcher(t *testing.T) {
	matcher := NewCoreProtocolMatcher()

	// 测试标签匹配
	protocol := matcher.MatchProtocol("<div>", 0)
	if protocol == nil {
		t.Error("Should match standard tag")
	} else if protocol.Name != "markit-standard-tag" {
		t.Errorf("Expected markit-standard-tag, got %s", protocol.Name)
	}

	// 测试注释匹配
	protocol = matcher.MatchProtocol("<!-- comment -->", 0)
	if protocol == nil {
		t.Error("Should match comment")
	} else if protocol.Name != "markit-comment" {
		t.Errorf("Expected markit-comment, got %s", protocol.Name)
	}

	// 测试不匹配
	protocol = matcher.MatchProtocol("plain text", 0)
	if protocol != nil {
		t.Error("Should not match plain text")
	}
}

func TestCoreParserBasic(t *testing.T) {
	input := "<root>Hello World</root>"
	parser := NewParser(input)

	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(doc.Children))
	}

	root, ok := doc.Children[0].(*Element)
	if !ok {
		t.Error("First child should be an element")
	} else if root.TagName != "root" {
		t.Errorf("Expected root element name 'root', got '%s'", root.TagName)
	}
}

func TestCoreParserWithComments(t *testing.T) {
	input := "<!-- comment --><root>content</root>"
	parser := NewParser(input)

	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(doc.Children))
	}

	// 第一个应该是注释
	comment, ok := doc.Children[0].(*Comment)
	if !ok {
		t.Error("First child should be a comment")
	} else if !strings.Contains(comment.Content, "comment") {
		t.Errorf("Comment content should contain 'comment', got '%s'", comment.Content)
	}

	// 第二个应该是元素
	element, ok := doc.Children[1].(*Element)
	if !ok {
		t.Error("Second child should be an element")
	} else if element.TagName != "root" {
		t.Errorf("Expected element name 'root', got '%s'", element.TagName)
	}
}

func TestCoreSelfClosingTagsAllowed(t *testing.T) {
	input := "<img src=\"test.jpg\" /><br/>"

	// 使用允许自封闭标签的配置
	config := DefaultConfig()
	config.AllowSelfCloseTags = true
	parser := NewParserWithConfig(input, config)

	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(doc.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(doc.Children))
	}

	// 第一个应该是自封闭的img标签
	img, ok := doc.Children[0].(*Element)
	if !ok {
		t.Error("First child should be an element")
	} else {
		if img.TagName != "img" {
			t.Errorf("Expected element name 'img', got '%s'", img.TagName)
		}
		if !img.SelfClose {
			t.Error("img element should be self-closing")
		}
		if img.Attributes["src"] != "test.jpg" {
			t.Errorf("Expected src attribute 'test.jpg', got '%s'", img.Attributes["src"])
		}
	}

	// 第二个应该是自封闭的br标签
	br, ok := doc.Children[1].(*Element)
	if !ok {
		t.Error("Second child should be an element")
	} else {
		if br.TagName != "br" {
			t.Errorf("Expected element name 'br', got '%s'", br.TagName)
		}
		if !br.SelfClose {
			t.Error("br element should be self-closing")
		}
	}
}

func TestCoreSelfClosingTagsNotAllowed(t *testing.T) {
	input := "<img src=\"test.jpg\" />"

	// 使用不允许自封闭标签的配置
	config := DefaultConfig()
	config.AllowSelfCloseTags = false
	parser := NewParserWithConfig(input, config)

	_, err := parser.Parse()
	if err == nil {
		t.Error("Expected parse error when self-closing tags are not allowed")
	}

	// 检查错误信息是否包含相关内容
	if !strings.Contains(err.Error(), "self-closing tags not allowed") {
		t.Errorf("Expected error about self-closing tags, got: %v", err)
	}
}
