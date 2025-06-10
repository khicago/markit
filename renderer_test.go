package markit

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// UnknownNode 用于测试的未知节点类型
type UnknownNode struct {
	pos Position
}

func (u *UnknownNode) Type() NodeType     { return NodeType(999) } // 返回一个无效的节点类型
func (u *UnknownNode) Position() Position { return u.pos }
func (u *UnknownNode) String() string     { return "UnknownNode" }

// TestNewRenderer 测试渲染器创建
func TestNewRenderer(t *testing.T) {
	t.Run("default renderer creation", func(t *testing.T) {
		renderer := NewRenderer()
		if renderer == nil {
			t.Fatal("renderer should not be nil")
		}
		if renderer.options == nil {
			t.Fatal("renderer options should not be nil")
		}
		if renderer.options.Indent != "  " {
			t.Errorf("expected default indent '  ', got %q", renderer.options.Indent)
		}
		if renderer.options.EscapeText != true {
			t.Error("expected EscapeText to be true by default")
		}
		if renderer.options.IncludeDeclaration != true {
			t.Error("expected IncludeDeclaration to be true by default")
		}
	})

	t.Run("renderer with custom options", func(t *testing.T) {
		opts := &RenderOptions{
			Indent:      "\t",
			CompactMode: true,
		}
		renderer := NewRendererWithOptions(opts)
		if renderer.options.Indent != "\t" {
			t.Errorf("expected indent '\t', got %q", renderer.options.Indent)
		}
		if !renderer.options.CompactMode {
			t.Error("expected CompactMode to be true")
		}
	})

	t.Run("renderer with nil options", func(t *testing.T) {
		renderer := NewRendererWithOptions(nil)
		if renderer == nil {
			t.Fatal("renderer should not be nil")
		}
		// 应该使用默认选项
		if renderer.options.Indent != "  " {
			t.Errorf("expected default indent when nil options, got %q", renderer.options.Indent)
		}
	})
}

// TestBasicRendering 测试基本渲染功能
func TestBasicRendering(t *testing.T) {
	t.Run("simple element", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "root",
					Children: []Node{
						&Text{Content: "Hello World"},
					},
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		expected := `<root>
  Hello World
</root>
`
		if result != expected {
			t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
		}
	})

	t.Run("element with attributes", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "div",
					Attributes: map[string]string{
						"id":    "main",
						"class": "container",
					},
					Children: []Node{
						&Text{Content: "Content"},
					},
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		// 检查是否包含属性
		if !strings.Contains(result, `id="main"`) {
			t.Error("result should contain id attribute")
		}
		if !strings.Contains(result, `class="container"`) {
			t.Error("result should contain class attribute")
		}
	})

	t.Run("self-closing element", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName:   "img",
					SelfClose: true,
					Attributes: map[string]string{
						"src": "test.jpg",
						"alt": "test",
					},
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<img") {
			t.Error("result should contain img tag")
		}
		if !strings.Contains(result, "/>") {
			t.Error("result should contain self-closing tag syntax")
		}
	})
}

// TestRenderOptions 测试渲染选项
func TestRenderOptions(t *testing.T) {
	doc := &Document{
		Children: []Node{
			&Element{
				TagName: "root",
				Children: []Node{
					&Element{
						TagName: "child",
						Children: []Node{
							&Text{Content: "text"},
						},
					},
				},
			},
		},
	}

	t.Run("custom indent", func(t *testing.T) {
		opts := &RenderOptions{
			Indent:             "\t",
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "\t<child>") {
			t.Error("result should contain tab indentation")
		}
	})

	t.Run("exclude declaration", func(t *testing.T) {
		opts := &RenderOptions{
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if strings.Contains(result, "<?") {
			t.Error("result should not contain processing instructions")
		}
	})

	t.Run("compact mode", func(t *testing.T) {
		smallDoc := &Document{
			Children: []Node{
				&Element{
					TagName: "small",
					Children: []Node{
						&Text{Content: "text"},
					},
				},
			},
		}

		opts := &RenderOptions{
			CompactMode:        true,
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(smallDoc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if strings.Count(result, "\n") > 1 {
			t.Error("compact mode should minimize newlines")
		}
	})

	t.Run("sort attributes", func(t *testing.T) {
		elemDoc := &Document{
			Children: []Node{
				&Element{
					TagName: "div",
					Attributes: map[string]string{
						"z-attr": "last",
						"a-attr": "first",
						"m-attr": "middle",
					},
				},
			},
		}

		opts := &RenderOptions{
			SortAttributes:     true,
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(elemDoc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		// 检查属性顺序
		aIndex := strings.Index(result, "a-attr")
		mIndex := strings.Index(result, "m-attr")
		zIndex := strings.Index(result, "z-attr")

		if aIndex == -1 || mIndex == -1 || zIndex == -1 {
			t.Fatal("all attributes should be present")
		}
		if !(aIndex < mIndex && mIndex < zIndex) {
			t.Error("attributes should be sorted alphabetically")
		}
	})
}

// TestEmptyElementStyles 测试空元素样式
func TestEmptyElementStyles(t *testing.T) {
	doc := &Document{
		Children: []Node{
			&Element{
				TagName:   "empty",
				SelfClose: true,
			},
		},
	}

	t.Run("self-closing style", func(t *testing.T) {
		opts := &RenderOptions{
			EmptyElementStyle:  SelfClosingStyle,
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<empty />") {
			t.Error("self-closing style should use self-closing syntax")
		}
	})

	t.Run("paired tag style", func(t *testing.T) {
		opts := &RenderOptions{
			EmptyElementStyle:  PairedTagStyle,
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<empty></empty>") {
			t.Error("paired tag style should use opening and closing tags")
		}
	})

	t.Run("void element style", func(t *testing.T) {
		htmlConfig := HTMLConfig()
		voidDoc := &Document{
			Children: []Node{
				&Element{
					TagName:   "br",
					SelfClose: true,
				},
			},
		}

		opts := &RenderOptions{
			EmptyElementStyle:  VoidElementStyle,
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithConfig(htmlConfig, opts)
		result, err := renderer.RenderToString(voidDoc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<br>") || strings.Contains(result, "</br>") {
			t.Error("void element should not have closing tag")
		}
	})
}

// TestRendererAllNodeTypes 测试所有节点类型的渲染（改名避免冲突）
func TestRendererAllNodeTypes(t *testing.T) {
	doc := &Document{
		Children: []Node{
			&ProcessingInstruction{
				Target:  "xml-stylesheet",
				Content: `type="text/xsl" href="style.xsl"`,
			},
			&Doctype{
				Content: "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\"",
			},
			&Comment{
				Content: " This is a comment ",
			},
			&Element{
				TagName: "root",
				Children: []Node{
					&Text{Content: "Hello World"},
					&CDATA{Content: "function() { return 'test'; }"},
					&Element{
						TagName:   "img",
						SelfClose: true,
						Attributes: map[string]string{
							"src": "image.png",
						},
					},
				},
			},
		},
	}

	renderer := NewRenderer()
	result, err := renderer.RenderToString(doc)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}

	// 验证每种节点类型都被正确渲染
	expectedContents := []string{
		"<?xml-stylesheet type=\"text/xsl\" href=\"style.xsl\"?>",
		"<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\">",
		"<!-- This is a comment -->",
		"<root>",
		"Hello World",
		"<![CDATA[function() { return 'test'; }]]>",
		"<img src=\"image.png\" />",
		"</root>",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(result, expected) {
			t.Errorf("expected result to contain %q, but it didn't.\nResult:\n%s", expected, result)
		}
	}
}

// TestStreamingSupport 测试流式支持
func TestStreamingSupport(t *testing.T) {
	doc := &Document{
		Children: []Node{
			&Element{
				TagName: "root",
				Children: []Node{
					&Text{Content: "Hello World"},
				},
			},
		},
	}

	t.Run("render to writer", func(t *testing.T) {
		var buf bytes.Buffer
		renderer := NewRenderer()
		err := renderer.RenderToWriter(doc, &buf)
		if err != nil {
			t.Fatalf("render to writer error: %v", err)
		}

		result := buf.String()
		if !strings.Contains(result, "Hello World") {
			t.Error("result should contain expected content")
		}
	})

	t.Run("render element to writer", func(t *testing.T) {
		elem := &Element{
			TagName: "test",
			Children: []Node{
				&Text{Content: "content"},
			},
		}

		var buf bytes.Buffer
		renderer := NewRenderer()
		err := renderer.RenderElementToWriter(elem, &buf)
		if err != nil {
			t.Fatalf("render element to writer error: %v", err)
		}

		result := buf.String()
		if !strings.Contains(result, "<test>") {
			t.Error("result should contain test element")
		}
	})
}

// TestElementLevelRendering 测试元素级渲染
func TestElementLevelRendering(t *testing.T) {
	elem := &Element{
		TagName: "div",
		Attributes: map[string]string{
			"class": "test",
		},
		Children: []Node{
			&Text{Content: "Element content"},
		},
	}

	renderer := NewRendererWithOptions(&RenderOptions{
		IncludeDeclaration: false,
	})
	result, err := renderer.RenderElement(elem)
	if err != nil {
		t.Fatalf("render element error: %v", err)
	}

	if !strings.Contains(result, `<div class="test">`) {
		t.Error("result should contain div element with class")
	}
	if !strings.Contains(result, "Element content") {
		t.Error("result should contain element content")
	}
	if !strings.Contains(result, "</div>") {
		t.Error("result should contain closing tag")
	}
}

// TestValidation 测试验证功能
func TestValidation(t *testing.T) {
	t.Run("valid document", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "valid-element",
					Attributes: map[string]string{
						"valid-attr": "value",
					},
					Children: []Node{
						&Text{Content: "Valid text content"},
					},
				},
			},
		}

		validation := &ValidationOptions{
			CheckWellFormed: true,
			CheckEncoding:   true,
		}

		renderer := NewRenderer()
		result, err := renderer.RenderWithValidation(doc, validation)
		if err != nil {
			t.Fatalf("validation should pass for valid document: %v", err)
		}

		if !strings.Contains(result, "valid-element") {
			t.Error("result should contain valid element")
		}
	})

	t.Run("invalid tag name", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "123invalid", // 无效标签名
					Children: []Node{
						&Text{Content: "content"},
					},
				},
			},
		}

		validation := &ValidationOptions{
			CheckWellFormed: true,
		}

		renderer := NewRenderer()
		_, err := renderer.RenderWithValidation(doc, validation)
		if err == nil {
			t.Error("validation should fail for invalid tag name")
		}

		validationErr, ok := err.(*ValidationError)
		if !ok {
			t.Error("error should be ValidationError type")
		} else if !strings.Contains(validationErr.Message, "invalid tag name") {
			t.Error("error message should mention invalid tag name")
		}
	})

	t.Run("invalid attribute name", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "div",
					Attributes: map[string]string{
						"123invalid": "value", // 无效属性名
					},
				},
			},
		}

		validation := &ValidationOptions{
			CheckWellFormed: true,
		}

		renderer := NewRenderer()
		_, err := renderer.RenderWithValidation(doc, validation)
		if err == nil {
			t.Error("validation should fail for invalid attribute name")
		}
	})
}

