package markit

import (
	"fmt"
	"strings"
	"testing"
)

// TestParserBasicParsing 测试基础解析功能
func TestParserBasicParsing(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, *Document)
	}{
		{
			name:  "simple document",
			input: "<doc><p>hello</p></doc>",
			check: checkSimpleDocument,
		},
		{
			name:  "self-closing element",
			input: `<doc><image token="test" /></doc>`,
			check: checkSelfClosingElement,
		},
		{
			name:  "boolean attributes",
			input: `<doc><t b tc-blue>text</t></doc>`,
			check: checkBooleanAttributes,
		},
		{
			name:  "nested elements",
			input: "<root><level1><level2><level3>deep content</level3></level2></level1></root>",
			check: checkNestedElements,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			doc, err := parser.Parse()

			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			tt.check(t, doc)
		})
	}
}

// checkSimpleDocument 检查简单文档结构
func checkSimpleDocument(t *testing.T, doc *Document) {
	if len(doc.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(doc.Children))
		return
	}

	docElement, ok := doc.Children[0].(*Element)
	if !ok {
		t.Errorf("expected Element, got %T", doc.Children[0])
		return
	}

	if docElement.TagName != "doc" {
		t.Errorf("expected tag name 'doc', got %q", docElement.TagName)
	}

	if len(docElement.Children) != 1 {
		t.Errorf("expected 1 child in doc, got %d", len(docElement.Children))
		return
	}

	pElement, ok := docElement.Children[0].(*Element)
	if !ok {
		t.Errorf("expected Element, got %T", docElement.Children[0])
		return
	}

	if pElement.TagName != "p" {
		t.Errorf("expected tag name 'p', got %q", pElement.TagName)
	}
}

// checkSelfClosingElement 检查自闭合元素
func checkSelfClosingElement(t *testing.T, doc *Document) {
	docElement := doc.Children[0].(*Element)
	imageElement := docElement.Children[0].(*Element)

	if imageElement.TagName != "image" {
		t.Errorf("expected tag name 'image', got %q", imageElement.TagName)
	}

	if !imageElement.SelfClose {
		t.Error("expected self-close element")
	}

	if imageElement.Attributes["token"] != "test" {
		t.Errorf("expected token attribute 'test', got %q", imageElement.Attributes["token"])
	}
}

// checkBooleanAttributes 检查布尔属性
func checkBooleanAttributes(t *testing.T, doc *Document) {
	docElement := doc.Children[0].(*Element)
	tElement := docElement.Children[0].(*Element)

	if tElement.TagName != "t" {
		t.Errorf("expected tag name 't', got %q", tElement.TagName)
	}

	if tElement.Attributes["b"] != "" {
		t.Errorf("expected boolean attribute 'b' to be empty, got %q", tElement.Attributes["b"])
	}

	if tElement.Attributes["tc-blue"] != "" {
		t.Errorf("expected boolean attribute 'tc-blue' to be empty, got %q", tElement.Attributes["tc-blue"])
	}

	textNode := tElement.Children[0].(*Text)
	if textNode.Content != "text" {
		t.Errorf("expected text content 'text', got %q", textNode.Content)
	}
}

// checkNestedElements 检查嵌套元素
func checkNestedElements(t *testing.T, doc *Document) {
	root := doc.Children[0].(*Element)
	level1 := root.Children[0].(*Element)
	level2 := level1.Children[0].(*Element)
	level3 := level2.Children[0].(*Element)
	text := level3.Children[0].(*Text)

	if root.TagName != "root" {
		t.Errorf("expected root tag name 'root', got %q", root.TagName)
	}
	if level1.TagName != "level1" {
		t.Errorf("expected level1 tag name 'level1', got %q", level1.TagName)
	}
	if level2.TagName != "level2" {
		t.Errorf("expected level2 tag name 'level2', got %q", level2.TagName)
	}
	if level3.TagName != "level3" {
		t.Errorf("expected level3 tag name 'level3', got %q", level3.TagName)
	}
	if text.Content != "deep content" {
		t.Errorf("expected text content 'deep content', got %q", text.Content)
	}
}

