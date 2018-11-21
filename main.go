package main

import _ "github.com/gocolly/colly"
import "github.com/milktea02/go-scrape-ygo/scraper"
import "github.com/milktea02/go-scrape-ygo/product"
import (
	"fmt"
	_ "io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r.URL.Path[1:])
	cardName := strings.Replace(r.URL.Path[1:], " ", "+", -1)
	f2fScraper := &scraper.F2FScraper{}
	cards, err := f2fScraper.Scrape(cardName)
	if err != nil {
		log.Printf("Error: '%s'", err)
	}
	fmt.Fprintf(w, "hi there, I love %s!\n", r.URL.Path[1:])
	printCardInfoForWeb(w, cards)
}

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	/*
		args := os.Args[1:]
		cardName := parseArgs(args)
		f2fScraper := &scraper.F2FScraper{}

		cards, err := f2fScraper.Scrape(cardName)
		if err != nil {
			log.Printf("Error: '%s'", err)
		}
		printCardInfo(cards)
	*/
}

func printCardInfoForWeb(w http.ResponseWriter, cards []*product.Info) {
	for _, card := range cards {
		fmt.Fprintf(w, "Card Name: %s \t Card Set: %s\n", card.Name, card.Set)
		for _, variant := range card.Variants {
			fmt.Fprintf(w, "%s \t %.2f \t %d\n", variant.Condition, variant.Price, variant.Quantity)
		}
	}

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
