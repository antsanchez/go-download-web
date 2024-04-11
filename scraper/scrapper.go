package scraper

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// New creates a new Scraper
func New(conf Conf) Scraper {
	return Scraper{
		OldDomain:  conf.OldDomain,
		NewDomain:  conf.NewDomain,
		Root:       conf.Root,
		Path:       conf.Path,
		UseQueries: conf.UseQueries,

		Scanning:    make(chan int, conf.Simultaneus), // Semaphore
		NewLinks:    make(chan []Links, 100000),       // New links to scan
		Pages:       make(chan Page, 100000),          // Pages scanned
		Attachments: make(chan []string, 100000),      // Attachments
		Started:     make(chan int, 100000),           // Crawls started
		Finished:    make(chan int, 100000),           // Crawls finished

		Indexed:    []string{},
		ForSitemap: []string{},
		Files:      []string{},
		StartTime:  time.Now(),

		Seen: make(map[string]bool),
	}
}

// Close closes the channels
func (s *Scraper) Close() {
	close(s.Scanning)
	close(s.NewLinks)
	close(s.Pages)
	close(s.Attachments)
	close(s.Started)
	close(s.Finished)
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
		if n.Type == html.ElementNode && n.Data == "a" {
			ok := false
			newLink := Links{}

			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)
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
func (s *Scraper) TakeLinks(link string) {
	s.Started <- 1
	s.Scanning <- 1
	defer func() {
		<-s.Scanning
		s.Finished <- 1
		fmt.Printf("Started: %6d - Finished %6d", len(s.Started), len(s.Finished))
	}()

	// Get links
	page, attached, err := s.getLinks(link)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Save Page
	s.Pages <- page

	s.Attachments <- attached

	// Save links
	s.NewLinks <- page.Links
}

// Scrape scrapes the site
func (s *Scraper) Scrape() {
	// Take the links from the startsite
	seen := make(map[string]bool)
	s.TakeLinks(s.OldDomain)
	seen[s.OldDomain] = true

	for {
		select {
		case links := <-s.NewLinks:
			for _, link := range links {
				if !seen[link.Href] {
					seen[link.Href] = true
					go s.TakeLinks(link.Href)
				}
			}
		case page := <-s.Pages:
			if !s.IsURLInSlice(page.URL, s.Indexed) {
				s.Indexed = append(s.Indexed, page.URL)
				go func() {
					err := s.SaveHTML(page.URL, page.HTML)
					if err != nil {
						log.Println(err)
					}
				}()
			}

			if !page.NoIndex {
				if !s.IsURLInSlice(page.URL, s.ForSitemap) {
					s.ForSitemap = append(s.ForSitemap, page.URL)
				}
			}
		case attachment := <-s.Attachments:
			for _, link := range attachment {
				if !s.IsURLInSlice(link, s.Files) {
					s.Files = append(s.Files, link)
				}
			}
		}

		// Break the for loop once all scans are finished
		if len(s.Started) > 1 && len(s.Scanning) == 0 && len(s.Started) == len(s.Finished) {
			break
		}
	}
}

// DownloadAttachments downloads the attachments
func (s *Scraper) DownloadAttachments() {
	for _, attachedFile := range s.Files {
		attachedFile := attachedFile

		// First, seek for more attachments on the CSS and JS files
		if strings.Contains(attachedFile, ".css") || strings.Contains(attachedFile, ".js") {
			moreAttachments := s.GetInsideAttachments(attachedFile)
			for _, link := range moreAttachments {
				link := link
				if !s.IsURLInSlice(link, s.Files) {
					log.Println("Appended: ", link)
					s.Files = append(s.Files, link)

					err := s.SaveAttachment(link)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

		err := s.SaveAttachment(attachedFile)
		if err != nil {
			log.Println(err)
		}
	}

}
