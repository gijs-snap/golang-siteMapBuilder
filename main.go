package main 

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	//"encoding/xml"
	htmlParser "github.com/gijs-snap/golang-htmlParser"
)

type URL struct {
    loc    string `xml:"loc"`
}

// type UrlArray struct {
// 	URLList []List
// }


func main() {
	url := "https://www.seltzers.co.nz/"
	html := getHtml(url)

	links := getLinksFromPage(html)

	var uniqueLinks []string
	for _, l := range links {
		isForeignSite := checkLinkDomain(l.Href)
		if isForeignSite != true {
			_, found := Find(uniqueLinks, l.Href)
			if !found {
				// create struct instead of slice
				// newl := URL{loc:l.Href}
				isMailTo := strings.HasPrefix(l.Href, "mailto");
				if isMailTo != true {
					uniqueLinks = append(uniqueLinks, l.Href)
					html := getHtml(url + l.Href)
					getLinksFromPage(html)
				}
			}			
		}
	}

	generateXML(uniqueLinks)
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

func generateXML(uniqueLinks []string) {
	fmt.Println(uniqueLinks)
	// xmlstring, err := xml.MarshalIndent()
}

func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}