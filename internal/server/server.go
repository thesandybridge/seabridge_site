package server

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sbxblog/internal/markdown"
)

var (
	indexTemplate        = template.Must(template.ParseFiles("templates/index.html", "templates/base.html", "templates/nav.html"))
	blogTemplate         = template.Must(template.ParseFiles("templates/blog.html", "templates/base.html", "templates/nav.html"))
	postTemplate         = template.Must(template.ParseFiles("templates/post.html", "templates/base.html", "templates/nav.html"))
	pageNotFoundTemplate = template.Must(template.ParseFiles("templates/404.html", "templates/base.html", "templates/nav.html"))
)

func StartServer() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/404", errorHandler)

	log.Println("Listening on http://localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title       string
		Description string
		Classes     string
		Date        string
		Path        string
	}{
		Title:       "sbx blog",
		Description: "A simple blog built with Go and Markdown.",
		Classes:     "home",
		Date:        time.Now().Format("2006"),
		Path:        r.URL.Path,
	}

	if err := indexTemplate.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	if err := pageNotFoundTemplate.ExecuteTemplate(w, "base", nil); err != nil {
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

	data := struct {
		Title           string
		Description     string
		Posts           []BlogPostSummary
		Classes         string
		ContentTemplate string
		Date            string
		Path            string
	}{
		Title:           "Your Blog Title",
		Description:     "Your Blog Description",
		Posts:           posts,
		Classes:         "blog",
		ContentTemplate: "blog.html",
		Date:            time.Now().Format("2006"),
		Path:            r.URL.Path,
	}

	if err := blogTemplate.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveBlogPost(w http.ResponseWriter, r *http.Request, slug string) {
	markdownFile := filepath.Join("content", slug+".md")
	if _, err := os.Stat(markdownFile); os.IsNotExist(err) {
		errorHandler(w, r)
		return
	}
	htmlContent, err := markdown.ConvertToHTML(markdownFile)
	if err != nil {
		http.Error(w, "Could not load blog post", http.StatusNotFound)
		return
	}

	data := struct {
		Title           string
		Content         template.HTML
		Description     string
		Classes         string
		ContentTemplate string
		Date            string
		Path            string
	}{
		Title:           strings.ReplaceAll(slug, "-", " "),
		Content:         template.HTML(htmlContent),
		Description:     "Your Blog Description",
		Classes:         "post",
		ContentTemplate: "post.html",
		Date:            time.Now().Format("2006"),
		Path:            r.URL.Path,
	}

	if err := postTemplate.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
