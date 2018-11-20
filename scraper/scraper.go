package scraper

import _ "github.com/gocolly/colly"
import "github.com/milktea02/go-scrape-ygo/product"

type Scraper interface {
	Scrape(cardName string) (products []*product.Info, err error)
}
