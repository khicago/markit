package markit

// CoreProtocol MarkIt 核心协议定义
// 这些是 MarkIt 的内置协议，不能被覆盖或移除
type CoreProtocol struct {
	Name        string
	OpenSeq     string
	CloseSeq    string
	SelfClose   string
	TokenType   TokenType
	Description string
}

// GetCoreProtocols 返回 MarkIt 的核心协议
// 这些协议是内置的，不需要注册，也不能被覆盖
func GetCoreProtocols() []CoreProtocol {
	return []CoreProtocol{
		{
			Name:        "markit-standard-tag",
			OpenSeq:     "<",
			CloseSeq:    ">",
			SelfClose:   "/",
			TokenType:   TokenOpenTag,
			Description: "MarkIt standard tags <tag>",
		},
		{
			Name:        "markit-comment",
			OpenSeq:     "<!--",
			CloseSeq:    "-->",
			SelfClose:   "",
			TokenType:   TokenComment,
			Description: "MarkIt comments <!-- -->",
		},
	}
}

// CoreProtocolMatcher MarkIt 核心协议匹配器
type CoreProtocolMatcher struct {
	protocols []CoreProtocol
	maxLen    int
}

// NewCoreProtocolMatcher 创建核心协议匹配器
func NewCoreProtocolMatcher() *CoreProtocolMatcher {
	protocols := GetCoreProtocols()
	matcher := &CoreProtocolMatcher{
		protocols: protocols,
		maxLen:    0,
	}

	// 计算最大长度
	for _, protocol := range protocols {
		if len(protocol.OpenSeq) > matcher.maxLen {
			matcher.maxLen = len(protocol.OpenSeq)
		}
	}

	return matcher
}

// MatchProtocol 匹配核心协议
func (cpm *CoreProtocolMatcher) MatchProtocol(input string, pos int) *CoreProtocol {
	// 按开始序列长度从长到短匹配，确保最长匹配优先
	for length := cpm.maxLen; length >= 1; length-- {
		if pos+length > len(input) {
			continue
		}

		candidate := input[pos : pos+length]
		for i := range cpm.protocols {
			if cpm.protocols[i].OpenSeq == candidate {
				return &cpm.protocols[i]
			}
		}
	}
	return nil
}
