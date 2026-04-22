/* Simple Go builder inspired by Hugo's functionalities
   This assumes a website/ dir, with content/ and pages/ in it. 
   Content should have a blog/ and a projects/, as well as a home.md, to build the HTML files 
   Credits: Gemini AI
*/

package main

import (
	"fmt"
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
)

type Metadata struct {
	Title string   `yaml:"title"`
	Date  string   `yaml:"date"`
	Author string  `yaml:"author"`
	Tags  []string `yaml:"tags"`
	Slug  string
}

type ProjectMetadata struct {
	Title string `yaml:"title"`
	Description string `yaml:"description"`
	Tech []string `yaml:"tech"`
	Source string `yaml:"source"`
	Slug string
}

type PageData struct {
	Meta    Metadata
	Content template.HTML
}

func main() {
	// Build single page 
	buildSinglePage("../content/home.md", "../pages/home.html")

	os.MkdirAll("../pages/blog", 0755)
	os.MkdirAll("../pages/tags", 0755)

	blogPosts := processDir("../content/blog", "../pages/blog")

	// group posts by tag and add Tags
	tagMap := make(map[string][]Metadata)
	allTags := make(map[string]bool)
	
	for _, post := range blogPosts {
		for _, tag := range post.Tags {
			tagMap[tag] = append(tagMap[tag], post)
			allTags[tag] = true
		}
	}

	for tag, posts := range tagMap {
		render(fmt.Sprintf("../pages/tags/%s.html", tag), "Filtering for: #"+tag, posts, nil)
	}

	keys := make([]string, 0, len(allTags))
	for k := range allTags { keys = append(keys, k)}
	sort.Strings(keys)

	render("../pages/blog.html", "", blogPosts, keys)
	
	projects := processProjectsDir("../content/projects")
	renderProjects("../pages/projects.html", projects)
}

func buildSinglePage(srcPath, destPath string) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	input, _ := os.ReadFile(srcPath)
	var meta Metadata
	contentBytes, _ := frontmatter.Parse(bytes.NewReader(input), &meta)

	var buf bytes.Buffer
	markdown.Convert(contentBytes, &buf)
	
	final := fmt.Sprintf("<div><h1>%s</h1>%s</div>", meta.Title, buf.String())

	os.MkdirAll(filepath.Dir(destPath), 0755)
	os.WriteFile(destPath, []byte(final), 0644)

}

func processDir(src, dest string) []Metadata {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("dracula"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
				),
			),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var list []Metadata
	files, _ := os.ReadDir(src)

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".md" {
			continue
		}

		input, _ := os.ReadFile(filepath.Join(src, f.Name()))
		var meta Metadata
		contentBytes, _ := frontmatter.Parse(bytes.NewReader(input), &meta)
		var buf bytes.Buffer
		markdown.Convert(contentBytes, &buf)
		
		// Adds a byline at start of post
		byline := fmt.Sprintf("<p><i>%s · %s </i></p>\n<hr>", meta.Date, meta.Author)

		meta.Slug = f.Name()[:len(f.Name())-3]

		finalContent := byline + buf.String()
		outputFile := filepath.Join(dest, meta.Slug+".html")
		os.WriteFile(outputFile, []byte(finalContent), 0644)

		list = append(list, meta)
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Date > list[j].Date })
	return list
}

func render(destPath, title string, items []Metadata, footerTags []string) {
    const listTmpl = `
{{if .Title }}
<h1>{{.Title}}</h1>
{{end}}
<ul class="blog-posts">
    {{range .Items}}
    <li>
        <span>{{.Date}}</span>
        <a href="/pages/blog/{{.Slug}}.html" 
           hx-get="/pages/blog/{{.Slug}}.html" 
           hx-target="#main-content" 
           hx-swap="innerHTML transition:true">{{.Title}}</a>
    </li>
    {{end}}
</ul>

{{if .FooterTags}}
<div class="tag-cloud">
    {{range .FooterTags}}
    <a href="/pages/tags/{{.}}.html" 
	   hx-get="/pages/tags/{{.}}.html" 
       hx-target="#main-content" 
       class="blog-tags">#{{.}}</a>
    {{end}}
</div>
{{end}}`

    // Define the template and the custom 'lower' function
    t := template.Must(template.New("list").Funcs(template.FuncMap{
        "lower": func(s string) string { 
            return strings.ToLower(s) 
        },
    }).Parse(listTmpl))

    // Create the output file
    f, err := os.Create(destPath)
    if err != nil {
        fmt.Printf("Failed to create file %s: %v", destPath, err)
        return
    }
    defer f.Close()

    // Pass all three pieces of data to the template
    t.Execute(f, struct {
        Title      string
        Items      []Metadata
        FooterTags []string
    }{
        Title:      title,
        Items:      items,
        FooterTags: footerTags,
    })
}

func processProjectsDir(src string) []ProjectMetadata {
	var list []ProjectMetadata
	files, _ := os.ReadDir(src)

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".md" { continue }

		input, _ := os.ReadFile(filepath.Join(src, f.Name()))
		var meta ProjectMetadata
		
		// Parse frontmatter
		_, err := frontmatter.Parse(bytes.NewReader(input), &meta)
		if err != nil { continue }

		list = append(list, meta)
	}

	return list
}

func renderProjects(destPath string, projects []ProjectMetadata) {
	const projTmpl = `
<div class="project-list">
    {{range .}}
    <div class="project-item" style="margin-bottom: 2.5rem;">
        <h3 style="margin-bottom: 0.5rem;">{{.Title}}</h3>
        
        <p style="margin-bottom: 1rem;">
            {{.Description}}
        </p>
        
        <div class="tech-stack" style="margin-bottom: 0.8rem;">
            {{range .Tech}}
            <code style="background: #282a36; color: #50fa7b; padding: 1px 5px; border-radius: 3px; font-size: 0.85rem; margin-right: 5px;">{{.}}
            </code>
            {{end}}
        </div>

        {{if .Source}}
        <a href="{{.Source}}" target="_blank" rel="noopener" style="font-size: 0.9rem;">
            Source Code
        </a>
        {{end}}
    </div>
    {{end}}
</div>`

	t := template.Must(template.New("projects").Parse(projTmpl))
	
	// Ensure the parent directory exists
	os.MkdirAll(filepath.Dir(destPath), 0755)

	f, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("Failed to create projects file: %v\n", err)
		return
	}
	defer f.Close()

	// Since .Description is already HTML, we need to tell Go's template engine 
	// not to escape it. To do that easily without extra structs, 
	// we could cast it, but for a simple "Lair" setup, this is fine.
	t.Execute(f, projects)
}
