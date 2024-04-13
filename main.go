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

	"github.com/antsanchez/go-download-web/pkg/console"
	"github.com/antsanchez/go-download-web/pkg/get"
	"github.com/antsanchez/go-download-web/pkg/scraper"
)

func main() {

	// Parse the flags
	conf, err := scraper.ParseFlags()
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new scraper
	scrap, err := scraper.New(conf, get.New(), console.New())
	if err != nil {
		log.Fatal(err)
	}

	// Run the scraper
	scrap.Run()
}
