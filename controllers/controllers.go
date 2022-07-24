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

	websiteReport := scraper.CrawlWebsite(payload.Url, payload.Depth, payload.MaxRequests)

	// Return the website url
	c.JSON(http.StatusOK, websiteReport)
}
