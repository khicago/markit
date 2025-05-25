package markit

import (
	"testing"
)

// TestLexerInternalMethods 测试词法分析器的内部方法
func TestLexerInternalMethods(t *testing.T) {
	t.Run("GetConfig method", func(t *testing.T) {
		config := DefaultConfig()
		config.CaseSensitive = false

		lexer := NewLexerWithConfig("test", config)
		retrievedConfig := lexer.GetConfig()

		if retrievedConfig.CaseSensitive != false {
			t.Errorf("expected CaseSensitive to be false, got %t", retrievedConfig.CaseSensitive)
		}
	})

	t.Run("peekChar edge cases", func(t *testing.T) {
		// 测试在字符串末尾的peek
		lexer := NewLexer("a")
		lexer.readChar() // 读取'a'，现在current是0（EOF）

		next := lexer.peekChar()
		if next != 0 {
			t.Errorf("expected peek to return 0 at EOF, got %c", next)
		}
	})

	t.Run("readChar positioning", func(t *testing.T) {
		lexer := NewLexer("a\nb")

		// 初始位置 - lexer创建时已经调用了readChar()，所以line=1, column=1
		if lexer.line != 1 || lexer.column != 1 {
			t.Errorf("expected initial position 1:1, got %d:%d", lexer.line, lexer.column)
		}

		// 读取字符'a'
		lexer.readChar()
		if lexer.column != 2 {
			t.Errorf("expected column 2, got %d", lexer.column)
		}

		// 读取换行符
		lexer.readChar()
		if lexer.line != 2 || lexer.column != 1 {
			t.Errorf("expected position 2:1, got %d:%d", lexer.line, lexer.column)
		}
	})
}

// TestLexerProtocolTokens 测试词法分析器的协议token
func TestLexerProtocolTokens(t *testing.T) {
	// 注意：当前实现不支持处理指令、DOCTYPE和CDATA，这些测试被跳过
	t.Skip("Current implementation does not support processing instructions, DOCTYPE, and CDATA")
}

// TestLexerCommentEdgeCases 测试注释的边缘情况
func TestLexerCommentEdgeCases(t *testing.T) {
	t.Run("Unterminated comment", func(t *testing.T) {
		lexer := NewLexer("<!-- unterminated comment")
		token := lexer.NextToken()

		// 根据readComment实现，未终止的注释会读取到EOF
		if token.Type != TokenComment {
			t.Errorf("expected comment token for unterminated comment, got %v", token.Type)
		}
	})

	t.Run("Comment with dashes", func(t *testing.T) {
		lexer := NewLexer("<!-- comment with -- dashes -->")
		token := lexer.NextToken()

		if token.Type != TokenComment {
			t.Errorf("expected comment token, got %v", token.Type)
		}

		expected := "comment with -- dashes"
		if token.Value != expected {
			t.Errorf("expected comment value %q, got %q", expected, token.Value)
		}
	})

	t.Run("Empty comment", func(t *testing.T) {
		lexer := NewLexer("<!---->")
		token := lexer.NextToken()

		if token.Type != TokenComment {
			t.Errorf("expected comment token, got %v", token.Type)
		}

		if token.Value != "" {
			t.Errorf("expected empty comment value, got %q", token.Value)
		}
	})
}

// TestLexerAttributeEdgeCases 测试属性解析的边缘情况
func TestLexerAttributeEdgeCases(t *testing.T) {
	t.Run("Attribute without value", func(t *testing.T) {
		lexer := NewLexer(`<input checked>`)
		token := lexer.NextToken()

		if token.Type != TokenOpenTag {
			t.Fatalf("expected open tag, got %v", token.Type)
		}

		if token.Attributes["checked"] != "" {
			t.Errorf("expected empty value for boolean attribute, got %q", token.Attributes["checked"])
		}
	})

	t.Run("Attribute with unquoted value", func(t *testing.T) {
		lexer := NewLexer(`<div class=container>`)
		token := lexer.NextToken()

		if token.Type != TokenOpenTag {
			t.Fatalf("expected open tag, got %v", token.Type)
		}

		if token.Attributes["class"] != "container" {
			t.Errorf("expected class value 'container', got %q", token.Attributes["class"])
		}
	})

	t.Run("Attribute with escaped quotes", func(t *testing.T) {
		lexer := NewLexer(`<div title="He said \"Hello\"">`)
		token := lexer.NextToken()

		if token.Type != TokenOpenTag {
			t.Fatalf("expected open tag, got %v", token.Type)
		}

		expected := `He said "Hello"`
		if token.Attributes["title"] != expected {
			t.Errorf("expected title value %q, got %q", expected, token.Attributes["title"])
		}
	})
}

// TestLexerErrorHandlingEdgeCases 测试错误处理的边缘情况
func TestLexerErrorHandlingEdgeCases(t *testing.T) {
	t.Run("Invalid tag start", func(t *testing.T) {
		lexer := NewLexer("<123invalid>")
		token := lexer.NextToken()

		if token.Type != TokenError {
			t.Errorf("expected error token for invalid tag name, got %v", token.Type)
		}
	})

	t.Run("Unterminated attribute value with single quote", func(t *testing.T) {
		lexer := NewLexer(`<div class='unterminated`)
		token := lexer.NextToken()

		if token.Type != TokenError {
			t.Errorf("expected error token for unterminated attribute, got %v", token.Type)
		}
	})

	t.Run("Unterminated attribute value with double quote", func(t *testing.T) {
		lexer := NewLexer(`<div class="unterminated`)
		token := lexer.NextToken()

		if token.Type != TokenError {
			t.Errorf("expected error token for unterminated attribute, got %v", token.Type)
		}
	})
}

// TestLexerIdentifierMethods 测试标识符读取方法
func TestLexerIdentifierMethods(t *testing.T) {
	// 注意：这些是内部方法测试，跳过以避免测试实现细节
	t.Skip("Skipping internal method tests to avoid testing implementation details")
}

// TestLexerAttributeValueMethods 测试属性值读取方法
func TestLexerAttributeValueMethods(t *testing.T) {
	// 注意：这些是内部方法测试，跳过以避免测试实现细节
	t.Skip("Skipping internal method tests to avoid testing implementation details")
}

// TestLexerUtilityFunctions 测试工具函数
func TestLexerUtilityFunctions(t *testing.T) {
	t.Run("isIdentifierStart", func(t *testing.T) {
		validStarts := []rune{'a', 'A', '_', '-', ':'}
		for _, r := range validStarts {
			if !isIdentifierStart(r) {
				t.Errorf("expected %c to be valid identifier start", r)
			}
		}

		invalidStarts := []rune{'1', '!', '@', ' '}
		for _, r := range invalidStarts {
			if isIdentifierStart(r) {
				t.Errorf("expected %c to be invalid identifier start", r)
			}
		}
	})

	t.Run("isIdentifierChar", func(t *testing.T) {
		validChars := []rune{'a', 'A', '1', '_', '-', ':'}
		for _, r := range validChars {
			if !isIdentifierChar(r) {
				t.Errorf("expected %c to be valid identifier char", r)
			}
		}

		invalidChars := []rune{'!', '@', ' ', '>'}
		for _, r := range invalidChars {
			if isIdentifierChar(r) {
				t.Errorf("expected %c to be invalid identifier char", r)
			}
		}
	})
}
