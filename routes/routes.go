package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mtrentz/Website-Text-Scraper/controllers"
)

func HandleRequest() {
	r := gin.Default()
	// r.POST("/api/v1/scrape", controllers.ScrapeWebsite)
	// r.GET("/api/v1/classify/:id", controllers.ClassifyWebsite)
	r.GET("/api/v1/hello", controllers.HelloController)
	r.POST("/api/v1/scrape", controllers.ScrapeUrl)
	r.Run()
}