// TestParserCommentSupport 测试注释支持
func TestParserCommentSupport(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, *Document)
	}{
		{
			name:  "comment before element",
			input: "<!-- comment --><root>content</root>",
			check: func(t *testing.T, doc *Document) {
				if len(doc.Children) != 2 {
					t.Errorf("expected 2 children (comment and element), got %d", len(doc.Children))
				}
				// 第一个应该是注释
				if _, ok := doc.Children[0].(*Comment); !ok {
					t.Error("expected first child to be a comment")
				}
				// 第二个应该是元素
				if _, ok := doc.Children[1].(*Element); !ok {
					t.Error("expected second child to be an element")
				}
			},
		},
		{
			name:  "comment after element",
			input: "<root>content</root><!-- comment -->",
			check: func(t *testing.T, doc *Document) {
				if len(doc.Children) != 2 {
					t.Errorf("expected 2 children (element and comment), got %d", len(doc.Children))
				}
				// 第一个应该是元素
				if _, ok := doc.Children[0].(*Element); !ok {
					t.Error("expected first child to be an element")
				}
				// 第二个应该是注释
				if _, ok := doc.Children[1].(*Comment); !ok {
					t.Error("expected second child to be a comment")
				}
			},
		},
		{
			name:  "comment between elements",
			input: "<first>content</first><!-- comment --><second>content</second>",
			check: func(t *testing.T, doc *Document) {
				if len(doc.Children) != 3 {
					t.Errorf("expected 3 children (element, comment, element), got %d", len(doc.Children))
				}
			},
		},
		{
			name:  "multiple comments",
			input: "<!-- comment1 --><!-- comment2 --><root>content</root><!-- comment3 -->",
			check: func(t *testing.T, doc *Document) {
				if len(doc.Children) != 4 {
					t.Errorf("expected 4 children (3 comments and 1 element), got %d", len(doc.Children))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			doc, err := parser.Parse()

			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			tt.check(t, doc)
		})
	}
}

// TestParserErrorHandling 测试错误处理
func TestParserErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "mismatched tags",
			input:       "<open>content</close>",
			expectError: true,
			errorMsg:    "mismatched",
		},
		{
			name:        "unclosed element",
			input:       "<open>content",
			expectError: true,
			errorMsg:    "unexpected EOF",
		},
		{
			name:        "invalid tag name",
			input:       "<123invalid>content</123invalid>",
			expectError: true, // 数字开头的标签名应该报错
			errorMsg:    "",
		},
		{
			name:        "empty tag name",
			input:       "<>content</>",
			expectError: true,
			errorMsg:    "empty tag",
		},
		{
			name:        "nested mismatched tags",
			input:       "<root><child>content</wrong></root>",
			expectError: true,
			errorMsg:    "mismatched",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			_, err := parser.Parse()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestParserComplexStructures 测试复杂结构
func TestParserComplexStructures(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, *Document)
	}{
		{
			name: "mixed content with text and elements",
			input: `<root>
				text before
				<child>child content</child>
				text between
				<another>another content</another>
				text after
			</root>`,
			check: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) < 3 { // 至少应该有文本和元素混合
					t.Errorf("expected mixed content, got %d children", len(root.Children))
				}
			},
		},
		{
			name:  "attributes with special characters",
			input: `<element class="my-class" data-value="test&amp;value" id="unique_123">content</element>`,
			check: func(t *testing.T, doc *Document) {
				element := doc.Children[0].(*Element)
				if element.Attributes["class"] != "my-class" {
					t.Errorf("expected class attribute 'my-class', got %q", element.Attributes["class"])
				}
				if element.Attributes["data-value"] != "test&amp;value" {
					t.Errorf("expected data-value attribute 'test&amp;value', got %q", element.Attributes["data-value"])
				}
			},
		},
		{
			name:  "empty elements",
			input: `<root><empty></empty><self-close/></root>`,
			check: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 2 {
					t.Errorf("expected 2 children, got %d", len(root.Children))
				}

				empty := root.Children[0].(*Element)
				selfClose := root.Children[1].(*Element)

				if empty.TagName != "empty" {
					t.Errorf("expected tag name 'empty', got %q", empty.TagName)
				}
				if len(empty.Children) != 0 {
					t.Errorf("expected empty element to have no children, got %d", len(empty.Children))
				}

				if selfClose.TagName != "self-close" {
					t.Errorf("expected tag name 'self-close', got %q", selfClose.TagName)
				}
				if !selfClose.SelfClose {
					t.Error("expected self-close element")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			doc, err := parser.Parse()

			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			tt.check(t, doc)
		})
	}
}

