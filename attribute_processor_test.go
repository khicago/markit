package markit

import (
	"testing"
)

// TestDefaultAttributeProcessorMethods 测试DefaultAttributeProcessor的所有方法
func TestDefaultAttributeProcessorMethods(t *testing.T) {
	processor := &DefaultAttributeProcessor{}

	t.Run("ProcessAttribute method", func(t *testing.T) {
		tests := []struct {
			name          string
			key           string
			value         string
			expectedKey   string
			expectedValue interface{}
			expectError   bool
		}{
			{
				name:          "Normal attribute",
				key:           "class",
				value:         "container",
				expectedKey:   "class",
				expectedValue: "container",
				expectError:   false,
			},
			{
				name:          "Empty key",
				key:           "",
				value:         "value",
				expectedKey:   "",
				expectedValue: "value",
				expectError:   false,
			},
			{
				name:          "Empty value",
				key:           "checked",
				value:         "",
				expectedKey:   "checked",
				expectedValue: true, // 空值被处理为布尔属性
				expectError:   false,
			},
			{
				name:          "Boolean attribute with empty value",
				key:           "disabled",
				value:         "",
				expectedKey:   "disabled",
				expectedValue: true,
				expectError:   false,
			},
			{
				name:          "Attribute with special characters",
				key:           "data-test",
				value:         "special@value#123",
				expectedKey:   "data-test",
				expectedValue: "special@value#123",
				expectError:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				key, value, err := processor.ProcessAttribute(tt.key, tt.value)

				if tt.expectError && err == nil {
					t.Errorf("expected error but got none")
				}
				if !tt.expectError && err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if key != tt.expectedKey {
					t.Errorf("expected key %q, got %q", tt.expectedKey, key)
				}

				if value != tt.expectedValue {
					t.Errorf("expected value %v, got %v", tt.expectedValue, value)
				}
			})
		}
	})

	t.Run("IsBooleanAttribute method", func(t *testing.T) {
		tests := []struct {
			name     string
			key      string
			expected bool
		}{
			{
				name:     "Standard boolean attribute - checked",
				key:      "checked",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - disabled",
				key:      "disabled",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - selected",
				key:      "selected",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - readonly",
				key:      "readonly",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - required",
				key:      "required",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - autofocus",
				key:      "autofocus",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - autoplay",
				key:      "autoplay",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - controls",
				key:      "controls",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - defer",
				key:      "defer",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - hidden",
				key:      "hidden",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - loop",
				key:      "loop",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - multiple",
				key:      "multiple",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - muted",
				key:      "muted",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - open",
				key:      "open",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - reversed",
				key:      "reversed",
				expected: true,
			},
			{
				name:     "Standard boolean attribute - scoped",
				key:      "scoped",
				expected: true,
			},
			{
				name:     "Non-boolean attribute - class",
				key:      "class",
				expected: false,
			},
			{
				name:     "Non-boolean attribute - id",
				key:      "id",
				expected: false,
			},
			{
				name:     "Non-boolean attribute - src",
				key:      "src",
				expected: false,
			},
			{
				name:     "Non-boolean attribute - href",
				key:      "href",
				expected: false,
			},
			{
				name:     "Empty key",
				key:      "",
				expected: false,
			},
			{
				name:     "Custom attribute",
				key:      "data-custom",
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := processor.IsBooleanAttribute(tt.key)
				if result != tt.expected {
					t.Errorf("expected IsBooleanAttribute(%q) to be %t, got %t", tt.key, tt.expected, result)
				}
			})
		}
	})
}

// TestDefaultAttributeProcessorIntegration 测试属性处理器的集成场景
func TestDefaultAttributeProcessorIntegration(t *testing.T) {
	processor := &DefaultAttributeProcessor{}

	t.Run("Processing boolean attributes", func(t *testing.T) {
		// 测试布尔属性的处理
		key, value, err := processor.ProcessAttribute("checked", "")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if key != "checked" {
			t.Errorf("expected key 'checked', got %q", key)
		}

		if value != true {
			t.Errorf("expected boolean value true, got %v", value)
		}
	})

	t.Run("Processing regular attributes", func(t *testing.T) {
		// 测试普通属性的处理
		key, value, err := processor.ProcessAttribute("class", "container")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if key != "class" {
			t.Errorf("expected key 'class', got %q", key)
		}

		if value != "container" {
			t.Errorf("expected string value 'container', got %v", value)
		}
	})

	t.Run("Processing attributes with special values", func(t *testing.T) {
		tests := []struct {
			key   string
			value string
		}{
			{"title", "Hello \"World\""},
			{"data-json", `{"key": "value"}`},
			{"style", "color: red; background: blue;"},
			{"onclick", "alert('test');"},
		}

		for _, tt := range tests {
			key, value, err := processor.ProcessAttribute(tt.key, tt.value)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", tt.key, err)
			}

			if key != tt.key {
				t.Errorf("expected key %q, got %q", tt.key, key)
			}

			if value != tt.value {
				t.Errorf("expected value %q, got %v", tt.value, value)
			}
		}
	})
}

// TestAttributeProcessorEdgeCasesAdvanced 测试属性处理器的高级边缘情况
func TestAttributeProcessorEdgeCasesAdvanced(t *testing.T) {
	processor := &DefaultAttributeProcessor{}

	t.Run("Unicode attribute names and values", func(t *testing.T) {
		key, value, err := processor.ProcessAttribute("数据", "测试值")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if key != "数据" {
			t.Errorf("expected key '数据', got %q", key)
		}

		if value != "测试值" {
			t.Errorf("expected value '测试值', got %v", value)
		}
	})

	t.Run("Very long attribute values", func(t *testing.T) {
		longValue := string(make([]byte, 1000))
		for i := range longValue {
			longValue = longValue[:i] + "a" + longValue[i+1:]
		}

		key, value, err := processor.ProcessAttribute("long", longValue)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if key != "long" {
			t.Errorf("expected key 'long', got %q", key)
		}

		if value != longValue {
			t.Errorf("expected long value to be preserved")
		}
	})

	t.Run("Attribute names with special characters", func(t *testing.T) {
		specialKeys := []string{
			"data-test",
			"aria-label",
			"ng:if",
			"v-model",
			"@click",
			":disabled",
		}

		for _, specialKey := range specialKeys {
			key, value, err := processor.ProcessAttribute(specialKey, "test")
			if err != nil {
				t.Errorf("unexpected error for key %q: %v", specialKey, err)
			}

			if key != specialKey {
				t.Errorf("expected key %q, got %q", specialKey, key)
			}

			if value != "test" {
				t.Errorf("expected value 'test', got %v", value)
			}
		}
	})
}
