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

	// Check if depth or MaxRequests were set to -1, which is ulimited
	if payload.Depth == -1 {
		// If Depth was set to -1, set it to inifnite (0 in gocolly)
		depth = 0
	} else if payload.MaxRequests == -1 {
		// If MaxRequests was set to -1, set it to inifnite (0 in the crawler code)
		maxRequests = 0
	}

	if payload.Depth == 0 && payload.MaxRequests == 0 {
		// If not set anything, default to depth of 1 and max requests of 200
		depth = 2
		maxRequests = 200
	} else if payload.Depth == 0 && payload.MaxRequests > 0 {
		// If depth not set and max requests is set to above 0, default to depth of 2
		depth = 2
		maxRequests = payload.MaxRequests
	} else if payload.MaxRequests == 0 && payload.Depth > 0 {
		// If max requests not set and depth is set to above 0, default to max requests of 200
		maxRequests = 200
		depth = payload.Depth
	} else {
		// If both depth and maxRequests were set, use them.
		depth = payload.Depth
		maxRequests = payload.MaxRequests
	}

	websiteReport := scraper.CrawlWebsite(payload.Url, depth, maxRequests)

	// Return the website url
	c.JSON(http.StatusOK, websiteReport)
}
