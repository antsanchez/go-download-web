package download

import (
	"os"
	"strings"

	"github.com/antsanchez/godownloadweb/commons"
)

func hasPaths(url string) bool {
	paths := strings.Split(url, "/")
	if len(paths) > 1 {
		return true
	}
	return false
}

func getOnlyPath(url string) (path string) {

	paths := strings.Split(url, "/")
	if len(paths) <= 1 {
		return url
	}

	total := paths[:len(paths)-1]

	return strings.Join(total[:], "/")
}

// GetPath returns only the path, without domain, from the given link
func GetPath(link string) string {
	return strings.Replace(link, commons.Root, "", 1)
}

// exists returns whether the given file or directory exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