// TestRendererErrorHandling 测试错误处理（改名避免冲突）
func TestRendererErrorHandling(t *testing.T) {
	renderer := NewRenderer()

	t.Run("nil document", func(t *testing.T) {
		_, err := renderer.RenderToString(nil)
		if err == nil {
			t.Error("should return error for nil document")
		}
	})

	t.Run("nil element", func(t *testing.T) {
		_, err := renderer.RenderElement(nil)
		if err == nil {
			t.Error("should return error for nil element")
		}
	})

	t.Run("nil writer", func(t *testing.T) {
		doc := &Document{}
		err := renderer.RenderToWriter(doc, nil)
		if err == nil {
			t.Error("should return error for nil writer")
		}
	})
}

// TestTextEscaping 测试文本转义
func TestTextEscaping(t *testing.T) {
	doc := &Document{
		Children: []Node{
			&Element{
				TagName: "root",
				Attributes: map[string]string{
					"attr": `<value & "quoted">`,
				},
				Children: []Node{
					&Text{Content: `<script>alert("XSS")</script>`},
				},
			},
		},
	}

	t.Run("with escaping", func(t *testing.T) {
		opts := &RenderOptions{
			EscapeText:         true,
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		// 检查文本内容是否被转义
		if strings.Contains(result, `<script>`) {
			t.Error("script tags should be escaped")
		}
		if !strings.Contains(result, `&lt;script&gt;`) {
			t.Error("text should contain escaped content")
		}
	})

	t.Run("without escaping", func(t *testing.T) {
		opts := &RenderOptions{
			EscapeText:         false,
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		// 检查文本内容是否未被转义
		if !strings.Contains(result, `<script>`) {
			t.Error("text should contain unescaped content")
		}
	})
}

// TestBackwardCompatibility 测试向后兼容性
func TestBackwardCompatibility(t *testing.T) {
	doc := &Document{
		Children: []Node{
			&Element{
				TagName: "root",
				Children: []Node{
					&Text{Content: "test"},
				},
			},
		},
	}

	renderer := NewRenderer()

	// 测试原有的 Render 方法
	result := renderer.Render(doc)
	if result == "" {
		t.Error("Render method should return non-empty result")
	}
	if !strings.Contains(result, "test") {
		t.Error("Render method should contain expected content")
	}
}

// TestComplexDocument 测试复杂文档渲染
func TestComplexDocument(t *testing.T) {
	// 创建一个复杂的文档结构
	doc := &Document{
		Children: []Node{
			&ProcessingInstruction{
				Target:  "xml",
				Content: `version="1.0" encoding="UTF-8"`,
			},
			&Comment{Content: " Document structure test "},
			&Element{
				TagName: "document",
				Attributes: map[string]string{
					"version": "1.0",
					"lang":    "en",
				},
				Children: []Node{
					&Element{
						TagName: "header",
						Children: []Node{
							&Element{
								TagName: "title",
								Children: []Node{
									&Text{Content: "Test Document"},
								},
							},
							&Element{
								TagName: "meta",
								Attributes: map[string]string{
									"charset": "UTF-8",
								},
								SelfClose: true,
							},
						},
					},
					&Element{
						TagName: "body",
						Children: []Node{
							&Element{
								TagName: "h1",
								Attributes: map[string]string{
									"id": "main-title",
								},
								Children: []Node{
									&Text{Content: "Main Title"},
								},
							},
							&Element{
								TagName: "p",
								Children: []Node{
									&Text{Content: "This is a paragraph with "},
									&Element{
										TagName: "strong",
										Children: []Node{
											&Text{Content: "bold"},
										},
									},
									&Text{Content: " text."},
								},
							},
							&Element{
								TagName: "ul",
								Attributes: map[string]string{
									"class": "list",
								},
								Children: []Node{
									&Element{
										TagName: "li",
										Children: []Node{
											&Text{Content: "Item 1"},
										},
									},
									&Element{
										TagName: "li",
										Children: []Node{
											&Text{Content: "Item 2"},
										},
									},
								},
							},
							&CDATA{Content: "Some CDATA content"},
						},
					},
				},
			},
		},
	}

	renderer := NewRenderer()
	result, err := renderer.RenderToString(doc)
	if err != nil {
		t.Fatalf("complex document render error: %v", err)
	}

	// 验证结构完整性 - 不依赖属性顺序
	expectedElements := []string{
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>",
		"<!-- Document structure test -->",
		"<header>",
		"<title>",
		"Test Document",
		"</title>",
		"<meta charset=\"UTF-8\" />",
		"</header>",
		"<body>",
		"<h1 id=\"main-title\">",
		"Main Title",
		"</h1>",
		"<p>",
		"This is a paragraph with ",
		"<strong>",
		"bold",
		"</strong>",
		" text.",
		"</p>",
		"<ul class=\"list\">",
		"<li>",
		"Item 1",
		"</li>",
		"<li>",
		"Item 2",
		"</li>",
		"</ul>",
		"<![CDATA[Some CDATA content]]>",
		"</body>",
		"</document>",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("complex document should contain: %q\nResult:\n%s", expected, result)
		}
	}

	// 验证document标签的属性（不依赖顺序）
	if !strings.Contains(result, `version="1.0"`) {
		t.Error("document should contain version attribute")
	}
	if !strings.Contains(result, `lang="en"`) {
		t.Error("document should contain lang attribute")
	}
	if !strings.Contains(result, "<document") {
		t.Error("document should contain document tag")
	}

	// 验证缩进结构
	lines := strings.Split(result, "\n")
	foundTitle := false
	for _, line := range lines {
		if strings.Contains(line, "<title>") {
			foundTitle = true
			// title 元素应该有正确的缩进（4个空格）
			if !strings.HasPrefix(line, "    ") {
				t.Error("title element should have proper indentation")
			}
			break
		}
	}
	if !foundTitle {
		t.Error("should find title element in output")
	}
}

// TestDeclarationControl 测试声明控制
func TestDeclarationControl(t *testing.T) {
	doc := &Document{
		Children: []Node{
			&ProcessingInstruction{
				Target:  "xml",
				Content: `version="1.0"`,
			},
			&Doctype{
				Content: "html",
			},
			&Element{
				TagName: "root",
				Children: []Node{
					&Text{Content: "content"},
				},
			},
		},
	}

	t.Run("include declarations", func(t *testing.T) {
		opts := &RenderOptions{
			IncludeDeclaration: true,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<?xml") {
			t.Error("should include processing instruction")
		}
		if !strings.Contains(result, "<!DOCTYPE") {
			t.Error("should include doctype")
		}
	})

	t.Run("exclude declarations", func(t *testing.T) {
		opts := &RenderOptions{
			IncludeDeclaration: false,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if strings.Contains(result, "<?xml") {
			t.Error("should not include processing instruction")
		}
		if strings.Contains(result, "<!DOCTYPE") {
			t.Error("should not include doctype")
		}
		if !strings.Contains(result, "<root>") {
			t.Error("should still include elements")
		}
	})
}

// TestUncoveredCodePaths 测试未覆盖的代码路径
func TestUncoveredCodePaths(t *testing.T) {
	t.Run("ValidationError methods", func(t *testing.T) {
		err := &ValidationError{
			Message:  "test error",
			Position: Position{Line: 1, Column: 5},
			NodeType: NodeTypeElement,
		}
		expected := "validation error at line 1, column 5: test error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Renderer setter methods", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试 SetOptions
		newOpts := &RenderOptions{
			Indent:      "\t",
			CompactMode: true,
		}
		renderer.SetOptions(newOpts)
		if renderer.options.Indent != "\t" {
			t.Error("SetOptions should update indent")
		}

		// 测试 SetOptions with nil
		renderer.SetOptions(nil)
		if renderer.options.Indent != "\t" {
			t.Error("SetOptions with nil should not change options")
		}

		// 测试 SetConfig
		config := DefaultConfig()
		renderer.SetConfig(config)
		if renderer.config != config {
			t.Error("SetConfig should update config")
		}

		// 测试 SetValidation
		validation := &ValidationOptions{
			CheckWellFormed: true,
		}
		renderer.SetValidation(validation)
		if renderer.validation != validation {
			t.Error("SetValidation should update validation")
		}
	})

	t.Run("renderDocument method", func(t *testing.T) {
		renderer := NewRenderer()
		doc := &Document{
			Children: []Node{
				&Text{Content: "test"},
			},
		}

		var buf strings.Builder
		err := renderer.renderDocument(doc, &buf, 0)
		if err != nil {
			t.Fatalf("renderDocument error: %v", err)
		}

		result := buf.String()
		if result != "test" {
			t.Errorf("expected 'test', got %q", result)
		}
	})

	t.Run("isSmallElement method", func(t *testing.T) {
		renderer := NewRenderer()

		// 空元素
		emptyElem := &Element{TagName: "empty"}
		if !renderer.isSmallElement(emptyElem) {
			t.Error("empty element should be small")
		}

		// 短文本元素
		shortTextElem := &Element{
			TagName: "short",
			Children: []Node{
				&Text{Content: "short"},
			},
		}
		if !renderer.isSmallElement(shortTextElem) {
			t.Error("short text element should be small")
		}

		// 长文本元素
		longTextElem := &Element{
			TagName: "long",
			Children: []Node{
				&Text{Content: strings.Repeat("x", 100)},
			},
		}
		if renderer.isSmallElement(longTextElem) {
			t.Error("long text element should not be small")
		}

		// 多子元素
		multiChildElem := &Element{
			TagName: "multi",
			Children: []Node{
				&Text{Content: "a"},
				&Text{Content: "b"},
			},
		}
		if renderer.isSmallElement(multiChildElem) {
			t.Error("multi-child element should not be small")
		}
	})

	t.Run("isOnlyTextChildren method", func(t *testing.T) {
		renderer := NewRenderer()

		// 只有文本子元素
		textOnlyElem := &Element{
			TagName: "text-only",
			Children: []Node{
				&Text{Content: "text1"},
				&Text{Content: "text2"},
			},
		}
		if !renderer.isOnlyTextChildren(textOnlyElem) {
			t.Error("element with only text children should return true")
		}

		// 混合子元素
		mixedElem := &Element{
			TagName: "mixed",
			Children: []Node{
				&Text{Content: "text"},
				&Element{TagName: "child"},
			},
		}
		if renderer.isOnlyTextChildren(mixedElem) {
			t.Error("element with mixed children should return false")
		}

		// 空元素
		emptyElem := &Element{TagName: "empty"}
		if !renderer.isOnlyTextChildren(emptyElem) {
			t.Error("empty element should return true")
		}
	})
}

// TestTextRenderingEdgeCases 测试文本渲染的边缘情况
func TestTextRenderingEdgeCases(t *testing.T) {
	t.Run("multiline text with indentation", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "root",
					Children: []Node{
						&Text{Content: "line1\nline2\nline3"},
					},
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		// 检查多行文本是否正确处理
		if !strings.Contains(result, "line1\n") {
			t.Error("multiline text should preserve line breaks")
		}
	})

	t.Run("text with tabs and spaces", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "root",
					Children: []Node{
						&Text{Content: "text\twith\ttabs\r\nand\rcarriage"},
					},
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "\t") {
			t.Error("text should preserve tabs")
		}
	})
}

