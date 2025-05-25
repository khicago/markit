package markit

import (
	"testing"
)

// TestTokenTypeStringComplete 测试所有TokenType的String方法
func TestTokenTypeStringComplete(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{TokenError, "ERROR"},
		{TokenEOF, "EOF"},
		{TokenText, "TEXT"},
		{TokenOpenTag, "OPEN_TAG"},
		{TokenCloseTag, "CLOSE_TAG"},
		{TokenSelfCloseTag, "SELF_CLOSE_TAG"},
		{TokenAttribute, "ATTRIBUTE"},
		{TokenComment, "COMMENT"},
		{TokenProcessingInstruction, "PROCESSING_INSTRUCTION"},
		{TokenDoctype, "DOCTYPE"},
		{TokenCDATA, "CDATA"},
		{TokenEntity, "ENTITY"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.tokenType.String()
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

// TestTokenTypeStringUnknown 测试未知TokenType的String方法
func TestTokenTypeStringUnknown(t *testing.T) {
	unknownType := TokenType(999)
	result := unknownType.String()
	expected := "UNKNOWN(999)"

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// TestTokenTypeStringEdgeCases 测试TokenType.String的边界情况
func TestTokenTypeStringEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		tokenType TokenType
		expected  string
	}{
		{"Negative value", TokenType(-1), "UNKNOWN(-1)"},
		{"Zero value", TokenType(0), "ERROR"}, // 0对应TokenError
		{"Large value", TokenType(1000), "UNKNOWN(1000)"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.tokenType.String()
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}

// TestTokenStringRepresentation 测试Token结构体的字符串表示
func TestTokenStringRepresentation(t *testing.T) {
	token := Token{
		Type:     TokenText,
		Value:    "Hello World",
		Position: Position{Line: 1, Column: 5},
	}

	// 测试token的各个字段
	if token.Type.String() != "TEXT" {
		t.Errorf("expected token type string 'TEXT', got %q", token.Type.String())
	}

	if token.Value != "Hello World" {
		t.Errorf("expected token value 'Hello World', got %q", token.Value)
	}

	if token.Position.Line != 1 {
		t.Errorf("expected line 1, got %d", token.Position.Line)
	}

	if token.Position.Column != 5 {
		t.Errorf("expected column 5, got %d", token.Position.Column)
	}
}

// TestTokenCreationCoverage 测试不同类型token的创建
func TestTokenCreationCoverage(t *testing.T) {
	tests := []struct {
		name      string
		tokenType TokenType
		value     string
		line      int
		column    int
	}{
		{"Error token", TokenError, "invalid", 6, 1},
		{"EOF token", TokenEOF, "", 10, 1},
		{"Text token", TokenText, "content", 2, 5},
		{"Open tag token", TokenOpenTag, "div", 1, 1},
		{"Close tag token", TokenCloseTag, "div", 1, 10},
		{"Self-close tag token", TokenSelfCloseTag, "br", 2, 1},
		{"Attribute token", TokenAttribute, "class=\"test\"", 1, 5},
		{"Comment token", TokenComment, " comment ", 3, 1},
		{"Processing instruction token", TokenProcessingInstruction, "xml version=\"1.0\"", 1, 1},
		{"Doctype token", TokenDoctype, "html", 1, 1},
		{"CDATA token", TokenCDATA, "script content", 4, 1},
		{"Entity token", TokenEntity, "&amp;", 5, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token := Token{
				Type:     test.tokenType,
				Value:    test.value,
				Position: Position{Line: test.line, Column: test.column},
			}

			if token.Type != test.tokenType {
				t.Errorf("expected token type %v, got %v", test.tokenType, token.Type)
			}

			if token.Value != test.value {
				t.Errorf("expected token value %q, got %q", test.value, token.Value)
			}

			if token.Position.Line != test.line {
				t.Errorf("expected line %d, got %d", test.line, token.Position.Line)
			}

			if token.Position.Column != test.column {
				t.Errorf("expected column %d, got %d", test.column, token.Position.Column)
			}

			// 测试token类型的字符串表示
			typeString := token.Type.String()
			if typeString == "" {
				t.Error("token type string should not be empty")
			}
		})
	}
}

// TestPositionStruct 测试Position结构体
func TestPositionStruct(t *testing.T) {
	pos := Position{Line: 42, Column: 13}

	if pos.Line != 42 {
		t.Errorf("expected line 42, got %d", pos.Line)
	}

	if pos.Column != 13 {
		t.Errorf("expected column 13, got %d", pos.Column)
	}
}

// TestTokenEquality 测试token的相等性比较
func TestTokenEquality(t *testing.T) {
	token1 := Token{
		Type:     TokenText,
		Value:    "test",
		Position: Position{Line: 1, Column: 1},
	}

	token2 := Token{
		Type:     TokenText,
		Value:    "test",
		Position: Position{Line: 1, Column: 1},
	}

	token3 := Token{
		Type:     TokenText,
		Value:    "different",
		Position: Position{Line: 1, Column: 1},
	}

	// 测试相同的token
	if token1.Type != token2.Type {
		t.Error("tokens with same type should be equal")
	}

	if token1.Value != token2.Value {
		t.Error("tokens with same value should be equal")
	}

	// 测试不同的token
	if token1.Value == token3.Value {
		t.Error("tokens with different values should not be equal")
	}
}

// TestTokenTypeConstants 测试所有TokenType常量都有定义
func TestTokenTypeConstants(t *testing.T) {
	// 确保所有的token类型都有合理的字符串表示
	tokenTypes := []TokenType{
		TokenError,
		TokenEOF,
		TokenText,
		TokenOpenTag,
		TokenCloseTag,
		TokenSelfCloseTag,
		TokenAttribute,
		TokenComment,
		TokenProcessingInstruction,
		TokenDoctype,
		TokenCDATA,
		TokenEntity,
	}

	for _, tokenType := range tokenTypes {
		str := tokenType.String()
		if str == "" {
			t.Errorf("token type %d should have a non-empty string representation", int(tokenType))
		}

		// 确保不是UNKNOWN格式
		if len(str) > 7 && str[:7] == "UNKNOWN" {
			t.Errorf("token type %d should not have UNKNOWN string representation: %s", int(tokenType), str)
		}
	}
}
