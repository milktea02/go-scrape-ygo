package scraper

import "github.com/gocolly/colly"

type Scraper interface {
	Scrape(cardName string) (htmlBody *colly.HTMLElement, err error)
}
