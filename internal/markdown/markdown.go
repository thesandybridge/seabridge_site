package markdown

import (
	"os"

	"github.com/russross/blackfriday/v2"
)

func ConvertToHTML(filePath string) (string, error) {
	markdownContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	htmlContent := blackfriday.Run(markdownContent)
	return string(htmlContent), nil
}
