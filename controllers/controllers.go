package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mtrentz/Website-Text-Scraper/scraper"
)

func ScrapeUrl(c *gin.Context) {
	payload := struct {
		Url         string `json:"url"`
		Depth       int    `json:"depth"`
		MaxRequests int    `json:"max_requests"`
	}{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Since I'm not allowing for inifnite requests.
	var depth int
	var maxRequests int

	defaultDepth := 2
	defaultMaxRequests := 200

	// If depth or max requests were not set, they will be 0 in payload.
	// If so, set them to the default values.
	if payload.Depth == 0 {
		depth = defaultDepth
	}

	if payload.MaxRequests == 0 {
		maxRequests = defaultMaxRequests
	}

	// If they were set to -1, means that the client wants unlimited.
	// If thats the case, set the variables to 0, which in gocolly means unlimited.
	if payload.Depth == -1 {
		depth = 0
	}

	if payload.MaxRequests == -1 {
		maxRequests = 0
	}

	// Now, if they were set to more than 0 in the payload, use them.
	if payload.Depth > 0 {
		depth = payload.Depth
	}

	if payload.MaxRequests > 0 {
		maxRequests = payload.MaxRequests
	}

	websiteReport := scraper.CrawlWebsite(payload.Url, depth, maxRequests)

	// Return the website url
	c.JSON(http.StatusOK, websiteReport)
}
