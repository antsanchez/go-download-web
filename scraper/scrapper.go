package scraper

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Scraper struct {
	// Original domain
	OldDomain string

	// New domain to rewrite the download HTML sites
	NewDomain string

	// Root domain
	Root string

	// Path where to save the downloads
	Path string

	// Use args on URLs
	UseQueries bool
}

// Links model
type Links struct {
	Href string
}

// Page model
type Page struct {
	URL       string
	Canonical string
	Links     []Links
	NoIndex   bool
	HTML      string
}

// getLinks Get the links from a HTML site
func (s *Scraper) getLinks(domain string) (page Page, attachments []string, err error) {
	resp, err := http.Get(domain)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	page.HTML = buf.String()

	doc, err := html.Parse(buf)
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
					found := s.getURLEmbeeded(a.Val)
					if found != "" {
						link, err := resp.Request.URL.Parse(found)
						if err == nil {
							foundLink := s.sanitizeURL(link.String())
							if s.isValidAttachment(foundLink) {
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
						foundLink := s.sanitizeURL(link.String())
						if s.isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						} else if s.isValidLink(foundLink) {
							page.Links = append(page.Links, Links{Href: foundLink})
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
						foundLink := s.sanitizeURL(link.String())
						if s.isValidAttachment(foundLink) {
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
						foundLink := s.sanitizeURL(link.String())
						if s.isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}
				if a.Key == "srcset" {
					links := strings.Split(a.Val, " ")
					for _, val := range links {
						link, err := resp.Request.URL.Parse(val)
						if err == nil {
							foundLink := s.sanitizeURL(link.String())
							if s.isValidAttachment(foundLink) {
								attachments = append(attachments, foundLink)
							}
						}
					}
				}
			}
		}

		// Get links
		// super lazy alphanumeric checking because I don't want to do regex.
		alphanumeric := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321"
		if n.Type == html.ElementNode && n.Data == "a" {
			ok := false
			newLink := Links{}

			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)

					if strings.Contains(alphanumeric, string(a.Val[0])) {
						link, _ = url.Parse(resp.Request.URL.Scheme + "://" + resp.Request.URL.Hostname() + "/" + a.Val)
					}
					if err == nil {
						foundLink := s.sanitizeURL(link.String())
						if s.isValidLink(foundLink) {
							ok = true
							newLink.Href = foundLink
						} else if s.isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}

			}

			if ok && !s.doesLinkExist(newLink, page.Links) {
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
func (s *Scraper) TakeLinks(
	toScan string,
	started chan int,
	finished chan int,
	scanning chan int,
	newLinks chan []Links,
	pages chan Page,
	attachments chan []string,
) {
	started <- 1
	scanning <- 1
	defer func() {
		<-scanning
		finished <- 1
		fmt.Printf("Started: %6d - Finished %6d", len(started), len(finished))
	}()

	// Get links
	page, attached, err := s.getLinks(toScan)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Save Page
	pages <- page

	attachments <- attached

	// Save links
	newLinks <- page.Links
}
