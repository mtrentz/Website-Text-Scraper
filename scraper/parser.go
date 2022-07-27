package scraper

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func removeHeadersAndFooters(doc *goquery.Document) (rest *goquery.Document, headers []*goquery.Selection, footers []*goquery.Selection) {
	headers = make([]*goquery.Selection, 0)
	footers = make([]*goquery.Selection, 0)

	// Find the header element, append to headers and remove it
	doc.Find("header, [class*=header], [id*=header]").Each(func(i int, s *goquery.Selection) {
		headers = append(headers, s)
		s.Remove()
	})

	// Find the footer element, append to footers and remove it
	doc.Find("footer, [class*=footer], [id*=footer]").Each(func(i int, s *goquery.Selection) {
		footers = append(footers, s)
		s.Remove()
	})

	// return the document and the headers and footers
	return doc, headers, footers
}
