package scrapper

import (
	"bytes"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/antsanchez/godownloadweb/commons"
)

var (
	extensions = []string{".png", ".jpg", ".jpeg", ".json", ".js", ".tiff", ".pdf", ".txt", ".gif", ".psd", ".ai", "dwg", ".bmp", ".zip", ".tar", ".gzip", ".svg", ".avi", ".mov", ".json", ".xml", ".mp3", ".wav", ".mid", ".ogg", ".acc", ".ac3", "mp4", ".ogm", ".cda", ".mpeg", ".avi", ".swf", ".acg", ".bat", ".ttf", ".msi", ".lnk", ".dll", ".db", ".css"}
	falseURLs  = []string{"mailto:", "javascript:", "tel:", "whatsapp:", "callto:", "wtai:", "sms:", "market:", "geopoint:", "ymsgr:", "msnim:", "gtalk:", "skype:"}
	validURL   = regexp.MustCompile(`\(([^()]*)\)`)
	validCSS   = regexp.MustCompile(`\{(\s*?.*?)*?\}`)
)

// isInternLink checks if a link is intern
func isInternLink(link string) bool {
	if strings.Index(link, commons.Root) == 0 {
		return true
	}
	return false
}

// removeQuery removes the query parameters from the given link
func removeQuery(link string) string {
	return strings.Split(link, "?")[0]
}

// isStart cheks if the site is the startsite
func isStart(link string) bool {
	if strings.Compare(link, commons.Root) == 0 {
		return true
	}
	return false
}

// sanitizeURL sanitizes a URL
func sanitizeURL(link string) string {
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

	if !*UseQueries {
		tram = removeQuery(tram)
	}

	return tram
}

// IsValidExtension check if an extension is valid
func IsValidExtension(link string) bool {
	for _, extension := range extensions {
		if strings.Contains(strings.ToLower(link), extension) {
			return false
		}
	}
	return true
}

// isValidLink checks if a link is valid
func isValidLink(link string) bool {
	if isInternLink(link) && !isStart(link) && IsValidExtension(link) {
		return true
	}

	return false
}

// isValidAttachment checks if the link is a valid extension
func isValidAttachment(link string) bool {
	if isInternLink(link) && !isStart(link) && !IsValidExtension(link) {
		return true
	}

	return false
}

// doesLinkExist checks if a link exists in a given slice
func doesLinkExist(newLink commons.Links, existingLinks []commons.Links) (exists bool) {
	for _, val := range existingLinks {
		if strings.Compare(newLink.Href, val.Href) == 0 {
			exists = true
		}
	}

	return
}

// IsURLInSlice checks if a URL is in a slice
func IsURLInSlice(search string, array []string) bool {
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
func IsLinkScanned(link string, scanned []string) (exists bool) {
	for _, val := range scanned {
		if strings.Compare(link, val) == 0 {
			exists = true
		}
	}

	return
}

// getURLEmbeeded from HTML or CSS
func getURLEmbeeded(body string) (url string) {
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
func GetInsideAttachments(url string) (attachments []string) {
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
				found := getURLEmbeeded(r)
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

	return
}
