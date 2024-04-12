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

	"github.com/antsanchez/go-download-web/pkg/console"
	"github.com/antsanchez/go-download-web/pkg/get"
	"github.com/antsanchez/go-download-web/pkg/scraper"
)

// ParseFlags parses the flags
func parseFlags() (conf scraper.Conf, err error) {
	flag.StringVar(&conf.OldDomain, "u", "", "URL to download")
	flag.StringVar(&conf.NewDomain, "new", "", "New URL")
	flag.StringVar(&conf.IncludedURLs, "r", "", "URL prefixes/root paths that should be included in the scraper, in addition to the domain")
	flag.IntVar(&conf.Simultaneous, "s", 3, "Number of concurrent connections")
	flag.BoolVar(&conf.UseQueries, "q", false, "Ignore queries on URLs")
	flag.StringVar(&conf.Path, "path", "./website", "Local path for downloaded files")
	flag.Parse()

	if conf.OldDomain == "" {
		err = errors.New("URL cannot be empty! Please, use '-u <URL>'")
		return
	}

	if conf.Simultaneous <= 0 {
		err = errors.New("the number of concurrent connections be at least 1'")
		return
	}

	return
}

func main() {
	// Parse the flags
	conf, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new scraper
	scrap := scraper.New(&conf, get.New(), console.New())

	// Run the scraper
	scrap.Run()
}