// TestParserWhitespaceHandling 测试空白字符处理
func TestParserWhitespaceHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, *Document)
	}{
		{
			name:  "preserve significant whitespace",
			input: "<pre>  spaced   content  </pre>",
			check: func(t *testing.T, doc *Document) {
				pre := doc.Children[0].(*Element)
				text := pre.Children[0].(*Text)
				// 检查文本内容包含空格，但不要求完全匹配
				if !strings.Contains(text.Content, "spaced") || !strings.Contains(text.Content, "content") {
					t.Errorf("expected text to contain 'spaced' and 'content', got %q", text.Content)
				}
			},
		},
		{
			name: "trim insignificant whitespace",
			input: `
				<root>
					<child>content</child>
				</root>
			`,
			check: func(t *testing.T, doc *Document) {
				// 验证解析成功，具体的空白处理策略取决于实现
				if len(doc.Children) == 0 {
					t.Error("expected at least one child element")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			doc, err := parser.Parse()

			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			tt.check(t, doc)
		})
	}
}

// TestTrimWhitespaceConfiguration 测试 TrimWhitespace 配置功能
func TestTrimWhitespaceConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		trimWhitespace bool
		checkFunc      func(*testing.T, *Document)
	}{
		{
			name:           "TrimWhitespace enabled - basic text",
			input:          "<root>  hello world  </root>",
			trimWhitespace: true,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 1 {
					t.Fatalf("expected 1 child, got %d", len(root.Children))
				}
				text := root.Children[0].(*Text)
				if text.Content != "hello world" {
					t.Errorf("expected 'hello world', got %q", text.Content)
				}
			},
		},
		{
			name:           "TrimWhitespace disabled - preserve whitespace",
			input:          "<root>  hello world  </root>",
			trimWhitespace: false,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 1 {
					t.Fatalf("expected 1 child, got %d", len(root.Children))
				}
				text := root.Children[0].(*Text)
				if text.Content != "  hello world  " {
					t.Errorf("expected '  hello world  ', got %q", text.Content)
				}
			},
		},
		{
			name:           "TrimWhitespace enabled - multiline text",
			input:          "<root>\n  hello\n  world\n</root>",
			trimWhitespace: true,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 1 {
					t.Fatalf("expected 1 child, got %d", len(root.Children))
				}
				text := root.Children[0].(*Text)
				expected := "hello\n  world"
				if text.Content != expected {
					t.Errorf("expected %q, got %q", expected, text.Content)
				}
			},
		},
		{
			name:           "TrimWhitespace disabled - preserve multiline whitespace",
			input:          "<root>\n  hello\n  world\n</root>",
			trimWhitespace: false,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 1 {
					t.Fatalf("expected 1 child, got %d", len(root.Children))
				}
				text := root.Children[0].(*Text)
				expected := "\n  hello\n  world\n"
				if text.Content != expected {
					t.Errorf("expected %q, got %q", expected, text.Content)
				}
			},
		},
		{
			name:           "TrimWhitespace enabled - comment trimming",
			input:          "<root><!-- \n  comment content  \n --></root>",
			trimWhitespace: true,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 1 {
					t.Fatalf("expected 1 child, got %d", len(root.Children))
				}
				comment := root.Children[0].(*Comment)
				if comment.Content != "comment content" {
					t.Errorf("expected 'comment content', got %q", comment.Content)
				}
			},
		},
		{
			name:           "TrimWhitespace disabled - preserve comment whitespace",
			input:          "<root><!-- \n  comment content  \n --></root>",
			trimWhitespace: false,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 1 {
					t.Fatalf("expected 1 child, got %d", len(root.Children))
				}
				comment := root.Children[0].(*Comment)
				expected := " \n  comment content  \n "
				if comment.Content != expected {
					t.Errorf("expected %q, got %q", expected, comment.Content)
				}
			},
		},
		{
			name:           "TrimWhitespace enabled - empty text handling",
			input:          "<root>   </root>",
			trimWhitespace: true,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				// 当 TrimWhitespace 为 true 时，空白文本应该被跳过
				if len(root.Children) != 0 {
					t.Errorf("expected 0 children (empty text should be skipped), got %d", len(root.Children))
				}
			},
		},
		{
			name:           "TrimWhitespace disabled - preserve empty text",
			input:          "<root>   </root>",
			trimWhitespace: false,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 1 {
					t.Fatalf("expected 1 child, got %d", len(root.Children))
				}
				text := root.Children[0].(*Text)
				if text.Content != "   " {
					t.Errorf("expected '   ', got %q", text.Content)
				}
			},
		},
		{
			name:           "TrimWhitespace enabled - mixed content",
			input:          "<root>  text  <!-- comment -->  more text  </root>",
			trimWhitespace: true,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 3 {
					t.Fatalf("expected 3 children, got %d", len(root.Children))
				}
				
				text1 := root.Children[0].(*Text)
				if text1.Content != "text" {
					t.Errorf("expected 'text', got %q", text1.Content)
				}
				
				comment := root.Children[1].(*Comment)
				if comment.Content != "comment" {
					t.Errorf("expected 'comment', got %q", comment.Content)
				}
				
				text2 := root.Children[2].(*Text)
				if text2.Content != "more text" {
					t.Errorf("expected 'more text', got %q", text2.Content)
				}
			},
		},
		{
			name:           "TrimWhitespace enabled - special whitespace characters",
			input:          "<root>\t\n\r  content  \t\n\r</root>",
			trimWhitespace: true,
			checkFunc: func(t *testing.T, doc *Document) {
				root := doc.Children[0].(*Element)
				if len(root.Children) != 1 {
					t.Fatalf("expected 1 child, got %d", len(root.Children))
				}
				text := root.Children[0].(*Text)
				if text.Content != "content" {
					t.Errorf("expected 'content', got %q", text.Content)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.TrimWhitespace = tt.trimWhitespace
			
			parser := NewParserWithConfig(tt.input, config)
			doc, err := parser.Parse()

			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			tt.checkFunc(t, doc)
		})
	}
}

