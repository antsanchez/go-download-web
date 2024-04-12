package scraper

import (
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	renderedExtensions = []string{".html", ".htm", ".php", ".asp", ".aspx", ".jsp", ".cfm", ".cgi", ".pl", ".py", ".rb", ".shtml", ".shtm", ".phtml"}
	extensions         = []string{
		".png", ".jpg", ".jpeg", ".json", ".js", ".tiff", ".pdf", ".txt", ".gif", ".psd", ".ai", "dwg", ".bmp", ".zip", ".tar", ".gzip", ".svg",
		".avi", ".mov", ".xml", ".mp3", ".wav", ".mid", ".ogg", ".acc", ".ac3", "mp4", ".ogm", ".cda", ".mpeg", ".avi", ".swf", ".acg",
		".bat", ".ttf", ".msi", ".lnk", ".dll", ".db", ".css", ".csv", ".parquet", ".tar", ".gz", ".bz2", ".xz", ".7z", ".rar", ".zip", ".tar.gz",
		".asc", ".pgp", ".sig", ".md5", ".sha1", ".sha256", ".sha512", ".asc", ".pgp", ".sig", ".md5", ".sha1", ".sha256", ".sha512", ".tgz", ".wmv",
		".flv", ".rm", ".rmvb", ".asf", ".mpg", ".mpeg", ".mpe", ".wmv", ".mp4", ".mkv", ".vob", ".mov", ".qt", ".avi", ".asf", ".rm", ".rmvb",
		".ico", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".ods", ".odp", ".odg", ".odf", ".txt", ".rtf", ".pdf", ".epub", ".mobi",
		".aiff", ".wav", ".mp3", ".aac", ".ogg", ".wma", ".flac", ".alac", ".ape", ".aif", ".mid", ".midi", ".mka", ".opus", ".ra", ".rm", ".sln",
		".webmanifest", ".manifest",
	}
	falseURLs = []string{
		"mailto:", "javascript:", "tel:", "whatsapp:", "callto:", "wtai:", "sms:", "market:", "geopoint:", "ymsgr:", "msnim:", "gtalk:", "skype:",
		"aim:", "icq:", "irc:", "ircs:", "mumble:", "sip:", "xmpp:", "aim:", "itms:", "itms-apps:", "itms-services:", "data:", "blob:", "about:",
		"chrome:", "chrome-extension:", "chrome-untrusted:", "chrome-search:", "chrome-native", "chrome-devtools:", "chrome-devtools:", "chrome-devtools:",
	}
	urlsInCSS      = regexp.MustCompile(`url\(['"]?(.*?)['"]?\)`)
	validJS        = regexp.MustCompile(`import\s+[\w\*\s]+\s+from\s+['"](.*?)['"]`)
	validJSImport  = regexp.MustCompile(`import\s+['"](.*?)['"]`)
	validJSRequire = regexp.MustCompile(`require\s*\(\s*['"](.*?)['"]\s*\)`)
)

// IsInternLink checks if a link is intern
func (s *Scraper) IsInternLink(link string) bool {
	for _, root := range s.Roots {
		if strings.Index(link, root) == 0 {
			return true
		}
	}
	return false
}

// IsStart cheks if the site is the startsite
func (s *Scraper) IsStart(link string) bool {
	for _, root := range s.Roots {
		if strings.Compare(link, root) == 0 {
			return true
		}
	}
	return false
}

// SanitizeURL sanitizes a URL
func (s *Scraper) SanitizeURL(link string) string {
	for _, fal := range falseURLs {
		if strings.Contains(link, fal) {
			return ""
		}
	}

	link = strings.TrimSpace(link)

	tram := strings.Split(link, "#")[0]

	if !s.UseQueries {
		tram = strings.Split(tram, "?")[0]
	}

	// avoid index out of range
	if len(tram) == 0 {
		return tram
	}

	if !s.IsValidExtension(tram) && !s.HasRenderedExtension(tram) {
		if string(tram[len(tram)-1]) != "/" {
			tram = tram + "/"
		}
	}

	return tram
}

// IsValidExtension check if an extension is valid
func (s *Scraper) IsValidExtension(link string) bool {
	if !strings.Contains(link, ".") {
		return false
	}

	found := link[strings.LastIndex(link, "."):]
	if found == "" {
		return false
	}

	for _, ext := range extensions {
		if strings.Compare(found, ext) == 0 {
			return true
		}
	}

	return false
}

func (s *Scraper) HasRenderedExtension(link string) bool {
	if !strings.Contains(link, ".") {
		return false
	}

	found := link[strings.LastIndex(link, "."):]
	if found == "" {
		return false
	}

	for _, ext := range renderedExtensions {
		if strings.Compare(found, ext) == 0 {
			return true
		}
	}

	return false
}

// IsValidLink checks if the link is a valid url and from the domain
func (s *Scraper) IsValidLink(link string) (ok bool) {
	if !s.IsInternLink(link) {
		return false
	}

	// check if is a valid url with the url package
	got, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	if got.Scheme == "" || got.Host == "" {
		return false
	}

	return true
}

