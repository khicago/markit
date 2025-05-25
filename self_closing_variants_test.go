package markit

import (
	"testing"
)

// TestSelfClosingTagVariants 测试不同格式的自封闭标签
func TestSelfClosingTagVariants(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "XML style with space",
			input:       `<br />`,
			expectError: false,
			description: "XML风格的自封闭标签，带空格",
		},
		{
			name:        "XML style without space",
			input:       `<br/>`,
			expectError: false,
			description: "XML风格的自封闭标签，不带空格",
		},
		{
			name:        "HTML void element style",
			input:       `<br>`,
			expectError: true, // 当前实现会期望结束标签
			description: "HTML风格的无效元素（void element），无结束标签",
		},
		{
			name:        "IMG with attributes and space",
			input:       `<img src="test.jpg" alt="test" />`,
			expectError: false,
			description: "带属性的自封闭IMG标签，带空格",
		},
		{
			name:        "IMG with attributes without space",
			input:       `<img src="test.jpg" alt="test"/>`,
			expectError: false,
			description: "带属性的自封闭IMG标签，不带空格",
		},
		{
			name:        "HTML void IMG without closing",
			input:       `<img src="test.jpg" alt="test">`,
			expectError: true, // 当前实现会期望结束标签
			description: "HTML风格的IMG无效元素，无结束标签",
		},
		{
			name:        "Multiple self-closing tags",
			input:       `<img src="1.jpg" /><br/><hr />`,
			expectError: false,
			description: "多个不同格式的自封闭标签",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			config.AllowSelfCloseTags = true

			parser := NewParserWithConfig(tt.input, config)
			doc, err := parser.Parse()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but parsing succeeded", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
				return
			}

			if doc == nil {
				t.Errorf("Expected non-nil document for %s", tt.description)
				return
			}

			// 验证所有元素都是自封闭的
			for i, child := range doc.Children {
				element, ok := child.(*Element)
				if !ok {
					t.Errorf("Expected element at index %d for %s", i, tt.description)
					continue
				}

				if !element.SelfClose {
					t.Errorf("Expected self-closing element at index %d for %s", i, tt.description)
				}

				if len(element.Children) != 0 {
					t.Errorf("Self-closing element should have no children at index %d for %s", i, tt.description)
				}
			}
		})
	}
}

// TestHTMLVoidElementsSupport 测试HTML无效元素的支持
func TestHTMLVoidElementsSupport(t *testing.T) {
	// HTML5 无效元素列表
	voidElements := []string{
		"area", "base", "br", "col", "embed", "hr", "img", "input",
		"link", "meta", "param", "source", "track", "wbr",
	}

	t.Run("XML style self-closing", func(t *testing.T) {
		for _, tagName := range voidElements {
			input := "<" + tagName + " />"

			config := DefaultConfig()
			config.AllowSelfCloseTags = true

			parser := NewParserWithConfig(input, config)
			doc, err := parser.Parse()

			if err != nil {
				t.Errorf("Error parsing XML-style self-closing %s: %v", tagName, err)
				continue
			}

			if len(doc.Children) != 1 {
				t.Errorf("Expected 1 child for %s, got %d", tagName, len(doc.Children))
				continue
			}

			element := doc.Children[0].(*Element)
			if element.TagName != tagName {
				t.Errorf("Expected tag name %s, got %s", tagName, element.TagName)
			}

			if !element.SelfClose {
				t.Errorf("Expected %s to be self-closing", tagName)
			}
		}
	})

	t.Run("HTML void style - current limitation", func(t *testing.T) {
		// 这个测试展示当前实现的限制：不支持HTML风格的无效元素
		for _, tagName := range []string{"br", "hr", "img"} {
			input := "<" + tagName + ">"

			config := DefaultConfig()
			config.AllowSelfCloseTags = true

			parser := NewParserWithConfig(input, config)
			_, err := parser.Parse()

			// 当前实现会报错，因为它期望找到结束标签
			if err == nil {
				t.Errorf("Current implementation should fail for HTML void element %s (this test documents current limitation)", tagName)
			}
		}
	})
}

// TestSelfClosingConfigurationControl 测试自封闭标签配置的控制
func TestSelfClosingConfigurationControl(t *testing.T) {
	input := `<img src="test.jpg" />`

	t.Run("AllowSelfCloseTags enabled", func(t *testing.T) {
		config := DefaultConfig()
		config.AllowSelfCloseTags = true

		parser := NewParserWithConfig(input, config)
		doc, err := parser.Parse()

		if err != nil {
			t.Errorf("Expected successful parsing when AllowSelfCloseTags is enabled: %v", err)
			return
		}

		element := doc.Children[0].(*Element)
		if !element.SelfClose {
			t.Error("Expected self-closing element when AllowSelfCloseTags is enabled")
		}
	})

	t.Run("AllowSelfCloseTags disabled", func(t *testing.T) {
		config := DefaultConfig()
		config.AllowSelfCloseTags = false

		parser := NewParserWithConfig(input, config)
		_, err := parser.Parse()

		if err == nil {
			t.Error("Expected parsing error when AllowSelfCloseTags is disabled")
		}
	})
}