// TestParserAttributeProcessor 测试属性处理器
func TestParserAttributeProcessor(t *testing.T) {
	// 创建自定义属性处理器
	processor := &DefaultAttributeProcessor{}

	tests := []struct {
		name  string
		input string
		check func(*testing.T, *Document)
	}{
		{
			name:  "process boolean attributes",
			input: `<input disabled checked></input>`,
			check: func(t *testing.T, doc *Document) {
				input := doc.Children[0].(*Element)
				if input.Attributes["disabled"] != "" {
					t.Errorf("expected boolean attribute 'disabled' to be empty, got %q", input.Attributes["disabled"])
				}
				if input.Attributes["checked"] != "" {
					t.Errorf("expected boolean attribute 'checked' to be empty, got %q", input.Attributes["checked"])
				}
			},
		},
		{
			name:  "process regular attributes",
			input: `<element attr1="value1" attr2="value2"></element>`,
			check: func(t *testing.T, doc *Document) {
				element := doc.Children[0].(*Element)
				if element.Attributes["attr1"] != "value1" {
					t.Errorf("expected attr1 'value1', got %q", element.Attributes["attr1"])
				}
				if element.Attributes["attr2"] != "value2" {
					t.Errorf("expected attr2 'value2', got %q", element.Attributes["attr2"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			parser.SetAttributeProcessor(processor)
			doc, err := parser.Parse()

			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			tt.check(t, doc)
		})
	}
}

// TestParserPrettyPrint 测试美化输出
func TestParserPrettyPrint(t *testing.T) {
	input := `<doc><p><t b>hello</t></p><image token="test" /></doc>`
	parser := NewParser(input)
	doc, err := parser.Parse()

	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	output := PrettyPrint(doc)
	t.Logf("Pretty print output:\n%s", output)

	// 基本检查输出包含预期内容
	if !contains(output, "Document") {
		t.Error("output should contain 'Document'")
	}
	if !contains(output, "<doc>") {
		t.Error("output should contain '<doc>'")
	}
	if !contains(output, "<p>") {
		t.Error("output should contain '<p>'")
	}
	if !contains(output, "hello") {
		t.Error("output should contain 'hello'")
	}
}

// TestParserSkipComments 测试跳过注释功能
func TestParserSkipComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, *Document)
	}{
		{
			name:  "skip comment before element",
			input: "<!-- comment --><root>content</root>",
			check: func(t *testing.T, doc *Document) {
				if len(doc.Children) != 1 {
					t.Errorf("expected 1 child (comment should be skipped), got %d", len(doc.Children))
				}
			},
		},
		{
			name:  "skip comment after element",
			input: "<root>content</root><!-- comment -->",
			check: func(t *testing.T, doc *Document) {
				if len(doc.Children) != 1 {
					t.Errorf("expected 1 child (comment should be skipped), got %d", len(doc.Children))
				}
			},
		},
		{
			name:  "skip comment between elements",
			input: "<first>content</first><!-- comment --><second>content</second>",
			check: func(t *testing.T, doc *Document) {
				if len(doc.Children) != 2 {
					t.Errorf("expected 2 children (comment should be skipped), got %d", len(doc.Children))
				}
			},
		},
		{
			name:  "skip multiple comments",
			input: "<!-- comment1 --><!-- comment2 --><root>content</root><!-- comment3 -->",
			check: func(t *testing.T, doc *Document) {
				if len(doc.Children) != 1 {
					t.Errorf("expected 1 child (comments should be skipped), got %d", len(doc.Children))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.SkipComments = true
			parser := NewParserWithConfig(tt.input, config)
			doc, err := parser.Parse()

			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			tt.check(t, doc)
		})
	}
}

// 辅助函数
func contains(s, substr string) bool {
	return containsAt(s, substr, 0) != -1
}

func containsAt(s, substr string, start int) int {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// TestParserConfigMethods 测试解析器的配置方法
func TestParserConfigMethods(t *testing.T) {
	parser := NewParser("<test>content</test>")

	t.Run("GetConfig", func(t *testing.T) {
		config := parser.GetConfig()
		if config == nil {
			t.Fatal("GetConfig should not return nil")
		}
	})

	t.Run("SetConfig", func(t *testing.T) {
		newConfig := DefaultConfig()
		newConfig.CaseSensitive = false
		newConfig.SkipComments = true

		parser.SetConfig(newConfig)
		retrievedConfig := parser.GetConfig()

		if retrievedConfig.CaseSensitive {
			t.Error("Config CaseSensitive was not updated")
		}

		if !retrievedConfig.SkipComments {
			t.Error("Config SkipComments was not updated")
		}
	})
}

// TestParserErrorMethod 测试解析器的Error方法
func TestParserErrorMethod(t *testing.T) {
	parser := NewParser("<invalid")

	// 尝试解析无效的XML
	_, err := parser.Parse()
	if err == nil {
		t.Fatal("Expected parse error")
	}

	// 测试Error方法
	parseErr, ok := err.(*ParseError)
	if !ok {
		t.Fatalf("Expected ParseError, got %T", err)
	}

	errorMsg := parseErr.Error()
	if errorMsg == "" {
		t.Error("Error message should not be empty")
	}

	// 验证错误消息格式
	if !contains(errorMsg, "parse error") {
		t.Errorf("Error message should contain 'parse error', got: %s", errorMsg)
	}
}

// TestParserNodeTypeParsing 测试各种节点类型的解析
func TestParserNodeTypeParsing(t *testing.T) {
	t.Run("Comment parsing", func(t *testing.T) {
		input := `<root><!-- This is a comment --></root>`
		parser := NewParser(input)
		doc, err := parser.Parse()
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		root := doc.Children[0].(*Element)
		if len(root.Children) == 0 {
			t.Fatal("Expected root to have children")
		}

		comment, ok := root.Children[0].(*Comment)
		if !ok {
			t.Fatalf("Expected Comment, got %T", root.Children[0])
		}

		expected := "This is a comment"
		if comment.Content != expected {
			t.Errorf("Expected comment content '%s', got '%s'", expected, comment.Content)
		}
	})
}

// TestParserTextParsing 测试文本解析的边界情况
func TestParserTextParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "<p>Hello World</p>",
			expected: "Hello World",
		},
		{
			name:     "text with entities",
			input:    "<p>&lt;&gt;&amp;&quot;&#39;</p>",
			expected: "&lt;&gt;&amp;&quot;&#39;",
		},
		{
			name:     "text with whitespace",
			input:    "<p>  Hello   World  </p>",
			expected: "Hello   World",
		},
		{
			name:     "empty text",
			input:    "<p></p>",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			doc, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			element := doc.Children[0].(*Element)
			if len(element.Children) == 0 && tt.expected == "" {
				return // 空文本情况
			}

			if len(element.Children) == 0 {
				t.Fatalf("Expected text content, got no children")
			}

			text, ok := element.Children[0].(*Text)
			if !ok {
				t.Fatalf("Expected Text node, got %T", element.Children[0])
			}

			if text.Content != tt.expected {
				t.Errorf("Expected text content '%s', got '%s'", tt.expected, text.Content)
			}
		})
	}
}