// TestSpecialNodeRendering 测试特殊节点的渲染
func TestSpecialNodeRendering(t *testing.T) {
	t.Run("comment with compact mode", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Comment{Content: " test comment "},
			},
		}

		opts := &RenderOptions{
			CompactMode: true,
		}
		renderer := NewRendererWithOptions(opts)
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<!-- test comment -->") {
			t.Error("comment should be rendered correctly")
		}
		if strings.Contains(result, "\n") {
			t.Error("compact mode should not add newlines")
		}
	})

	t.Run("processing instruction with empty content", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&ProcessingInstruction{
					Target:  "xml-target",
					Content: "",
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		expected := "<?xml-target?>"
		if !strings.Contains(result, expected) {
			t.Error("processing instruction with empty content should render correctly")
		}
	})
}

// TestValidationEdgeCases 测试验证的边缘情况
func TestValidationEdgeCases(t *testing.T) {
	t.Run("validation with invalid UTF-8", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Text{Content: string([]byte{0xff, 0xfe, 0xfd})}, // invalid UTF-8
			},
		}

		validation := &ValidationOptions{
			CheckEncoding: true,
		}

		renderer := NewRenderer()
		_, err := renderer.RenderWithValidation(doc, validation)
		if err == nil {
			t.Error("validation should fail for invalid UTF-8")
		}

		validationErr, ok := err.(*ValidationError)
		if !ok {
			t.Error("error should be ValidationError type")
		} else if validationErr.NodeType != NodeTypeText {
			t.Error("error should be for text node")
		}
	})

	t.Run("validation with invalid tag names", func(t *testing.T) {
		tests := []struct {
			name    string
			tagName string
			valid   bool
		}{
			{"empty tag name", "", false},
			{"valid tag name", "valid-tag_name.ext", true},
			{"starts with number", "123invalid", false},
			{"starts with hyphen", "-invalid", false},
			{"contains invalid chars", "invalid@tag", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				doc := &Document{
					Children: []Node{
						&Element{TagName: tt.tagName},
					},
				}

				validation := &ValidationOptions{
					CheckWellFormed: true,
				}

				renderer := NewRenderer()
				_, err := renderer.RenderWithValidation(doc, validation)

				if tt.valid && err != nil {
					t.Errorf("valid tag name %q should not cause error", tt.tagName)
				}
				if !tt.valid && err == nil {
					t.Errorf("invalid tag name %q should cause error", tt.tagName)
				}
			})
		}
	})
}

