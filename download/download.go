package download

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/antsanchez/go-download-web/commons"
)

// Settings needed
type Settings struct {
	OldDomain string
	NewDomain string // NewDomain to rewrite the download HTML sites
}

type Downloader struct {
	Conf Settings
}

func New(oldDomain, newDomain string) Downloader {
	return Downloader{
		Conf: Settings{
			OldDomain: oldDomain,
			NewDomain: newDomain,
		},
	}
}

// Download a single link
func (d *Downloader) download(url string, filename string, changeDomain bool) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()

	// Change domain
	if changeDomain && d.Conf.NewDomain != "" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		newStr := buf.String()

		newStr = strings.ReplaceAll(newStr, d.Conf.OldDomain, d.Conf.NewDomain)

		newContent := bytes.NewBufferString(newStr)
		_, err = io.Copy(f, newContent)

	} else {
		_, err = io.Copy(f, resp.Body)
	}

	return
}

// All Download the complete website
func (d *Downloader) All(indexed []string) {
	for _, url := range indexed {
		filepath := d.GetPath(url)
		if filepath == "" {
			filepath = "/index.html"
		}

		// Get last path
		if d.hasPaths(filepath) {
			if commons.IsFinal(filepath) {
				// if the url is a final url in a folder, like example.com/path
				// this will create the folder "path" and, inside, the index.html file
				if !d.exists(commons.PATH + filepath) {
					os.MkdirAll(commons.PATH+filepath, 0755) // first create directory
					filepath = filepath + "index.html"
				}
			} else {
				// if the url is not a final url in a folder, like example.com/path/bum.html
				// this will create the folder "path" and, inside, the bum.html file
				path := d.getOnlyPath(filepath)
				if !d.exists(commons.PATH + path) {
					os.MkdirAll(commons.PATH+path, 0755) // first create directory
				}
			}

		}

		d.download(url, commons.PATH+filepath, true)
	}
}

// Attachments download all attachments
func (d *Downloader) Attachments(attachments []string) {
	for _, url := range attachments {
		filepath := d.GetPath(url)
		if filepath == "" {
			continue
		}

		// Get last path
		if d.hasPaths(filepath) {
			if commons.IsFinal(filepath) {
				// if the url is a final url in a folder, like example.com/path/
				// this will create the folder "path" and, inside, the index.html file
				filepath = commons.RemoveLastSlash(filepath)
				url = commons.RemoveLastSlash(url)
			}

			path := d.getOnlyPath(filepath)
			if !d.exists(commons.PATH + path) {
				os.MkdirAll(commons.PATH+path, 0755) // first create directory
			}
		}

		d.download(url, commons.PATH+filepath, false)
	}
}
