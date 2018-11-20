package scraper

import "github.com/gocolly/colly"
import (
	"fmt"
	"github.com/milktea02/go-scrape-ygo/product"
	"log"
	"sort"
	"strconv"
	"strings"
)

type F2FScraper struct {
}

func (*F2FScraper) Scrape(cardName string) (htmlBody *colly.HTMLElement, err error) {

	log.Printf("Scraping F2F for card '%s'", cardName)
	c := colly.NewCollector(
		colly.AllowedDomains("facetofacegames.com", "www.facetofacegames.com"),
	)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		htmlBody = e
	})
	if err := c.Visit("http://www.facetofacegames.com/products/search?query=" + cardName); err != nil {
		return nil, fmt.Errorf("Error while visiting facetoface - '%s'\n", err)

	}

	return htmlBody, nil

}

func (*F2FScraper) ProcessBody(args []string) {
	fmt.Printf("Args: '%+v'", args)
	cardName := strings.Join(args, "+")

	fmt.Printf("Querying for '%s'\n", cardName)

	products := []*product.Info{}
	c := colly.NewCollector(
		colly.AllowedDomains("facetofacegames.com", "www.facetofacegames.com"),
	)

	c.OnHTML("table.products_table", func(e *colly.HTMLElement) {
		e.ForEach("td.meta", func(_ int, tdMeta *colly.HTMLElement) {
			productName := tdMeta.ChildText("a.name")
			productSet := tdMeta.ChildText("small.category")
			variants := []*product.Variant{}

			tdMeta.ForEach("tr.variantRow", func(_ int, trVariantRow *colly.HTMLElement) {
				variantCondition := trVariantRow.ChildText("td.variantInfo")
				variantPriceString := strings.Trim(trVariantRow.ChildText("td.price"), "CAD$ ")
				variantPrice, err := strconv.ParseFloat(variantPriceString, 64)
				if err != nil {
					log.Fatalf("Error wile parsing string to float: '%s'", variantPriceString)
				}
				variants = append(variants, &product.Variant{
					Condition: variantCondition,
					Price:     variantPrice,
				})
			})
			if len(variants) == 0 {
				variantCondition := tdMeta.ChildText("td.variantInfo")
				variantPriceString := tdMeta.ChildText("table > tbody > tr > td:nth-child(2)")
				variantPriceString = strings.Trim(variantPriceString, "CAD$ ")
				variantPrice, err := strconv.ParseFloat(variantPriceString, 64)
				if err != nil {
					log.Fatalf("Error wile parsing string to float: '%s'", variantPriceString)
				}
				variants = append(variants, &product.Variant{
					Condition: variantCondition,
					Price:     variantPrice,
				})
			}
			sort.Slice(variants, func(first, second int) bool {
				return variants[first].Price < variants[second].Price
			})
			products = append(products, &product.Info{
				Name:     productName,
				Set:      productSet,
				Variants: variants,
			})

		})
	})
	sort.Slice(products, func(first, second int) bool {
		return products[first].Variants[0].Price < products[second].Variants[0].Price
	})

	for _, product := range products {
		fmt.Printf("%+v\n", product)
		for _, variant := range product.Variants {
			fmt.Printf("%+v\n", variant)
		}
	}
}