// errorWriter 用于测试 Writer 错误的辅助结构
type errorWriter struct{}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("write error")
}

// TestWriterErrorHandling 测试 Writer 错误处理
func TestWriterErrorHandling(t *testing.T) {
	t.Run("writeIndent with error", func(t *testing.T) {
		renderer := NewRenderer()

		// 创建一个会产生错误的 Writer
		errorWriter := &errorWriter{}

		err := renderer.writeIndent(errorWriter, 2)
		if err == nil {
			t.Error("writeIndent should return error when writer fails")
		}
	})
}

// TestRenderingErrorConditions 测试渲染错误条件
func TestRenderingErrorConditions(t *testing.T) {
	t.Run("RenderToString with validation error", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{TagName: ""}, // 空标签名
			},
		}

		renderer := NewRenderer()
		validation := &ValidationOptions{
			CheckWellFormed: true,
		}
		renderer.SetValidation(validation)

		_, err := renderer.RenderToString(doc)
		if err == nil {
			t.Error("RenderToString should fail with validation error")
		}
	})

	t.Run("RenderToWriter with write error", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{TagName: "test"},
			},
		}

		renderer := NewRenderer()
		errorWriter := &errorWriter{}

		err := renderer.RenderToWriter(doc, errorWriter)
		if err == nil {
			t.Error("RenderToWriter should fail with write error")
		}
	})

	t.Run("RenderElement with nil element", func(t *testing.T) {
		renderer := NewRenderer()

		_, err := renderer.RenderElement(nil)
		if err == nil {
			t.Error("RenderElement should fail with nil element")
		}
	})

	t.Run("RenderElementToWriter with error writer", func(t *testing.T) {
		elem := &Element{TagName: "test"}
		renderer := NewRenderer()
		errorWriter := &errorWriter{}

		err := renderer.RenderElementToWriter(elem, errorWriter)
		if err == nil {
			t.Error("RenderElementToWriter should fail with write error")
		}
	})

	t.Run("RenderWithValidation validation errors", func(t *testing.T) {
		tests := []struct {
			name string
			doc  *Document
		}{
			{
				"empty tag name",
				&Document{Children: []Node{&Element{TagName: ""}}},
			},
			{
				"invalid character in tag name",
				&Document{Children: []Node{&Element{TagName: "inv@lid"}}},
			},
			{
				"tag name starts with number",
				&Document{Children: []Node{&Element{TagName: "123tag"}}},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				renderer := NewRenderer()
				validation := &ValidationOptions{
					CheckWellFormed: true,
				}

				_, err := renderer.RenderWithValidation(tt.doc, validation)
				if err == nil {
					t.Errorf("RenderWithValidation should fail for %s", tt.name)
				}
			})
		}
	})
}

// TestSpecialTextCases 测试特殊文本情况
func TestSpecialTextCases(t *testing.T) {
	t.Run("text with control characters", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "root",
					Children: []Node{
						&Text{Content: "text\x00with\x01control"},
					},
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "text") {
			t.Error("should preserve normal text parts")
		}
	})

	t.Run("empty text nodes", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "root",
					Children: []Node{
						&Text{Content: ""},
						&Text{Content: "  "},
						&Text{Content: "\n\t"},
					},
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<root>") {
			t.Error("should render root element")
		}
	})

	t.Run("text with line breaks and indentation", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "root",
					Children: []Node{
						&Text{Content: "\n  indented\n  text\n"},
					},
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		// 测试多行文本的处理
		lines := strings.Split(result, "\n")
		if len(lines) < 2 {
			t.Error("should preserve line structure in text")
		}
	})
}

// TestSpecialNodesEdgeCases 测试特殊节点的边缘情况
func TestSpecialNodesEdgeCases(t *testing.T) {
	t.Run("comment with special characters", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Comment{Content: " comment with -- dashes and \n newlines "},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<!--") && !strings.Contains(result, "-->") {
			t.Error("should render comment markers")
		}
	})

	t.Run("processing instruction with complex content", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&ProcessingInstruction{
					Target:  "xml-stylesheet",
					Content: "type=\"text/xsl\" href=\"style.xsl\"",
				},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<?xml-stylesheet") {
			t.Error("should render processing instruction target")
		}
		if !strings.Contains(result, "type=\"text/xsl\"") {
			t.Error("should render processing instruction content")
		}
	})

	t.Run("doctype with public and system id", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Doctype{Content: "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\""},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<!DOCTYPE") {
			t.Error("should render doctype declaration")
		}
		if !strings.Contains(result, "PUBLIC") {
			t.Error("should render public identifier")
		}
	})

	t.Run("CDATA with complex content", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&CDATA{Content: "function() { return '<test>'; /* comment */ }"},
			},
		}

		renderer := NewRenderer()
		result, err := renderer.RenderToString(doc)
		if err != nil {
			t.Fatalf("render error: %v", err)
		}

		if !strings.Contains(result, "<![CDATA[") {
			t.Error("should render CDATA start marker")
		}
		if !strings.Contains(result, "]]>") {
			t.Error("should render CDATA end marker")
		}
		if !strings.Contains(result, "function()") {
			t.Error("should preserve CDATA content")
		}
	})
}

