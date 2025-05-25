package markit

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer 词法分析器
type Lexer struct {
	input    string
	position int
	line     int
	column   int
	current  rune
	config   *ParserConfig
}

// NewLexer 创建新的词法分析器（使用默认配置）
func NewLexer(input string) *Lexer {
	return NewLexerWithConfig(input, DefaultConfig())
}

// NewLexerWithConfig 创建带配置的词法分析器
func NewLexerWithConfig(input string, config *ParserConfig) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
		config: config,
	}
	l.readChar()
	return l
}

// SetConfig 设置词法分析器配置
func (l *Lexer) SetConfig(config *ParserConfig) {
	l.config = config
}

// GetConfig 获取词法分析器配置
func (l *Lexer) GetConfig() *ParserConfig {
	return l.config
}

// NextToken 获取下一个 token
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	pos := Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.position,
	}

	if l.position >= len(l.input) {
		return Token{Type: TokenEOF, Value: "", Position: pos}
	}

	// 计算当前字符的位置（因为 readChar 已经移动了位置）
	currentPos := l.position
	if l.current != 0 {
		// 回退到当前字符的位置
		_, size := utf8.DecodeRuneInString(l.input[l.position-1:])
		currentPos = l.position - size
	}

	// 使用核心协议匹配器检查是否是标签开始
	if protocol := l.config.CoreMatcher.MatchProtocol(l.input, currentPos); protocol != nil {
		return l.readProtocolToken(protocol)
	}

	// 读取文本内容
	token := l.readText(pos)
	return token
}

// readChar 读取下一个字符
func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.current = 0 // EOF
	} else {
		if l.current == '\n' {
			l.line++
			l.column = 0
		}
		// 正确解码UTF-8字符
		r, size := utf8.DecodeRuneInString(l.input[l.position:])
		l.current = r
		l.position += size
		l.column++
	}
}

// peekChar 查看下一个字符但不移动位置
func (l *Lexer) peekChar() rune {
	if l.position >= len(l.input) {
		return 0
	}
	// 正确解码UTF-8字符
	r, _ := utf8.DecodeRuneInString(l.input[l.position:])
	return r
}

// skipWhitespace 跳过空白字符
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.current) {
		l.readChar()
	}
}

// readText 读取文本内容
func (l *Lexer) readText(pos Position) Token {
	var text strings.Builder

	for l.current != '<' && l.current != 0 {
		text.WriteRune(l.current)
		l.readChar()
	}

	content := strings.TrimSpace(text.String())

	return Token{
		Type:     TokenText,
		Value:    content,
		Position: pos,
	}
}

// readIdentifier 读取标识符（标签名或属性名）
func (l *Lexer) readIdentifier() string {
	var identifier strings.Builder

	// 第一个字符必须是字母、下划线或连字符
	if !isIdentifierStart(l.current) {
		return ""
	}

	for isIdentifierChar(l.current) {
		identifier.WriteRune(l.current)
		l.readChar()
	}

	return identifier.String()
}

// readAttribute 读取属性
func (l *Lexer) readAttribute() (string, string, error) {
	// 读取属性名
	name := l.readIdentifier()
	if name == "" {
		return "", "", fmt.Errorf("invalid attribute name")
	}

	l.skipWhitespace()

	// 检查是否有等号
	if l.current != '=' {
		// 布尔属性，没有值
		return name, "", nil
	}

	l.readChar() // 跳过 '='
	l.skipWhitespace()

	// 读取属性值
	value, err := l.readAttributeValue()
	if err != nil {
		return "", "", err
	}

	return name, value, nil
}

// readAttributeValue 读取属性值
func (l *Lexer) readAttributeValue() (string, error) {
	if l.current == '"' || l.current == '\'' {
		// 带引号的值
		quote := l.current
		l.readChar() // 跳过开始引号

		var value strings.Builder
		for l.current != quote && l.current != 0 {
			if l.current == '\\' {
				l.readChar()
				if l.current != 0 {
					value.WriteRune(l.current)
					l.readChar()
				}
			} else {
				value.WriteRune(l.current)
				l.readChar()
			}
		}

		if l.current != quote {
			return "", fmt.Errorf("unterminated quoted string")
		}
		l.readChar() // 跳过结束引号

		return value.String(), nil
	} else {
		// 不带引号的值
		var value strings.Builder
		for !unicode.IsSpace(l.current) && l.current != '>' && l.current != '/' && l.current != 0 {
			value.WriteRune(l.current)
			l.readChar()
		}
		return value.String(), nil
	}
}

