package scraper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antsanchez/go-download-web/pkg/get"
	"golang.org/x/net/html"
)

// New creates a new Scraper
func New(conf *Conf, getter get.GetInterface) Scraper {

	// Get the root domain
	resp, _, err := getter.Get(conf.OldDomain)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the roots
	conf.Roots = append(conf.Roots, resp.Request.URL.String())
	if len(conf.IncludedURLs) > 0 {
		var urls = strings.Split(conf.IncludedURLs, ",")
		for _, url := range urls {
			if len(url) == 0 {
				continue
			}
			conf.Roots = append(conf.Roots, url)
		}
	}

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

func (s *Scraper) ExtractURLs(link string) (page Page, attachments []string, err error) {
	resp, buf, err := s.Get.Get(link)
	if err != nil {
		log.Println(err)
		return
	}

	// if not 200 status code, return
	if resp.StatusCode != 200 {
		return page, attachments, fmt.Errorf("Status code: %d", resp.StatusCode)
	}

	page.HTML = buf.String()
	page.URL = link

	urls := make([]string, 0)
	tokenizer := html.NewTokenizer(strings.NewReader(page.HTML))

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}
		token := tokenizer.Token()

		// Check all tokens for relevant attributes
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			for _, attr := range token.Attr {
				if attr.Key == "href" || attr.Key == "src" || attr.Key == "srcset" || attr.Key == "data-src" || attr.Key == "data-href" {
					urls = append(urls, attr.Val)
				}
			}
			if tokenType == html.StartTagToken && token.Data == "meta" {
				// Check meta tags for content attribute
				for _, attr := range token.Attr {
					if attr.Key == "content" {
						urls = append(urls, attr.Val)
					}
				}
			}
		}
	}

	for _, url := range urls {
		// Parse the URL to check if it's valid
		found, err := resp.Request.URL.Parse(url)
		if err != nil {
			continue
		}
		foundLink := s.SanitizeURL(found.String())
		if s.IsValidSite(foundLink) {
			page.Links = append(page.Links, Links{Href: foundLink})
		} else if s.IsValidAttachment(foundLink) {
			attachments = append(attachments, foundLink)
		}
	}

	return page, attachments, nil
}

// getLinks Get the links from a HTML site
func (s *Scraper) getLinks(domain string) (page Page, attachments []string, err error) {
	resp, buf, err := s.Get.Get(domain)
	if err != nil {
		log.Println(err)
		return
	}

	page.HTML = buf.String()

	doc, err := html.Parse(buf)
	if err != nil {
		log.Println(err)
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
						fmt.Printf("%q is not a valid URL\n", urlStr)
						continue
					}

					foundLink := s.SanitizeURL(found.String())
					if s.IsValidAttachment(foundLink) {
						fmt.Println("Adding inline CSS attachment:", foundLink)
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
							fmt.Println("Adding node attachments:", foundLink, a.Val)
							attachments = append(attachments, s.RemoveTrailingSlash(foundLink))
						} else if s.IsValidSite(foundLink) {
							fmt.Println("Adding node links:", foundLink, a.Val)
							page.Links = append(page.Links, Links{Href: foundLink})
						}
					}
				}
			}
		}

		// Get CSS and JS from no script
		if n.Type == html.ElementNode && n.Data == "noscript" {
			fmt.Println("Found noscript")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				fmt.Println("Found noscript child")
				if c.Type == html.ElementNode && c.Data == "link" {
					fmt.Println("Found noscript link")
					for _, a := range c.Attr {
						if a.Key == "href" {
							link, err := resp.Request.URL.Parse(a.Val)
							if err == nil {
								foundLink := s.SanitizeURL(link.String())
								if s.IsValidAttachment(foundLink) {
									fmt.Println("Adding noscript attachments:", foundLink)
									attachments = append(attachments, foundLink)
								}
							}
						}
					}
				}
				if c.Type == html.ElementNode && c.Data == "script" {
					fmt.Println("Found noscript script")
					for _, a := range c.Attr {
						if a.Key == "src" {
							link, err := resp.Request.URL.Parse(a.Val)
							if err == nil {
								foundLink := s.SanitizeURL(link.String())
								if s.IsValidAttachment(foundLink) {
									fmt.Println("Adding noscript attachments:", foundLink)
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
							fmt.Println("Adding JS script:", foundLink)
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
							fmt.Println("Adding img:", foundLink)
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
								fmt.Println("Adding srcset:", foundLink)
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
							fmt.Println("Adding link:", foundLink)
							ok = true
							newLink.Href = foundLink
						} else if s.IsValidAttachment(foundLink) {
							fmt.Println("Adding attachment from links:", foundLink)
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
	defer func() {
		<-s.Scanning
		s.Finished <- 1
		fmt.Printf("Started: %6d - Finished %6d\n", len(s.Started), len(s.Finished))
	}()

	// Get links
	page, attached, err := s.ExtractURLs(link)
	if err != nil {
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
			moreAttachments, err := s.GetInsideAttachments(attachedFile)
			if err != nil {
				log.Println(err)
				continue
			}
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
