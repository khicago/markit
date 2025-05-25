package markit

import (
	"testing"
)

// TestVoidElementsConfiguration 测试 void element 配置功能
func TestVoidElementsConfiguration(t *testing.T) {
	t.Run("默认配置不包含 void elements", func(t *testing.T) {
		config := DefaultConfig()

		if config.IsVoidElement("br") {
			t.Error("默认配置不应该包含 br 作为 void element")
		}

		if config.IsVoidElement("img") {
			t.Error("默认配置不应该包含 img 作为 void element")
		}
	})

	t.Run("HTML 配置包含标准 void elements", func(t *testing.T) {
		config := HTMLConfig()

		// 测试标准 HTML5 void elements
		voidElements := []string{"br", "hr", "img", "input", "area", "base", "col", "embed", "link", "meta", "param", "source", "track", "wbr"}

		for _, element := range voidElements {
			if !config.IsVoidElement(element) {
				t.Errorf("HTML 配置应该包含 %s 作为 void element", element)
			}
		}

		// 测试非 void element
		if config.IsVoidElement("div") {
			t.Error("div 不应该是 void element")
		}

		if config.IsVoidElement("p") {
			t.Error("p 不应该是 void element")
		}
	})

	t.Run("手动添加和移除 void elements", func(t *testing.T) {
		config := DefaultConfig()

		// 添加自定义 void element
		config.AddVoidElement("custom-tag")

		if !config.IsVoidElement("custom-tag") {
			t.Error("添加的 custom-tag 应该是 void element")
		}

		// 移除 void element
		config.RemoveVoidElement("custom-tag")

		if config.IsVoidElement("custom-tag") {
			t.Error("移除后的 custom-tag 不应该是 void element")
		}
	})

	t.Run("设置完整的 void elements 列表", func(t *testing.T) {
		config := DefaultConfig()

		customVoidElements := []string{"foo", "bar", "baz"}
		config.SetVoidElements(customVoidElements)

		for _, element := range customVoidElements {
			if !config.IsVoidElement(element) {
				t.Errorf("%s 应该是 void element", element)
			}
		}

		// 原来的不应该存在
		if config.IsVoidElement("br") {
			t.Error("设置新列表后，br 不应该是 void element")
		}
	})

	t.Run("大小写敏感性测试", func(t *testing.T) {
		// 大小写敏感配置
		config := DefaultConfig()
		config.CaseSensitive = true
		config.AddVoidElement("BR")

		if !config.IsVoidElement("BR") {
			t.Error("大小写敏感模式下，BR 应该是 void element")
		}

		if config.IsVoidElement("br") {
			t.Error("大小写敏感模式下，br 不应该是 void element")
		}

		// 大小写不敏感配置
		config.CaseSensitive = false
		config.SetVoidElements([]string{"BR"})

		if !config.IsVoidElement("br") {
			t.Error("大小写不敏感模式下，br 应该是 void element")
		}

		if !config.IsVoidElement("BR") {
			t.Error("大小写不敏感模式下，BR 应该是 void element")
		}
	})
}

