package logging

import (
	"log"
	"os"
)

// Define global logger
var Logger *log.Logger

func init() {
	// Set up logger to file
	logFile, err := os.OpenFile("scraper.logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	Logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

}
