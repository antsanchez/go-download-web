package commons

type Conf struct {
	// Root domain Root
	Root string

	// PATH Path where to save the website
	Path string
}

func New(root, path string) Conf {
	return Conf{
		Root: root,
		Path: path,
	}
}

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
	return string(url[len(url)-1]) == "/"
}

// RemoveLastSlash removes the last slash
func RemoveLastSlash(url string) string {
	if string(url[len(url)-1]) == "/" {
		return url[:len(url)-1]
	}
	return url
}
