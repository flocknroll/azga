package utils

import (
	"bufio"
	"crypto/md5"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Returns a MD5 digest of the lines.
// Strips the new line character if present.
func HashLines(lines []string) (int, []byte) {
	nb := 0
	md5 := md5.New()

	for ix, line := range lines {
		io.WriteString(md5, line)
		nb = ix + 1
	}

	return nb, md5.Sum(nil)
}

// Reads all the lines from a file.
func ReadEntireFile(path string) []string {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err == nil {
		defer f.Close()
	} else {
		log.Fatal(err)
	}

	result := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	return result
}

// Copy the remote content into a temporary file
func DownloadContent(url string) string {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "azga-installer-*.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer tmpFile.Close()

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return tmpFile.Name()
}
