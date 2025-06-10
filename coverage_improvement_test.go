package markit

import (
	"testing"
)

// TestDefaultHTMLConfig 测试 DefaultHTMLConfig 函数
func TestDefaultHTMLConfig(t *testing.T) {
	config := DefaultHTMLConfig()

	if config == nil {
		t.Fatal("DefaultHTMLConfig should not return nil")
	}

	// 验证HTML配置的基本属性
	if config.CaseSensitive {
		t.Error("HTML config should be case insensitive")
	}

	if !config.TrimWhitespace {
		t.Error("HTML config should trim whitespace")
	}

	// 验证void elements是否正确设置
	if !config.IsVoidElement("br") {
		t.Error("br should be a void element in HTML config")
	}

	if !config.IsVoidElement("img") {
		t.Error("img should be a void element in HTML config")
	}
}

// TestReadTextRecursiveCall 测试 readText 函数的递归调用行为
func TestReadTextRecursiveCall(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "whitespace only text followed by tag",
			input: "   <tag>",
			expected: []Token{
				{Type: TokenOpenTag, Value: "tag"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "multiple whitespace segments",
			input: "  \n\t  <tag>content</tag>",
			expected: []Token{
				{Type: TokenOpenTag, Value: "tag"},
				{Type: TokenText, Value: "content"},
				{Type: TokenCloseTag, Value: "tag"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "whitespace between tags",
			input: "<start>  \n  </start>",
			expected: []Token{
				{Type: TokenOpenTag, Value: "start"},
				{Type: TokenCloseTag, Value: "start"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "consecutive empty text segments",
			input: "   \n   \t   <tag/>",
			expected: []Token{
				{Type: TokenSelfCloseTag, Value: "tag"},
				{Type: TokenEOF, Value: ""},
			},
		},
		{
			name:  "empty text at end of input",
			input: "<tag>content</tag>   \n   ",
			expected: []Token{
				{Type: TokenOpenTag, Value: "tag"},
				{Type: TokenText, Value: "content"},
				{Type: TokenCloseTag, Value: "tag"},
				{Type: TokenEOF, Value: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.TrimWhitespace = true
			lexer := NewLexerWithConfig(tt.input, config)

			var tokens []Token
			for {
				token := lexer.NextToken()
				tokens = append(tokens, token)
				if token.Type == TokenEOF {
					break
				}
			}

			if len(tokens) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokens))
				return
			}

			for i, expected := range tt.expected {
				if tokens[i].Type != expected.Type {
					t.Errorf("Token %d: expected type %v, got %v", i, expected.Type, tokens[i].Type)
				}
				if expected.Value != "" && tokens[i].Value != expected.Value {
					t.Errorf("Token %d: expected value %q, got %q", i, expected.Value, tokens[i].Value)
				}
			}
		})
	}
}

// TestVoidElementEdgeCases 测试 AddVoidElement 和 RemoveVoidElement 的边界情况
func TestVoidElementEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "AddVoidElement with nil VoidElements",
			testFunc: func(t *testing.T) {
				config := &ParserConfig{
					VoidElements: nil,
				}
				config.AddVoidElement("test")
				if !config.IsVoidElement("test") {
					t.Error("Should be able to add void element when VoidElements is nil")
				}
			},
		},
		{
			name: "AddVoidElement case sensitivity",
			testFunc: func(t *testing.T) {
				config := &ParserConfig{
					CaseSensitive: true,
					VoidElements:  make(map[string]bool),
				}
				config.AddVoidElement("Test")
				if !config.IsVoidElement("Test") {
					t.Error("Should find exact case match")
				}
				if config.IsVoidElement("test") {
					t.Error("Should not find different case when case sensitive")
				}
			},
		},
		{
			name: "AddVoidElement case insensitive",
			testFunc: func(t *testing.T) {
				config := &ParserConfig{
					CaseSensitive: false,
					VoidElements:  make(map[string]bool),
				}
				config.AddVoidElement("Test")
				if !config.IsVoidElement("test") {
					t.Error("Should find case insensitive match")
				}
			},
		},
		{
			name: "RemoveVoidElement with nil VoidElements",
			testFunc: func(t *testing.T) {
				config := &ParserConfig{
					VoidElements: nil,
				}
				// Should not panic
				config.RemoveVoidElement("test")
			},
		},
		{
			name: "RemoveVoidElement existing element",
			testFunc: func(t *testing.T) {
				config := &ParserConfig{
					VoidElements: map[string]bool{"test": true},
				}
				config.RemoveVoidElement("test")
				if config.IsVoidElement("test") {
					t.Error("Element should be removed")
				}
			},
		},
		{
			name: "RemoveVoidElement nonexistent element",
			testFunc: func(t *testing.T) {
				config := &ParserConfig{
					VoidElements: map[string]bool{"other": true},
				}
				// Should not panic
				config.RemoveVoidElement("test")
				if !config.IsVoidElement("other") {
					t.Error("Other elements should remain")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestEmptyTextHandlingEdgeCases 测试空文本处理的边界情况
func TestEmptyTextHandlingEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "consecutive whitespace segments",
			input:    "   \n   \t   <tag>content</tag>   \n   ",
			expected: []TokenType{TokenOpenTag, TokenText, TokenCloseTag, TokenEOF},
		},
		{
			name:     "whitespace at end of input",
			input:    "<tag>content</tag>   \n   \t   ",
			expected: []TokenType{TokenOpenTag, TokenText, TokenCloseTag, TokenEOF},
		},
		{
			name:     "only whitespace input",
			input:    "   \n   \t   ",
			expected: []TokenType{TokenEOF},
		},
		{
			name:     "whitespace between nested tags",
			input:    "<outer>   \n   <inner/>   \n   </outer>",
			expected: []TokenType{TokenOpenTag, TokenSelfCloseTag, TokenCloseTag, TokenEOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.TrimWhitespace = true
			lexer := NewLexerWithConfig(tt.input, config)

			var tokenTypes []TokenType
			for {
				token := lexer.NextToken()
				tokenTypes = append(tokenTypes, token.Type)
				if token.Type == TokenEOF {
					break
				}
			}

			if len(tokenTypes) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokenTypes))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", tokenTypes)
				return
			}

			for i, expected := range tt.expected {
				if tokenTypes[i] != expected {
					t.Errorf("Token %d: expected type %v, got %v", i, expected, tokenTypes[i])
				}
			}
		})
	}
}