// TestValidationComprehensive 测试全面的验证功能
func TestValidationComprehensive(t *testing.T) {
	t.Run("text validation with invalid UTF-8", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Text{Content: string([]byte{0xff, 0xfe, 0xfd})},
			},
		}

		renderer := NewRenderer()
		validation := &ValidationOptions{
			CheckEncoding: true,
		}

		_, err := renderer.RenderWithValidation(doc, validation)
		if err == nil {
			t.Error("should fail validation for invalid UTF-8")
		}

		validationErr, ok := err.(*ValidationError)
		if !ok {
			t.Error("error should be ValidationError")
		} else {
			if validationErr.NodeType != NodeTypeText {
				t.Error("error should be for text node")
			}
		}
	})

	t.Run("element validation with invalid tag names", func(t *testing.T) {
		testCases := []struct {
			tagName string
			valid   bool
		}{
			{"", false},
			{"123invalid", false},
			{"-invalid", false},
			{"in@valid", false},
			{"valid", true},
			{"valid-name", true},
			{"valid_name", true},
			{"valid.name", true},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("tag_%s", tc.tagName), func(t *testing.T) {
				doc := &Document{
					Children: []Node{
						&Element{TagName: tc.tagName},
					},
				}

				renderer := NewRenderer()
				validation := &ValidationOptions{
					CheckWellFormed: true,
				}

				_, err := renderer.RenderWithValidation(doc, validation)

				if tc.valid && err != nil {
					t.Errorf("valid tag name %q should not cause error: %v", tc.tagName, err)
				}
				if !tc.valid && err == nil {
					t.Errorf("invalid tag name %q should cause error", tc.tagName)
				}
			})
		}
	})

	t.Run("nested validation errors", func(t *testing.T) {
		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "valid",
					Children: []Node{
						&Element{TagName: "123invalid"},
					},
				},
			},
		}

		renderer := NewRenderer()
		validation := &ValidationOptions{
			CheckWellFormed: true,
		}

		_, err := renderer.RenderWithValidation(doc, validation)
		if err == nil {
			t.Error("should fail validation for nested invalid tag name")
		}
	})
}

// TestWriteIndentEdgeCases 测试缩进写入的边缘情况
func TestWriteIndentEdgeCases(t *testing.T) {
	t.Run("writeIndent with zero level", func(t *testing.T) {
		renderer := NewRenderer()
		var buf strings.Builder

		err := renderer.writeIndent(&buf, 0)
		if err != nil {
			t.Errorf("writeIndent should not error with zero level: %v", err)
		}

		result := buf.String()
		if result != "" {
			t.Error("zero indent level should write nothing")
		}
	})

	t.Run("writeIndent with custom indent string", func(t *testing.T) {
		opts := &RenderOptions{
			Indent: "\t\t",
		}
		renderer := NewRendererWithOptions(opts)
		var buf strings.Builder

		err := renderer.writeIndent(&buf, 2)
		if err != nil {
			t.Errorf("writeIndent should not error: %v", err)
		}

		result := buf.String()
		expected := "\t\t\t\t" // 2 levels * 2 tabs each
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})
}

// TestReadTextEdgeCases 测试读取文本的边缘情况
func TestReadTextEdgeCases(t *testing.T) {
	t.Run("lexer readText with complex whitespace", func(t *testing.T) {
		input := "   \n\t\r   <tag>"
		lexer := NewLexer(input)

		// 读取直到遇到标签
		for {
			token := lexer.NextToken()
			if token.Type == TokenOpenTag {
				break
			}
			if token.Type == TokenEOF {
				t.Error("should find OPEN_TAG before EOF")
				break
			}
		}
	})
}

// TestRendererComprehensiveCoverage 全面覆盖测试
func TestRendererComprehensiveCoverage(t *testing.T) {
	t.Run("ValidationError Error method", func(t *testing.T) {
		err := &ValidationError{
			Message:  "test validation error",
			Position: Position{Line: 10, Column: 5},
			NodeType: NodeTypeElement,
		}
		expected := "validation error at line 10, column 5: test validation error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("SetOptions comprehensive", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试设置新选项
		newOpts := &RenderOptions{
			Indent:      "\t",
			CompactMode: true,
		}
		renderer.SetOptions(newOpts)
		if renderer.options.Indent != "\t" {
			t.Errorf("expected indent '\t', got %q", renderer.options.Indent)
		}
		if !renderer.options.CompactMode {
			t.Error("expected CompactMode to be true")
		}

		// 测试设置 nil 选项（应该不改变当前选项）
		renderer.SetOptions(nil)
		if renderer.options.Indent != "\t" {
			t.Error("SetOptions with nil should not change current options")
		}
	})

	t.Run("SetConfig method", func(t *testing.T) {
		renderer := NewRenderer()
		config := DefaultConfig()
		config.TrimWhitespace = false

		renderer.SetConfig(config)
		if renderer.config != config {
			t.Error("SetConfig should set the config")
		}
		if renderer.config.TrimWhitespace {
			t.Error("config should have TrimWhitespace=false")
		}
	})

	t.Run("SetValidation method", func(t *testing.T) {
		renderer := NewRenderer()
		validation := &ValidationOptions{
			CheckWellFormed: true,
			CheckEncoding:   true,
		}

		renderer.SetValidation(validation)
		if renderer.validation != validation {
			t.Error("SetValidation should set the validation")
		}
		if !renderer.validation.CheckWellFormed {
			t.Error("validation should have CheckWellFormed=true")
		}
	})
}

// TestUtilityFunctions 测试工具函数
func TestUtilityFunctions(t *testing.T) {
	t.Run("isValidTagName", func(t *testing.T) {
		tests := []struct {
			name     string
			tagName  string
			expected bool
		}{
			{"empty string", "", false},
			{"valid simple", "div", true},
			{"valid with dash", "my-tag", true},
			{"valid with underscore", "my_tag", true},
			{"valid with dot", "my.tag", true},
			{"valid with numbers", "tag123", true},
			{"invalid starts with number", "123tag", false},
			{"invalid starts with dash", "-tag", false},
			{"invalid starts with dot", ".tag", false},
			{"invalid with space", "my tag", false},
			{"invalid with @", "my@tag", false},
			{"valid uppercase", "MyTag", true},
			{"valid mixed case", "myTag", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := isValidTagName(tt.tagName)
				if result != tt.expected {
					t.Errorf("isValidTagName(%q) = %v, expected %v", tt.tagName, result, tt.expected)
				}
			})
		}
	})

	t.Run("isValidAttributeName", func(t *testing.T) {
		tests := []struct {
			name     string
			attrName string
			expected bool
		}{
			{"empty string", "", false},
			{"valid simple", "class", true},
			{"valid with dash", "data-value", true},
			{"valid with underscore", "ng_model", true},
			{"invalid starts with number", "123attr", false},
			{"invalid with space", "my attr", false},
			{"valid uppercase", "ID", true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := isValidAttributeName(tt.attrName)
				if result != tt.expected {
					t.Errorf("isValidAttributeName(%q) = %v, expected %v", tt.attrName, result, tt.expected)
				}
			})
		}
	})

	t.Run("escapeText", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"no escaping needed", "hello world", "hello world"},
			{"escape ampersand", "Tom & Jerry", "Tom &amp; Jerry"},
			{"escape less than", "3 < 5", "3 &lt; 5"},
			{"escape greater than", "5 > 3", "5 &gt; 3"},
			{"escape double quote", `say "hello"`, "say &quot;hello&quot;"},
			{"escape single quote", "don't", "don&#39;t"},
			{"escape all", `<script>alert("XSS & 'attack'");</script>`, "&lt;script&gt;alert(&quot;XSS &amp; &#39;attack&#39;&quot;);&lt;/script&gt;"},
			{"empty string", "", ""},
			{"multiple escapes", "<<>>&&\"\"''", "&lt;&lt;&gt;&gt;&amp;&amp;&quot;&quot;&#39;&#39;"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := escapeText(tt.input)
				if result != tt.expected {
					t.Errorf("escapeText(%q) = %q, expected %q", tt.input, result, tt.expected)
				}
			})
		}
	})
}

