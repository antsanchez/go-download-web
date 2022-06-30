package sitemap

import (
	"encoding/xml"
	"io/ioutil"
)

const SitemapFile = "sitemap.xml"

// URLSitemap is the model for every url entry on the sitemap
type URLSitemap struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
}

func appendBytes(appendTo []byte, toAppend []byte) []byte {
	return append(appendTo, toAppend...)
}

func sitemapPath(filaneme string) string {
	if filaneme != "" {
		return filaneme + "/" + SitemapFile
	}
	return SitemapFile
}

// CreateSitemap creates the sitemap
func CreateSitemap(links []string, filename string) error {
	filename = sitemapPath(filename)

	var total = []byte(xml.Header)
	total = appendBytes(total, []byte(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`))
	total = appendBytes(total, []byte("\n"))

	for _, val := range links {
		pos := URLSitemap{Loc: val}
		output, err := xml.MarshalIndent(pos, "  ", "    ")
		if err != nil {
			return err
		}
		total = appendBytes(total, output)
		total = appendBytes(total, []byte("\n"))
	}

	total = appendBytes(total, []byte(`</urlset>`))

	return ioutil.WriteFile(filename, total, 0644)
}
