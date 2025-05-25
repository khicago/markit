package markit

import (
	"testing"
)

// TestLexerBasicFunctionality 测试词法分析器的基本功能
func TestLexerBasicFunctionality(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "Simple tag",
			input:    "<div>content</div>",
			expected: []TokenType{TokenOpenTag, TokenText, TokenCloseTag, TokenEOF},
		},
		{
			name:     "Self-closing tag",
			input:    "<img />",
			expected: []TokenType{TokenSelfCloseTag, TokenEOF},
		},
		{
			name:     "Comment",
			input:    "<!-- comment -->",
			expected: []TokenType{TokenComment, TokenEOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			for i, expectedType := range tt.expected {
				token := lexer.NextToken()
				if token.Type != expectedType {
					t.Errorf("token %d: expected type %v, got %v", i, expectedType, token.Type)
				}
			}
		})
	}
}

// TestLexerPositionTracking 测试词法分析器的位置跟踪
func TestLexerPositionTracking(t *testing.T) {
	input := `<div>
	<span>text</span>
</div>`

	lexer := NewLexer(input)

	// 第一个token应该在第1行第1列
	token := lexer.NextToken()
	if token.Position.Line != 1 || token.Position.Column != 1 {
		t.Errorf("expected position (1,1), got (%d,%d)", token.Position.Line, token.Position.Column)
	}

	// 继续读取tokens并验证位置
	for token.Type != TokenEOF {
		token = lexer.NextToken()
		// 验证位置信息存在且合理
		if token.Position.Line < 1 || token.Position.Column < 1 {
			t.Errorf("invalid position (%d,%d) for token %v", token.Position.Line, token.Position.Column, token.Type)
		}
	}
}

// TestLexerAttributeParsing 测试属性解析
func TestLexerAttributeParsing(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tagName    string
		attributes map[string]string
	}{
		{
			name:    "Single attribute",
			input:   `<div class="container">`,
			tagName: "div",
			attributes: map[string]string{
				"class": "container",
			},
		},
		{
			name:    "Multiple attributes",
			input:   `<img src="image.jpg" alt="description" width="100">`,
			tagName: "img",
			attributes: map[string]string{
				"src":   "image.jpg",
				"alt":   "description",
				"width": "100",
			},
		},
		{
			name:    "Boolean attribute",
			input:   `<input type="checkbox" checked>`,
			tagName: "input",
			attributes: map[string]string{
				"type":    "checkbox",
				"checked": "",
			},
		},
		{
			name:    "Single quoted attributes",
			input:   `<div class='container' id='main'>`,
			tagName: "div",
			attributes: map[string]string{
				"class": "container",
				"id":    "main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			token := lexer.NextToken()

			if token.Type != TokenOpenTag && token.Type != TokenSelfCloseTag {
				t.Fatalf("expected tag token, got %v", token.Type)
			}

			if token.Value != tt.tagName {
				t.Errorf("expected tag name %q, got %q", tt.tagName, token.Value)
			}

			if len(token.Attributes) != len(tt.attributes) {
				t.Errorf("expected %d attributes, got %d", len(tt.attributes), len(token.Attributes))
			}

			for key, expectedValue := range tt.attributes {
				if actualValue, exists := token.Attributes[key]; !exists {
					t.Errorf("missing attribute %q", key)
				} else if actualValue != expectedValue {
					t.Errorf("attribute %q: expected %q, got %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}

// TestLexerCommentParsing 测试注释解析
func TestLexerCommentParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple comment",
			input:    "<!-- Hello World -->",
			expected: "Hello World",
		},
		{
			name:     "Comment with special characters",
			input:    "<!-- <script>alert('test');</script> -->",
			expected: "<script>alert('test');</script>",
		},
		{
			name:     "Multiline comment",
			input:    "<!-- Line 1\nLine 2\nLine 3 -->",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "Empty comment",
			input:    "<!---->",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			token := lexer.NextToken()

			if token.Type != TokenComment {
				t.Fatalf("expected comment token, got %v", token.Type)
			}

			if token.Value != tt.expected {
				t.Errorf("expected comment %q, got %q", tt.expected, token.Value)
			}
		})
	}
}

// TestLexerTextParsing 测试文本解析
func TestLexerTextParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Text with whitespace",
			input:    "  Hello   World  ",
			expected: "Hello   World",
		},
		{
			name:     "Text with special characters",
			input:    "Hello & World © 2023",
			expected: "Hello & World © 2023",
		},
		{
			name:     "Unicode text",
			input:    "你好世界 🌍",
			expected: "你好世界 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			token := lexer.NextToken()

			if token.Type != TokenText {
				t.Fatalf("expected text token, got %v", token.Type)
			}

			if token.Value != tt.expected {
				t.Errorf("expected text %q, got %q", tt.expected, token.Value)
			}
		})
	}
}