// IsValidLink checks if the link is a site and not an attachment
func (s *Scraper) IsValidSite(link string) bool {
	if !s.IsInternLink(link) {
		return false
	}

	if s.IsStart(link) {
		return false
	}

	if s.IsValidExtension(link) {
		return false
	}

	return true
}

// IsValidAttachment checks if the link is a valid extension, not a site
func (s *Scraper) IsValidAttachment(link string) bool {
	if !s.IsInternLink(link) {
		return false
	}

	return s.IsValidExtension(s.RemoveTrailingSlash(link))
}

// RemoveTrailingSlash removes the trailing slash from a link
func (s *Scraper) RemoveTrailingSlash(link string) string {
	if len(link) == 0 {
		return link
	}

	if string(link[len(link)-1]) == "/" {
		return link[:len(link)-1]
	}

	return link
}

// DoesLinkExist checks if a link exists in a given slice
func (s *Scraper) DoesLinkExist(newLink Links, existingLinks []Links) (exists bool) {
	for _, val := range existingLinks {
		if strings.Compare(newLink.Href, val.Href) == 0 {
			exists = true
		}
	}

	return
}

// IsURLInSlice checks if a URL is in a slice
func (s *Scraper) IsURLInSlice(search string, array []string) bool {
	if len(search) == 0 {
		return false
	}
	if len(array) == 0 {
		return false
	}
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

// GetInsideAttachments gets inside CSS and JS Files
func (s *Scraper) GetInsideAttachments(link string) (attachments []string, err error) {
	if IsFinal(link) {
		// if the url is a final url in a folder, like example.com/path/
		// this will create the folder "path" and, inside, the index.html file
		link = RemoveLastSlash(link)
	}

	resp, buf, err := s.Get.Get(link)
	if err != nil {
		return
	}

	got := resp.Request.URL.String()
	body := buf.String()

	// First, search for JavaScript
	if strings.Contains(got, ".js") {
		blocks := validJS.FindAll([]byte(body), -1)
		for _, b := range blocks {
			// Extract the URL from the import statement or require function
			found := s.getJSURLEmbedded(string(b))
			if found != "" {
				link, err := resp.Request.URL.Parse(found)
				if err == nil {
					foundLink := s.SanitizeURL(link.String())
					if s.IsValidAttachment(foundLink) {
						attachments = append(attachments, foundLink)
					}
				}
			}
		}
	}

	// Second, search for CSS
	if strings.Contains(got, ".css") {
		matches := urlsInCSS.FindAllStringSubmatch(body, -1)
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

	return
}

// getJSURLEmbedded from JavaScript
func (s *Scraper) getJSURLEmbedded(body string) (url string) {
	// Use a regular expression to find import statements or require functions
	valid := validJSImport.Find([]byte(body))
	if valid == nil {
		valid = validJSRequire.Find([]byte(body))
	}
	if valid == nil {
		return
	}

	// Extract the URL from the import statement or require function
	url = string(valid)

	// Remove surrounding quotes
	if string(url[0]) == `'` || string(url[0]) == `"` {
		url = url[1:]
	}
	if string(url[len(url)-1]) == `'` || string(url[len(url)-1]) == `"` {
		url = url[:len(url)-1]
	}

	return url
}

// GetPath returns the path of a given URL
func (s *Scraper) GetPath(url string) (path string) {
	paths := strings.Split(url, "/")
	if len(paths) <= 1 {
		return url
	}

	total := paths[:len(paths)-1]
	path = strings.Join(total[:], "/")

	if len(path) == 0 {
		return "/"
	}

	if string(path[0]) != "/" {
		path = "/" + path
	}

	if string(path[len(path)-1]) != "/" {
		path = path + "/"
	}

	return
}

// GetLastFolder returns the last folder of a path
func (s *Scraper) GetLastFolder(path string) string {
	paths := strings.Split(path, "/")
	if len(paths) == 0 {
		return ""
	}

	return paths[len(paths)-2]
}

// RemoveDomain returns only the path, without domain, from the given link
func (s *Scraper) RemoveDomain(link string) string {
	for _, root := range s.Roots {
		if strings.Index(link, root) == 0 {
			return strings.Replace(link, root, "", 1)
		}
	}

	return link
}

// exists returns whether the given file or directory exists
func (s *Scraper) exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsInSlice check if the given link is in a slice
func IsInSlice(search string, array []string) bool {
	for _, val := range array {
		if val == search {
			return true
		}
	}

	return false
}

// IsFinal check if the url is a folder-like path, like example.com/path/
func IsFinal(url string) bool {
	return string(url[len(url)-1]) == "/"
}

// RemoveLastSlash removes the last slash
func RemoveLastSlash(url string) string {
	if len(url) == 0 {
		return url
	}
	if string(url[len(url)-1]) == "/" {
		return url[:len(url)-1]
	}
	return url
}