// TestValidateText 测试文本验证
func TestValidateText(t *testing.T) {
	t.Run("valid UTF-8 text", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(&ValidationOptions{CheckEncoding: true})

		text := &Text{Content: "Hello 世界 🌍"}
		err := renderer.validateText(text)
		if err != nil {
			t.Errorf("valid UTF-8 text should not cause error: %v", err)
		}
	})

	t.Run("invalid UTF-8 text", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(&ValidationOptions{CheckEncoding: true})

		// 创建无效的 UTF-8 字符串
		invalidUTF8 := string([]byte{0xff, 0xfe, 0xfd})
		text := &Text{Content: invalidUTF8}

		err := renderer.validateText(text)
		if err == nil {
			t.Error("invalid UTF-8 text should cause validation error")
		}

		validationErr, ok := err.(*ValidationError)
		if !ok {
			t.Error("error should be ValidationError type")
		} else {
			if validationErr.NodeType != NodeTypeText {
				t.Error("error should be for text node")
			}
			if !strings.Contains(validationErr.Message, "UTF-8") {
				t.Error("error message should mention UTF-8")
			}
		}
	})

	t.Run("validation disabled", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(&ValidationOptions{CheckEncoding: false})

		invalidUTF8 := string([]byte{0xff, 0xfe, 0xfd})
		text := &Text{Content: invalidUTF8}

		err := renderer.validateText(text)
		if err != nil {
			t.Error("should not validate when CheckEncoding is false")
		}
	})

	t.Run("no validation options", func(t *testing.T) {
		renderer := NewRenderer()

		text := &Text{Content: "any content"}
		err := renderer.validateText(text)
		if err != nil {
			t.Error("should not validate when validation is nil")
		}
	})
}

