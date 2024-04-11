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
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/antsanchez/go-download-web/scraper"
	"github.com/antsanchez/go-download-web/sitemap"
)

func main() {
	conf, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	// Do First call to domain
	resp, err := http.Get(conf.OldDomain)
	if err != nil {
		log.Println("Domain could not be reached!")
		return
	}
	defer resp.Body.Close()

	// Prepare the root domains
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

	// Create directory for downloaded website
	err = os.MkdirAll(conf.Path, 0755)
	if err != nil {
		log.Println(conf.Path)
		log.Fatal(err)
	}

	scrap := scraper.New(conf)
	defer scrap.Close()

	scrap.Scrape()

	log.Println("\nFinished scraping the site...")

	scrap.DownloadAttachments()

	log.Println("Creating Sitemap...")
	err = sitemap.CreateSitemap(scrap.ForSitemap, scrap.Path)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished.")
}
