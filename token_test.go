package markit

import (
	"strings"
	"testing"
)

// TestTokenTypeString 测试TokenType的String方法
func TestTokenTypeString(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{TokenOpenTag, "OPEN_TAG"},
		{TokenCloseTag, "CLOSE_TAG"},
		{TokenSelfCloseTag, "SELF_CLOSE_TAG"},
		{TokenText, "TEXT"},
		{TokenComment, "COMMENT"},
		{TokenAttribute, "ATTRIBUTE"},
		{TokenError, "ERROR"},
		{TokenEOF, "EOF"},
		{TokenProcessingInstruction, "PROCESSING_INSTRUCTION"},
		{TokenDoctype, "DOCTYPE"},
		{TokenCDATA, "CDATA"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.tokenType.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}

	// 测试未知的TokenType
	t.Run("unknown token type", func(t *testing.T) {
		unknownType := TokenType(999)
		result := unknownType.String()
		if !strings.Contains(result, "UNKNOWN") {
			t.Errorf("Unknown token type should contain 'UNKNOWN', got '%s'", result)
		}
	})
}

// TestTokenString 测试Token的String方法
func TestTokenString(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected string
	}{
		{
			name: "error token",
			token: Token{
				Type:  TokenError,
				Value: "unexpected character",
			},
			expected: "ERROR(unexpected character)",
		},
		{
			name: "EOF token",
			token: Token{
				Type: TokenEOF,
			},
			expected: "EOF",
		},
		{
			name: "text token",
			token: Token{
				Type:  TokenText,
				Value: "hello world",
			},
			expected: "TEXT(\"hello world\")",
		},
		{
			name: "open tag token",
			token: Token{
				Type:  TokenOpenTag,
				Value: "div",
			},
			expected: "OPEN_TAG(div)",
		},
		{
			name: "close tag token",
			token: Token{
				Type:  TokenCloseTag,
				Value: "div",
			},
			expected: "CLOSE_TAG(div)",
		},
		{
			name: "self-close tag token",
			token: Token{
				Type:  TokenSelfCloseTag,
				Value: "img",
			},
			expected: "SELF_CLOSE_TAG(img)",
		},
		{
			name: "attribute token",
			token: Token{
				Type:  TokenAttribute,
				Value: "class=container",
			},
			expected: "ATTR(class=container)",
		},
		{
			name: "comment token",
			token: Token{
				Type:  TokenComment,
				Value: " this is a comment ",
			},
			expected: "COMMENT( this is a comment )",
		},
		{
			name: "unknown token type",
			token: Token{
				Type:  TokenType(999),
				Value: "unknown",
			},
			expected: "UNKNOWN(999)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestPositionString 测试Position的String方法
func TestPositionString(t *testing.T) {
	tests := []struct {
		name     string
		position Position
		expected string
	}{
		{
			name: "basic position",
			position: Position{
				Line:   1,
				Column: 1,
				Offset: 0,
			},
			expected: "1:1",
		},
		{
			name: "middle position",
			position: Position{
				Line:   10,
				Column: 25,
				Offset: 150,
			},
			expected: "10:25",
		},
		{
			name: "large position",
			position: Position{
				Line:   1000,
				Column: 500,
				Offset: 50000,
			},
			expected: "1000:500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.position.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestTokenCreation 测试Token的创建和属性
func TestTokenCreation(t *testing.T) {
	t.Run("token with attributes", func(t *testing.T) {
		token := Token{
			Type:  TokenOpenTag,
			Value: "div",
			Attributes: map[string]string{
				"class": "container",
				"id":    "main",
			},
			Position: Position{Line: 1, Column: 1, Offset: 0},
		}

		if token.Type != TokenOpenTag {
			t.Errorf("Expected TokenOpenTag, got %v", token.Type)
		}

		if token.Value != "div" {
			t.Errorf("Expected 'div', got '%s'", token.Value)
		}

		if len(token.Attributes) != 2 {
			t.Errorf("Expected 2 attributes, got %d", len(token.Attributes))
		}

		if token.Attributes["class"] != "container" {
			t.Errorf("Expected class='container', got '%s'", token.Attributes["class"])
		}

		if token.Position.Line != 1 {
			t.Errorf("Expected line 1, got %d", token.Position.Line)
		}
	})

	t.Run("empty token", func(t *testing.T) {
		token := Token{}

		if token.Type != TokenType(0) {
			t.Errorf("Expected zero value TokenType, got %v", token.Type)
		}

		if token.Value != "" {
			t.Errorf("Expected empty value, got '%s'", token.Value)
		}

		if token.Attributes != nil {
			t.Errorf("Expected nil attributes, got %v", token.Attributes)
		}
	})
}
