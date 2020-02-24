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
	maxDepth := flag.Int("depth", 10, "the maximum number of links depth to traverse")
	flag.Parse()

	// fmt.Println("The urlFlag: ", *urlFlag) // TESTING
	pages := bfs(*urlFlag, *maxDepth)
	for _, page := range pages {
		fmt.Println("The page is: ", page)
	}

	// pages := get(*urlFlag)
	// for _, page := range pages {
	// 	fmt.Println("The page is: ", page) // TESTING
	// }
}

func bfs(urlStr string, maxDepth int) []string {
	// A collection of url that we've seen
	seen := make(map[string]struct{})

	// To get the links on the page that we're in
	var q map[string]struct{}

	// To get the links on the next page
	// or the page tha twe're not in yet
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url, _ := range q {

			// If there is such "key" inside that map
			// Or tell me if there's a value inside the key within this map
			// True if there is such key, otherwise it'll be False
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	var ret []string
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}

// 1. GET the webpage.
func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body) // for testing purpose

	// Get the URL from a requested URL
	// Example:
	// 		URL: "https://abc.com/efg/"
	// 		Requested URL: "https://abc.com/efg/"
	reqURL := resp.Request.URL
	// fmt.Println("Request URL: ", reqURL.String()) // TESTING

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

	links := hrefs(resp.Body, base)
	return filter(links, withPrefix(base))
}

// 2. Parse all the links from an url
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

// 4. filter out any links with a different domain
func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {

		// Tests if link begins with the "base" prefix
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