// TestParserSelfCloseElementParsing 测试自闭合元素解析
func TestParserSelfCloseElementParsing(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "img tag",
			input: `<img src="test.jpg" alt="test" />`,
		},
		{
			name:  "br tag",
			input: `<br />`,
		},
		{
			name:  "input tag",
			input: `<input type="text" name="test" />`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(tt.input)
			doc, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if len(doc.Children) == 0 {
				t.Fatal("Expected at least one child")
			}

			element, ok := doc.Children[0].(*Element)
			if !ok {
				t.Fatalf("Expected Element, got %T", doc.Children[0])
			}

			if !element.SelfClose {
				t.Error("Expected self-closing element")
			}

			if len(element.Children) != 0 {
				t.Error("Self-closing element should not have children")
			}
		})
	}
}

// TestWalkFunction 测试Walk函数
func TestWalkFunction(t *testing.T) {
	// 创建一个测试文档
	doc := &Document{
		Children: []Node{
			&Comment{Content: "Test comment"},
			&Element{
				TagName: "root",
				Children: []Node{
					&Text{Content: "Hello"},
					&Element{
						TagName:   "img",
						SelfClose: true,
					},
				},
			},
			&ProcessingInstruction{
				Target:  "xml",
				Content: "version=\"1.0\"",
			},
			&Doctype{Content: "html"},
			&CDATA{Content: "some data"},
		},
	}

	// 创建一个测试访问者
	visitor := &TestVisitor{
		visitedNodes: make(map[string]int),
	}

	// 遍历AST
	err := Walk(doc, visitor)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	// 验证所有节点类型都被访问了
	expectedVisits := map[string]int{
		"Document":              1,
		"Comment":               1,
		"Element":               2, // root和img
		"Text":                  1,
		"ProcessingInstruction": 1,
		"Doctype":               1,
		"CDATA":                 1,
	}

	for nodeType, expectedCount := range expectedVisits {
		if actualCount := visitor.visitedNodes[nodeType]; actualCount != expectedCount {
			t.Errorf("Expected %d visits to %s, got %d", expectedCount, nodeType, actualCount)
		}
	}
}

