package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bimonestle/go-exercise-projects/04.HTML-Link-Parser/link"
)

// 1. GET the webpage.
// 2. parse all the links on the page.
// 3. build proper urls with our links
// 4. filter out any links with a different domain
// 5. Find all pages (BFS)
// 6. Print out XML
func main() {
	urlFlag := flag.String("url", "https://gophercises.com/", "The url that you want to build a sitemap for")
	flag.Parse()

	fmt.Println("The urlFlag: ", *urlFlag) // TESTING

	// 1. GET the webpage.
	resp, err := http.Get(*urlFlag)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body) // for testing purpose

	// 2. Parse all the links from an url
	// Get the URL from a requested URL
	// Example:
	// 		URL: "https://abc.com/efg/"
	// 		Requested URL: "https://abc.com/efg/"
	reqURL := resp.Request.URL
	fmt.Println("Request URL: ", reqURL.String()) // TESTING

	// Get the base form of URL from a URL
	// Example:
	// 		URL : "https://abc.com/efg/hi"
	// 		Requested URL: "https://abc.com/efg/hi"
	// 		Base URL: "https://abc.com"
	baseURL := &url.URL{
		Scheme: reqURL.Scheme, // "https:""
		Host:   reqURL.Host,   // "//some-domain.com"
	}
	base := baseURL.String()
	// fmt.Println("Base URL: ", base) // TESTING

	pages := hrefs(resp.Body, base)
	for _, page := range pages {
		fmt.Println(page) // TESTING
	}
}

func hrefs(body io.Reader, base string) []string {
	links, _ := link.Parse(body)
	var ret []string
	for _, l := range links {
		// fmt.Println("The parsed link from a response Body: ", l) // TESTING
		switch {
		// for links without base URL
		// example: "/about-us"
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		// for links which begins with "http"
		// example: "https://twitter.com/...."
		// example: "https://thisdomain.com/..."
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}
