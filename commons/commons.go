package commons

// Page model
type Page struct {
	URL       string
	Canonical string
	Links     []Links
	NoIndex   bool
}

// Links model
type Links struct {
	Href string
}

// Root domain Root
var Root string

// PATH Path where to save the website
const PATH = "website/"

// IsInSlice check if the given link is in a slice
func IsInSlice(search string, array []string) bool {
	for _, val := range array {
		if val == search {
			return true
		}
	}

	return false
}

// IsFinal check if the url is a folder-like path, like example.com/path/
func IsFinal(url string) bool {
	if string(url[len(url)-1]) == "/" {
		return true
	}

	return false
}

// RemoveLastSlash removes the last slash
func RemoveLastSlash(url string) string {
	if string(url[len(url)-1]) == "/" {
		return url[:len(url)-1]
	}
	return url
}