// TestLexerErrorHandling 测试错误处理
func TestLexerErrorHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Unterminated tag",
			input: "<div",
		},
		{
			name:  "Invalid tag name",
			input: "<123invalid>",
		},
		{
			name:  "Unterminated attribute value",
			input: `<div class="unterminated`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)

			// 读取所有tokens，确保不会崩溃
			for {
				token := lexer.NextToken()
				if token.Type == TokenEOF || token.Type == TokenError {
					break
				}
			}

			// 如果到达这里，说明没有崩溃
			t.Logf("Successfully handled error case: %s", tt.name)
		})
	}
}

// TestLexerComplexDocument 测试复杂文档
func TestLexerComplexDocument(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<!-- This is a comment -->
<html lang="en">
	<head>
		<title>Test Document</title>
	</head>
	<body>
		<h1>Hello World</h1>
		<p>This is a <strong>test</strong> document.</p>
		<img src="image.jpg" alt="Test Image" />
	</body>
</html>`

	lexer := NewLexer(input)
	tokenCount := 0

	for {
		token := lexer.NextToken()
		tokenCount++

		// 验证每个token都有有效的位置信息
		if token.Position.Line < 1 || token.Position.Column < 1 {
			t.Errorf("invalid position for token %d: (%d,%d)", tokenCount, token.Position.Line, token.Position.Column)
		}

		if token.Type == TokenEOF {
			break
		}

		// 防止无限循环
		if tokenCount > 100 {
			t.Fatal("too many tokens, possible infinite loop")
		}
	}

	if tokenCount < 10 {
		t.Errorf("expected more tokens for complex document, got %d", tokenCount)
	}
}

// TestLexerSelfClosingTags 测试自封闭标签
func TestLexerSelfClosingTags(t *testing.T) {
	// 创建允许自封闭标签的配置
	config := DefaultConfig()
	config.AllowSelfCloseTags = true

	tests := []struct {
		name     string
		input    string
		expected TokenType
	}{
		{
			name:     "Self-closing img tag",
			input:    "<img />",
			expected: TokenSelfCloseTag,
		},
		{
			name:     "Self-closing br tag",
			input:    "<br/>",
			expected: TokenSelfCloseTag,
		},
		{
			name:     "Self-closing input tag",
			input:    "<input type='text' />",
			expected: TokenSelfCloseTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexerWithConfig(tt.input, config)
			token := lexer.NextToken()

			if token.Type != tt.expected {
				t.Errorf("expected token type %v, got %v", tt.expected, token.Type)
			}
		})
	}
}

// TestLexerConfigurationEffects 测试配置对词法分析的影响
func TestLexerConfigurationEffects(t *testing.T) {
	input := "<img />"

	// 测试不允许自封闭标签的配置
	configNoSelfClose := DefaultConfig()
	configNoSelfClose.AllowSelfCloseTags = false

	lexer := NewLexerWithConfig(input, configNoSelfClose)
	token := lexer.NextToken()

	// 应该返回错误token
	if token.Type != TokenError {
		t.Errorf("expected error token when self-closing tags disabled, got %v", token.Type)
	}

	// 测试允许自封闭标签的配置
	configAllowSelfClose := DefaultConfig()
	configAllowSelfClose.AllowSelfCloseTags = true

	lexer = NewLexerWithConfig(input, configAllowSelfClose)
	token = lexer.NextToken()

	// 应该返回自封闭标签token
	if token.Type != TokenSelfCloseTag {
		t.Errorf("expected self-close tag token when enabled, got %v", token.Type)
	}
}
