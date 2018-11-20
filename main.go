package main

import _ "github.com/gocolly/colly"
import "github.com/milktea02/go-scrape-ygo/scraper"
import (
	"log"
)

func main() {
	f2fScraper := &scraper.F2FScraper{}

	htmlBody, err := f2fScraper.Scrape("mystical+space+typhoon")
	if err != nil {
		log.Printf("Error: '%s'", err)
	}

	log.Printf("HTML Body: '%+v'", htmlBody)

}
