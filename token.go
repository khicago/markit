package markit

import "fmt"

// TokenType 表示标记类型
type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenText
	TokenOpenTag
	TokenCloseTag
	TokenSelfCloseTag
	TokenAttribute
	TokenComment
	// 新增的协议类型
	TokenProcessingInstruction
	TokenDoctype
	TokenCDATA
	TokenEntity
)

// String 返回 TokenType 的字符串表示
func (t TokenType) String() string {
	switch t {
	case TokenError:
		return "ERROR"
	case TokenEOF:
		return "EOF"
	case TokenText:
		return "TEXT"
	case TokenOpenTag:
		return "OPEN_TAG"
	case TokenCloseTag:
		return "CLOSE_TAG"
	case TokenSelfCloseTag:
		return "SELF_CLOSE_TAG"
	case TokenAttribute:
		return "ATTRIBUTE"
	case TokenComment:
		return "COMMENT"
	case TokenProcessingInstruction:
		return "PROCESSING_INSTRUCTION"
	case TokenDoctype:
		return "DOCTYPE"
	case TokenCDATA:
		return "CDATA"
	case TokenEntity:
		return "ENTITY"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(t))
	}
}

// Token 表示一个词法标记
type Token struct {
	Type       TokenType
	Value      string
	Attributes map[string]string
	Position   Position
}

// Position 表示源码中的位置信息
type Position struct {
	Line   int
	Column int
	Offset int
}

// String 返回 Token 的字符串表示
func (t Token) String() string {
	switch t.Type {
	case TokenError:
		return fmt.Sprintf("ERROR(%s)", t.Value)
	case TokenEOF:
		return "EOF"
	case TokenText:
		return fmt.Sprintf("TEXT(%q)", t.Value)
	case TokenOpenTag:
		return fmt.Sprintf("OPEN_TAG(%s)", t.Value)
	case TokenCloseTag:
		return fmt.Sprintf("CLOSE_TAG(%s)", t.Value)
	case TokenSelfCloseTag:
		return fmt.Sprintf("SELF_CLOSE_TAG(%s)", t.Value)
	case TokenAttribute:
		return fmt.Sprintf("ATTR(%s)", t.Value)
	case TokenComment:
		return fmt.Sprintf("COMMENT(%s)", t.Value)
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(t.Type))
	}
}

// String 返回 Position 的字符串表示
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}
