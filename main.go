package main 

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"encoding/xml"
	htmlParser "github.com/gijs-snap/golang-htmlParser"
)

var siteMapXmlStart string = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`
var siteMapXmlEnd string = `</urlset>`

type Url struct {
    Loc    string `xml:"loc"`
}


func main() {
	url := "https://www.seltzers.co.nz/"
	html := getHtml(url)

	links := getLinksFromPage(html)

	var uniqueLinks []string
	var allLinks []Url

	for _, l := range links {
		isForeignSite := checkLinkDomain(l.Href)
		if isForeignSite != true {
			_, found := Find(uniqueLinks, l.Href)
			if !found {
				isMailTo := strings.HasPrefix(l.Href, "mailto");
				if isMailTo != true {
					uniqueLinks = append(uniqueLinks, l.Href)
					newLink := Url{Loc: l.Href}
					allLinks = append(allLinks, newLink)
					html := getHtml(url + l.Href)
					getLinksFromPage(html)
				}

			}			
		}
	}
	generateXML(allLinks)
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


func generateXML(allLinks []Url) {
	fmt.Println(allLinks)
	if xmlstring, err := xml.MarshalIndent(allLinks, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + siteMapXmlStart  + "\n" + string(xmlstring) + "\n" + siteMapXmlEnd)
		_ = ioutil.WriteFile("map.xml", xmlstring, 0644)
	}
}

func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}