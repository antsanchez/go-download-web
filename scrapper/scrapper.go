package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/antsanchez/godownloadweb/commons"
	"golang.org/x/net/html"
)

// UseQueries true if query parameters should be also crawled
var UseQueries *bool

// getLinks Get the links from a HTML site
func getLinks(domain string) (page commons.Page, attachments []string, err error) {
	resp, err := http.Get(domain)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	page.URL = domain
	foundMeta := false

	var f func(*html.Node)
	f = func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key == "style" {
				if strings.Contains(a.Val, "url(") {
					found := getURLEmbeeded(a.Val)
					if found != "" {
						link, err := resp.Request.URL.Parse(found)
						if err == nil {
							foundLink := sanitizeURL(link.String())
							if isValidAttachment(foundLink) {
								attachments = append(attachments, foundLink)
							}
						}
					}

				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == "robots" {
					foundMeta = true
				}
				if foundMeta {
					if a.Key == "content" && strings.Contains(a.Val, "noindex") {
						page.NoIndex = true
					}
				}
			}
		}

		// Get CSS and AMP
		if n.Type == html.ElementNode && n.Data == "link" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := sanitizeURL(link.String())
						if isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						} else if isValidLink(foundLink) {
							page.Links = append(page.Links, commons.Links{Href: foundLink})
						}
					}
				}
			}
		}

		// Get JS Scripts
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := sanitizeURL(link.String())
						if isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}
			}
		}

		// Get Images
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := sanitizeURL(link.String())
						if isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}
				if a.Key == "srcset" {
					links := strings.Split(a.Val, " ")
					for _, val := range links {
						link, err := resp.Request.URL.Parse(val)
						if err == nil {
							foundLink := sanitizeURL(link.String())
							if isValidAttachment(foundLink) {
								attachments = append(attachments, foundLink)
							}
						}
					}
				}
			}
		}

		// Get links
		if n.Type == html.ElementNode && n.Data == "a" {
			ok := false
			newLink := commons.Links{}

			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := sanitizeURL(link.String())
						if isValidLink(foundLink) {
							ok = true
							newLink.Href = foundLink
						} else if isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}

			}

			if ok && !doesLinkExist(newLink, page.Links) {
				page.Links = append(page.Links, newLink)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return
}

// TakeLinks take links from the given site
func TakeLinks(toScan string, started chan int, finished chan int, scanning chan int, newLinks chan []commons.Links, pages chan commons.Page, attachments chan []string) {
	started <- 1
	scanning <- 1
	defer func() {
		<-scanning
		finished <- 1
		fmt.Printf("\rStarted: %6d - Finished %6d", len(started), len(finished))
	}()

	// Get links
	page, attached, err := getLinks(toScan)
	if err != nil {
		return
	}

	// Save Page
	pages <- page

	attachments <- attached

	// Save links
	newLinks <- page.Links
}