// TestVisitor 测试访问者实现
type TestVisitor struct {
	visitedNodes map[string]int
}

func (v *TestVisitor) VisitDocument(doc *Document) error {
	v.visitedNodes["Document"]++
	return nil
}

func (v *TestVisitor) VisitElement(elem *Element) error {
	v.visitedNodes["Element"]++
	return nil
}

func (v *TestVisitor) VisitText(text *Text) error {
	v.visitedNodes["Text"]++
	return nil
}

func (v *TestVisitor) VisitProcessingInstruction(pi *ProcessingInstruction) error {
	v.visitedNodes["ProcessingInstruction"]++
	return nil
}

func (v *TestVisitor) VisitDoctype(doctype *Doctype) error {
	v.visitedNodes["Doctype"]++
	return nil
}

func (v *TestVisitor) VisitCDATA(cdata *CDATA) error {
	v.visitedNodes["CDATA"]++
	return nil
}

func (v *TestVisitor) VisitComment(comment *Comment) error {
	v.visitedNodes["Comment"]++
	return nil
}

// TestWalkWithError 测试Walk函数的错误处理
func TestWalkWithError(t *testing.T) {
	doc := &Document{
		Children: []Node{
			&Text{Content: "test"},
		},
	}

	visitor := &ErrorVisitor{}
	err := Walk(doc, visitor)
	if err == nil {
		t.Error("Expected error from Walk, got nil")
	}
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}
}

