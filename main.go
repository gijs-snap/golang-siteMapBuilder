package main 

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	htmlParser "github.com/gijs-snap/golang-htmlParser"
)

func main() {
	url := "https://www.seltzers.co.nz/"
	html := getHtml(url)

	links := getLinksFromPage(html)

	for _, l := range links {
		isForeignSite := checkLinkDomain(l.Href)
		if isForeignSite != true {
			fmt.Println(l.Href)
			html := getHtml(url + l.Href)
			fmt.Println(html)
		}
	}

}

func getHtml(url string) string{
	fmt.Printf("HTML code of %s ...\n", url)

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(html)
}

func checkLinkDomain(l string) bool{
	return strings.HasPrefix(l, "http");
}

func getLinksFromPage(html string) []htmlParser.Link {
	r := strings.NewReader(html)

	parsed, err := htmlParser.Parse(r)
	if err != nil {
		fmt.Println("Error getting links from page")
		panic(err)
	}
	return parsed
}