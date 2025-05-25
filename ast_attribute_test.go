package markit

import (
	"testing"
)

// TestASTAttributeHandling 测试属性处理
func TestASTAttributeHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "single attribute",
			input: `<element id="test">content</element>`,
			expected: map[string]string{
				"id": "test",
			},
		},
		{
			name:  "multiple attributes",
			input: `<element id="test" class="main" data-value="123">content</element>`,
			expected: map[string]string{
				"id":         "test",
				"class":      "main",
				"data-value": "123",
			},
		},
		{
			name:  "boolean attributes",
			input: `<input type="checkbox" checked disabled />`,
			expected: map[string]string{
				"type":     "checkbox",
				"checked":  "",
				"disabled": "",
			},
		},
		{
			name:  "mixed attributes",
			input: `<input type="text" name="username" required placeholder="Enter username" />`,
			expected: map[string]string{
				"type":        "text",
				"name":        "username",
				"required":    "",
				"placeholder": "Enter username",
			},
		},
		{
			name:  "attributes with special characters",
			input: `<element data-test="value with spaces" data-json='{"key": "value"}' onclick="alert('hello')">content</element>`,
			expected: map[string]string{
				"data-test": "value with spaces",
				"data-json": `{"key": "value"}`,
				"onclick":   "alert('hello')",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser(test.input)
			doc, err := parser.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			element := doc.Children[0].(*Element)

			// 验证属性数量
			if len(element.Attributes) != len(test.expected) {
				t.Errorf("expected %d attributes, got %d", len(test.expected), len(element.Attributes))
			}

			// 验证每个属性
			for key, expectedValue := range test.expected {
				actualValue, exists := element.Attributes[key]
				if !exists {
					t.Errorf("expected attribute %q not found", key)
					continue
				}
				if actualValue != expectedValue {
					t.Errorf("expected attribute %s=%q, got %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}

// TestDefaultAttributeProcessor 测试默认属性处理器
func TestDefaultAttributeProcessor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "quoted attributes",
			input: `<element id="test" class='main'>content</element>`,
			expected: map[string]string{
				"id":    "test",
				"class": "main",
			},
		},
		{
			name:  "unquoted attributes",
			input: `<element id=test class=main>content</element>`,
			expected: map[string]string{
				"id":    "test",
				"class": "main",
			},
		},
		{
			name:  "boolean attributes",
			input: `<input checked disabled readonly />`,
			expected: map[string]string{
				"checked":  "",
				"disabled": "",
				"readonly": "",
			},
		},
		{
			name:  "mixed quote styles",
			input: `<element id="double" class='single' data=unquoted>content</element>`,
			expected: map[string]string{
				"id":    "double",
				"class": "single",
				"data":  "unquoted",
			},
		},
		{
			name:  "empty attribute values",
			input: `<element id="" class=''>content</element>`,
			expected: map[string]string{
				"id":    "",
				"class": "",
			},
		},
		{
			name:  "attributes with special characters",
			input: `<element data-test="value-with-dashes" data_underscore="value_with_underscores">content</element>`,
			expected: map[string]string{
				"data-test":       "value-with-dashes",
				"data_underscore": "value_with_underscores",
			},
		},
		{
			name:  "attributes with numbers",
			input: `<element id="test123" data-value="456" tabindex="1">content</element>`,
			expected: map[string]string{
				"id":         "test123",
				"data-value": "456",
				"tabindex":   "1",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser(test.input)
			doc, err := parser.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			element := doc.Children[0].(*Element)

			// 验证属性数量
			if len(element.Attributes) != len(test.expected) {
				t.Errorf("expected %d attributes, got %d", len(test.expected), len(element.Attributes))
			}

			// 验证每个属性
			for key, expectedValue := range test.expected {
				actualValue, exists := element.Attributes[key]
				if !exists {
					t.Errorf("expected attribute %q not found", key)
					continue
				}
				if actualValue != expectedValue {
					t.Errorf("expected attribute %s=%q, got %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}

// TestAttributeProcessorEdgeCases 测试属性处理器的边界情况
func TestAttributeProcessorEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:  "attribute with equals in value",
			input: `<element data-formula="x=y+z">content</element>`,
			expected: map[string]string{
				"data-formula": "x=y+z",
			},
		},
		{
			name:  "attribute with quotes in value",
			input: `<element title='He said "Hello"'>content</element>`,
			expected: map[string]string{
				"title": `He said "Hello"`,
			},
		},
		{
			name:  "attribute with escaped quotes",
			input: `<element data-json="{\"key\": \"value\"}">content</element>`,
			expected: map[string]string{
				"data-json": `{"key": "value"}`,
			},
		},
		{
			name:  "attribute with whitespace",
			input: `<element title="  spaced value  ">content</element>`,
			expected: map[string]string{
				"title": "  spaced value  ",
			},
		},
		{
			name:  "multiple spaces between attributes",
			input: `<element    id="test"     class="main"    >content</element>`,
			expected: map[string]string{
				"id":    "test",
				"class": "main",
			},
		},
		{
			name:  "attribute with newlines",
			input: "<element\n  id=\"test\"\n  class=\"main\"\n>content</element>",
			expected: map[string]string{
				"id":    "test",
				"class": "main",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser(test.input)
			doc, err := parser.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			element := doc.Children[0].(*Element)

			// 验证属性数量
			if len(element.Attributes) != len(test.expected) {
				t.Errorf("expected %d attributes, got %d", len(test.expected), len(element.Attributes))
			}

			// 验证每个属性
			for key, expectedValue := range test.expected {
				actualValue, exists := element.Attributes[key]
				if !exists {
					t.Errorf("expected attribute %q not found", key)
					continue
				}
				if actualValue != expectedValue {
					t.Errorf("expected attribute %s=%q, got %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}
