package markit

import (
	"testing"
)

// trimWhitespace 辅助函数，用于去除字符串两端的空白字符
func trimWhitespace(s string) string {
	start := 0
	end := len(s)

	// 找到第一个非空白字符
	for start < end && isWhitespace(s[start]) {
		start++
	}

	// 找到最后一个非空白字符
	for end > start && isWhitespace(s[end-1]) {
		end--
	}

	return s[start:end]
}

// isWhitespace 检查字符是否为空白字符
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// TestASTConstruction 测试通过解析构建AST
func TestASTConstruction(t *testing.T) {
	input := `<document>
		<header id="main-header">
			<title>Test Document</title>
		</header>
		<content class="main-content">
			<p>First paragraph</p>
			<p>Second paragraph</p>
		</content>
	</document>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// 验证文档结构
	if len(doc.Children) != 1 {
		t.Errorf("expected 1 child in document, got %d", len(doc.Children))
	}

	documentElement := doc.Children[0].(*Element)
	if documentElement.TagName != "document" {
		t.Errorf("expected tag name 'document', got %q", documentElement.TagName)
	}

	if len(documentElement.Children) != 2 {
		t.Errorf("expected 2 children in document element, got %d", len(documentElement.Children))
	}

	// 验证header元素
	headerElement := documentElement.Children[0].(*Element)
	if headerElement.TagName != "header" {
		t.Errorf("expected tag name 'header', got %q", headerElement.TagName)
	}

	if headerElement.Attributes["id"] != "main-header" {
		t.Errorf("expected id attribute 'main-header', got %q", headerElement.Attributes["id"])
	}

	// 验证title元素
	if len(headerElement.Children) != 1 {
		t.Errorf("expected 1 child in header, got %d", len(headerElement.Children))
	}

	titleElement := headerElement.Children[0].(*Element)
	if titleElement.TagName != "title" {
		t.Errorf("expected tag name 'title', got %q", titleElement.TagName)
	}

	// 验证content元素
	contentElement := documentElement.Children[1].(*Element)
	if contentElement.TagName != "content" {
		t.Errorf("expected tag name 'content', got %q", contentElement.TagName)
	}

	if contentElement.Attributes["class"] != "main-content" {
		t.Errorf("expected class attribute 'main-content', got %q", contentElement.Attributes["class"])
	}

	// 验证p元素
	if len(contentElement.Children) != 2 {
		t.Errorf("expected 2 children in content, got %d", len(contentElement.Children))
	}

	for i, expectedText := range []string{"First paragraph", "Second paragraph"} {
		pElement := contentElement.Children[i].(*Element)
		if pElement.TagName != "p" {
			t.Errorf("expected tag name 'p' at index %d, got %q", i, pElement.TagName)
		}

		if len(pElement.Children) != 1 {
			t.Errorf("expected 1 text child in p element at index %d, got %d", i, len(pElement.Children))
			continue
		}

		textNode := pElement.Children[0].(*Text)
		actualText := trimWhitespace(textNode.Content)
		if actualText != expectedText {
			t.Errorf("expected text %q at index %d, got %q", expectedText, i, actualText)
		}
	}
}

// TestASTSelfClosingElements 测试自闭合元素
func TestASTSelfClosingElements(t *testing.T) {
	input := `<container>
		<img src="image1.jpg" alt="Image 1" />
		<br />
		<input type="text" name="username" />
		<hr />
	</container>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	containerElement := doc.Children[0].(*Element)
	if containerElement.TagName != "container" {
		t.Errorf("expected tag name 'container', got %q", containerElement.TagName)
	}

	// 验证自闭合元素
	expectedElements := []struct {
		tagName    string
		attributes map[string]string
	}{
		{"img", map[string]string{"src": "image1.jpg", "alt": "Image 1"}},
		{"br", map[string]string{}},
		{"input", map[string]string{"type": "text", "name": "username"}},
		{"hr", map[string]string{}},
	}

	if len(containerElement.Children) != len(expectedElements) {
		t.Errorf("expected %d children, got %d", len(expectedElements), len(containerElement.Children))
	}

	for i, expected := range expectedElements {
		element := containerElement.Children[i].(*Element)
		if element.TagName != expected.tagName {
			t.Errorf("expected tag name %q at index %d, got %q", expected.tagName, i, element.TagName)
		}

		if !element.SelfClose {
			t.Errorf("expected self-close element at index %d", i)
		}

		for key, value := range expected.attributes {
			if element.Attributes[key] != value {
				t.Errorf("expected attribute %s=%q at index %d, got %q", key, value, i, element.Attributes[key])
			}
		}
	}
}

