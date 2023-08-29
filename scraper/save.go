package scraper

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/antsanchez/go-download-web/commons"
)

// Download a single link
func (s *Scraper) SaveAttachment(url string) (err error) {
	filepath := s.GetPath(url)
	if filepath == "" {
		return
	}

	// Get last path
	if s.hasPaths(filepath) {
		if commons.IsFinal(filepath) {
			// if the url is a final url in a folder, like example.com/path/
			// this will create the folder "path" and, inside, the file
			filepath = commons.RemoveLastSlash(filepath)
			url = commons.RemoveLastSlash(url)
		}

		path := s.getOnlyPath(filepath)
		if !s.exists(s.Path + path) {
			os.MkdirAll(s.Path+path, 0755) // first create directory
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed retrieving '", url, "' - ", err.Error())
		return
	}
	defer resp.Body.Close()

	absoluteFile := strings.ReplaceAll(s.Path+filepath, "%20", " ")
	f, err := os.Create(absoluteFile)
	if err != nil {
		log.Println("Failed creating '", absoluteFile, "' - ", err.Error())
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return
}

// Download a single link
func (s *Scraper) SaveHTML(url string, html string) (err error) {
	filepath := s.GetPath(url)
	if filepath == "" {
		filepath = "/index.html"
	}

	if s.hasPaths(filepath) {
		if commons.IsFinal(filepath) {
			// if the url is a final url in a folder, like example.com/path
			// this will create the folder "path" and, inside, the index.html file
			if !s.exists(s.Path + filepath) {
				os.MkdirAll(s.Path+filepath, 0755) // first create directory
				filepath = filepath + "index.html"
			}
		} else {
			// if the url is not a final url in a folder, like example.com/path/bum.html
			// this will create the folder "path" and, inside, the bum.html file
			path := s.getOnlyPath(filepath)
			if !s.exists(s.Path + path) {
				os.MkdirAll(s.Path+path, 0755) // first create directory
			}
		}
	}

	f, err := os.Create(s.Path + filepath)
	if err != nil {
		log.Println("Failed creating '", s.Path+filepath, "' - [", err.Error(), "]")
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
