package sitemap

import (
	"encoding/xml"
	"io/ioutil"
)

// URLSitemap is the model for every url entry on the sitemap
type URLSitemap struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

func appendBytes(appendTo []byte, toAppend []byte) []byte {
	for _, val := range toAppend {
		appendTo = append(appendTo, val)
	}
	return appendTo
}

// CreateSitemap create the sitemap
func CreateSitemap(links []string, filename string) {
	var total = []byte(xml.Header)
	total = appendBytes(total, []byte(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`))
	total = appendBytes(total, []byte("\n"))

	for _, val := range links {
		pos := URLSitemap{Loc: val}
		output, err := xml.MarshalIndent(pos, "  ", "    ")
		if err != nil {
			panic(err)
		}
		total = appendBytes(total, output)
		total = appendBytes(total, []byte("\n"))
	}

	total = appendBytes(total, []byte(`</urlset>`))

	err := ioutil.WriteFile(filename, total, 0644)
	if err != nil {
		panic(err)
	}
}
