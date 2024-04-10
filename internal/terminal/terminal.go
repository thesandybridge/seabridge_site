package terminal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

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

func CommandHandler(w http.ResponseWriter, r *http.Request) {
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
