package markit

import (
	"strings"
	"testing"
)

// BenchmarkLexerSimple 基准测试：简单词法分析
func BenchmarkLexerSimple(b *testing.B) {
	input := `<element attr="value">text</element>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			token := lexer.NextToken()
			if token.Type == TokenEOF {
				break
			}
		}
	}
}

// BenchmarkLexerComplex 基准测试：复杂词法分析
func BenchmarkLexerComplex(b *testing.B) {
	input := `<root>
		<element id="test" class="example" disabled>
			<child>Some text content</child>
			<self-close attr="value" />
		</element>
		<!-- comment -->
		<another>More content</another>
	</root>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			token := lexer.NextToken()
			if token.Type == TokenEOF {
				break
			}
		}
	}
}

// BenchmarkParserSimple 基准测试：简单解析
func BenchmarkParserSimple(b *testing.B) {
	input := `<element attr="value">text</element>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(input)
		_, err := parser.Parse()
		if err != nil {
			b.Fatalf("parsing failed: %v", err)
		}
	}
}

// BenchmarkParserComplex 基准测试：复杂解析
func BenchmarkParserComplex(b *testing.B) {
	input := `<root>
		<element id="test" class="example" disabled>
			<child>Some text content</child>
			<self-close attr="value" />
		</element>
		<!-- comment -->
		<another>More content</another>
	</root>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(input)
		_, err := parser.Parse()
		if err != nil {
			b.Fatalf("parsing failed: %v", err)
		}
	}
}

// BenchmarkParserLarge 基准测试：大文档解析
func BenchmarkParserLarge(b *testing.B) {
	// 生成大型文档
	var builder strings.Builder
	builder.WriteString("<root>")

	for i := 0; i < 1000; i++ {
		builder.WriteString(`<item id="`)
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString(`" class="test">`)
		builder.WriteString("Content ")
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString("</item>")
	}

	builder.WriteString("</root>")
	input := builder.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(input)
		_, err := parser.Parse()
		if err != nil {
			b.Fatalf("parsing failed: %v", err)
		}
	}
}

// BenchmarkAttributeProcessing 测试属性处理性能
func BenchmarkAttributeProcessing(b *testing.B) {
	input := `<element attr1="value1" attr2="value2" attr3="value3"></element>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(input)
		_, err := parser.Parse()
		if err != nil {
			b.Fatalf("parsing failed: %v", err)
		}
	}
}

// BenchmarkPrettyPrint 基准测试：格式化输出
func BenchmarkPrettyPrint(b *testing.B) {
	input := `<root><element attr="value"><child>text</child></element></root>`
	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		b.Fatalf("parsing failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PrettyPrint(doc)
	}
}

// BenchmarkMemoryAllocation 基准测试：内存分配
func BenchmarkMemoryAllocation(b *testing.B) {
	input := `<root><child>text</child><child>text</child><child>text</child></root>`

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(input)
		_, err := parser.Parse()
		if err != nil {
			b.Fatalf("parsing failed: %v", err)
		}
	}
}

// BenchmarkNestedElements 基准测试：嵌套元素
func BenchmarkNestedElements(b *testing.B) {
	// 生成深度嵌套的文档
	var builder strings.Builder
	depth := 100

	// 开始标签
	for i := 0; i < depth; i++ {
		builder.WriteString("<level")
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString(">")
	}

	builder.WriteString("content")

	// 结束标签
	for i := depth - 1; i >= 0; i-- {
		builder.WriteString("</level")
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString(">")
	}

	input := builder.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(input)
		_, err := parser.Parse()
		if err != nil {
			b.Fatalf("parsing failed: %v", err)
		}
	}
}

// BenchmarkSelfClosingElements 基准测试：自闭合元素
func BenchmarkSelfClosingElements(b *testing.B) {
	var builder strings.Builder
	builder.WriteString("<root>")

	for i := 0; i < 100; i++ {
		builder.WriteString(`<img src="image`)
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString(`.jpg" alt="Image `)
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString(`" />`)
	}

	builder.WriteString("</root>")
	input := builder.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(input)
		_, err := parser.Parse()
		if err != nil {
			b.Fatalf("parsing failed: %v", err)
		}
	}
}

// BenchmarkMixedContent 基准测试：混合内容
func BenchmarkMixedContent(b *testing.B) {
	input := `<article>
		<h1>Title</h1>
		<p>This is a paragraph with <strong>bold text</strong> and <em>italic text</em>.</p>
		<ul>
			<li>Item 1</li>
			<li>Item 2 with <a href="link">link</a></li>
			<li>Item 3</li>
		</ul>
		<img src="image.jpg" alt="Description" />
		<p>Another paragraph with more <span class="highlight">highlighted</span> content.</p>
	</article>`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := NewParser(input)
		_, err := parser.Parse()
		if err != nil {
			b.Fatalf("parsing failed: %v", err)
		}
	}
}

// BenchmarkComparison 基准测试：与标准库比较
func BenchmarkComparison(b *testing.B) {
	input := `<root>
		<element id="test" class="example">
			<child>Some text content</child>
			<self-close attr="value" />
		</element>
	</root>`

	b.Run("MarkIt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parser := NewParser(input)
			_, err := parser.Parse()
			if err != nil {
				b.Fatalf("parsing failed: %v", err)
			}
		}
	})
}

// BenchmarkTokenTypes 基准测试：不同token类型的处理
func BenchmarkTokenTypes(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "OpenTag",
			input: `<element>`,
		},
		{
			name:  "CloseTag",
			input: `</element>`,
		},
		{
			name:  "SelfCloseTag",
			input: `<element />`,
		},
		{
			name:  "Text",
			input: `plain text content`,
		},
		{
			name:  "AttributeTag",
			input: `<element attr="value" class="test">`,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				lexer := NewLexer(tt.input)
				for {
					token := lexer.NextToken()
					if token.Type == TokenEOF {
						break
					}
				}
			}
		})
	}
}