// TestASTNestedStructure 测试嵌套结构
func TestASTNestedStructure(t *testing.T) {
	input := `<html>
		<head>
			<title>Nested Test</title>
			<meta charset="utf-8" />
		</head>
		<body>
			<div class="container">
				<div class="row">
					<div class="col">Column 1</div>
					<div class="col">Column 2</div>
				</div>
			</div>
		</body>
	</html>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	htmlElement := doc.Children[0].(*Element)
	if htmlElement.TagName != "html" {
		t.Errorf("expected tag name 'html', got %q", htmlElement.TagName)
	}

	if len(htmlElement.Children) != 2 {
		t.Errorf("expected 2 children in html, got %d", len(htmlElement.Children))
	}

	// 验证head元素
	headElement := htmlElement.Children[0].(*Element)
	if headElement.TagName != "head" {
		t.Errorf("expected tag name 'head', got %q", headElement.TagName)
	}

	// 验证body元素
	bodyElement := htmlElement.Children[1].(*Element)
	if bodyElement.TagName != "body" {
		t.Errorf("expected tag name 'body', got %q", bodyElement.TagName)
	}

	// 验证嵌套的div结构
	containerDiv := bodyElement.Children[0].(*Element)
	if containerDiv.TagName != "div" || containerDiv.Attributes["class"] != "container" {
		t.Errorf("expected container div, got %q with class %q", containerDiv.TagName, containerDiv.Attributes["class"])
	}

	rowDiv := containerDiv.Children[0].(*Element)
	if rowDiv.TagName != "div" || rowDiv.Attributes["class"] != "row" {
		t.Errorf("expected row div, got %q with class %q", rowDiv.TagName, rowDiv.Attributes["class"])
	}

	if len(rowDiv.Children) != 2 {
		t.Errorf("expected 2 column divs, got %d", len(rowDiv.Children))
	}

	for i, expectedText := range []string{"Column 1", "Column 2"} {
		colDiv := rowDiv.Children[i].(*Element)
		if colDiv.TagName != "div" || colDiv.Attributes["class"] != "col" {
			t.Errorf("expected col div at index %d, got %q with class %q", i, colDiv.TagName, colDiv.Attributes["class"])
		}

		if len(colDiv.Children) != 1 {
			t.Errorf("expected 1 text child in col div at index %d, got %d", i, len(colDiv.Children))
			continue
		}

		textNode := colDiv.Children[0].(*Text)
		actualText := trimWhitespace(textNode.Content)
		if actualText != expectedText {
			t.Errorf("expected text %q at index %d, got %q", expectedText, i, actualText)
		}
	}
}

// TestASTPositionTracking 测试位置跟踪
func TestASTPositionTracking(t *testing.T) {
	input := `<root>
	<child>text</child>
</root>`

	parser := NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// 验证根元素位置
	rootElement := doc.Children[0].(*Element)
	rootPos := rootElement.Position()
	if rootPos.Line != 1 || rootPos.Column != 1 {
		t.Errorf("expected root element at line 1, column 1, got line %d, column %d", rootPos.Line, rootPos.Column)
	}

	// 验证子元素位置
	childElement := rootElement.Children[0].(*Element)
	childPos := childElement.Position()
	if childPos.Line != 2 || childPos.Column != 2 {
		t.Errorf("expected child element at line 2, column 2, got line %d, column %d", childPos.Line, childPos.Column)
	}
}
