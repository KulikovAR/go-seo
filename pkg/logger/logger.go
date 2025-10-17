package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	logFile := fmt.Sprintf("logs/track_site_%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogTrackSiteParams(siteID int, source, device, os string, ads bool, country, lang string, pages int, subdomains bool) {
	InfoLogger.Printf("TrackSite request - SiteID: %d, Source: %s, Device: %s, OS: %s, Ads: %t, Country: %s, Lang: %s, Pages: %d, Subdomains: %t",
		siteID, source, device, os, ads, country, lang, pages, subdomains)
}

func LogXMLRiverURL(url string, params map[string]string) {
	InfoLogger.Printf("XMLRiver request - URL: %s, Params: %+v", url, params)
}

func LogXMLRiverResponse(statusCode int, responseBody string) {
	if statusCode != 200 {
		ErrorLogger.Printf("XMLRiver error response - Status: %d, Body: %s", statusCode, responseBody)
	} else {
		InfoLogger.Printf("XMLRiver success response - Status: %d", statusCode)
	}
}
