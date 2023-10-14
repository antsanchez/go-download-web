package scraper

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/CalderWhite/go-download-web/commons"
)

var (
	extensions = []string{".png", ".jpg", ".jpeg", ".json", ".js", ".tiff", ".pdf", ".txt", ".gif", ".psd", ".ai", "dwg", ".bmp", ".zip", ".tar", ".gzip", ".svg", ".avi", ".mov", ".json", ".xml", ".mp3", ".wav", ".mid", ".ogg", ".acc", ".ac3", "mp4", ".ogm", ".cda", ".mpeg", ".avi", ".swf", ".acg", ".bat", ".ttf", ".msi", ".lnk", ".dll", ".db", ".css", ".csv", ".parquet"}
	falseURLs  = []string{"mailto:", "javascript:", "tel:", "whatsapp:", "callto:", "wtai:", "sms:", "market:", "geopoint:", "ymsgr:", "msnim:", "gtalk:", "skype:"}
	validURL   = regexp.MustCompile(`\(([^()]*)\)`)
	validCSS   = regexp.MustCompile(`\{(\s*?.*?)*?\}`)
)

// isInternLink checks if a link is intern
func (s *Scraper) isInternLink(link string) bool {
	return strings.Index(link, s.Root) == 0
}

// removeQuery removes the query parameters from the given link
func (s *Scraper) removeQuery(link string) string {
	return strings.Split(link, "?")[0]
}

// isStart cheks if the site is the startsite
func (s *Scraper) isStart(link string) bool {
	return strings.Compare(link, s.Root) == 0
}

// sanitizeURL sanitizes a URL
func (s *Scraper) sanitizeURL(link string) string {
	for _, fal := range falseURLs {
		if strings.Contains(link, fal) {
			return ""
		}
	}

	link = strings.TrimSpace(link)

	if string(link[len(link)-1]) != "/" {
		link = link + "/"
	}

	tram := strings.Split(link, "#")[0]

	if !s.UseQueries {
		tram = s.removeQuery(tram)
	}

	return tram
}

// IsValidExtension check if an extension is valid
func (s *Scraper) IsValidExtension(link string) bool {
	for _, extension := range extensions {
		if strings.Contains(strings.ToLower(link), extension) {
			return false
		}
	}
	return true
}

// isValidLink checks if a link is valid
func (s *Scraper) isValidLink(link string) bool {
	if s.isInternLink(link) && !s.isStart(link) && s.IsValidExtension(link) {
		return true
	}

	return false
}

// isValidAttachment checks if the link is a valid extension
func (s *Scraper) isValidAttachment(link string) bool {
	if s.isInternLink(link) && !s.isStart(link) && !s.IsValidExtension(link) {
		return true
	}

	return false
}

// doesLinkExist checks if a link exists in a given slice
func (s *Scraper) doesLinkExist(newLink Links, existingLinks []Links) (exists bool) {
	for _, val := range existingLinks {
		if strings.Compare(newLink.Href, val.Href) == 0 {
			exists = true
		}
	}

	return
}

// IsURLInSlice checks if a URL is in a slice
func (s *Scraper) IsURLInSlice(search string, array []string) bool {
	withSlash := search[:len(search)-1]
	withoutSlash := search

	if string(search[len(search)-1]) == "/" {
		withSlash = search
		withoutSlash = search[:len(search)-1]
	}

	for _, val := range array {
		if val == withSlash || val == withoutSlash {
			return true
		}
	}

	return false
}

// IsLinkScanned checks if a link has already been scanned
func (s *Scraper) IsLinkScanned(link string, scanned []string) (exists bool) {
	for _, val := range scanned {
		if strings.Compare(link, val) == 0 {
			exists = true
		}
	}

	return
}

// getURLEmbeeded from HTML or CSS
func (s *Scraper) getURLEmbeeded(body string) (url string) {
	valid := validURL.Find([]byte(body))
	if valid == nil {
		return
	}

	url = string(valid)

	// Remove ()
	if string(url[0]) == `(` {
		url = url[1:]
	}
	if string(url[len(url)-1]) == `)` {
		url = url[:len(url)-1]
	}

	// Remove "
	if string(url[0]) == `"` {
		url = url[1:]
	}
	if string(url[len(url)-1]) == `"` {
		url = url[:len(url)-1]
	}

	// Remove '
	if string(url[0]) == `'` {
		url = url[1:]
	}
	if string(url[len(url)-1]) == `'` {
		url = url[:len(url)-1]
	}

	// To do: check if this is a valid url

	return url
}

// GetInsideAttachments gets inside CSS Files
func (s *Scraper) GetInsideAttachments(url string) (attachments []string) {
	if commons.IsFinal(url) {
		// if the url is a final url in a folder, like example.com/path/
		// this will create the folder "path" and, inside, the index.html file
		url = commons.RemoveLastSlash(url)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	if strings.Contains(body, "url(") {
		// Second, search for backgrounds
		blocks := validCSS.FindAll([]byte(body), -1)
		for _, b := range blocks {
			rules := strings.Split(string(b), ";")
			for _, r := range rules {
				found := s.getURLEmbeeded(r)
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

	return
}

func (s *Scraper) hasPaths(url string) bool {
	return len(strings.Split(url, "/")) > 1
}

func (s *Scraper) getOnlyPath(url string) (path string) {
	paths := strings.Split(url, "/")
	if len(paths) <= 1 {
		return url
	}

	total := paths[:len(paths)-1]
	return strings.Join(total[:], "/")
}

// GetPath returns only the path, without domain, from the given link
func (s *Scraper) GetPath(link string) string {
	return strings.Replace(link, s.Root, "", 1)
}

// exists returns whether the given file or directory exists
func (s *Scraper) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
