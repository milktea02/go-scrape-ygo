package main

import _ "github.com/gocolly/colly"
import "github.com/milktea02/go-scrape-ygo/scraper"
import "github.com/milktea02/go-scrape-ygo/product"
import (
	"log"
)

func main() {
	f2fScraper := &scraper.F2FScraper{}

	cards, err := f2fScraper.Scrape("mystical+space+typhoon")
	if err != nil {
		log.Printf("Error: '%s'", err)
	}
	printCardInfo(cards)
}

func printCardInfo(cards []*product.Info) {
	for _, card := range cards {
		log.Printf("Card Name: %s \t Card Set: %s", card.Name, card.Set)
		for _, variant := range card.Variants {
			log.Printf("%s \t %.2f \t %d", variant.Condition, variant.Price, variant.Quantity)
		}
	}

}
