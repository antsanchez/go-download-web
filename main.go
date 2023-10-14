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
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/antsanchez/go-download-web/scraper"
	"github.com/antsanchez/go-download-web/sitemap"
)

type Flags struct {
	// Domain to be scraped
	Domain *string

	// New Domain to be set
	NewDomain *string

	// Number of concurrent queries
	Simultaneus *int

	// Use query parameters on URLs
	UseQueries *bool

	// Path where to download the files to
	Path *string
}

func parseFlags() (flags Flags, err error) {
	flags.Domain = flag.String("u", "", "URL to copy")
	flags.NewDomain = flag.String("new", "", "New URL")
	flags.Simultaneus = flag.Int("s", 3, "Number of concurrent connections")
	flags.UseQueries = flag.Bool("q", false, "Ignore queries on URLs")
	flags.Path = flag.String("path", "./website", "Local path for downloaded files")
	flag.Parse()

	if *flags.Domain == "" {
		err = errors.New("URL cannot be empty! Please, use '-u <URL>'")
		return
	}

	if *flags.Simultaneus <= 0 {
		err = errors.New("the number of concurrent connections be at least 1'")
		return
	}

	log.Println("Domain:", *flags.Domain)
	if *flags.NewDomain != "" {
		log.Println("New Domain: ", *flags.NewDomain)
	}
	log.Println("Simultaneus:", *flags.Simultaneus)
	log.Println("Use Queries:", *flags.UseQueries)

	return
}

func main() {
	flags, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	// Create directory for downloaded website
	err = os.MkdirAll(*flags.Path, 0755)
	if err != nil {
		log.Println(*flags.Path)
		log.Fatal(err)
	}

	scanning := make(chan int, *flags.Simultaneus) // Semaphore
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
	resp, err := http.Get(*flags.Domain)
	if err != nil {
		log.Println("Domain could not be reached!")
		return
	}
	defer resp.Body.Close()

	s := scraper.Scraper{
		OldDomain:  *flags.Domain,
		NewDomain:  *flags.NewDomain,
		Root:       resp.Request.URL.String(),
		Path:       *flags.Path,
		UseQueries: *flags.UseQueries,
	}

	// Take the links from the startsite
	s.TakeLinks(*flags.Domain, started, finished, scanning, newLinks, pages, attachments)
	seen[*flags.Domain] = true

	for {
		select {
		case links := <-newLinks:
			for _, link := range links {
				if !seen[link.Href] {
					seen[link.Href] = true
					go s.TakeLinks(link.Href, started, finished, scanning, newLinks, pages, attachments)
				}
			}
		case page := <-pages:
			if !s.IsURLInSlice(page.URL, indexed) {
				indexed = append(indexed, page.URL)
				go func() {
					err := s.SaveHTML(page.URL, page.HTML)
					if err != nil {
						log.Println(err)
					}
				}()
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

	log.Println("\nFinished scraping the site...")

	// Semaphore for downloading attachments
	attachmentJobs := make(chan string, len(files))
	var wg sync.WaitGroup
	log.Println("\nDownloading attachments...")
	for i := 1; i <= *flags.Simultaneus; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for attachedFile := range attachmentJobs {
				if strings.Contains(attachedFile, ".css") || strings.Contains(attachedFile, ".csv") || strings.Contains(attachedFile, ".parquet") {
					log.Println(attachedFile)
					s.SaveAttachment(attachedFile)

					moreAttachments := s.GetInsideAttachments(attachedFile)
					for _, link := range moreAttachments {
						if !s.IsURLInSlice(link, files) {
							log.Println("Appended: ", link)
							files = append(files, link)
							go func() {
								err := s.SaveAttachment(link)
								if err != nil {
									log.Println(err)
								}
							}()
						}
					}
				}
			}
		}()
	}

	for _, attachedFile := range files {
		attachmentJobs <- attachedFile
	}

	wg.Wait()

	log.Println("Creating Sitemap...")
	err = sitemap.CreateSitemap(forSitemap, *flags.Path)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished.")
}
