package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/antsanchez/go-download-web/pkg/console"
	"github.com/antsanchez/go-download-web/pkg/get"
	"golang.org/x/net/html"
)

// New creates a new Scraper
func New(conf *Conf, getter get.GetInterface, con console.TerminalInterface) Scraper {

	con.AddStatus("Checking domain")

	// Get the root domain
	resp, _, err := getter.Get(conf.OldDomain)
	if err != nil {
		log.Fatal(err)
	}

	correct := RemoveLastSlash(resp.Request.URL.String())

	// Prepare the roots
	conf.Roots = append(conf.Roots, correct)
	if len(conf.IncludedURLs) > 0 {
		var urls = strings.Split(conf.IncludedURLs, ",")
		for _, url := range urls {
			if len(url) == 0 {
				continue
			}
			url = strings.TrimSpace(url)
			url = RemoveLastSlash(url)
			conf.Roots = append(conf.Roots, url)
		}
	}

	con.AddDomain(correct)

	con.AddStatus("Initiating scraper")

	return Scraper{
		OldDomain:  conf.OldDomain,
		NewDomain:  conf.NewDomain,
		Roots:      conf.Roots,
		Path:       conf.Path,
		UseQueries: conf.UseQueries,

		Scanning:    make(chan int, conf.Simultaneous), // Semaphore
		NewLinks:    make(chan []Links, 100000),        // New links to scan
		Pages:       make(chan Page, 100000),           // Pages scanned
		Attachments: make(chan []string, 100000),       // Attachments
		Started:     make(chan int, 100000),            // Crawls started
		Finished:    make(chan int, 100000),            // Crawls finished

		Indexed:    []string{},
		ForSitemap: []string{},
		Files:      []string{},
		StartTime:  time.Now(),

		Seen: make(map[string]bool),

		Get: getter,
		Con: con,
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

func (s *Scraper) Run() {
	defer s.Close()
	s.Scrape()
	s.DownloadAttachments()
}

// getLinks Get the links from a HTML site
func (s *Scraper) getLinks(domain string) (page Page, attachments []string, err error) {
	resp, buf, err := s.Get.Get(domain)
	if err != nil {
		s.Con.AddErrors(err.Error())
		return
	}

	// If rediection, get the new domain
	if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound {
		s.NewDomain = resp.Header.Get("Location")
	} else if resp.StatusCode != http.StatusOK {
		return page, attachments, fmt.Errorf("status code error: %d %s on %s", resp.StatusCode, resp.Status, domain)
	}

	page.HTML = buf.String()

	doc, err := html.Parse(buf)
	if err != nil {
		s.Con.AddErrors(err.Error())
		return
	}

	page.URL = domain

	var f func(*html.Node)
	f = func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key == "style" {

				// Get the attachments from the CSS
				matches := urlsInCSS.FindAllStringSubmatch(a.Val, -1)
				for _, match := range matches {
					if len(match) < 2 {
						continue
					}

					// Trim quotes and whitespace
					urlStr := strings.TrimSpace(match[1])

					// Parse the URL to check if it's valid
					found, err := resp.Request.URL.Parse(urlStr)
					if err != nil {
						continue
					}

					foundLink := s.SanitizeURL(found.String())
					if s.IsValidAttachment(foundLink) {
						attachments = append(attachments, foundLink)
					}
				}
			}
		}

		// Get CSS Links
		if n.Type == html.ElementNode && n.Data == "link" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := s.SanitizeURL(link.String())
						if s.IsValidAttachment(foundLink) {
							attachments = append(attachments, s.RemoveTrailingSlash(foundLink))
						} else if s.IsValidSite(foundLink) {
							page.Links = append(page.Links, Links{Href: foundLink})
						}
					}
				}
			}
		}

		// Get CSS and JS from no script
		if n.Type == html.ElementNode && n.Data == "noscript" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "link" {
					for _, a := range c.Attr {
						if a.Key == "href" {
							link, err := resp.Request.URL.Parse(a.Val)
							if err == nil {
								foundLink := s.SanitizeURL(link.String())
								if s.IsValidAttachment(foundLink) {
									attachments = append(attachments, foundLink)
								}
							}
						}
					}
				}
				if c.Type == html.ElementNode && c.Data == "script" {
					for _, a := range c.Attr {
						if a.Key == "src" {
							link, err := resp.Request.URL.Parse(a.Val)
							if err == nil {
								foundLink := s.SanitizeURL(link.String())
								if s.IsValidAttachment(foundLink) {
									attachments = append(attachments, foundLink)
								}
							}
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
						foundLink := s.SanitizeURL(link.String())
						if s.IsValidAttachment(foundLink) {
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
						foundLink := s.SanitizeURL(link.String())
						if s.IsValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}
				if a.Key == "srcset" {
					links := strings.Split(a.Val, " ")
					for _, val := range links {
						link, err := resp.Request.URL.Parse(val)
						if err == nil {
							foundLink := s.SanitizeURL(link.String())
							if s.IsValidAttachment(foundLink) {
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
						foundLink := s.SanitizeURL(link.String())
						if s.IsValidSite(foundLink) {
							ok = true
							newLink.Href = foundLink
						} else if s.IsValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}

			}

			if ok && !s.DoesLinkExist(newLink, page.Links) {
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

	s.Con.AddStarted()
	s.Con.AddStatus("Scraping " + link)

	defer func() {
		<-s.Scanning
		s.Finished <- 1
		s.Con.AddFinished()
	}()

	// Get links
	page, attached, err := s.getLinks(link)
	if err != nil {
		s.Con.AddErrors(err.Error())
	} else {
		s.Pages <- page
		s.Attachments <- attached
		s.NewLinks <- page.Links
	}
}

// Scrape scrapes the site
func (s *Scraper) Scrape() {

	s.Con.AddStatus("Scraping " + s.OldDomain)

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
						s.Con.AddErrors(err.Error())
					}
				}()
			}
		case attachment := <-s.Attachments:
			for _, link := range attachment {
				if !s.IsURLInSlice(link, s.Files) {
					s.Con.AddAttachments()
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
			moreAttachments, err := s.GetInsideAttachments(attachedFile)
			if err != nil {
				s.Con.AddErrors(err.Error())
				continue
			}
			for _, link := range moreAttachments {
				link := link
				if !s.IsURLInSlice(link, s.Files) {
					s.Files = append(s.Files, link)
					s.Con.AddDownloading()

					err := s.SaveAttachment(link)
					if err != nil {
						s.Con.AddErrors(err.Error())
					}
					s.Con.AddDownloaded()
				}
			}
		}

		err := s.SaveAttachment(attachedFile)
		if err != nil {
			s.Con.AddErrors(err.Error())
		}

		s.Con.AddDownloaded()
	}

}
