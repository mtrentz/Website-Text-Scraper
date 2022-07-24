package scraper

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/mtrentz/Website-Text-Scraper/logging"
)

var mu = &sync.Mutex{}

func CrawlWebsite(websiteUrl string, depth int, maxRequests int) WebsiteResult {
	// If not set depth and maxRequests, default to depth = 2 without maxRequests.
	// Since I'm not allowing for inifnite requests.
	if depth == 0 && maxRequests == 0 {
		depth = 2
		maxRequests = 0
	} else if depth == 0 {
		// If depth not set, set it to 0 which means infinite for gocolly
		depth = 0
	} else if maxRequests == 0 {
		// If maxRequests set it to 0, which means inifnite
		maxRequests = 0
	}

	requestCount := 0

	// Instantiate the website result
	websiteResult := &WebsiteResult{
		Url:        websiteUrl,
		PageAmount: 0,
		VisitedAt:  time.Now().Format("2006-01-02 15:04:05"),
		Pages:      []PageResult{},
	}

	// Get the domain of the website, since it will not be allowed to scrape
	// other domains.
	domain, err := parseDomain(websiteUrl)
	if err != nil {
		logging.Logger.Fatal(err)
	}

	// Log that crawler started on website
	logging.Logger.Printf("Crawler started on website %s with a depth of %d and a maxRequests of %d", websiteUrl, depth, maxRequests)

	// Create and configure colly collector
	c := colly.NewCollector()
	c.MaxDepth = depth
	c.AllowedDomains = []string{domain}
	c.AllowURLRevisit = false
	c.Async = true
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4})

	c.OnHTML("html body", func(e *colly.HTMLElement) {
		// Page URL (since it recursively goes into all hrefs)
		pageUrl := e.Request.URL.String()
		// Get the HTML string
		pageHtml, _ := e.DOM.Html()
		// Parse only the text
		pageText, err := ParseHtmlText(pageHtml)
		if err != nil {
			// This will exit this function, but not the crawler
			logging.Logger.Printf("Error parsing HTML text: %s\n", err)
			return
		}

		// Create a pageResult
		pageResult := PageResult{
			Url:       pageUrl,
			Text:      pageText,
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
			if requestCount >= maxRequests {
				fmt.Println("Stoping request")
				r.Abort()
			}
			fmt.Println("Requesting", r.URL)
			requestCount++
			mu.Unlock()
		}

		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(websiteUrl)

	c.Wait()

	logging.Logger.Printf("Finished crawling website %s\n", websiteUrl)

	return *websiteResult
}

func parseDomain(website string) (domain string, err error) {
	// Parse URL
	u, err := url.Parse(website)
	if err != nil {
		return "", err
	}
	domain = u.Hostname()
	return domain, nil
}
