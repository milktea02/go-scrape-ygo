package main

import "github.com/gocolly/colly"
import "log"
import "fmt"
import "strings"
import "strconv"
import "os"
import "sort"

type ProductInfo struct {
	Name     string
	Set      string
	Variants []*VariantInfo
}

type VariantInfo struct {
	Condition string
	Price     float64
}

func main() {

	args := os.Args[1:]
	fmt.Printf("Args: '%+v'", os.Args)
	cardName := strings.Join(args, "+")

	fmt.Printf("Querying for '%s'\n", cardName)

	c := colly.NewCollector(
		colly.AllowedDomains("facetofacegames.com", "www.facetofacegames.com"),
	)

	products := []*ProductInfo{}

	c.OnHTML("table.products_table", func(e *colly.HTMLElement) {
		e.ForEach("td.meta", func(_ int, tdMeta *colly.HTMLElement) {
			productName := tdMeta.ChildText("a.name")
			productSet := tdMeta.ChildText("small.category")
			variants := []*VariantInfo{}
			tdMeta.ForEach("tr.variantRow", func(_ int, trVariantRow *colly.HTMLElement) {
				variantCondition := trVariantRow.ChildText("td.variantInfo")
				variantPriceString := strings.Trim(trVariantRow.ChildText("td.price"), "CAD$ ")
				variantPrice, err := strconv.ParseFloat(variantPriceString, 64)
				if err != nil {
					log.Fatalf("Error wile parsing string to float: '%s'", variantPriceString)
				}
				variants = append(variants, &VariantInfo{
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
				variants = append(variants, &VariantInfo{
					Condition: variantCondition,
					Price:     variantPrice,
				})
			}
			sort.Slice(variants, func(first, second int) bool {
				return variants[first].Price < variants[second].Price
			})
			products = append(products, &ProductInfo{
				Name:     productName,
				Set:      productSet,
				Variants: variants,
			})

		})
	})
	if err := c.Visit("http://www.facetofacegames.com/products/search?query=" + cardName); err != nil {
		log.Fatal("Error while visiting facetoface ", err)

	}

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
