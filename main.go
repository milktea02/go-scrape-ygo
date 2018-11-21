package main

import _ "github.com/gocolly/colly"
import "github.com/milktea02/go-scrape-ygo/scraper"
import "github.com/milktea02/go-scrape-ygo/product"
import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	args := os.Args[1:]
	cardName, _ := parseArgs(args)
	f2fScraper := &scraper.F2FScraper{}

	cards, err := f2fScraper.Scrape(cardName)
	if err != nil {
		log.Printf("Error: '%s'", err)
	}
	printCardInfo(cards)
}

func printCardInfo(cards []*product.Info) {
	for _, card := range cards {
		fmt.Printf("Card Name: %s \t Card Set: %s\n", card.Name, card.Set)
		for _, variant := range card.Variants {
			fmt.Printf("%s \t %.2f \t %d\n", variant.Condition, variant.Price, variant.Quantity)
		}
	}

}

func parseArgs(args []string) (cardName string) {

	for i, arg := range args {
		args[i] = strings.TrimFunc(strings.ToLower(arg), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
	}

	cardName = strings.Join(args, "+")
	return cardName
}