// TestAdvancedRenderingScenarios 测试高级渲染场景
func TestAdvancedRenderingScenarios(t *testing.T) {
	t.Run("renderDocument direct call", func(t *testing.T) {
		renderer := NewRenderer()
		doc := &Document{
			Children: []Node{
				&Text{Content: "document content"},
				&Element{TagName: "div", Children: []Node{&Text{Content: "element content"}}},
			},
		}

		var buf strings.Builder
		err := renderer.renderDocument(doc, &buf, 0)
		if err != nil {
			t.Fatalf("renderDocument error: %v", err)
		}

		result := buf.String()
		if !strings.Contains(result, "document content") {
			t.Error("should contain document content")
		}
		if !strings.Contains(result, "<div>") {
			t.Error("should contain div element")
		}
	})

	t.Run("renderAttributes edge cases", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试空属性
		elem1 := &Element{TagName: "div"}
		var buf1 strings.Builder
		err := renderer.renderAttributes(elem1, &buf1)
		if err != nil {
			t.Errorf("renderAttributes with no attributes should not error: %v", err)
		}
		if buf1.String() != "" {
			t.Error("no attributes should render nothing")
		}

		// 测试 nil 属性
		elem2 := &Element{TagName: "div", Attributes: nil}
		var buf2 strings.Builder
		err = renderer.renderAttributes(elem2, &buf2)
		if err != nil {
			t.Errorf("renderAttributes with nil attributes should not error: %v", err)
		}

		// 测试空值属性（布尔属性）
		elem3 := &Element{
			TagName: "input",
			Attributes: map[string]string{
				"checked": "",
				"type":    "checkbox",
			},
		}
		var buf3 strings.Builder
		err = renderer.renderAttributes(elem3, &buf3)
		if err != nil {
			t.Errorf("renderAttributes should not error: %v", err)
		}
		result := buf3.String()
		if !strings.Contains(result, "checked") {
			t.Error("should contain boolean attribute")
		}
		if !strings.Contains(result, `type="checkbox"`) {
			t.Error("should contain regular attribute")
		}

		// 测试属性排序
		renderer.SetOptions(&RenderOptions{SortAttributes: true, EscapeText: false})
		elem4 := &Element{
			TagName: "div",
			Attributes: map[string]string{
				"z-attr": "last",
				"a-attr": "first",
				"m-attr": "middle",
			},
		}
		var buf4 strings.Builder
		err = renderer.renderAttributes(elem4, &buf4)
		if err != nil {
			t.Errorf("renderAttributes should not error: %v", err)
		}
		result4 := buf4.String()
		aIndex := strings.Index(result4, "a-attr")
		mIndex := strings.Index(result4, "m-attr")
		zIndex := strings.Index(result4, "z-attr")
		if !(aIndex < mIndex && mIndex < zIndex) {
			t.Error("attributes should be sorted alphabetically")
		}

		// 测试属性值转义
		renderer.SetOptions(&RenderOptions{EscapeText: true})
		elem5 := &Element{
			TagName: "div",
			Attributes: map[string]string{
				"data": `<value & "quoted">`,
			},
		}
		var buf5 strings.Builder
		err = renderer.renderAttributes(elem5, &buf5)
		if err != nil {
			t.Errorf("renderAttributes should not error: %v", err)
		}
		result5 := buf5.String()
		if !strings.Contains(result5, "&lt;value") {
			t.Error("attribute values should be escaped")
		}
	})

	t.Run("renderText edge cases", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试空文本
		text1 := &Text{Content: ""}
		var buf1 strings.Builder
		err := renderer.renderText(text1, &buf1, 0)
		if err != nil {
			t.Errorf("renderText with empty content should not error: %v", err)
		}

		// 测试只有空白字符的文本
		text2 := &Text{Content: "   \t\n   "}
		var buf2 strings.Builder
		err = renderer.renderText(text2, &buf2, 1)
		if err != nil {
			t.Errorf("renderText with whitespace should not error: %v", err)
		}

		// 测试不转义模式
		renderer.SetOptions(&RenderOptions{EscapeText: false, CompactMode: true})
		text3 := &Text{Content: "<script>alert('test')</script>"}
		var buf3 strings.Builder
		err = renderer.renderText(text3, &buf3, 0)
		if err != nil {
			t.Errorf("renderText should not error: %v", err)
		}
		result3 := buf3.String()
		if !strings.Contains(result3, "<script>") {
			t.Error("unescaped text should contain original content")
		}

		// 测试紧凑模式下的多行文本
		renderer.SetOptions(&RenderOptions{CompactMode: true})
		text4 := &Text{Content: "line1\nline2\nline3"}
		var buf4 strings.Builder
		err = renderer.renderText(text4, &buf4, 2)
		if err != nil {
			t.Errorf("renderText should not error: %v", err)
		}

		// 测试非紧凑模式下的简单文本
		renderer.SetOptions(&RenderOptions{CompactMode: false, Indent: "  "})
		text5 := &Text{Content: "simple text"}
		var buf5 strings.Builder
		err = renderer.renderText(text5, &buf5, 0)
		if err != nil {
			t.Errorf("renderText should not error: %v", err)
		}
	})

	t.Run("renderComment edge cases", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试紧凑模式的注释
		renderer.SetOptions(&RenderOptions{CompactMode: true})
		comment1 := &Comment{Content: " compact comment "}
		var buf1 strings.Builder
		err := renderer.renderComment(comment1, &buf1, 0)
		if err != nil {
			t.Errorf("renderComment should not error: %v", err)
		}
		result1 := buf1.String()
		if strings.Contains(result1, "\n") {
			t.Error("compact mode should not add newlines")
		}

		// 测试非紧凑模式，深度为0的注释
		renderer.SetOptions(&RenderOptions{CompactMode: false})
		comment2 := &Comment{Content: " root level comment "}
		var buf2 strings.Builder
		err = renderer.renderComment(comment2, &buf2, 0)
		if err != nil {
			t.Errorf("renderComment should not error: %v", err)
		}
		result2 := buf2.String()
		if !strings.Contains(result2, "<!--") {
			t.Error("should contain comment start marker")
		}

		// 测试有缩进的注释
		comment3 := &Comment{Content: " indented comment "}
		var buf3 strings.Builder
		err = renderer.renderComment(comment3, &buf3, 2)
		if err != nil {
			t.Errorf("renderComment should not error: %v", err)
		}
	})

	t.Run("renderProcessingInstruction edge cases", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试不包含声明的情况
		renderer.SetOptions(&RenderOptions{IncludeDeclaration: false})
		pi1 := &ProcessingInstruction{Target: "xml", Content: "version=\"1.0\""}
		var buf1 strings.Builder
		err := renderer.renderProcessingInstruction(pi1, &buf1, 0)
		if err != nil {
			t.Errorf("renderProcessingInstruction should not error: %v", err)
		}
		if buf1.String() != "" {
			t.Error("should not render when IncludeDeclaration is false")
		}

		// 测试包含声明的情况
		renderer.SetOptions(&RenderOptions{IncludeDeclaration: true, CompactMode: false})
		pi2 := &ProcessingInstruction{Target: "xml", Content: "version=\"1.0\""}
		var buf2 strings.Builder
		err = renderer.renderProcessingInstruction(pi2, &buf2, 0)
		if err != nil {
			t.Errorf("renderProcessingInstruction should not error: %v", err)
		}
		result2 := buf2.String()
		if !strings.Contains(result2, "<?xml") {
			t.Error("should contain processing instruction")
		}

		// 测试空内容的处理指令
		pi3 := &ProcessingInstruction{Target: "target-only", Content: ""}
		var buf3 strings.Builder
		err = renderer.renderProcessingInstruction(pi3, &buf3, 0)
		if err != nil {
			t.Errorf("renderProcessingInstruction should not error: %v", err)
		}
		result3 := buf3.String()
		expected := "<?target-only?>\n"
		if result3 != expected {
			t.Errorf("expected %q, got %q", expected, result3)
		}

		// 测试紧凑模式
		renderer.SetOptions(&RenderOptions{IncludeDeclaration: true, CompactMode: true})
		pi4 := &ProcessingInstruction{Target: "xml", Content: "version=\"1.0\""}
		var buf4 strings.Builder
		err = renderer.renderProcessingInstruction(pi4, &buf4, 0)
		if err != nil {
			t.Errorf("renderProcessingInstruction should not error: %v", err)
		}
		result4 := buf4.String()
		if strings.Contains(result4, "\n") {
			t.Error("compact mode should not add newlines")
		}

		// 测试有缩进的处理指令
		renderer.SetOptions(&RenderOptions{IncludeDeclaration: true, CompactMode: false})
		pi5 := &ProcessingInstruction{Target: "xml", Content: "version=\"1.0\""}
		var buf5 strings.Builder
		err = renderer.renderProcessingInstruction(pi5, &buf5, 2)
		if err != nil {
			t.Errorf("renderProcessingInstruction should not error: %v", err)
		}
	})

	t.Run("renderDoctype comprehensive", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试不包含声明的情况
		renderer.SetOptions(&RenderOptions{IncludeDeclaration: false})
		doctype1 := &Doctype{Content: "html"}
		var buf1 strings.Builder
		err := renderer.renderDoctype(doctype1, &buf1, 0)
		if err != nil {
			t.Errorf("renderDoctype should not error: %v", err)
		}
		if buf1.String() != "" {
			t.Error("should not render when IncludeDeclaration is false")
		}

		// 测试包含声明的情况
		renderer.SetOptions(&RenderOptions{IncludeDeclaration: true, CompactMode: false})
		doctype2 := &Doctype{Content: "html"}
		var buf2 strings.Builder
		err = renderer.renderDoctype(doctype2, &buf2, 0)
		if err != nil {
			t.Errorf("renderDoctype should not error: %v", err)
		}
		result2 := buf2.String()
		expected2 := "<!DOCTYPE html>\n"
		if result2 != expected2 {
			t.Errorf("expected %q, got %q", expected2, result2)
		}

		// 测试复杂的 DOCTYPE
		doctype3 := &Doctype{Content: "html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\""}
		var buf3 strings.Builder
		err = renderer.renderDoctype(doctype3, &buf3, 0)
		if err != nil {
			t.Errorf("renderDoctype should not error: %v", err)
		}
		result3 := buf3.String()
		if !strings.Contains(result3, "PUBLIC") {
			t.Error("should contain PUBLIC identifier")
		}

		// 测试紧凑模式
		renderer.SetOptions(&RenderOptions{IncludeDeclaration: true, CompactMode: true})
		doctype4 := &Doctype{Content: "html"}
		var buf4 strings.Builder
		err = renderer.renderDoctype(doctype4, &buf4, 0)
		if err != nil {
			t.Errorf("renderDoctype should not error: %v", err)
		}
		result4 := buf4.String()
		if strings.Contains(result4, "\n") {
			t.Error("compact mode should not add newlines")
		}

		// 测试有缩进的 DOCTYPE
		renderer.SetOptions(&RenderOptions{IncludeDeclaration: true, CompactMode: false})
		doctype5 := &Doctype{Content: "html"}
		var buf5 strings.Builder
		err = renderer.renderDoctype(doctype5, &buf5, 2)
		if err != nil {
			t.Errorf("renderDoctype should not error: %v", err)
		}
	})

	t.Run("renderCDATA edge cases", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试紧凑模式
		renderer.SetOptions(&RenderOptions{CompactMode: true})
		cdata1 := &CDATA{Content: "alert('test');"}
		var buf1 strings.Builder
		err := renderer.renderCDATA(cdata1, &buf1, 0)
		if err != nil {
			t.Errorf("renderCDATA should not error: %v", err)
		}
		result1 := buf1.String()
		if strings.Contains(result1, "\n") {
			t.Error("compact mode should not add newlines")
		}

		// 测试非紧凑模式，深度为0
		renderer.SetOptions(&RenderOptions{CompactMode: false})
		cdata2 := &CDATA{Content: "some data"}
		var buf2 strings.Builder
		err = renderer.renderCDATA(cdata2, &buf2, 0)
		if err != nil {
			t.Errorf("renderCDATA should not error: %v", err)
		}
		result2 := buf2.String()
		expected2 := "<![CDATA[some data]]>\n"
		if result2 != expected2 {
			t.Errorf("expected %q, got %q", expected2, result2)
		}

		// 测试有缩进的 CDATA
		cdata3 := &CDATA{Content: "indented data"}
		var buf3 strings.Builder
		err = renderer.renderCDATA(cdata3, &buf3, 2)
		if err != nil {
			t.Errorf("renderCDATA should not error: %v", err)
		}

		// 测试空内容 CDATA
		cdata4 := &CDATA{Content: ""}
		var buf4 strings.Builder
		err = renderer.renderCDATA(cdata4, &buf4, 0)
		if err != nil {
			t.Errorf("renderCDATA with empty content should not error: %v", err)
		}
		result4 := buf4.String()
		expected4 := "<![CDATA[]]>\n"
		if result4 != expected4 {
			t.Errorf("expected %q, got %q", expected4, result4)
		}
	})
}

