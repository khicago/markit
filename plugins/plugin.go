package plugins

import (
	"fmt"
)

// ExtendedProtocol 扩展协议定义
// 用于XML等扩展功能
type ExtendedProtocol struct {
	Name        string
	OpenSeq     string
	CloseSeq    string
	SelfClose   string
	TokenType   int // 使用int避免循环依赖
	Description string
}

// ProtocolMatcher 扩展协议匹配器
type ProtocolMatcher struct {
	protocols []ExtendedProtocol
	maxLen    int
}

// NewProtocolMatcher 创建扩展协议匹配器
func NewProtocolMatcher() *ProtocolMatcher {
	return &ProtocolMatcher{
		protocols: make([]ExtendedProtocol, 0),
		maxLen:    0,
	}
}

// RegisterProtocol 注册新的扩展协议
func (pm *ProtocolMatcher) RegisterProtocol(protocol ExtendedProtocol) {
	pm.protocols = append(pm.protocols, protocol)
	if len(protocol.OpenSeq) > pm.maxLen {
		pm.maxLen = len(protocol.OpenSeq)
	}
}

// MatchProtocol 匹配扩展协议
func (pm *ProtocolMatcher) MatchProtocol(input string, pos int) *ExtendedProtocol {
	// 按开始序列长度从长到短匹配，确保最长匹配优先
	for length := pm.maxLen; length >= 1; length-- {
		if pos+length > len(input) {
			continue
		}

		candidate := input[pos : pos+length]
		for i := range pm.protocols {
			if pm.protocols[i].OpenSeq == candidate {
				return &pm.protocols[i]
			}
		}
	}
	return nil
}

// Plugin 插件接口
type Plugin interface {
	// Name 返回插件名称
	Name() string

	// Install 安装插件到配置中
	Install(matcher *ProtocolMatcher) error

	// Description 返回插件描述
	Description() string
}

// PluginManager 插件管理器
type PluginManager struct {
	plugins map[string]Plugin
}

// NewPluginManager 创建插件管理器
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

// RegisterPlugin 注册插件
func (pm *PluginManager) RegisterPlugin(plugin Plugin) error {
	name := plugin.Name()
	if _, exists := pm.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}
	pm.plugins[name] = plugin
	return nil
}

// LoadPlugins 加载插件到匹配器
func (pm *PluginManager) LoadPlugins(matcher *ProtocolMatcher, pluginNames ...string) error {
	for _, name := range pluginNames {
		plugin, exists := pm.plugins[name]
		if !exists {
			return fmt.Errorf("plugin %s not found", name)
		}
		if err := plugin.Install(matcher); err != nil {
			return fmt.Errorf("failed to install plugin %s: %w", name, err)
		}
	}
	return nil
}