// isIdentifierStart 检查字符是否可以作为标识符的开始
func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || r == '-' || r == ':'
}

// isIdentifierChar 检查字符是否可以作为标识符的一部分
func isIdentifierChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == ':'
}

// readComment 读取 XML 注释 <!-- ... -->
func (l *Lexer) readComment(pos Position) Token {
	// 跳过 "<!--" 序列（已经被协议匹配器识别）
	for i := 0; i < 4; i++ { // "<!--" 长度为4
		l.readChar()
	}

	var comment strings.Builder

	// 读取注释内容直到找到 -->
	for l.current != 0 {
		if l.current == '-' && l.peekChar() == '-' {
			// 检查是否是注释结束
			l.readChar() // 跳过第一个 '-'
			if l.current == '-' && l.peekChar() == '>' {
				l.readChar() // 跳过第二个 '-'
				l.readChar() // 跳过 '>'
				break
			} else {
				// 不是注释结束，将 '-' 添加到内容中
				comment.WriteRune('-')
			}
		} else {
			comment.WriteRune(l.current)
			l.readChar()
		}
	}

	return Token{
		Type:     TokenComment,
		Value:    strings.TrimSpace(comment.String()), // 去除前后空格
		Position: pos,
	}
}

// readProtocolToken 读取协议token
func (l *Lexer) readProtocolToken(protocol *CoreProtocol) Token {
	pos := Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.position,
	}

	if protocol.Name == "markit-standard-tag" {
		return l.readTag(pos)
	} else if protocol.Name == "markit-comment" {
		return l.readComment(pos)
	}

	// 对于其他协议，使用原来的逻辑
	start := l.position
	if l.current != 0 {
		_, size := utf8.DecodeRuneInString(l.input[l.position-1:])
		start = l.position - size
	}

	// 跳过开始序列
	for i := 0; i < len(protocol.OpenSeq); i++ {
		l.readChar()
	}

	// 查找结束序列
	closeSeq := protocol.CloseSeq
	for l.position < len(l.input) {
		if strings.HasPrefix(l.input[l.position:], closeSeq) {
			content := l.input[start : l.position+len(closeSeq)]
			// 跳过结束序列
			for i := 0; i < len(closeSeq); i++ {
				l.readChar()
			}
			return Token{Type: protocol.TokenType, Value: content, Position: pos}
		}
		l.readChar()
	}

	// 如果没有找到结束序列，返回到文件末尾
	content := l.input[start:]
	return Token{Type: protocol.TokenType, Value: content, Position: pos}
}

// readTag 读取标签
func (l *Lexer) readTag(pos Position) Token {
	l.readChar() // 跳过 '<'

	// 检查是否是结束标签
	isCloseTag := false
	if l.current == '/' {
		isCloseTag = true
		l.readChar() // 跳过 '/'
	}

	// 读取标签名
	tagName := l.readIdentifier()
	if tagName == "" {
		return Token{Type: TokenError, Value: "invalid tag name", Position: pos}
	}

	// 跳过空白
	l.skipWhitespace()

	// 读取属性
	attributes := make(map[string]string)
	if !isCloseTag {
		for l.current != '>' && l.current != '/' && l.current != 0 {
			name, value, err := l.readAttribute()
			if err != nil {
				return Token{Type: TokenError, Value: err.Error(), Position: pos}
			}
			attributes[name] = value
			l.skipWhitespace()
		}
	}

	// 检查自封闭标签
	isSelfClose := false
	if l.current == '/' {
		// 检查配置是否允许自封闭标签
		if l.config != nil && l.config.AllowSelfCloseTags {
			isSelfClose = true
			l.readChar() // 跳过 '/'
		} else {
			// 如果不允许自封闭标签，将 '/' 视为普通字符
			// 这里可以选择报错或者继续处理
			return Token{Type: TokenError, Value: "self-closing tags not allowed", Position: pos}
		}
	}

	// 跳过 '>'
	if l.current != '>' {
		return Token{Type: TokenError, Value: "expected '>'", Position: pos}
	}
	l.readChar()

	// 确定token类型
	var tokenType TokenType
	if isCloseTag {
		tokenType = TokenCloseTag
	} else if isSelfClose {
		tokenType = TokenSelfCloseTag
	} else {
		tokenType = TokenOpenTag
	}

	return Token{
		Type:       tokenType,
		Value:      tagName,
		Attributes: attributes,
		Position:   pos,
	}
}
