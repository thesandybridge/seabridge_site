package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"

	"sbxblog/internal/markdown"
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
	http.HandleFunc("/commands", commandHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("Listening on http://localhost:" + port + "...")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

type CommandResponse struct {
	Action  string `json:"action,omitempty"`
	Message string `json:"message,omitempty"`
	URL     string `json:"url,omitempty"`
}

func getItems() []string {
	items := []string{
		"blog",
		"home",
	}

	return items
}

func buildRelativeURL(currentURL string, target string) (string, error) {
	parsedURL, err := url.Parse(currentURL)
	if err != nil {
		return "", err
	}

	appPath := strings.TrimPrefix(parsedURL.Path, "/commands")

	segments := strings.Split(strings.Trim(appPath, "/"), "/")

	switch target {
	case "home", "":
		parsedURL.Path = "/"
	case "..":
		if len(segments) > 1 {
			parsedURL.Path = "/" + strings.Join(segments[:len(segments)-1], "/")
		} else {
			parsedURL.Path = "/"
		}
	default:
		if strings.HasPrefix(target, "/") {
			parsedURL.Path = target
		} else {
			if appPath == "/" {
				parsedURL.Path += target
			} else {
				newTargetPath := strings.Join(segments, "/") + "/" + target
				parsedURL.Path = "/" + newTargetPath
			}
		}
	}

	newURL := parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path

	return newURL, nil
}

func getFullURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	} else if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	}

	host := r.Host

	path := r.URL.Path

	rawQuery := r.URL.RawQuery
	if rawQuery != "" {
		path += "?" + rawQuery
	}

	return scheme + "://" + host + path
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	p := bluemonday.StrictPolicy()
	command := p.Sanitize(r.FormValue("cmd"))
	args := strings.Fields(command)

	var response CommandResponse
	var message string

	items := getItems()

	switch args[0] {
	case "help":
		message = `Commands available:
    help    - Show this help message
    cd      - Navigate to another page
    ls      - List available pages
    clear   - Clear the screen
    github  - Open the GitHub page in a new tab
    echo    - Echo back the input
    contact - Show contact information`
	case "clear":
		response.Action = "clear"
	case "ls":
		message = strings.Join(items, " ")
	case "github":
		response.Action = "open-url"
		response.URL = "https://github.com/thesandybridge"
	case "echo":
		if len(args) > 1 {
			message = strings.Join(args[1:], " ")
		} else {
			message = "echo: no message provided"
		}
	case "contact":
		message = `Contact info:
    email:      matt@mattmillerdev.io
    linkedin:   /in/mattmillerdev/
    github:     /thesandybridge`
	case "cd":
		if len(args) > 1 {
			currentURL := getFullURL(r)
			targetPath := args[1]
			newPath, err := buildRelativeURL(currentURL, targetPath)
			if err != nil {
				message = fmt.Sprintf("Error: %v", err)
			} else {
				response.Action = "navigate"
				response.URL = newPath
			}
		} else {
			message = "cd: path required"
		}
	case "rotate":
		response.Action = "rotate"
	case "malware":
		response.Action = "malware"
	default:
		message = args[0] + ": command not found"
	}

	log.Printf("Args: %s, User-Agent: %s, Remote Address: %s",
		args, r.UserAgent(), r.RemoteAddr)

	if message != "" {
		response.Message = fmt.Sprintf("<pre class='ignore'>&gt; %s\n%s</pre>", strings.Join(args, " "), message)
	}

	if response.Action != "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(response.Message))
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
