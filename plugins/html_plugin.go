package plugins

import (
	"strings"
)

// HTMLPlugin HTML扩展插件
// 提供HTML特有的功能：void elements、case insensitive、boolean attributes等
type HTMLPlugin struct{}

// Name 返回插件名称
func (p *HTMLPlugin) Name() string {
	return "html"
}

// Description 返回插件描述
func (p *HTMLPlugin) Description() string {
	return "HTML extension plugin providing HTML5 void elements, case insensitive parsing, and boolean attributes support"
}

// Install 安装HTML插件到协议匹配器
func (p *HTMLPlugin) Install(matcher *ProtocolMatcher) error {
	// HTML插件主要提供配置而不是新的协议
	// 实际的协议注册在配置创建时处理
	return nil
}

// NewHTMLPlugin 创建新的HTML插件实例
func NewHTMLPlugin() *HTMLPlugin {
	return &HTMLPlugin{}
}

// GetHTML5VoidElements 返回HTML5标准的void elements列表
func (p *HTMLPlugin) GetHTML5VoidElements() []string {
	return []string{
		"area", "base", "br", "col", "embed", "hr", "img", "input",
		"link", "meta", "param", "source", "track", "wbr",
	}
}

// GetHTML5VoidElementsMap 返回HTML5标准的void elements映射
func (p *HTMLPlugin) GetHTML5VoidElementsMap() map[string]bool {
	voidElements := make(map[string]bool)
	for _, element := range p.GetHTML5VoidElements() {
		voidElements[element] = true
	}
	return voidElements
}

// IsHTML5VoidElement 检查是否是HTML5标准的void element
func (p *HTMLPlugin) IsHTML5VoidElement(tagName string) bool {
	voidElements := map[string]bool{
		"area": true, "base": true, "br": true, "col": true,
		"embed": true, "hr": true, "img": true, "input": true,
		"link": true, "meta": true, "param": true, "source": true,
		"track": true, "wbr": true,
	}
	return voidElements[strings.ToLower(tagName)]
}

// HTMLAttributeProcessor HTML特定的属性处理器
// 实现了核心的 AttributeProcessor 接口
type HTMLAttributeProcessor struct{}

// ProcessAttribute 处理单个属性 - 实现核心接口
func (hap *HTMLAttributeProcessor) ProcessAttribute(key, value string) (string, interface{}, error) {
	// HTML属性名不区分大小写
	normalizedKey := strings.ToLower(key)

	// 处理布尔属性
	if hap.IsBooleanAttribute(normalizedKey) {
		// HTML布尔属性：如果存在就是true，值可以是空字符串或属性名本身
		if value == "" || value == normalizedKey {
			return normalizedKey, true, nil
		}
	}

	return normalizedKey, value, nil
}

// IsBooleanAttribute 检查是否是HTML布尔属性 - 实现核心接口
func (hap *HTMLAttributeProcessor) IsBooleanAttribute(key string) bool {
	// HTML5标准布尔属性
	booleanAttrs := map[string]bool{
		"autofocus": true,
		"autoplay":  true,
		"checked":   true,
		"controls":  true,
		"defer":     true,
		"disabled":  true,
		"hidden":    true,
		"loop":      true,
		"multiple":  true,
		"muted":     true,
		"open":      true,
		"readonly":  true,
		"required":  true,
		"reversed":  true,
		"scoped":    true,
		"selected":  true,
	}

	return booleanAttrs[strings.ToLower(key)]
}

// NewHTMLAttributeProcessor 创建HTML属性处理器
func NewHTMLAttributeProcessor() *HTMLAttributeProcessor {
	return &HTMLAttributeProcessor{}
}
