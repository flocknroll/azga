package utils

import (
	"bufio"
	"crypto/md5"
	"io"
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
		panic(err)
	}

	result := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	return result
}
