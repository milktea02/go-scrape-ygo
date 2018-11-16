package main

import (
	"fmt"
	_ "github.com/PuerkitoBio/goquery"
	_ "io"
	"log"
	"net/http"
	_ "os"
	_ "strings"
	"time"
)

func main() {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Create and modify HTTP requeset before sending
	request, err := http.NewRequest("GET", "http://www.facetofacegames.com/products/search?query=cocoon+of+ultra+evolution", nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "Testing Go-lang Application")

	// Make HTTP GET request

	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	// Find all links and process them with the function

	document.Find("a").Each(processElement)

}

// This will get called for each HTML element found
func processElement(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	href, exists := element.Attr("href")
	if exists {
		fmt.Println(href)
	}
}
