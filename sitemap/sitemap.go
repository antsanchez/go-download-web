package sitemap

import (
	"encoding/xml"
	"os"
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

	var data = []byte(xml.Header)
	data = appendBytes(data, []byte(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`))
	data = appendBytes(data, []byte("\n"))

	for _, val := range links {
		pos := URLSitemap{Loc: val}
		output, err := xml.MarshalIndent(pos, "  ", "    ")
		if err != nil {
			return err
		}
		data = appendBytes(data, output)
		data = appendBytes(data, []byte("\n"))
	}

	data = appendBytes(data, []byte(`</urlset>`))

	return os.WriteFile(filename, data, 0644)
}
