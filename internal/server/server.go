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
	indexTemplate        = template.Must(template.ParseFiles("templates/index.html", "templates/base.html", "templates/nav.html"))
	blogTemplate         = template.Must(template.ParseFiles("templates/blog.html", "templates/base.html", "templates/nav.html"))
	postTemplate         = template.Must(template.ParseFiles("templates/post.html", "templates/base.html", "templates/nav.html"))
	pageNotFoundTemplate = template.Must(template.ParseFiles("templates/404.html", "templates/base.html", "templates/nav.html"))
)

func StartServer() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/", catchAllHandler)

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
		Title       string
		Description string
		Classes     string
	}{
		Title:       "sbx blog",
		Description: "A simple blog built with Go and Markdown.",
		Classes:     "home",
	}

	if err := indexTemplate.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		data := struct {
			Title       string
			Description string
			Classes     string
		}{
			Title:       "404 Not Found",
			Description: "The page you're looking for doesn't exist.",
			Classes:     "not-found",
		}

		if err := pageNotFoundTemplate.ExecuteTemplate(w, "base.html", data); err != nil {
			log.Printf("Failed to execute 404 template: %v", err)
			http.Error(w, "404 Not Found", http.StatusNotFound)
		}
	})
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		viewHandler(w, r)
		return
	}

	notFoundHandler().ServeHTTP(w, r)
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
	}{
		Title:           "Your Blog Title",
		Description:     "Your Blog Description",
		Posts:           posts,
		Classes:         "blog",
		ContentTemplate: "blog.html",
	}

	if err := blogTemplate.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		Title           string
		Content         template.HTML
		Description     string
		Classes         string
		ContentTemplate string
	}{
		Title:           strings.ReplaceAll(slug, "-", " "),
		Content:         template.HTML(htmlContent),
		Description:     "Your Blog Description",
		Classes:         "post",
		ContentTemplate: "post.html",
	}

	if err := postTemplate.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
