package markit

import (
	"testing"
)

func TestDebugLexer(t *testing.T) {
	input := "<root>Hello World</root>"
	lexer := NewLexer(input)

	t.Logf("Input: %s", input)

	// 获取所有 tokens
	for i := 0; i < 10; i++ {
		token := lexer.NextToken()
		t.Logf("Token %d: Type=%s, Value=%s, Position=%+v", i, token.Type, token.Value, token.Position)
		if token.Type == TokenEOF {
			break
		}
	}
}
