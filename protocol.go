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

	// Void Elements 配置
	VoidElements map[string]bool // 定义哪些标签是 void element（如 HTML 的 br, hr, img 等）
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
		AllowSelfCloseTags: true,                  // 默认允许自封闭标签
		VoidElements:       make(map[string]bool), // 默认不定义任何 void element
	}

	return config
}

// IsVoidElement 检查指定标签是否是 void element
func (config *ParserConfig) IsVoidElement(tagName string) bool {
	if config.VoidElements == nil {
		return false
	}

	// 根据大小写敏感性配置进行标准化
	normalizedTagName := config.NormalizeCase(tagName)
	return config.VoidElements[normalizedTagName]
}

// AddVoidElement 添加 void element
func (config *ParserConfig) AddVoidElement(tagName string) {
	if config.VoidElements == nil {
		config.VoidElements = make(map[string]bool)
	}
	normalizedTagName := config.NormalizeCase(tagName)
	config.VoidElements[normalizedTagName] = true
}

// RemoveVoidElement 移除 void element
func (config *ParserConfig) RemoveVoidElement(tagName string) {
	if config.VoidElements == nil {
		return
	}
	normalizedTagName := config.NormalizeCase(tagName)
	delete(config.VoidElements, normalizedTagName)
}

// SetVoidElements 设置完整的 void elements 列表
func (config *ParserConfig) SetVoidElements(elements []string) {
	config.VoidElements = make(map[string]bool)
	for _, element := range elements {
		normalizedElement := config.NormalizeCase(element)
		config.VoidElements[normalizedElement] = true
	}
}

// NormalizeCase 根据配置标准化大小写
func (config *ParserConfig) NormalizeCase(s string) string {
	if !config.CaseSensitive {
		return strings.ToLower(s)
	}
	return s
}
