package plugins

// XMLPlugin XML扩展插件
// 提供XML特有的功能：CDATA、DOCTYPE、ProcessingInstruction等
type XMLPlugin struct{}

// Name 返回插件名称
func (p *XMLPlugin) Name() string {
	return "xml"
}

// Description 返回插件描述
func (p *XMLPlugin) Description() string {
	return "XML extension plugin providing CDATA, DOCTYPE, and ProcessingInstruction support"
}

// Install 安装XML插件到协议匹配器
func (p *XMLPlugin) Install(matcher *ProtocolMatcher) error {
	// XML 处理指令 <?xml ... ?>
	matcher.RegisterProtocol(ExtendedProtocol{
		Name:        "xml-processing-instruction",
		OpenSeq:     "<?",
		CloseSeq:    "?>",
		SelfClose:   "",
		TokenType:   9, // TokenProcessingInstruction
		Description: "XML processing instructions",
	})

	// CDATA 节 <![CDATA[...]]>
	matcher.RegisterProtocol(ExtendedProtocol{
		Name:        "xml-cdata",
		OpenSeq:     "<![CDATA[",
		CloseSeq:    "]]>",
		SelfClose:   "",
		TokenType:   11, // TokenCDATA
		Description: "XML CDATA sections",
	})

	// DOCTYPE 声明 <!DOCTYPE ...>
	matcher.RegisterProtocol(ExtendedProtocol{
		Name:        "xml-doctype",
		OpenSeq:     "<!",
		CloseSeq:    ">",
		SelfClose:   "",
		TokenType:   10, // TokenDoctype
		Description: "XML DOCTYPE declarations",
	})

	return nil
}

// NewXMLPlugin 创建新的XML插件实例
func NewXMLPlugin() *XMLPlugin {
	return &XMLPlugin{}
}
