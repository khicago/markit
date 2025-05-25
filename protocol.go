package markit

import (
	"strings"
)

// ParserConfig 解析器配置
type ParserConfig struct {
	// 大小写敏感性配置
	CaseSensitive bool

	// 核心协议匹配器（内置，不可修改）
	CoreMatcher *CoreProtocolMatcher

	// 属性处理器
	AttributeProcessor AttributeProcessor

	// 其他配置选项
	TrimWhitespace     bool
	SkipComments       bool
	AllowEmptyElements bool
	AllowSelfCloseTags bool // 是否允许自封闭标签
}

// DefaultConfig 创建默认配置
func DefaultConfig() *ParserConfig {
	config := &ParserConfig{
		CaseSensitive:      true,
		CoreMatcher:        NewCoreProtocolMatcher(),
		AttributeProcessor: &DefaultAttributeProcessor{},
		TrimWhitespace:     true,
		SkipComments:       false,
		AllowEmptyElements: true,
		AllowSelfCloseTags: true, // 默认允许自封闭标签
	}

	return config
}

// NormalizeCase 根据配置标准化大小写
func (config *ParserConfig) NormalizeCase(s string) string {
	if !config.CaseSensitive {
		return strings.ToLower(s)
	}
	return s
}
