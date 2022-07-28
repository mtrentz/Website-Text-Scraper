package scraper

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/mtrentz/Website-Text-Scraper/logging"
)

var mu = &sync.Mutex{}

// Crawls many pages of a website.
// The crawler is constrained to the domains of the website.
// Passing a max requests doesn't guarantee that the crawler will return that many amount of pages.
func CrawlWebsite(websiteUrl string, depth int, maxRequests int) WebsiteResult {
	requestCount := 0

	// Instantiate the website result
	websiteResult := &WebsiteResult{
		Url:        websiteUrl,
		PageAmount: 0,
		VisitedAt:  time.Now().Format("2006-01-02 15:04:05"),
		Pages:      []PageResult{},
	}

	// Get the domain of the website, since it will not be allowed to scrape
	// other domains. Domains is a slice because it adds 'www' as a variation.
	domains, err := getAllowedDomains(websiteUrl)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	// Log that crawler started on website
	logging.Logger.Printf("Crawler started on website %s with a depth of %d and a maxRequests of %d", websiteUrl, depth, maxRequests)

	// Create and configure colly collector
	c := colly.NewCollector()
	c.MaxDepth = depth
	c.AllowedDomains = domains
	c.AllowURLRevisit = false
	c.Async = true
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})

	c.OnHTML("html body", func(e *colly.HTMLElement) {
		// Page URL (since it recursively goes into all hrefs)
		pageUrl := e.Request.URL.String()
		// Get the HTML string
		pageHtml, _ := e.DOM.Html()

		// Put the string into a Reader so I can parse into a goquery document
		reader := strings.NewReader(pageHtml)
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			logging.Logger.Fatal(err)
		}

		// Separate the headers, footers and the bodu (rest) of the html
		rest, headers, footers := removeHeadersAndFooters(doc)

		// Get the text from the rest of the html
		restHtml, err := rest.Html()
		if err != nil {
			logging.Logger.Fatalf("Error parsing rest Html for page %s: %e", pageUrl, err)
		}

		restText, err := ParseHtmlText(restHtml)
		if err != nil {
			logging.Logger.Fatalf("Error parsing rest Text for page %s: %e", pageUrl, err)
		}

		// Get the text of the footers and headers
		var headersText string
		var footersText string

		for _, header := range headers {
			headerHtml, err := header.Html()
			if err != nil {
				logging.Logger.Fatalf("Error parsing header Html for page %s: %e", pageUrl, err)
			}
			headerText, err := ParseHtmlText(headerHtml)
			if err != nil {
				logging.Logger.Fatalf("Error parsing header Text for page %s: %e", pageUrl, err)
			}
			headersText += headerText + "\n"
		}
		for _, footer := range footers {
			footerHtml, err := footer.Html()
			if err != nil {
				logging.Logger.Fatalf("Error parsing footer Html for page %s: %e", pageUrl, err)
			}
			footerText, err := ParseHtmlText(footerHtml)
			if err != nil {
				logging.Logger.Fatalf("Error parsing footer Text for page %s: %e", pageUrl, err)
			}
			footersText += footerText + "\n"
		}

		// Create a pageResult
		pageResult := PageResult{
			Url:       pageUrl,
			Header:    headersText,
			Body:      restText,
			Footer:    footersText,
			VisitedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// Log saying that page was scraped
		logging.Logger.Printf("Scraped page: %s\n", pageUrl)

		// Add the pageResult to the websiteResult
		websiteResult.AddPage(pageResult)
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		// Only check for max requests if is not set to 0
		if maxRequests != 0 {
			// Before request, check if already sent the max amount
			// of requests, if not, increment counter and continue
			mu.Lock()
			if requestCount > maxRequests {
				logging.Logger.Printf("Stoping request because of max requests")
				r.Abort()
			}
			fmt.Println("Requesting", r.URL)
			requestCount++
			mu.Unlock()
		}
	})

	c.Visit(websiteUrl)

	c.Wait()

	logging.Logger.Printf("Finished crawling website %s\n", websiteUrl)

	return *websiteResult
}

func getAllowedDomains(website string) (domains []string, err error) {
	// Parse URL
	u, err := url.Parse(website)
	if err != nil {
		return nil, err
	}
	domains = append(domains, u.Hostname())
	// Add the variation with www
	domains = append(domains, "www."+u.Hostname())
	return domains, nil
}