// TestVoidElementsParsing 测试 void element 解析功能
func TestVoidElementsParsing(t *testing.T) {
	t.Run("HTML style void elements", func(t *testing.T) {
		config := HTMLConfig()

		testCases := []struct {
			name    string
			input   string
			tagName string
		}{
			{"br tag", "<br>", "br"},
			{"hr tag", "<hr>", "hr"},
			{"img tag", `<img src="test.jpg">`, "img"},
			{"input tag", `<input type="text" name="test">`, "input"},
			{"meta tag", `<meta charset="utf-8">`, "meta"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				parser := NewParserWithConfig(tc.input, config)
				doc, err := parser.Parse()

				if err != nil {
					t.Fatalf("解析 %s 失败: %v", tc.input, err)
				}

				if len(doc.Children) != 1 {
					t.Fatalf("期望 1 个子节点，得到 %d", len(doc.Children))
				}

				element, ok := doc.Children[0].(*Element)
				if !ok {
					t.Fatalf("期望 Element 节点，得到 %T", doc.Children[0])
				}

				if element.TagName != tc.tagName {
					t.Errorf("期望标签名 %s，得到 %s", tc.tagName, element.TagName)
				}

				if !element.SelfClose {
					t.Errorf("void element %s 应该是自闭合的", tc.tagName)
				}

				if len(element.Children) != 0 {
					t.Errorf("void element %s 不应该有子节点", tc.tagName)
				}
			})
		}
	})

	t.Run("Mixed content with void elements", func(t *testing.T) {
		config := HTMLConfig()
		input := `<div>
			<p>Some text</p>
			<br>
			<img src="test.jpg" alt="test">
			<hr>
			<p>More text</p>
		</div>`

		parser := NewParserWithConfig(input, config)
		doc, err := parser.Parse()

		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		if len(doc.Children) != 1 {
			t.Fatalf("期望 1 个子节点，得到 %d", len(doc.Children))
		}

		divElement := doc.Children[0].(*Element)
		if divElement.TagName != "div" {
			t.Errorf("期望标签名 div，得到 %s", divElement.TagName)
		}

		// 检查 div 内的子元素
		var elements []*Element
		for _, child := range divElement.Children {
			if element, ok := child.(*Element); ok {
				elements = append(elements, element)
			}
		}

		// 应该有 p, br, img, hr, p 五个元素
		expectedTags := []struct {
			tagName   string
			selfClose bool
		}{
			{"p", false},
			{"br", true},
			{"img", true},
			{"hr", true},
			{"p", false},
		}

		if len(elements) != len(expectedTags) {
			t.Fatalf("期望 %d 个元素，得到 %d", len(expectedTags), len(elements))
		}

		for i, expected := range expectedTags {
			if elements[i].TagName != expected.tagName {
				t.Errorf("元素 %d: 期望标签名 %s，得到 %s", i, expected.tagName, elements[i].TagName)
			}

			if elements[i].SelfClose != expected.selfClose {
				t.Errorf("元素 %d (%s): 期望 SelfClose=%v，得到 %v", i, expected.tagName, expected.selfClose, elements[i].SelfClose)
			}
		}
	})

	t.Run("Custom void elements", func(t *testing.T) {
		config := DefaultConfig()
		config.SetVoidElements([]string{"my-icon", "my-separator"})

		input := `<container>
			<my-icon name="star">
			<p>Text content</p>
			<my-separator type="dotted">
		</container>`

		parser := NewParserWithConfig(input, config)
		doc, err := parser.Parse()

		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		containerElement := doc.Children[0].(*Element)

		var elements []*Element
		for _, child := range containerElement.Children {
			if element, ok := child.(*Element); ok {
				elements = append(elements, element)
			}
		}

		expectedElements := []struct {
			tagName   string
			selfClose bool
		}{
			{"my-icon", true},
			{"p", false},
			{"my-separator", true},
		}

		if len(elements) != len(expectedElements) {
			t.Fatalf("期望 %d 个元素，得到 %d", len(expectedElements), len(elements))
		}

		for i, expected := range expectedElements {
			if elements[i].TagName != expected.tagName {
				t.Errorf("元素 %d: 期望标签名 %s，得到 %s", i, expected.tagName, elements[i].TagName)
			}

			if elements[i].SelfClose != expected.selfClose {
				t.Errorf("元素 %d (%s): 期望 SelfClose=%v，得到 %v", i, expected.tagName, expected.selfClose, elements[i].SelfClose)
			}
		}
	})

	t.Run("Void elements with attributes", func(t *testing.T) {
		config := HTMLConfig()
		input := `<img src="test.jpg" alt="Test Image" width="100" height="200">`

		parser := NewParserWithConfig(input, config)
		doc, err := parser.Parse()

		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		element := doc.Children[0].(*Element)

		expectedAttributes := map[string]string{
			"src":    "test.jpg",
			"alt":    "Test Image",
			"width":  "100",
			"height": "200",
		}

		for key, expectedValue := range expectedAttributes {
			if actualValue, exists := element.Attributes[key]; !exists {
				t.Errorf("缺少属性 %s", key)
			} else if actualValue != expectedValue {
				t.Errorf("属性 %s: 期望值 %s，得到 %s", key, expectedValue, actualValue)
			}
		}

		if !element.SelfClose {
			t.Error("带属性的 img 元素应该是自闭合的")
		}
	})
}

// TestVoidElementsEdgeCases 测试 void element 的边界情况
func TestVoidElementsEdgeCases(t *testing.T) {
	t.Run("Void element 配置为 nil", func(t *testing.T) {
		config := DefaultConfig()
		config.VoidElements = nil

		input := "<br>"
		parser := NewParserWithConfig(input, config)
		_, err := parser.Parse()

		// 应该报错，因为没有找到结束标签
		if err == nil {
			t.Error("VoidElements 为 nil 时，应该报错")
		}
	})

	t.Run("XML style 和 HTML style 混合", func(t *testing.T) {
		config := HTMLConfig()
		input := `<div>
			<br />
			<br>
			<img src="test.jpg" />
			<img src="test2.jpg">
		</div>`

		parser := NewParserWithConfig(input, config)
		doc, err := parser.Parse()

		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		divElement := doc.Children[0].(*Element)

		var elements []*Element
		for _, child := range divElement.Children {
			if element, ok := child.(*Element); ok {
				elements = append(elements, element)
			}
		}

		// 所有的 br 和 img 都应该是自闭合的
		for i, element := range elements {
			if element.TagName == "br" || element.TagName == "img" {
				if !element.SelfClose {
					t.Errorf("元素 %d (%s) 应该是自闭合的", i, element.TagName)
				}
			}
		}
	})

	t.Run("Non-void element 不应该受影响", func(t *testing.T) {
		config := HTMLConfig()
		input := `<div><p>content</p></div>`

		parser := NewParserWithConfig(input, config)
		doc, err := parser.Parse()

		if err != nil {
			t.Fatalf("解析失败: %v", err)
		}

		divElement := doc.Children[0].(*Element)
		pElement := divElement.Children[0].(*Element)

		if divElement.SelfClose {
			t.Error("div 元素不应该是自闭合的")
		}

		if pElement.SelfClose {
			t.Error("p 元素不应该是自闭合的")
		}
	})
}
