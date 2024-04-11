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
	"sbxblog/internal/terminal"
)

func parseTemplates(templateFiles ...string) *template.Template {
	baseTemplate := "templates/base.html"
	navTemplate := "templates/nav.html"
	terminalTemplate := "templates/terminal.html"

	files := append([]string{baseTemplate, navTemplate, terminalTemplate}, templateFiles...)
	return template.Must(template.ParseFiles(files...))
}

var (
	indexTemplate        = parseTemplates("templates/index.html")
	blogTemplate         = parseTemplates("templates/blog.html")
	postTemplate         = parseTemplates("templates/post.html")
	pageNotFoundTemplate = parseTemplates("templates/404.html")
)

func StartServer() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/404", errorHandler)
	http.HandleFunc("/commands", terminal.CommandHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("Listening on http://localhost:" + port + "...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

type BlogPostSummary struct {
	Title string
	Slug  string
}

type Page struct {
	Title           string
	Description     string
	Posts           *[]BlogPostSummary
	Content         *template.HTML
	Classes         string
	ContentTemplate string
	Date            string
	Path            string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := Page{
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

	data := Page{
		Title:           "The Blog",
		Description:     "This is THE blog",
		Posts:           &posts,
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

	content := template.HTML(htmlContent)

	data := Page{
		Title:           strings.ReplaceAll(slug, "-", " "),
		Content:         &content,
		Description:     strings.ReplaceAll(slug, "-", " "),
		Classes:         "post",
		ContentTemplate: "post.html",
		Date:            time.Now().Format("2006"),
		Path:            r.URL.Path,
	}

	if err := postTemplate.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