// TestReadTextComplexRecursion 测试 readText 函数的复杂递归场景
func TestReadTextComplexRecursion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "deeply nested empty text",
			input:    "   \n   \t   <a>   \n   <b>   \n   </b>   \n   </a>   \n   ",
			expected: []TokenType{TokenOpenTag, TokenOpenTag, TokenCloseTag, TokenCloseTag, TokenEOF},
		},
		{
			name:     "alternating content and whitespace",
			input:    "content1   \n   <tag>   \n   content2   \n   </tag>   \n   content3",
			expected: []TokenType{TokenText, TokenOpenTag, TokenText, TokenCloseTag, TokenText, TokenEOF},
		},
		{
			name:     "unicode whitespace",
			input:    "\u00A0\u2000\u2001<tag>\u00A0\u2000content\u2001</tag>\u00A0\u2000\u2001",
			expected: []TokenType{TokenOpenTag, TokenText, TokenCloseTag, TokenEOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.TrimWhitespace = true
			lexer := NewLexerWithConfig(tt.input, config)

			var tokenTypes []TokenType
			for {
				token := lexer.NextToken()
				tokenTypes = append(tokenTypes, token.Type)
				if token.Type == TokenEOF {
					break
				}
			}

			if len(tokenTypes) != len(tt.expected) {
				t.Errorf("Expected %d tokens, got %d", len(tt.expected), len(tokenTypes))
				t.Errorf("Expected: %v", tt.expected)
				t.Errorf("Got: %v", tokenTypes)
				return
			}

			for i, expected := range tt.expected {
				if tokenTypes[i] != expected {
					t.Errorf("Token %d: expected type %v, got %v", i, expected, tokenTypes[i])
				}
			}
		})
	}
}

