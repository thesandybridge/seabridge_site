package server

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"sbxblog/internal/markdown"
)

var (
	tmplIndex    = template.Must(template.ParseFiles("templates/index.html"))
	tmplBlogList = template.Must(template.ParseFiles("templates/blog.html"))
	tmplBlogPost = template.Must(template.ParseFiles("templates/post.html"))
)

func StartServer() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/blog/", blogHandler)

	log.Println("Listening on http://localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Message string
	}{
		Message: "Welcome to the Go webserver!",
	}

	if err := tmplIndex.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type BlogPostSummary struct {
	Title string
	Slug  string
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/blog/")
	if slug == "" {
		serveBlogList(w, r)
	} else {
		serveBlogPost(w, r, slug)
	}
}

func serveBlogList(w http.ResponseWriter, r *http.Request) {
	var posts []BlogPostSummary
	err := filepath.WalkDir("content", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			slug := strings.TrimSuffix(d.Name(), ".md")
			title := strings.ReplaceAll(slug, "-", " ")
			posts = append(posts, BlogPostSummary{Title: title, Slug: slug})
		}
		return nil
	})
	if err != nil {
		log.Println("Failed to walk content directory:", err)
		http.Error(w, "Failed to load blog list", http.StatusInternalServerError)
		return
	}

	if err := tmplBlogList.Execute(w, posts); err != nil {
		http.Error(w, "Error rendering blog list", http.StatusInternalServerError)
	}
}

func serveBlogPost(w http.ResponseWriter, r *http.Request, slug string) {
	markdownFile := filepath.Join("content", slug+".md")
	htmlContent, err := markdown.ConvertToHTML(markdownFile)
	if err != nil {
		http.Error(w, "Could not load blog post", http.StatusNotFound)
		return
	}

	data := struct {
		Title   string
		Content template.HTML
	}{
		Title:   strings.ReplaceAll(slug, "-", " "),
		Content: template.HTML(htmlContent),
	}

	if err := tmplBlogPost.Execute(w, data); err != nil {
		http.Error(w, "Error rendering blog post", http.StatusInternalServerError)
	}
}
