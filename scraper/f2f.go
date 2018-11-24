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

func (f2fScraper *F2FScraper) Scrape(cardName string) (products []*product.Info, err error) {

	log.Printf("Beginning the scrape on F2F for card '%s'", cardName)

	log.Printf("Getting the html body")
	htmlBody, err := f2fScraper.getHTMLBody(cardName)
	if err != nil {
		return nil, err
	}
	log.Printf("Got the htmlBody: '%s'", htmlBody.Name)
	log.Printf("processing the htmlBody")
	products, err = f2fScraper.processBody(htmlBody)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (*F2FScraper) getHTMLBody(cardName string) (htmlBody *colly.HTMLElement, err error) {

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

func (*F2FScraper) processBody(htmlBody *colly.HTMLElement) (products []*product.Info, err error) {

	products = []*product.Info{}

	htmlBody.ForEach("table.products_table", func(_ int, e *colly.HTMLElement) {
		log.Printf("Found the table with class name products_table: '%s'", e.Name)
		e.ForEach("td.meta", func(_ int, tdMeta *colly.HTMLElement) {
			productName := tdMeta.ChildText("a.name")
			productSet := tdMeta.ChildText("small.category")
			variants := []*product.Variant{}

			tdMeta.ForEach("tr.variantRow", func(_ int, trVariantRow *colly.HTMLElement) {
				variantCondition := trVariantRow.ChildText("td.variantInfo")
				variantPriceString := strings.Trim(trVariantRow.ChildText("td.price"), "CAD$ ")
				variantPriceString = strings.Split(variantPriceString, " ")[0]
				variantPrice, err := strconv.ParseFloat(variantPriceString, 64)
				if err != nil {
					log.Printf("Error wile parsing string to float: '%s', - '%s'", variantPriceString, err)
					variantPrice = float64(0)
				}
				variantQuantityString := strings.Trim(trVariantRow.ChildText("td:nth-child(3)"), "x ")
				variantQuantity, err := strconv.ParseInt(variantQuantityString, 0, 64)
				if err != nil {
					log.Printf("Error wile parsing string to int: '%s', - '%s'", variantQuantityString, err)
					variantQuantity = 0
				}

				variants = append(variants, &product.Variant{
					Condition: variantCondition,
					Price:     variantPrice,
					Quantity:  variantQuantity,
				})
			})
			// This usually means theres none in stock so get the "speculated" pricing
			if len(variants) == 0 {
				variantCondition := tdMeta.ChildText("td.variantInfo")
				variantPriceString := tdMeta.ChildText("table > tbody > tr > td:nth-child(2)")
				variantPriceString = strings.Trim(variantPriceString, "CAD$ ")
				variantPrice, err := strconv.ParseFloat(variantPriceString, 64)
				if err != nil {
					log.Printf("Error wile parsing string to float: '%s', '%s'", variantPriceString, err)
					variantPrice = float64(0)
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

	return products, nil
}
