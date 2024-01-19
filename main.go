// Copyright 2021 Antonio Sanchez (asanchez.dev). All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/antsanchez/go-download-web/scraper"
	"github.com/antsanchez/go-download-web/sitemap"
)

type Flags struct {
	// Domain to be scraped
	Domain string

	// New Domain to be set
	NewDomain string

	// URL prefixes/roots that should be included in the scraper
	IncludedURLs string

	// Number of concurrent queries
	Simultaneous int

	// Use query parameters on URLs
	UseQueries bool

	// Path where to download the files to
	Path string
}

func parseFlags() (Flags, error) {
	flags := Flags{}
	flag.StringVar(&flags.Domain, "u", "", "URL to download")
	flag.StringVar(&flags.NewDomain, "new", "", "New URL")
	flag.StringVar(&flags.IncludedURLs, "r", "", "URL prefixes/root paths that should be included in the scraper, in addition to the domain")
	flag.IntVar(&flags.Simultaneous, "s", 3, "Number of concurrent connections")
	flag.BoolVar(&flags.UseQueries, "q", false, "Ignore queries on URLs")
	flag.StringVar(&flags.Path, "path", "./website", "Local path for downloaded files")
	flag.Parse()

	if flags.Domain == "" {
		return flags, errors.New("URL cannot be empty! Please, use '-u <URL>'")
	}

	if flags.Simultaneous <= 0 {
		return flags, errors.New("the number of concurrent connections be at least 1'")
	}

	log.Println("Simultaneous: ", flags.Simultaneous)
	log.Println("Use Queries:  ", flags.UseQueries)
	log.Println("Domain:       ", flags.Domain)
	if flags.NewDomain != "" {
		log.Println("New Domain:   ", flags.NewDomain)
	}

	return flags, nil
}

func main() {
	flags, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	// Create directory for downloaded website
	err = os.MkdirAll(flags.Path, 0755)
	if err != nil {
		log.Fatal(fmt.Errorf("Failed creating target path '%s': %w", flags.Path, err))
	}

	scanning := make(chan int, flags.Simultaneous) // Semaphore
	newLinks := make(chan []scraper.Links, 100000) // New links to scan
	pages := make(chan scraper.Page, 100000)       // Pages scanned
	attachments := make(chan []string, 100000)     // Attachments
	started := make(chan int, 100000)              // Crawls started
	finished := make(chan int, 100000)             // Crawls finished

	var indexed, forSitemap, files []string

	seen := make(map[string]bool)

	start := time.Now()

	defer func() {
		close(newLinks)
		close(pages)
		close(started)
		close(finished)
		close(scanning)

		log.Printf("\nDuration: %s\n", time.Since(start))
		log.Printf("Number of pages: %6d\n", len(indexed))
	}()

	// Do First call to domain
	resp, err := http.Get(flags.Domain)
	if err != nil {
		log.Println("Domain could not be reached!")
		return
	}
	defer resp.Body.Close()

	roots := []string{}
	roots = append(roots, resp.Request.URL.String())

	if len(flags.IncludedURLs) > 0 {
		var urls = strings.Split(flags.IncludedURLs, ",")
		for _, url := range urls {
			if len(url) == 0 {
				continue
			}
			roots = append(roots, url)
		}
	}

	s := scraper.Scraper{
		OldDomain:  flags.Domain,
		NewDomain:  flags.NewDomain,
		Roots:      roots,
		Path:       flags.Path,
		UseQueries: flags.UseQueries,
	}

	// Take the links from the startsite
	s.TakeLinks(flags.Domain, started, finished, scanning, newLinks, pages, attachments)
	seen[flags.Domain] = true

	for {
		select {
		case links := <-newLinks:
			for _, link := range links {
				link := link
				if !seen[link.Href] {
					seen[link.Href] = true
					go s.TakeLinks(link.Href, started, finished, scanning, newLinks, pages, attachments)
				}
			}

		case page := <-pages:
			if !s.IsURLInSlice(page.URL, indexed) {
				indexed = append(indexed, page.URL)

				err := s.SaveHTML(page.URL, page.HTML)
				if err != nil {
					log.Println(err)
				}
			}

			if !page.NoIndex {
				if !s.IsURLInSlice(page.URL, forSitemap) {
					forSitemap = append(forSitemap, page.URL)
				}
			}

		case attachment := <-attachments:
			for _, link := range attachment {
				if !s.IsURLInSlice(link, files) {
					files = append(files, link)
				}
			}
		}

		// Break the for loop once all scans are finished
		if len(started) > 1 && len(scanning) == 0 && len(started) == len(finished) {
			break
		}
	}

	log.Println("\nFinished crawling the site...")

	log.Println("\nDownloading attachments... ", len(files))
	for _, attachedFile := range files {
		attachedFile := attachedFile
		if strings.Contains(attachedFile, ".css") {
			moreAttachments := s.GetInsideAttachments(attachedFile)
			for _, link := range moreAttachments {
				link := link
				if !s.IsURLInSlice(link, files) {
					log.Println("Appended: ", link)
					files = append(files, link)

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

	log.Println("Creating Sitemap...")
	err = sitemap.CreateSitemap(forSitemap, flags.Path)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished.")
}
