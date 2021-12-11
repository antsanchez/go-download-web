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
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/antsanchez/godownloadweb/commons"
	"github.com/antsanchez/godownloadweb/download"
	"github.com/antsanchez/godownloadweb/scrapper"
	"github.com/antsanchez/godownloadweb/sitemap"
)

func main() {
	var domain, newDomain string
	if len(os.Args) == 1 {
		fmt.Println("URL can not be empty")
		os.Exit(1)
	}
	domain = os.Args[1]

	if len(os.Args) == 3 {
		newDomain = os.Args[2]
	}

	simultaneus := flag.Int("s", 3, "Number of concurrent connections")
	scrapper.UseQueries = flag.Bool("q", false, "Ignore queries on URLs")
	flag.Parse()

	download.Conf.OldDomain = domain
	download.Conf.NewDomain = newDomain

	fmt.Println("Domain:", domain)
	fmt.Println("New Domain: ", newDomain)
	fmt.Println("Simultaneus:", *simultaneus)
	fmt.Println("Use Queries:", *scrapper.UseQueries)

	if *simultaneus < 1 {
		fmt.Println("There can't be less than 1 simulataneous conexions")
		os.Exit(1)
	}

	scanning := make(chan int, *simultaneus)       // Semaphore
	newLinks := make(chan []commons.Links, 100000) // New links to scan
	pages := make(chan commons.Page, 100000)       // Pages scanned
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

		fmt.Printf("\nTime finished sitemap %s\n", time.Since(start))
		fmt.Printf("Index: %6d\n", len(indexed))
	}()

	// Do First call to domain
	resp, err := http.Get(domain)
	if err != nil {
		fmt.Println("Domain could not be reached!")
		return
	}
	// Todo: get favourite version of URL here
	defer resp.Body.Close()

	// Detected root domain
	commons.Root = resp.Request.URL.String()

	// Take the links from the startsite
	scrapper.TakeLinks(domain, started, finished, scanning, newLinks, pages, attachments)
	seen[domain] = true

	for {
		select {
		case links := <-newLinks:
			for _, link := range links {
				if !seen[link.Href] {
					seen[link.Href] = true
					go scrapper.TakeLinks(link.Href, started, finished, scanning, newLinks, pages, attachments)
				}
			}
		case page := <-pages:
			if !scrapper.IsURLInSlice(page.URL, indexed) {
				indexed = append(indexed, page.URL)
			}

			if !page.NoIndex {
				if !scrapper.IsURLInSlice(page.URL, forSitemap) {
					forSitemap = append(forSitemap, page.URL)
				}
			}
		case attachment := <-attachments:
			for _, link := range attachment {
				if !scrapper.IsURLInSlice(link, files) {
					files = append(files, link)
				}
			}
		}

		// Break the for loop once all scans are finished
		if len(started) > 1 && len(scanning) == 0 && len(started) == len(finished) {
			break
		}
	}

	// Get Inside Attachments
	for _, attachedFile := range files {
		if strings.Contains(attachedFile, ".css") {
			moreAttachments := scrapper.GetInsideAttachments(attachedFile)
			for _, link := range moreAttachments {
				if !scrapper.IsURLInSlice(link, files) {
					fmt.Println("Appended: ", link)
					files = append(files, link)
				}
			}
		}
	}

	// Create directory for website
	os.MkdirAll("website", 0755)

	// Create the sitemap
	sitemap.CreateSitemap(forSitemap, "website/sitemap.xml")

	// Download the complete website
	download.All(indexed)

	// Download the complete attachments
	download.Attachments(files)
}