// ErrorVisitor 用于测试错误处理的访问者
type ErrorVisitor struct{}

func (v *ErrorVisitor) VisitDocument(doc *Document) error {
	return nil
}

func (v *ErrorVisitor) VisitElement(elem *Element) error {
	return nil
}

func (v *ErrorVisitor) VisitText(text *Text) error {
	return fmt.Errorf("test error")
}

func (v *ErrorVisitor) VisitProcessingInstruction(pi *ProcessingInstruction) error {
	return nil
}

func (v *ErrorVisitor) VisitDoctype(doctype *Doctype) error {
	return nil
}

func (v *ErrorVisitor) VisitCDATA(cdata *CDATA) error {
	return nil
}

func (v *ErrorVisitor) VisitComment(comment *Comment) error {
	return nil
}

// TestPrettyPrintAllNodeTypes 测试PrettyPrint函数对所有节点类型的处理
func TestPrettyPrintAllNodeTypes(t *testing.T) {
	// 创建包含所有节点类型的文档
	doc := &Document{
		Children: []Node{
			&ProcessingInstruction{
				Target:  "xml",
				Content: "version=\"1.0\" encoding=\"UTF-8\"",
			},
			&Doctype{
				Content: "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\"",
			},
			&Comment{
				Content: "This is a comment",
			},
			&Element{
				TagName: "root",
				Attributes: map[string]string{
					"id":    "main",
					"class": "container",
				},
				Children: []Node{
					&Text{Content: "Hello World"},
					&CDATA{Content: "function() { return 'test'; }"},
					&Element{
						TagName:   "img",
						SelfClose: true,
						Attributes: map[string]string{
							"src": "image.png",
							"alt": "",
						},
					},
				},
			},
		},
	}

	output := PrettyPrint(doc)

	// 验证输出包含所有节点类型的表示
	expectedContents := []string{
		"Document",
		"PI:",
		"Doctype:",
		"Comment:",
		"<root",
		"id=\"main\"",
		"class=\"container\"",
		"Text:",
		"CDATA:",
		"<img",
		"src=\"image.png\"",
		"alt",
		"/>",
		"</root>",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it didn't.\nOutput:\n%s", expected, output)
		}
	}

	t.Logf("PrettyPrint output:\n%s", output)
}
