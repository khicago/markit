package markit

import (
	"github.com/khicago/markit/plugins"
)

// HTMLConfig 创建适用于 HTML 的配置
// 使用HTML插件提供的功能来配置解析器
func HTMLConfig() *ParserConfig {
	htmlPlugin := plugins.NewHTMLPlugin()

	config := &ParserConfig{
		CaseSensitive:      false,                    // HTML不区分大小写
		CoreMatcher:        NewCoreProtocolMatcher(), // 必须设置核心协议匹配器
		AttributeProcessor: plugins.NewHTMLAttributeProcessor(),
		TrimWhitespace:     true,
		SkipComments:       false,
		AllowEmptyElements: true,
		AllowSelfCloseTags: true,
		VoidElements:       htmlPlugin.GetHTML5VoidElementsMap(),
	}

	return config
}

// DefaultHTMLConfig HTMLConfig的别名，提供更明确的命名
func DefaultHTMLConfig() *ParserConfig {
	return HTMLConfig()
}