// TestValidationComprehensiveEdgeCases 测试验证的综合边缘情况
func TestValidationComprehensiveEdgeCases(t *testing.T) {
	t.Run("validateDocument with multiple errors", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(&ValidationOptions{CheckWellFormed: true})

		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "123invalid", // 第一个错误
					Children: []Node{
						&Element{TagName: "-invalid"}, // 第二个错误，但只返回第一个
					},
				},
			},
		}

		err := renderer.validateDocument(doc)
		if err == nil {
			t.Error("should return validation error")
		}

		validationErr, ok := err.(*ValidationError)
		if !ok {
			t.Error("error should be ValidationError type")
		} else {
			if !strings.Contains(validationErr.Message, "123invalid") {
				t.Error("should return first error")
			}
		}
	})

	t.Run("validateDocument without validation options", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(nil)

		doc := &Document{
			Children: []Node{
				&Element{TagName: "123invalid"}, // 这个无效但不应该被检查
			},
		}

		err := renderer.validateDocument(doc)
		if err != nil {
			t.Error("should not validate when validation is nil")
		}
	})

	t.Run("validateDocument with no errors", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(&ValidationOptions{CheckWellFormed: true})

		doc := &Document{
			Children: []Node{
				&Element{
					TagName: "valid-element",
					Attributes: map[string]string{
						"valid-attr": "value",
					},
					Children: []Node{
						&Text{Content: "valid text"},
					},
				},
			},
		}

		err := renderer.validateDocument(doc)
		if err != nil {
			t.Errorf("valid document should not cause error: %v", err)
		}
	})

	t.Run("validateNode with unknown node type", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(&ValidationOptions{CheckWellFormed: true})

		// 测试注释节点（应该不返回错误）
		comment := &Comment{Content: " test comment "}
		err := renderer.validateNode(comment)
		if err != nil {
			t.Error("comment node should not cause validation error")
		}

		// 测试处理指令节点（应该不返回错误）
		pi := &ProcessingInstruction{Target: "xml", Content: "version=\"1.0\""}
		err = renderer.validateNode(pi)
		if err != nil {
			t.Error("processing instruction node should not cause validation error")
		}
	})

	t.Run("validateElement with invalid attributes", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(&ValidationOptions{CheckWellFormed: true})

		elem := &Element{
			TagName: "valid-element",
			Attributes: map[string]string{
				"123invalid-attr": "value", // 无效属性名
				"valid-attr":      "value",
			},
		}

		err := renderer.validateElement(elem)
		if err == nil {
			t.Error("should return error for invalid attribute name")
		}

		validationErr, ok := err.(*ValidationError)
		if !ok {
			t.Error("error should be ValidationError type")
		} else {
			if !strings.Contains(validationErr.Message, "attribute name") {
				t.Error("error should mention attribute name")
			}
		}
	})

	t.Run("validateElement without CheckWellFormed", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetValidation(&ValidationOptions{CheckWellFormed: false})

		elem := &Element{
			TagName: "123invalid", // 这个无效但不应该被检查
		}

		err := renderer.validateElement(elem)
		if err != nil {
			t.Error("should not validate when CheckWellFormed is false")
		}
	})
}

// TestRenderNodeEdgeCases 测试 renderNode 的边缘情况
func TestRenderNodeEdgeCases(t *testing.T) {
	t.Run("renderNode with nil node", func(t *testing.T) {
		renderer := NewRenderer()
		var buf strings.Builder

		err := renderer.renderNode(nil, &buf, 0)
		if err != nil {
			t.Error("renderNode with nil should not error")
		}
		if buf.String() != "" {
			t.Error("nil node should render nothing")
		}
	})

	t.Run("renderNode with unknown node type", func(t *testing.T) {
		renderer := NewRenderer()
		var buf strings.Builder

		// 创建一个不支持的节点类型
		unknownNode := &UnknownNode{pos: Position{Line: 0, Column: 0}}

		err := renderer.renderNode(unknownNode, &buf, 0)
		if err == nil {
			t.Error("should return error for unknown node type")
		}
		if !strings.Contains(err.Error(), "unknown node type") {
			t.Error("error should mention unknown node type")
		}
	})

	t.Run("renderNode with Document node", func(t *testing.T) {
		renderer := NewRenderer()
		var buf strings.Builder

		doc := &Document{
			Children: []Node{
				&Text{Content: "test content"},
			},
		}

		err := renderer.renderNode(doc, &buf, 0)
		if err != nil {
			t.Errorf("renderNode with Document should not error: %v", err)
		}
		result := buf.String()
		if !strings.Contains(result, "test content") {
			t.Error("should render document content")
		}
	})
}

// TestRenderingMethodsErrorBranches 测试渲染方法的错误分支
func TestRenderingMethodsErrorBranches(t *testing.T) {
	t.Run("RenderToString error branches", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试 validation 错误分支
		renderer.SetValidation(&ValidationOptions{CheckWellFormed: true})
		invalidDoc := &Document{
			Children: []Node{
				&Element{TagName: ""}, // 无效标签名
			},
		}

		_, err := renderer.RenderToString(invalidDoc)
		if err == nil {
			t.Error("should return validation error")
		}
	})

	t.Run("RenderToWriter error branches", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试 validation 错误分支
		renderer.SetValidation(&ValidationOptions{CheckWellFormed: true})
		invalidDoc := &Document{
			Children: []Node{
				&Element{TagName: ""}, // 无效标签名
			},
		}

		var buf strings.Builder
		err := renderer.RenderToWriter(invalidDoc, &buf)
		if err == nil {
			t.Error("should return validation error")
		}
	})

	t.Run("RenderWithValidation error branches", func(t *testing.T) {
		renderer := NewRenderer()

		// 测试恢复原始验证设置的分支
		originalValidation := &ValidationOptions{CheckEncoding: true}
		renderer.SetValidation(originalValidation)

		doc := &Document{
			Children: []Node{
				&Element{TagName: "test"},
			},
		}

		newValidation := &ValidationOptions{CheckWellFormed: true}
		_, err := renderer.RenderWithValidation(doc, newValidation)
		if err != nil {
			t.Errorf("should not error for valid document: %v", err)
		}

		// 验证原始验证设置已恢复
		if renderer.validation != originalValidation {
			t.Error("should restore original validation settings")
		}
	})
}

// TestWriteIndentComprehensive 测试 writeIndent 的全面场景
func TestWriteIndentComprehensive(t *testing.T) {
	t.Run("writeIndent with various depths", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetOptions(&RenderOptions{Indent: "    "}) // 4 spaces

		tests := []struct {
			depth    int
			expected string
		}{
			{0, ""},
			{1, "    "},
			{2, "        "},
			{3, "            "},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("depth_%d", tt.depth), func(t *testing.T) {
				var buf strings.Builder
				err := renderer.writeIndent(&buf, tt.depth)
				if err != nil {
					t.Errorf("writeIndent should not error: %v", err)
				}
				result := buf.String()
				if result != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			})
		}
	})

	t.Run("writeIndent with tab indentation", func(t *testing.T) {
		renderer := NewRenderer()
		renderer.SetOptions(&RenderOptions{Indent: "\t"})

		var buf strings.Builder
		err := renderer.writeIndent(&buf, 3)
		if err != nil {
			t.Errorf("writeIndent should not error: %v", err)
		}
		result := buf.String()
		expected := "\t\t\t"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("writeIndent with error writer", func(t *testing.T) {
		renderer := NewRenderer()
		errorWriter := &errorWriter{}

		err := renderer.writeIndent(errorWriter, 1)
		if err == nil {
			t.Error("should return error when writer fails")
		}
	})
}