// TestReadTextDeepRecursion 测试深度递归的边界情况
func TestReadTextDeepRecursion(t *testing.T) {
	// 创建一个包含大量连续空白的输入，测试递归调用的深度
	input := ""
	for i := 0; i < 100; i++ {
		input += "   \n   \t   "
	}
	input += "<tag>content</tag>"

	config := DefaultConfig()
	config.TrimWhitespace = true
	lexer := NewLexerWithConfig(input, config)

	// 应该跳过所有空白，直接得到第一个有效token
	token := lexer.NextToken()
	if token.Type != TokenOpenTag {
		t.Errorf("Expected TokenOpenTag, got %v", token.Type)
	}
	if token.Value != "tag" {
		t.Errorf("Expected 'tag', got %q", token.Value)
	}
}

// TestReadTextRecursionWithComments 测试包含注释的递归情况
func TestReadTextRecursionWithComments(t *testing.T) {
	input := "   \n   <!-- comment -->   \n   <tag>content</tag>"

	config := DefaultConfig()
	config.TrimWhitespace = true
	lexer := NewLexerWithConfig(input, config)

	var tokens []Token
	for {
		token := lexer.NextToken()
		tokens = append(tokens, token)
		if token.Type == TokenEOF {
			break
		}
	}

	expectedTypes := []TokenType{TokenComment, TokenOpenTag, TokenText, TokenCloseTag, TokenEOF}
	if len(tokens) != len(expectedTypes) {
		t.Errorf("Expected %d tokens, got %d", len(expectedTypes), len(tokens))
		return
	}

	for i, expected := range expectedTypes {
		if tokens[i].Type != expected {
			t.Errorf("Token %d: expected type %v, got %v", i, expected, tokens[i].Type)
		}
	}
}

func TestReadTextWithNilConfig(t *testing.T) {
	// 测试当 config 为 nil 时的 readText 行为
	// 注意：我们不能完全设置 config 为 nil，因为 NextToken 需要 CoreMatcher
	// 所以我们创建一个配置，但不启用 TrimWhitespace
	config := DefaultConfig()
	config.TrimWhitespace = false

	lexer := NewLexerWithConfig("some text content", config)

	token := lexer.NextToken()

	if token.Type != TokenText {
		t.Errorf("Expected TokenText, got %v", token.Type)
	}

	if token.Value != "some text content" {
		t.Errorf("Expected 'some text content', got '%s'", token.Value)
	}
}

func TestReadTextConfigNilWithWhitespace(t *testing.T) {
	// 测试当 TrimWhitespace 为 false 时，空白字符不会被修剪
	config := DefaultConfig()
	config.TrimWhitespace = false

	lexer := NewLexerWithConfig("  whitespace text  ", config)

	token := lexer.NextToken()

	if token.Type != TokenText {
		t.Errorf("Expected TokenText, got %v", token.Type)
	}

	// 当 TrimWhitespace 为 false 时，空白字符应该保留
	if token.Value != "  whitespace text  " {
		t.Errorf("Expected '  whitespace text  ', got '%s'", token.Value)
	}
}

func TestReadTextTrimWhitespaceDisabled(t *testing.T) {
	// 测试 TrimWhitespace 为 false 时的各种情况
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "leading whitespace",
			input:    "   text",
			expected: "   text",
		},
		{
			name:     "trailing whitespace",
			input:    "text   ",
			expected: "text   ",
		},
		{
			name:     "both leading and trailing whitespace",
			input:    "   text   ",
			expected: "   text   ",
		},
		{
			name:     "only whitespace",
			input:    "   ",
			expected: "   ",
		},
		{
			name:     "tabs and spaces",
			input:    "\t  text  \t",
			expected: "\t  text  \t",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := DefaultConfig()
			config.TrimWhitespace = false

			lexer := NewLexerWithConfig(tc.input, config)
			token := lexer.NextToken()

			if token.Type != TokenText {
				t.Errorf("Expected TokenText, got %v", token.Type)
			}

			if token.Value != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, token.Value)
			}
		})
	}
}
