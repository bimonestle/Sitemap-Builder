package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// 1. GET the webpage.
// 2. parse all the links on the page.
// 3. build proper urls with our links
// 4. filter out any links with a different domain
// 5. Find all pages (BFS)
// 6. Print out XML
func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "The url that you want to build a sitemap for")
	flag.Parse()

	fmt.Println(*urlFlag)

	// 1. GET the webpage.
	resp, err := http.Get(*urlFlag)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
