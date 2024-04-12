package scraper

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
)

// PreparePathsFile prepares the folder and filename for a given URL, assuming it's a file
func (s *Scraper) PreparePathsFile(url string) (folder, filename string) {
	url = s.RemoveDomain(url)
	if url == "" {
		return "", ""
	}

	folder = s.GetPath(url)
	if folder == "" {
		folder = "/"
	}

	filename = url[strings.LastIndex(url, "/")+1:]

	// If filename is empty, get the last folder as filename
	if filename == "" && folder != "/" {
		folder = folder[:len(folder)-1]
		filename = folder[strings.LastIndex(folder, "/")+1:]
		folder = folder[:strings.LastIndex(folder, "/")+1]
	}

	return
}

// PreparePathsPage prepares the folder and filename for a given URL, assuming it's a page
func (s *Scraper) PreparePathsPage(url string) (folder, filename string) {
	url = s.RemoveDomain(url)
	if url == "" {
		return "/", "index.html"
	}

	folder = s.GetPath(url)
	if folder == "" {
		folder = "/"
	}

	// check if last path has a renderable extension
	if len(folder) > 1 && IsFinal(url) {
		last := s.GetLastFolder(folder)
		if s.HasRenderedExtension(RemoveLastSlash(last)) {
			filename = last
			folder = folder[:len(folder)-1]
			folder = strings.Replace(folder, last, "", 1)
			return
		}
	}

	filename = url[strings.LastIndex(url, "/")+1:]
	if filename == "" {
		filename = "index.html"
		return
	}

	if !s.HasRenderedExtension(filename) {
		folder = folder + filename + "/"
		filename = "index.html"
	}

	return

}

// Download a single link
func (s *Scraper) SaveAttachment(url string) (err error) {
	folder, filename := s.PreparePathsFile(url)
	folder = s.Path + folder
	final := folder + filename

	if !s.exists(folder) {
		os.MkdirAll(folder, 0755) // first create directory
	}

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(final)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return
}

// Download a single link
func (s *Scraper) SaveHTML(url string, html string) (err error) {
	folder, filename := s.PreparePathsPage(url)
	folder = s.Path + folder
	final := folder + filename

	if !s.exists(folder) {
		os.MkdirAll(folder, 0755) // first create directory
	}

	f, err := os.Create(final)
	if err != nil {
		return
	}
	defer f.Close()

	for _, root := range s.Roots {
		html = strings.ReplaceAll(html, root, "")
	}

	if s.NewDomain != "" && s.OldDomain != s.NewDomain {
		newStr := strings.ReplaceAll(html, s.OldDomain, s.NewDomain)
		newContent := bytes.NewBufferString(newStr)
		_, err = io.Copy(f, newContent)
	} else {
		_, err = io.Copy(f, bytes.NewBufferString(html))
	}

	return
}
