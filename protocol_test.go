package markit

import (
	"testing"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig should not return nil")
	}

	if config.CaseSensitive != true {
		t.Errorf("Expected CaseSensitive to be true, got %v", config.CaseSensitive)
	}

	if config.SkipComments != false {
		t.Errorf("Expected SkipComments to be false, got %v", config.SkipComments)
	}

	if config.AllowSelfCloseTags != true {
		t.Errorf("Expected AllowSelfCloseTags to be true, got %v", config.AllowSelfCloseTags)
	}

	if config.CoreMatcher == nil {
		t.Error("Expected CoreMatcher to be initialized")
	}

	if config.AttributeProcessor == nil {
		t.Error("Expected AttributeProcessor to be initialized")
	}
}

// TestNormalizeCase 测试大小写规范化函数
func TestNormalizeCase(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		caseSensitive bool
		expected      string
	}{
		{
			name:          "case sensitive - no change",
			input:         "TestString",
			caseSensitive: true,
			expected:      "TestString",
		},
		{
			name:          "case insensitive - to lower",
			input:         "TestString",
			caseSensitive: false,
			expected:      "teststring",
		},
		{
			name:          "case insensitive - already lower",
			input:         "teststring",
			caseSensitive: false,
			expected:      "teststring",
		},
		{
			name:          "case insensitive - all upper",
			input:         "TESTSTRING",
			caseSensitive: false,
			expected:      "teststring",
		},
		{
			name:          "case sensitive - preserve case",
			input:         "MixedCaseString",
			caseSensitive: true,
			expected:      "MixedCaseString",
		},
		{
			name:          "empty string",
			input:         "",
			caseSensitive: false,
			expected:      "",
		},
		{
			name:          "unicode characters",
			input:         "测试String",
			caseSensitive: false,
			expected:      "测试string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ParserConfig{
				CaseSensitive: tt.caseSensitive,
			}
			result := config.NormalizeCase(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
