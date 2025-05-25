package main

import (
	"fmt"
	"log"

	"github.com/khicago/markit"
)

func main() {
	// 使用简单的HTML内容测试
	fmt.Println("=== 基本HTML解析测试 ===")
	simpleHTML := `<div class="container">
    <h1>Welcome</h1>
    <p>This is a <strong>test</strong> page.</p>
    <img src="image.jpg" alt="Test Image">
    <br>
    <input type="text" name="username" required>
    <hr>
</div>`

	// 创建HTML配置的解析器
	parser := markit.NewParserWithConfig(simpleHTML, markit.HTMLConfig())

	// 解析文档
	doc, err := parser.Parse()
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	// 打印解析结果
	fmt.Println(markit.PrettyPrint(doc))

	// 演示void elements的处理
	fmt.Println("\n=== Void Elements 演示 ===")
	voidElementsHTML := `<div>
    <p>Before image</p>
    <img src="test.jpg" alt="Test">
    <br>
    <hr>
    <input type="text" name="test">
    <p>After input</p>
</div>`

	voidParser := markit.NewParserWithConfig(voidElementsHTML, markit.HTMLConfig())
	voidDoc, err := voidParser.Parse()
	if err != nil {
		log.Fatalf("解析void elements失败: %v", err)
	}

	fmt.Println(markit.PrettyPrint(voidDoc))

	// 演示属性处理（HTML不区分大小写，布尔属性）
	fmt.Println("\n=== 属性处理演示 ===")
	attributesHTML := `<form>
    <input TYPE="text" NAME="username" REQUIRED>
    <input type="checkbox" name="agree" checked>
    <button type="submit" disabled>Submit</button>
</form>`

	attrParser := markit.NewParserWithConfig(attributesHTML, markit.HTMLConfig())
	attrDoc, err := attrParser.Parse()
	if err != nil {
		log.Fatalf("解析属性失败: %v", err)
	}

	fmt.Println(markit.PrettyPrint(attrDoc))
}
