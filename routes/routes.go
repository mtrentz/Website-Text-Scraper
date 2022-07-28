package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mtrentz/Website-Text-Scraper/controllers"
)

func HandleRequest() {
	r := gin.Default()
	r.POST("/api/scrape/", controllers.ScrapeUrl)
	r.Run(":8080")
}
