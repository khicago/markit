# MarkIt Documentation

This directory contains the complete documentation for MarkIt parser, built with Jekyll for GitHub Pages.

## 📁 Structure

```
docs/
├── _config.yml              # Jekyll configuration
├── _layouts/
│   └── default.html         # Main layout template
├── Gemfile                  # Ruby dependencies
├── index.md                 # Homepage
├── getting-started.md       # Quick start guide
├── configuration.md         # Configuration options
├── void-elements.md         # Void elements feature
├── api-reference.md         # Complete API documentation
├── examples.md              # Practical examples
├── contributing.md          # Contribution guidelines
├── faq.md                   # Frequently asked questions
└── README.md               # This file
```

## 🚀 Local Development

### Prerequisites

- Ruby 2.7+
- Bundler

### Setup

1. **Install dependencies**:
```bash
cd docs
bundle install
```

2. **Start development server**:
```bash
bundle exec jekyll serve
```

3. **Open in browser**:
```
http://localhost:4000
```

### Building for Production

```bash
bundle exec jekyll build
```

The built site will be in the `_site` directory.

## 📝 Writing Documentation

### Front Matter

Each documentation page should include front matter:

```yaml
---
layout: default
title: "Page Title"
description: "Page description for SEO"
keywords: "relevant, keywords, for, seo"
author: "MarkIt Team"
---
```

### Style Guide

- Use clear, concise language
- Include practical code examples
- Add table of contents for long pages
- Use proper markdown formatting
- Include cross-references between pages

### Code Examples

Use fenced code blocks with language specification:

````markdown
```go
package main

import "github.com/khicago/markit"

func main() {
    parser := markit.NewParser("<div>Hello</div>")
    doc, err := parser.Parse()
    // ...
}
```
````

### Internal Links

Use relative links for internal navigation:

```markdown
[Getting Started](getting-started)
[API Reference](api-reference)
[Examples](examples)
```

## 🎨 Styling

The documentation uses a custom CSS framework built into the default layout:

- **Responsive design** - Works on all devices
- **Syntax highlighting** - Code blocks are highlighted
- **Copy buttons** - Easy code copying
- **Smooth scrolling** - Better navigation experience
- **Mobile menu** - Collapsible navigation on mobile

## 📊 SEO Optimization

The documentation is optimized for search engines:

- **Meta tags** - Title, description, keywords
- **Open Graph** - Social media sharing
- **Twitter Cards** - Twitter sharing
- **Structured data** - Better search results
- **Sitemap** - Automatic generation
- **SEO plugin** - Jekyll SEO tag

## 🔧 Configuration

Key configuration options in `_config.yml`:

```yaml
# Site settings
title: "MarkIt Documentation"
description: "Next-generation extensible markup parser for Go"
url: "https://khicago.github.io"
baseurl: "/markit"

# Build settings
markdown: kramdown
highlighter: rouge
theme: minima

# Plugins
plugins:
  - jekyll-feed
  - jekyll-sitemap
  - jekyll-seo-tag
```

## 📱 Mobile Support

The documentation is fully responsive:

- **Mobile-first design**
- **Touch-friendly navigation**
- **Optimized typography**
- **Fast loading times**

## 🔍 Search

Search functionality can be added using:

- **Jekyll search plugins**
- **Algolia DocSearch**
- **Custom JavaScript search**

## 📈 Analytics

Google Analytics can be enabled by adding:

```yaml
google_analytics: "GA_TRACKING_ID"
```

## 🤝 Contributing

To contribute to the documentation:

1. **Fork the repository**
2. **Create a feature branch**
3. **Make your changes**
4. **Test locally**
5. **Submit a pull request**

### Guidelines

- Follow the existing style and structure
- Test all links and code examples
- Ensure mobile compatibility
- Update table of contents if needed
- Add appropriate front matter

## 📄 License

This documentation is part of the MarkIt project and follows the same license terms.

---

For more information about MarkIt, visit the [main repository](https://github.com/khicago/markit). 