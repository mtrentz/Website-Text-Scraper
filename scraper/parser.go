package scraper

import (
	"io"
	"strings"

	"github.com/mtrentz/Website-Text-Scraper/logging"
	"golang.org/x/net/html"
)

func ParseHtmlText(pageHtml string) (pageText string, err error) {
	textTags := []string{
		"a",
		"p", "span", "em", "string", "blockquote", "q", "cite",
		"h1", "h2", "h3", "h4", "h5", "h6",
	}

	tag := ""
	enter := false

	var text string

	r := strings.NewReader(pageHtml)
	tokenizer := html.NewTokenizer(r)
	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			logging.Logger.Fatal(err)
		case html.StartTagToken, html.SelfClosingTagToken:
			enter = false

			tag = token.Data
			for _, ttt := range textTags {
				if tag == ttt {
					enter = true
					break
				}
			}
		case html.TextToken:
			if enter {
				data := strings.TrimSpace(token.Data)

				if len(data) > 0 {
					text += data + "\n"
				}
			}
		}
	}

	return text, nil
}
