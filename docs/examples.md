---
layout: default
title: "Examples - Real-World Usage"
description: "Practical examples of using MarkIt parser for web scraping, template processing, API parsing, and more."
keywords: "markit examples, go xml parser examples, html parser usage, web scraping"
author: "MarkIt Team"
---

# Examples

> **Learn MarkIt through practical, real-world examples**

This page provides comprehensive examples showing how to use MarkIt in various scenarios.

## 📋 Table of Contents

- [Basic Usage](#basic-usage)
- [Web Scraping](#web-scraping)
- [Template Processing](#template-processing)
- [API Response Parsing](#api-response-parsing)
- [Configuration Files](#configuration-files)
- [Custom Transformations](#custom-transformations)
- [Performance Optimization](#performance-optimization)
- [Error Handling](#error-handling)

## Basic Usage

### Simple XML Parsing

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/khicago/markit"
)

func main() {
    xmlContent := `
    <bookstore>
        <book id="1" category="fiction">
            <title>The Great Gatsby</title>
            <author>F. Scott Fitzgerald</author>
            <price currency="USD">12.99</price>
        </book>
        <book id="2" category="science">
            <title>A Brief History of Time</title>
            <author>Stephen Hawking</author>
            <price currency="USD">15.99</price>
        </book>
    </bookstore>`
    
    // Parse with default configuration
    parser := markit.NewParser(xmlContent)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Access the root element
    bookstore := doc.Root
    fmt.Printf("Root element: %s\n", bookstore.TagName)
    
    // Find all books
    books := bookstore.FindChildrenByTag("book")
    fmt.Printf("Found %d books:\n", len(books))
    
    for _, book := range books {
        id, _ := book.GetAttribute("id")
        category, _ := book.GetAttribute("category")
        
        title := book.FindChildByTag("title")
        author := book.FindChildByTag("author")
        price := book.FindChildByTag("price")
        
        fmt.Printf("Book %s (%s):\n", id, category)
        if title != nil && len(title.Children) > 0 {
            if textNode, ok := title.Children[0].(*markit.TextNode); ok {
                fmt.Printf("  Title: %s\n", textNode.Content)
            }
        }
        if author != nil && len(author.Children) > 0 {
            if textNode, ok := author.Children[0].(*markit.TextNode); ok {
                fmt.Printf("  Author: %s\n", textNode.Content)
            }
        }
        if price != nil && len(price.Children) > 0 {
            if textNode, ok := price.Children[0].(*markit.TextNode); ok {
                currency, _ := price.GetAttribute("currency")
                fmt.Printf("  Price: %s %s\n", textNode.Content, currency)
            }
        }
        fmt.Println()
    }
}
```

### HTML5 Document Parsing

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/khicago/markit"
)

func main() {
    htmlContent := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Sample Page</title>
    </head>
    <body>
        <header>
            <h1>Welcome to My Site</h1>
            <nav>
                <ul>
                    <li><a href="/">Home</a></li>
                    <li><a href="/about">About</a></li>
                    <li><a href="/contact">Contact</a></li>
                </ul>
            </nav>
        </header>
        <main>
            <article>
                <h2>Article Title</h2>
                <p>This is the first paragraph.</p>
                <p>This is the second paragraph with <strong>bold text</strong>.</p>
                <img src="image.jpg" alt="Sample Image">
            </article>
        </main>
        <footer>
            <p>&copy; 2024 My Website</p>
        </footer>
    </body>
    </html>`
    
    // Use HTML configuration for proper HTML5 parsing
    config := markit.HTMLConfig()
    parser := markit.NewParserWithConfig(htmlContent, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Extract page title
    title := doc.Root.FindDescendantByTag("title")
    if title != nil && len(title.Children) > 0 {
        if textNode, ok := title.Children[0].(*markit.TextNode); ok {
            fmt.Printf("Page Title: %s\n", strings.TrimSpace(textNode.Content))
        }
    }
    
    // Extract navigation links
    nav := doc.Root.FindDescendantByTag("nav")
    if nav != nil {
        links := nav.FindDescendantsByTag("a")
        fmt.Println("Navigation Links:")
        for _, link := range links {
            href, _ := link.GetAttribute("href")
            if len(link.Children) > 0 {
                if textNode, ok := link.Children[0].(*markit.TextNode); ok {
                    fmt.Printf("  %s -> %s\n", strings.TrimSpace(textNode.Content), href)
                }
            }
        }
    }
    
    // Extract article content
    article := doc.Root.FindDescendantByTag("article")
    if article != nil {
        h2 := article.FindChildByTag("h2")
        if h2 != nil && len(h2.Children) > 0 {
            if textNode, ok := h2.Children[0].(*markit.TextNode); ok {
                fmt.Printf("\nArticle: %s\n", strings.TrimSpace(textNode.Content))
            }
        }
        
        paragraphs := article.FindChildrenByTag("p")
        fmt.Println("Paragraphs:")
        for i, p := range paragraphs {
            fmt.Printf("  %d. %s\n", i+1, extractTextContent(p))
        }
    }
}

// Helper function to extract all text content from an element
func extractTextContent(element *markit.Element) string {
    var texts []string
    
    var extractText func(node markit.Node)
    extractText = func(node markit.Node) {
        switch n := node.(type) {
        case *markit.TextNode:
            texts = append(texts, n.Content)
        case *markit.Element:
            for _, child := range n.Children {
                extractText(child)
            }
        }
    }
    
    extractText(element)
    return strings.TrimSpace(strings.Join(texts, ""))
}
```

## Web Scraping

### Extracting Product Information

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "io"
    "strings"
    
    "github.com/khicago/markit"
)

type Product struct {
    Name        string
    Price       string
    Description string
    ImageURL    string
    Rating      string
}

func main() {
    // Simulate fetching HTML content (replace with actual HTTP request)
    htmlContent := `
    <div class="product-list">
        <div class="product" data-id="1">
            <img src="/images/product1.jpg" alt="Product 1">
            <h3 class="product-name">Wireless Headphones</h3>
            <p class="price">$99.99</p>
            <p class="description">High-quality wireless headphones with noise cancellation.</p>
            <div class="rating" data-rating="4.5">★★★★☆</div>
        </div>
        <div class="product" data-id="2">
            <img src="/images/product2.jpg" alt="Product 2">
            <h3 class="product-name">Smart Watch</h3>
            <p class="price">$199.99</p>
            <p class="description">Feature-rich smartwatch with health monitoring.</p>
            <div class="rating" data-rating="4.8">★★★★★</div>
        </div>
    </div>`
    
    products := scrapeProducts(htmlContent)
    
    fmt.Printf("Found %d products:\n\n", len(products))
    for i, product := range products {
        fmt.Printf("Product %d:\n", i+1)
        fmt.Printf("  Name: %s\n", product.Name)
        fmt.Printf("  Price: %s\n", product.Price)
        fmt.Printf("  Description: %s\n", product.Description)
        fmt.Printf("  Image: %s\n", product.ImageURL)
        fmt.Printf("  Rating: %s\n", product.Rating)
        fmt.Println()
    }
}

func scrapeProducts(htmlContent string) []Product {
    // Use HTML configuration for web scraping
    config := markit.HTMLConfig()
    config.SkipComments = true // Skip comments for performance
    
    parser := markit.NewParserWithConfig(htmlContent, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    var products []Product
    
    // Find all product divs
    productElements := findElementsByClass(doc.Root, "product")
    
    for _, productEl := range productElements {
        product := Product{}
        
        // Extract product name
        nameEl := findElementByClass(productEl, "product-name")
        if nameEl != nil {
            product.Name = extractTextContent(nameEl)
        }
        
        // Extract price
        priceEl := findElementByClass(productEl, "price")
        if priceEl != nil {
            product.Price = extractTextContent(priceEl)
        }
        
        // Extract description
        descEl := findElementByClass(productEl, "description")
        if descEl != nil {
            product.Description = extractTextContent(descEl)
        }
        
        // Extract image URL
        imgEl := productEl.FindDescendantByTag("img")
        if imgEl != nil {
            src, _ := imgEl.GetAttribute("src")
            product.ImageURL = src
        }
        
        // Extract rating
        ratingEl := findElementByClass(productEl, "rating")
        if ratingEl != nil {
            rating, _ := ratingEl.GetAttribute("data-rating")
            product.Rating = rating
        }
        
        products = append(products, product)
    }
    
    return products
}

// Helper function to find elements by class name
func findElementsByClass(root *markit.Element, className string) []*markit.Element {
    var result []*markit.Element
    
    var search func(element *markit.Element)
    search = func(element *markit.Element) {
        // Check if current element has the class
        if hasClass(element, className) {
            result = append(result, element)
        }
        
        // Search in children
        for _, child := range element.Children {
            if childEl, ok := child.(*markit.Element); ok {
                search(childEl)
            }
        }
    }
    
    search(root)
    return result
}

func findElementByClass(root *markit.Element, className string) *markit.Element {
    elements := findElementsByClass(root, className)
    if len(elements) > 0 {
        return elements[0]
    }
    return nil
}

func hasClass(element *markit.Element, className string) bool {
    class, exists := element.GetAttribute("class")
    if !exists {
        return false
    }
    
    classes := strings.Fields(class)
    for _, c := range classes {
        if c == className {
            return true
        }
    }
    return false
}

func extractTextContent(element *markit.Element) string {
    var texts []string
    
    var extractText func(node markit.Node)
    extractText = func(node markit.Node) {
        switch n := node.(type) {
        case *markit.TextNode:
            texts = append(texts, n.Content)
        case *markit.Element:
            for _, child := range n.Children {
                extractText(child)
            }
        }
    }
    
    extractText(element)
    return strings.TrimSpace(strings.Join(texts, ""))
}
```

### RSS Feed Parser

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "time"
    
    "github.com/khicago/markit"
)

type RSSItem struct {
    Title       string
    Link        string
    Description string
    PubDate     string
    GUID        string
}

type RSSFeed struct {
    Title       string
    Link        string
    Description string
    Items       []RSSItem
}

func main() {
    rssContent := `<?xml version="1.0" encoding="UTF-8"?>
    <rss version="2.0">
        <channel>
            <title>Tech News</title>
            <link>https://example.com</link>
            <description>Latest technology news and updates</description>
            <item>
                <title>New Programming Language Released</title>
                <link>https://example.com/news/1</link>
                <description>A revolutionary new programming language has been announced...</description>
                <pubDate>Mon, 01 Jan 2024 10:00:00 GMT</pubDate>
                <guid>https://example.com/news/1</guid>
            </item>
            <item>
                <title>AI Breakthrough in Machine Learning</title>
                <link>https://example.com/news/2</link>
                <description>Researchers have achieved a significant breakthrough in AI...</description>
                <pubDate>Sun, 31 Dec 2023 15:30:00 GMT</pubDate>
                <guid>https://example.com/news/2</guid>
            </item>
        </channel>
    </rss>`
    
    feed := parseRSSFeed(rssContent)
    
    fmt.Printf("Feed: %s\n", feed.Title)
    fmt.Printf("Description: %s\n", feed.Description)
    fmt.Printf("Link: %s\n\n", feed.Link)
    
    fmt.Printf("Items (%d):\n", len(feed.Items))
    for i, item := range feed.Items {
        fmt.Printf("%d. %s\n", i+1, item.Title)
        fmt.Printf("   Link: %s\n", item.Link)
        fmt.Printf("   Date: %s\n", item.PubDate)
        fmt.Printf("   Description: %s\n", truncateString(item.Description, 100))
        fmt.Println()
    }
}

func parseRSSFeed(rssContent string) RSSFeed {
    parser := markit.NewParser(rssContent)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    feed := RSSFeed{}
    
    // Find the channel element
    channel := doc.Root.FindDescendantByTag("channel")
    if channel == nil {
        log.Fatal("No channel element found")
    }
    
    // Extract feed metadata
    if title := channel.FindChildByTag("title"); title != nil {
        feed.Title = extractTextContent(title)
    }
    if link := channel.FindChildByTag("link"); link != nil {
        feed.Link = extractTextContent(link)
    }
    if desc := channel.FindChildByTag("description"); desc != nil {
        feed.Description = extractTextContent(desc)
    }
    
    // Extract items
    items := channel.FindChildrenByTag("item")
    for _, itemEl := range items {
        item := RSSItem{}
        
        if title := itemEl.FindChildByTag("title"); title != nil {
            item.Title = extractTextContent(title)
        }
        if link := itemEl.FindChildByTag("link"); link != nil {
            item.Link = extractTextContent(link)
        }
        if desc := itemEl.FindChildByTag("description"); desc != nil {
            item.Description = extractTextContent(desc)
        }
        if pubDate := itemEl.FindChildByTag("pubDate"); pubDate != nil {
            item.PubDate = extractTextContent(pubDate)
        }
        if guid := itemEl.FindChildByTag("guid"); guid != nil {
            item.GUID = extractTextContent(guid)
        }
        
        feed.Items = append(feed.Items, item)
    }
    
    return feed
}

func truncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen] + "..."
}

func extractTextContent(element *markit.Element) string {
    var texts []string
    
    var extractText func(node markit.Node)
    extractText = func(node markit.Node) {
        switch n := node.(type) {
        case *markit.TextNode:
            texts = append(texts, n.Content)
        case *markit.Element:
            for _, child := range n.Children {
                extractText(child)
            }
        }
    }
    
    extractText(element)
    return strings.TrimSpace(strings.Join(texts, ""))
}
```

## Template Processing

### Simple Template Engine

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/khicago/markit"
)

type TemplateData struct {
    Title   string
    User    string
    Items   []string
    ShowNav bool
}

func main() {
    templateContent := `
    <template>
        <html>
            <head>
                <title>{{title}}</title>
            </head>
            <body>
                <if condition="showNav">
                    <nav>
                        <ul>
                            <li><a href="/">Home</a></li>
                            <li><a href="/profile">Profile</a></li>
                        </ul>
                    </nav>
                </if>
                
                <main>
                    <h1>Welcome, {{user}}!</h1>
                    
                    <if condition="items">
                        <h2>Your Items:</h2>
                        <ul>
                            <for each="items" as="item">
                                <li>{{item}}</li>
                            </for>
                        </ul>
                    </if>
                    
                    <unless condition="items">
                        <p>No items found.</p>
                    </unless>
                </main>
            </body>
        </html>
    </template>`
    
    data := TemplateData{
        Title:   "My Dashboard",
        User:    "John Doe",
        Items:   []string{"Task 1", "Task 2", "Task 3"},
        ShowNav: true,
    }
    
    result := processTemplate(templateContent, data)
    fmt.Println(result)
}

func processTemplate(templateContent string, data TemplateData) string {
    // Configure parser for template processing
    config := markit.DefaultConfig()
    config.CaseSensitive = false
    config.SetVoidElements([]string{"include", "import", "placeholder"})
    
    parser := markit.NewParserWithConfig(templateContent, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Process the template
    processor := &TemplateProcessor{Data: data}
    err = markit.WalkDepthFirst(doc.Root, processor)
    if err != nil {
        log.Fatal(err)
    }
    
    return doc.ToHTML()
}

type TemplateProcessor struct {
    Data TemplateData
}

func (tp *TemplateProcessor) VisitEnter(node markit.Node) error {
    switch n := node.(type) {
    case *markit.Element:
        return tp.processElement(n)
    case *markit.TextNode:
        return tp.processTextNode(n)
    }
    return nil
}

func (tp *TemplateProcessor) VisitLeave(node markit.Node) error {
    return nil
}

func (tp *TemplateProcessor) processElement(element *markit.Element) error {
    switch element.TagName {
    case "if":
        return tp.processIf(element)
    case "unless":
        return tp.processUnless(element)
    case "for":
        return tp.processFor(element)
    }
    return nil
}

func (tp *TemplateProcessor) processTextNode(textNode *markit.TextNode) error {
    // Replace template variables
    content := textNode.Content
    content = strings.ReplaceAll(content, "{{title}}", tp.Data.Title)
    content = strings.ReplaceAll(content, "{{user}}", tp.Data.User)
    textNode.Content = content
    return nil
}

func (tp *TemplateProcessor) processIf(element *markit.Element) error {
    condition, _ := element.GetAttribute("condition")
    shouldShow := tp.evaluateCondition(condition)
    
    if !shouldShow {
        // Remove this element from its parent
        if parent := element.Parent(); parent != nil {
            if parentEl, ok := parent.(*markit.Element); ok {
                parentEl.RemoveChild(element)
            }
        }
    }
    
    return nil
}

func (tp *TemplateProcessor) processUnless(element *markit.Element) error {
    condition, _ := element.GetAttribute("condition")
    shouldShow := !tp.evaluateCondition(condition)
    
    if !shouldShow {
        // Remove this element from its parent
        if parent := element.Parent(); parent != nil {
            if parentEl, ok := parent.(*markit.Element); ok {
                parentEl.RemoveChild(element)
            }
        }
    }
    
    return nil
}

func (tp *TemplateProcessor) processFor(element *markit.Element) error {
    each, _ := element.GetAttribute("each")
    as, _ := element.GetAttribute("as")
    
    if each == "items" {
        // Clone the element for each item
        parent := element.Parent()
        if parentEl, ok := parent.(*markit.Element); ok {
            // Remove the original for element
            parentEl.RemoveChild(element)
            
            // Add cloned elements for each item
            for _, item := range tp.Data.Items {
                clone := tp.cloneElement(element)
                tp.replaceInElement(clone, "{{"+as+"}}", item)
                parentEl.AddChild(clone)
            }
        }
    }
    
    return nil
}

func (tp *TemplateProcessor) evaluateCondition(condition string) bool {
    switch condition {
    case "showNav":
        return tp.Data.ShowNav
    case "items":
        return len(tp.Data.Items) > 0
    }
    return false
}

func (tp *TemplateProcessor) cloneElement(element *markit.Element) *markit.Element {
    clone := &markit.Element{
        TagName:    element.TagName,
        Attributes: make(map[string]string),
        SelfClosed: element.SelfClosed,
    }
    
    // Copy attributes
    for k, v := range element.Attributes {
        clone.Attributes[k] = v
    }
    
    // Clone children
    for _, child := range element.Children {
        switch c := child.(type) {
        case *markit.Element:
            childClone := tp.cloneElement(c)
            clone.AddChild(childClone)
        case *markit.TextNode:
            textClone := &markit.TextNode{Content: c.Content}
            clone.AddChild(textClone)
        }
    }
    
    return clone
}

func (tp *TemplateProcessor) replaceInElement(element *markit.Element, placeholder, value string) {
    for _, child := range element.Children {
        switch c := child.(type) {
        case *markit.Element:
            tp.replaceInElement(c, placeholder, value)
        case *markit.TextNode:
            c.Content = strings.ReplaceAll(c.Content, placeholder, value)
        }
    }
}
```

## API Response Parsing

### JSON-like XML API Response

```go
package main

import (
    "fmt"
    "log"
    "strconv"
    
    "github.com/khicago/markit"
)

type APIResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    UserData `json:"data"`
}

type UserData struct {
    ID       int      `json:"id"`
    Name     string   `json:"name"`
    Email    string   `json:"email"`
    Active   bool     `json:"active"`
    Tags     []string `json:"tags"`
    Profile  Profile  `json:"profile"`
}

type Profile struct {
    Bio     string `json:"bio"`
    Website string `json:"website"`
    Age     int    `json:"age"`
}

func main() {
    xmlResponse := `<?xml version="1.0" encoding="UTF-8"?>
    <response>
        <status>success</status>
        <message>User data retrieved successfully</message>
        <data>
            <user>
                <id>12345</id>
                <name>John Doe</name>
                <email>john.doe@example.com</email>
                <active>true</active>
                <tags>
                    <tag>developer</tag>
                    <tag>golang</tag>
                    <tag>api</tag>
                </tags>
                <profile>
                    <bio>Software developer with 5+ years experience</bio>
                    <website>https://johndoe.dev</website>
                    <age>30</age>
                </profile>
            </user>
        </data>
    </response>`
    
    response := parseAPIResponse(xmlResponse)
    
    fmt.Printf("Status: %s\n", response.Status)
    fmt.Printf("Message: %s\n", response.Message)
    fmt.Printf("\nUser Data:\n")
    fmt.Printf("  ID: %d\n", response.Data.ID)
    fmt.Printf("  Name: %s\n", response.Data.Name)
    fmt.Printf("  Email: %s\n", response.Data.Email)
    fmt.Printf("  Active: %t\n", response.Data.Active)
    fmt.Printf("  Tags: %v\n", response.Data.Tags)
    fmt.Printf("\nProfile:\n")
    fmt.Printf("  Bio: %s\n", response.Data.Profile.Bio)
    fmt.Printf("  Website: %s\n", response.Data.Profile.Website)
    fmt.Printf("  Age: %d\n", response.Data.Profile.Age)
}

func parseAPIResponse(xmlContent string) APIResponse {
    // Use default config for strict API parsing
    config := markit.DefaultConfig()
    config.SkipComments = true
    
    parser := markit.NewParserWithConfig(xmlContent, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    response := APIResponse{}
    
    root := doc.Root
    
    // Parse status
    if status := root.FindChildByTag("status"); status != nil {
        response.Status = extractTextContent(status)
    }
    
    // Parse message
    if message := root.FindChildByTag("message"); message != nil {
        response.Message = extractTextContent(message)
    }
    
    // Parse data
    if data := root.FindChildByTag("data"); data != nil {
        if user := data.FindChildByTag("user"); user != nil {
            response.Data = parseUserData(user)
        }
    }
    
    return response
}

func parseUserData(userElement *markit.Element) UserData {
    userData := UserData{}
    
    // Parse ID
    if id := userElement.FindChildByTag("id"); id != nil {
        if idStr := extractTextContent(id); idStr != "" {
            if idInt, err := strconv.Atoi(idStr); err == nil {
                userData.ID = idInt
            }
        }
    }
    
    // Parse name
    if name := userElement.FindChildByTag("name"); name != nil {
        userData.Name = extractTextContent(name)
    }
    
    // Parse email
    if email := userElement.FindChildByTag("email"); email != nil {
        userData.Email = extractTextContent(email)
    }
    
    // Parse active
    if active := userElement.FindChildByTag("active"); active != nil {
        if activeStr := extractTextContent(active); activeStr != "" {
            userData.Active = activeStr == "true"
        }
    }
    
    // Parse tags
    if tags := userElement.FindChildByTag("tags"); tags != nil {
        tagElements := tags.FindChildrenByTag("tag")
        for _, tagEl := range tagElements {
            userData.Tags = append(userData.Tags, extractTextContent(tagEl))
        }
    }
    
    // Parse profile
    if profile := userElement.FindChildByTag("profile"); profile != nil {
        userData.Profile = parseProfile(profile)
    }
    
    return userData
}

func parseProfile(profileElement *markit.Element) Profile {
    profile := Profile{}
    
    if bio := profileElement.FindChildByTag("bio"); bio != nil {
        profile.Bio = extractTextContent(bio)
    }
    
    if website := profileElement.FindChildByTag("website"); website != nil {
        profile.Website = extractTextContent(website)
    }
    
    if age := profileElement.FindChildByTag("age"); age != nil {
        if ageStr := extractTextContent(age); ageStr != "" {
            if ageInt, err := strconv.Atoi(ageStr); err == nil {
                profile.Age = ageInt
            }
        }
    }
    
    return profile
}

func extractTextContent(element *markit.Element) string {
    var texts []string
    
    var extractText func(node markit.Node)
    extractText = func(node markit.Node) {
        switch n := node.(type) {
        case *markit.TextNode:
            texts = append(texts, n.Content)
        case *markit.Element:
            for _, child := range n.Children {
                extractText(child)
            }
        }
    }
    
    extractText(element)
    return strings.TrimSpace(strings.Join(texts, ""))
}
```

## Configuration Files

### Application Configuration Parser

```go
package main

import (
    "fmt"
    "log"
    "strconv"
    "strings"
    
    "github.com/khicago/markit"
)

type AppConfig struct {
    Server   ServerConfig   `xml:"server"`
    Database DatabaseConfig `xml:"database"`
    Logging  LoggingConfig  `xml:"logging"`
    Features FeatureFlags   `xml:"features"`
}

type ServerConfig struct {
    Host         string `xml:"host"`
    Port         int    `xml:"port"`
    ReadTimeout  int    `xml:"readTimeout"`
    WriteTimeout int    `xml:"writeTimeout"`
    TLS          TLSConfig `xml:"tls"`
}

type TLSConfig struct {
    Enabled  bool   `xml:"enabled"`
    CertFile string `xml:"certFile"`
    KeyFile  string `xml:"keyFile"`
}

type DatabaseConfig struct {
    Driver   string `xml:"driver"`
    Host     string `xml:"host"`
    Port     int    `xml:"port"`
    Name     string `xml:"name"`
    Username string `xml:"username"`
    Password string `xml:"password"`
    Pool     PoolConfig `xml:"pool"`
}

type PoolConfig struct {
    MaxOpen     int `xml:"maxOpen"`
    MaxIdle     int `xml:"maxIdle"`
    MaxLifetime int `xml:"maxLifetime"`
}

type LoggingConfig struct {
    Level  string `xml:"level"`
    Format string `xml:"format"`
    Output string `xml:"output"`
}

type FeatureFlags struct {
    EnableMetrics   bool `xml:"enableMetrics"`
    EnableProfiling bool `xml:"enableProfiling"`
    EnableCaching   bool `xml:"enableCaching"`
}

func main() {
    configXML := `<?xml version="1.0" encoding="UTF-8"?>
    <config>
        <!-- Server Configuration -->
        <server>
            <host>localhost</host>
            <port>8080</port>
            <readTimeout>30</readTimeout>
            <writeTimeout>30</writeTimeout>
            <tls>
                <enabled>true</enabled>
                <certFile>/etc/ssl/server.crt</certFile>
                <keyFile>/etc/ssl/server.key</keyFile>
            </tls>
        </server>
        
        <!-- Database Configuration -->
        <database>
            <driver>postgres</driver>
            <host>db.example.com</host>
            <port>5432</port>
            <name>myapp</name>
            <username>dbuser</username>
            <password>secretpassword</password>
            <pool>
                <maxOpen>25</maxOpen>
                <maxIdle>5</maxIdle>
                <maxLifetime>300</maxLifetime>
            </pool>
        </database>
        
        <!-- Logging Configuration -->
        <logging>
            <level>info</level>
            <format>json</format>
            <output>stdout</output>
        </logging>
        
        <!-- Feature Flags -->
        <features>
            <enableMetrics>true</enableMetrics>
            <enableProfiling>false</enableProfiling>
            <enableCaching>true</enableCaching>
        </features>
    </config>`
    
    config := parseAppConfig(configXML)
    
    fmt.Println("Application Configuration:")
    fmt.Printf("\nServer:\n")
    fmt.Printf("  Host: %s\n", config.Server.Host)
    fmt.Printf("  Port: %d\n", config.Server.Port)
    fmt.Printf("  Read Timeout: %d seconds\n", config.Server.ReadTimeout)
    fmt.Printf("  Write Timeout: %d seconds\n", config.Server.WriteTimeout)
    fmt.Printf("  TLS Enabled: %t\n", config.Server.TLS.Enabled)
    if config.Server.TLS.Enabled {
        fmt.Printf("  TLS Cert: %s\n", config.Server.TLS.CertFile)
        fmt.Printf("  TLS Key: %s\n", config.Server.TLS.KeyFile)
    }
    
    fmt.Printf("\nDatabase:\n")
    fmt.Printf("  Driver: %s\n", config.Database.Driver)
    fmt.Printf("  Host: %s\n", config.Database.Host)
    fmt.Printf("  Port: %d\n", config.Database.Port)
    fmt.Printf("  Name: %s\n", config.Database.Name)
    fmt.Printf("  Username: %s\n", config.Database.Username)
    fmt.Printf("  Pool Max Open: %d\n", config.Database.Pool.MaxOpen)
    fmt.Printf("  Pool Max Idle: %d\n", config.Database.Pool.MaxIdle)
    fmt.Printf("  Pool Max Lifetime: %d seconds\n", config.Database.Pool.MaxLifetime)
    
    fmt.Printf("\nLogging:\n")
    fmt.Printf("  Level: %s\n", config.Logging.Level)
    fmt.Printf("  Format: %s\n", config.Logging.Format)
    fmt.Printf("  Output: %s\n", config.Logging.Output)
    
    fmt.Printf("\nFeature Flags:\n")
    fmt.Printf("  Metrics: %t\n", config.Features.EnableMetrics)
    fmt.Printf("  Profiling: %t\n", config.Features.EnableProfiling)
    fmt.Printf("  Caching: %t\n", config.Features.EnableCaching)
}

func parseAppConfig(configXML string) AppConfig {
    // Use default config for strict configuration parsing
    config := markit.DefaultConfig()
    config.SkipComments = false // Keep comments for documentation
    
    parser := markit.NewParserWithConfig(configXML, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    appConfig := AppConfig{}
    
    root := doc.Root
    
    // Parse server config
    if server := root.FindChildByTag("server"); server != nil {
        appConfig.Server = parseServerConfig(server)
    }
    
    // Parse database config
    if database := root.FindChildByTag("database"); database != nil {
        appConfig.Database = parseDatabaseConfig(database)
    }
    
    // Parse logging config
    if logging := root.FindChildByTag("logging"); logging != nil {
        appConfig.Logging = parseLoggingConfig(logging)
    }
    
    // Parse feature flags
    if features := root.FindChildByTag("features"); features != nil {
        appConfig.Features = parseFeatureFlags(features)
    }
    
    return appConfig
}

func parseServerConfig(serverEl *markit.Element) ServerConfig {
    server := ServerConfig{}
    
    if host := serverEl.FindChildByTag("host"); host != nil {
        server.Host = extractTextContent(host)
    }
    
    if port := serverEl.FindChildByTag("port"); port != nil {
        if portStr := extractTextContent(port); portStr != "" {
            if portInt, err := strconv.Atoi(portStr); err == nil {
                server.Port = portInt
            }
        }
    }
    
    if readTimeout := serverEl.FindChildByTag("readTimeout"); readTimeout != nil {
        if timeoutStr := extractTextContent(readTimeout); timeoutStr != "" {
            if timeoutInt, err := strconv.Atoi(timeoutStr); err == nil {
                server.ReadTimeout = timeoutInt
            }
        }
    }
    
    if writeTimeout := serverEl.FindChildByTag("writeTimeout"); writeTimeout != nil {
        if timeoutStr := extractTextContent(writeTimeout); timeoutStr != "" {
            if timeoutInt, err := strconv.Atoi(timeoutStr); err == nil {
                server.WriteTimeout = timeoutInt
            }
        }
    }
    
    if tls := serverEl.FindChildByTag("tls"); tls != nil {
        server.TLS = parseTLSConfig(tls)
    }
    
    return server
}

func parseTLSConfig(tlsEl *markit.Element) TLSConfig {
    tls := TLSConfig{}
    
    if enabled := tlsEl.FindChildByTag("enabled"); enabled != nil {
        tls.Enabled = extractTextContent(enabled) == "true"
    }
    
    if certFile := tlsEl.FindChildByTag("certFile"); certFile != nil {
        tls.CertFile = extractTextContent(certFile)
    }
    
    if keyFile := tlsEl.FindChildByTag("keyFile"); keyFile != nil {
        tls.KeyFile = extractTextContent(keyFile)
    }
    
    return tls
}

func parseDatabaseConfig(dbEl *markit.Element) DatabaseConfig {
    db := DatabaseConfig{}
    
    if driver := dbEl.FindChildByTag("driver"); driver != nil {
        db.Driver = extractTextContent(driver)
    }
    
    if host := dbEl.FindChildByTag("host"); host != nil {
        db.Host = extractTextContent(host)
    }
    
    if port := dbEl.FindChildByTag("port"); port != nil {
        if portStr := extractTextContent(port); portStr != "" {
            if portInt, err := strconv.Atoi(portStr); err == nil {
                db.Port = portInt
            }
        }
    }
    
    if name := dbEl.FindChildByTag("name"); name != nil {
        db.Name = extractTextContent(name)
    }
    
    if username := dbEl.FindChildByTag("username"); username != nil {
        db.Username = extractTextContent(username)
    }
    
    if password := dbEl.FindChildByTag("password"); password != nil {
        db.Password = extractTextContent(password)
    }
    
    if pool := dbEl.FindChildByTag("pool"); pool != nil {
        db.Pool = parsePoolConfig(pool)
    }
    
    return db
}

func parsePoolConfig(poolEl *markit.Element) PoolConfig {
    pool := PoolConfig{}
    
    if maxOpen := poolEl.FindChildByTag("maxOpen"); maxOpen != nil {
        if maxOpenStr := extractTextContent(maxOpen); maxOpenStr != "" {
            if maxOpenInt, err := strconv.Atoi(maxOpenStr); err == nil {
                pool.MaxOpen = maxOpenInt
            }
        }
    }
    
    if maxIdle := poolEl.FindChildByTag("maxIdle"); maxIdle != nil {
        if maxIdleStr := extractTextContent(maxIdle); maxIdleStr != "" {
            if maxIdleInt, err := strconv.Atoi(maxIdleStr); err == nil {
                pool.MaxIdle = maxIdleInt
            }
        }
    }
    
    if maxLifetime := poolEl.FindChildByTag("maxLifetime"); maxLifetime != nil {
        if maxLifetimeStr := extractTextContent(maxLifetime); maxLifetimeStr != "" {
            if maxLifetimeInt, err := strconv.Atoi(maxLifetimeStr); err == nil {
                pool.MaxLifetime = maxLifetimeInt
            }
        }
    }
    
    return pool
}

func parseLoggingConfig(loggingEl *markit.Element) LoggingConfig {
    logging := LoggingConfig{}
    
    if level := loggingEl.FindChildByTag("level"); level != nil {
        logging.Level = extractTextContent(level)
    }
    
    if format := loggingEl.FindChildByTag("format"); format != nil {
        logging.Format = extractTextContent(format)
    }
    
    if output := loggingEl.FindChildByTag("output"); output != nil {
        logging.Output = extractTextContent(output)
    }
    
    return logging
}

func parseFeatureFlags(featuresEl *markit.Element) FeatureFlags {
    features := FeatureFlags{}
    
    if enableMetrics := featuresEl.FindChildByTag("enableMetrics"); enableMetrics != nil {
        features.EnableMetrics = extractTextContent(enableMetrics) == "true"
    }
    
    if enableProfiling := featuresEl.FindChildByTag("enableProfiling"); enableProfiling != nil {
        features.EnableProfiling = extractTextContent(enableProfiling) == "true"
    }
    
    if enableCaching := featuresEl.FindChildByTag("enableCaching"); enableCaching != nil {
        features.EnableCaching = extractTextContent(enableCaching) == "true"
    }
    
    return features
}

func extractTextContent(element *markit.Element) string {
    var texts []string
    
    var extractText func(node markit.Node)
    extractText = func(node markit.Node) {
        switch n := node.(type) {
        case *markit.TextNode:
            texts = append(texts, n.Content)
        case *markit.Element:
            for _, child := range n.Children {
                extractText(child)
            }
        }
    }
    
    extractText(element)
    return strings.TrimSpace(strings.Join(texts, ""))
}
```

## Custom Transformations

### HTML to Markdown Converter

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/khicago/markit"
)

func main() {
    htmlContent := `
    <article>
        <h1>Getting Started with Go</h1>
        <p>Go is a <strong>powerful</strong> programming language developed by Google.</p>
        
        <h2>Key Features</h2>
        <ul>
            <li>Fast compilation</li>
            <li>Garbage collection</li>
            <li>Concurrency support</li>
        </ul>
        
        <h2>Code Example</h2>
        <pre><code>package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}</code></pre>
        
        <p>For more information, visit the <a href="https://golang.org">official website</a>.</p>
        
        <blockquote>
            <p>Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.</p>
        </blockquote>
    </article>`
    
    markdown := convertHTMLToMarkdown(htmlContent)
    fmt.Println(markdown)
}

func convertHTMLToMarkdown(htmlContent string) string {
    config := markit.HTMLConfig()
    parser := markit.NewParserWithConfig(htmlContent, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    converter := &HTMLToMarkdownConverter{}
    err = markit.WalkDepthFirst(doc.Root, converter)
    if err != nil {
        log.Fatal(err)
    }
    
    return strings.TrimSpace(converter.Result)
}

type HTMLToMarkdownConverter struct {
    Result     string
    ListDepth  int
    InCodeBlock bool
}

func (c *HTMLToMarkdownConverter) VisitEnter(node markit.Node) error {
    switch n := node.(type) {
    case *markit.Element:
        return c.handleElementEnter(n)
    case *markit.TextNode:
        return c.handleTextNode(n)
    }
    return nil
}

func (c *HTMLToMarkdownConverter) VisitLeave(node markit.Node) error {
    if element, ok := node.(*markit.Element); ok {
        return c.handleElementLeave(element)
    }
    return nil
}

func (c *HTMLToMarkdownConverter) handleElementEnter(element *markit.Element) error {
    switch element.TagName {
    case "h1":
        c.Result += "# "
    case "h2":
        c.Result += "## "
    case "h3":
        c.Result += "### "
    case "h4":
        c.Result += "#### "
    case "h5":
        c.Result += "##### "
    case "h6":
        c.Result += "###### "
    case "p":
        // Paragraph will be handled in leave
    case "strong", "b":
        c.Result += "**"
    case "em", "i":
        c.Result += "*"
    case "code":
        if !c.InCodeBlock {
            c.Result += "`"
        }
    case "pre":
        c.Result += "```\n"
        c.InCodeBlock = true
    case "ul":
        c.ListDepth++
    case "ol":
        c.ListDepth++
    case "li":
        indent := strings.Repeat("  ", c.ListDepth-1)
        c.Result += indent + "- "
    case "a":
        c.Result += "["
    case "blockquote":
        c.Result += "> "
    case "br":
        c.Result += "\n"
    }
    return nil
}

func (c *HTMLToMarkdownConverter) handleElementLeave(element *markit.Element) error {
    switch element.TagName {
    case "h1", "h2", "h3", "h4", "h5", "h6":
        c.Result += "\n\n"
    case "p":
        c.Result += "\n\n"
    case "strong", "b":
        c.Result += "**"
    case "em", "i":
        c.Result += "*"
    case "code":
        if !c.InCodeBlock {
            c.Result += "`"
        }
    case "pre":
        c.Result += "\n```\n\n"
        c.InCodeBlock = false
    case "ul", "ol":
        c.ListDepth--
        if c.ListDepth == 0 {
            c.Result += "\n"
        }
    case "li":
        c.Result += "\n"
    case "a":
        href, _ := element.GetAttribute("href")
        c.Result += "](" + href + ")"
    case "blockquote":
        c.Result += "\n\n"
    }
    return nil
}

func (c *HTMLToMarkdownConverter) handleTextNode(textNode *markit.TextNode) error {
    content := textNode.Content
    if !c.InCodeBlock {
        // Clean up whitespace for regular text
        content = strings.TrimSpace(content)
        if content != "" {
            c.Result += content
        }
    } else {
        // Preserve formatting in code blocks
        c.Result += content
    }
    return nil
}
```

## Performance Optimization

### Batch Processing with Reusable Configuration

```go
package main

import (
    "fmt"
    "log"
    "runtime"
    "sync"
    "time"
    
    "github.com/khicago/markit"
)

type ProcessingResult struct {
    Index    int
    Success  bool
    Error    error
    Duration time.Duration
    Elements int
}

func main() {
    // Simulate multiple documents to process
    documents := generateTestDocuments(1000)
    
    fmt.Printf("Processing %d documents...\n", len(documents))
    
    // Sequential processing
    start := time.Now()
    sequentialResults := processDocumentsSequential(documents)
    sequentialDuration := time.Since(start)
    
    // Parallel processing
    start = time.Now()
    parallelResults := processDocumentsParallel(documents)
    parallelDuration := time.Since(start)
    
    // Print results
    fmt.Printf("\nSequential Processing:\n")
    fmt.Printf("  Duration: %v\n", sequentialDuration)
    fmt.Printf("  Success: %d/%d\n", countSuccessful(sequentialResults), len(sequentialResults))
    
    fmt.Printf("\nParallel Processing:\n")
    fmt.Printf("  Duration: %v\n", parallelDuration)
    fmt.Printf("  Success: %d/%d\n", countSuccessful(parallelResults), len(parallelResults))
    fmt.Printf("  Speedup: %.2fx\n", float64(sequentialDuration)/float64(parallelDuration))
}

func generateTestDocuments(count int) []string {
    documents := make([]string, count)
    
    template := `<root>
        <metadata>
            <id>%d</id>
            <timestamp>%s</timestamp>
        </metadata>
        <content>
            <title>Document %d</title>
            <body>
                <p>This is paragraph 1 of document %d.</p>
                <p>This is paragraph 2 with <strong>bold</strong> text.</p>
                <ul>
                    <li>Item 1</li>
                    <li>Item 2</li>
                    <li>Item 3</li>
                </ul>
            </body>
        </content>
    </root>`
    
    for i := 0; i < count; i++ {
        documents[i] = fmt.Sprintf(template, i, time.Now().Format(time.RFC3339), i, i)
    }
    
    return documents
}

func processDocumentsSequential(documents []string) []ProcessingResult {
    // Reuse configuration for better performance
    config := createOptimizedConfig()
    
    results := make([]ProcessingResult, len(documents))
    
    for i, doc := range documents {
        start := time.Now()
        result := ProcessingResult{Index: i}
        
        parser := markit.NewParserWithConfig(doc, config)
        ast, err := parser.Parse()
        
        result.Duration = time.Since(start)
        
        if err != nil {
            result.Error = err
            result.Success = false
        } else {
            result.Success = true
            result.Elements = countElements(ast.Root)
        }
        
        results[i] = result
    }
    
    return results
}

func processDocumentsParallel(documents []string) []ProcessingResult {
    numWorkers := runtime.NumCPU()
    jobs := make(chan int, len(documents))
    results := make([]ProcessingResult, len(documents))
    
    var wg sync.WaitGroup
    
    // Start workers
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            // Each worker gets its own config to avoid race conditions
            config := createOptimizedConfig()
            
            for i := range jobs {
                start := time.Now()
                result := ProcessingResult{Index: i}
                
                parser := markit.NewParserWithConfig(documents[i], config)
                ast, err := parser.Parse()
                
                result.Duration = time.Since(start)
                
                if err != nil {
                    result.Error = err
                    result.Success = false
                } else {
                    result.Success = true
                    result.Elements = countElements(ast.Root)
                }
                
                results[i] = result
            }
        }()
    }
    
    // Send jobs
    for i := range documents {
        jobs <- i
    }
    close(jobs)
    
    // Wait for completion
    wg.Wait()
    
    return results
}

func createOptimizedConfig() *markit.ParserConfig {
    config := markit.DefaultConfig()
    
    // Optimize for performance
    config.CaseSensitive = true  // Faster string comparisons
    config.SkipComments = true   // Skip comment processing
    
    return config
}

func countElements(element *markit.Element) int {
    count := 1 // Count this element
    
    for _, child := range element.Children {
        if childEl, ok := child.(*markit.Element); ok {
            count += countElements(childEl)
        }
    }
    
    return count
}

func countSuccessful(results []ProcessingResult) int {
    count := 0
    for _, result := range results {
        if result.Success {
            count++
        }
    }
    return count
}
```

### Memory-Efficient Streaming Parser

```go
package main

import (
    "fmt"
    "log"
    "runtime"
    "strings"
    
    "github.com/khicago/markit"
)

type StreamProcessor struct {
    ElementCount map[string]int
    TextLength   int
    MaxDepth     int
    currentDepth int
}

func (sp *StreamProcessor) VisitEnter(node markit.Node) error {
    switch n := node.(type) {
    case *markit.Element:
        sp.currentDepth++
        if sp.currentDepth > sp.MaxDepth {
            sp.MaxDepth = sp.currentDepth
        }
        
        if sp.ElementCount == nil {
            sp.ElementCount = make(map[string]int)
        }
        sp.ElementCount[n.TagName]++
        
    case *markit.TextNode:
        sp.TextLength += len(strings.TrimSpace(n.Content))
    }
    
    return nil
}

func (sp *StreamProcessor) VisitLeave(node markit.Node) error {
    if _, ok := node.(*markit.Element); ok {
        sp.currentDepth--
    }
    return nil
}

func main() {
    // Large document simulation
    largeDocument := generateLargeDocument(10000) // 10k elements
    
    fmt.Printf("Document size: %d characters\n", len(largeDocument))
    
    // Memory usage before parsing
    var m1 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    // Parse with streaming processor
    processor := &StreamProcessor{}
    parseWithStreamProcessor(largeDocument, processor)
    
    // Memory usage after parsing
    var m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    // Print results
    fmt.Printf("\nParsing Results:\n")
    fmt.Printf("  Max Depth: %d\n", processor.MaxDepth)
    fmt.Printf("  Total Text Length: %d\n", processor.TextLength)
    fmt.Printf("  Element Counts:\n")
    for tag, count := range processor.ElementCount {
        fmt.Printf("    %s: %d\n", tag, count)
    }
    
    fmt.Printf("\nMemory Usage:\n")
    fmt.Printf("  Before: %d KB\n", m1.Alloc/1024)
    fmt.Printf("  After: %d KB\n", m2.Alloc/1024)
    fmt.Printf("  Difference: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
}

func generateLargeDocument(elementCount int) string {
    var builder strings.Builder
    builder.WriteString("<root>\n")
    
    for i := 0; i < elementCount; i++ {
        builder.WriteString(fmt.Sprintf(`  <item id="%d">
    <title>Item %d</title>
    <description>This is a description for item %d with some content.</description>
    <metadata>
      <created>2024-01-01</created>
      <category>test</category>
    </metadata>
  </item>
`, i, i, i))
    }
    
    builder.WriteString("</root>")
    return builder.String()
}

func parseWithStreamProcessor(document string, processor *StreamProcessor) {
    // Use memory-optimized configuration
    config := markit.DefaultConfig()
    config.SkipComments = true // Reduce memory usage
    
    parser := markit.NewParserWithConfig(document, config)
    doc, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }
    
    // Process with streaming visitor
    err = markit.WalkDepthFirst(doc.Root, processor)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Error Handling

### Robust Error Handling and Recovery

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/khicago/markit"
)

type ParseAttempt struct {
    Content string
    Config  *markit.ParserConfig
    Name    string
}

func main() {
    // Test various problematic documents
    problemDocuments := []ParseAttempt{
        {
            Content: `<root><unclosed>content`,
            Config:  markit.DefaultConfig(),
            Name:    "Unclosed tag",
        },
        {
            Content: `<root><tag attr="unclosed value>content</tag></root>`,
            Config:  markit.DefaultConfig(),
            Name:    "Unclosed attribute",
        },
        {
            Content: `<root><tag><nested><deep></deep></nested></tag></root>`,
            Config:  markit.DefaultConfig(),
            Name:    "Valid nested document",
        },
        {
            Content: `<root><self-closed /></root>`,
            Config:  createStrictConfig(),
            Name:    "Self-closed with strict config",
        },
        {
            Content: `<html><br><img src="test.jpg"></html>`,
            Config:  markit.HTMLConfig(),
            Name:    "HTML void elements",
        },
    }
    
    fmt.Println("Testing error handling and recovery:")
    fmt.Println(strings.Repeat("=", 50))
    
    for i, attempt := range problemDocuments {
        fmt.Printf("\nTest %d: %s\n", i+1, attempt.Name)
        fmt.Printf("Content: %s\n", attempt.Content)
        
        result := safeParseWithRecovery(attempt.Content, attempt.Config)
        
        if result.Success {
            fmt.Printf("✅ Success: Parsed %d elements\n", result.ElementCount)
        } else {
            fmt.Printf("❌ Failed: %s\n", result.Error)
            
            if result.RecoveryAttempted {
                fmt.Printf("🔄 Recovery attempted with result: %s\n", result.RecoveryResult)
            }
        }
    }
}

type ParseResult struct {
    Success           bool
    Error             string
    ElementCount      int
    RecoveryAttempted bool
    RecoveryResult    string
}

func safeParseWithRecovery(content string, config *markit.ParserConfig) ParseResult {
    result := ParseResult{}
    
    // First attempt with original configuration
    parser := markit.NewParserWithConfig(content, config)
    doc, err := parser.Parse()
    
    if err != nil {
        result.Error = err.Error()
        
        // Attempt recovery with more lenient configuration
        result.RecoveryAttempted = true
        recoveryResult := attemptRecovery(content, err)
        result.RecoveryResult = recoveryResult.Message
        
        if recoveryResult.Success {
            result.Success = true
            result.ElementCount = recoveryResult.ElementCount
        }
    } else {
        result.Success = true
        result.ElementCount = countElements(doc.Root)
    }
    
    return result
}

type RecoveryResult struct {
    Success      bool
    Message      string
    ElementCount int
}

func attemptRecovery(content string, originalError error) RecoveryResult {
    // Try different recovery strategies
    strategies := []struct {
        name   string
        config *markit.ParserConfig
        modify func(string) string
    }{
        {
            name:   "Lenient HTML config",
            config: createLenientHTMLConfig(),
            modify: func(s string) string { return s },
        },
        {
            name:   "Auto-close tags",
            config: markit.DefaultConfig(),
            modify: autoCloseTags,
        },
        {
            name:   "Strip problematic content",
            config: markit.DefaultConfig(),
            modify: stripProblematicContent,
        },
    }
    
    for _, strategy := range strategies {
        modifiedContent := strategy.modify(content)
        parser := markit.NewParserWithConfig(modifiedContent, strategy.config)
        doc, err := parser.Parse()
        
        if err == nil {
            return RecoveryResult{
                Success:      true,
                Message:      fmt.Sprintf("Recovered using strategy: %s", strategy.name),
                ElementCount: countElements(doc.Root),
            }
        }
    }
    
    return RecoveryResult{
        Success: false,
        Message: "All recovery strategies failed",
    }
}

func createStrictConfig() *markit.ParserConfig {
    config := markit.DefaultConfig()
    config.AllowSelfCloseTags = false
    return config
}

func createLenientHTMLConfig() *markit.ParserConfig {
    config := markit.HTMLConfig()
    config.CaseSensitive = false
    return config
}

func autoCloseTags(content string) string {
    // Simple auto-close strategy - wrap in a root element
    if !strings.HasPrefix(strings.TrimSpace(content), "<") {
        return content
    }
    
    // If content doesn't end with a closing tag, try to wrap it
    trimmed := strings.TrimSpace(content)
    if !strings.HasSuffix(trimmed, ">") {
        return "<root>" + content + "</root>"
    }
    
    return content
}

func stripProblematicContent(content string) string {
    // Remove unclosed attributes
    result := content
    
    // Fix unclosed quotes in attributes
    lines := strings.Split(result, "\n")
    for i, line := range lines {
        if strings.Contains(line, `="`) && !strings.Contains(line, `">`) {
            // Try to close unclosed attribute quotes
            if idx := strings.LastIndex(line, `="`); idx != -1 {
                if !strings.Contains(line[idx:], `"`) || strings.Count(line[idx:], `"`) == 1 {
                    lines[i] = line + `"`
                }
            }
        }
    }
    
    return strings.Join(lines, "\n")
}

func countElements(element *markit.Element) int {
    count := 1
    for _, child := range element.Children {
        if childEl, ok := child.(*markit.Element); ok {
            count += countElements(childEl)
        }
    }
    return count
}
```

---

## Next Steps

- [🏠 Back to Home](/) - Return to the main documentation
- [🚀 Getting Started](getting-started) - Learn the basics
- [⚙️ Configuration](configuration) - Master parser configuration
- [📚 API Reference](api-reference) - Complete API documentation

---

<div align="center">

**[📋 Report Issues](https://github.com/khicago/markit/issues)** • **[💬 Discussions](https://github.com/khicago/markit/discussions)** • **[🤝 Contributing](https://github.com/khicago/markit/blob/main/CONTRIBUTING.md)**

</div> 