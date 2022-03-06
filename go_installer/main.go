package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// Returns a MD5 digest of the lines.
// Strips the new line character if present.
func hashLines(lines []string) (int, []byte) {
	nb := 0
	md5 := md5.New()

	for ix, line := range lines {
		io.WriteString(md5, line)
		nb = ix + 1
	}

	return nb, md5.Sum(nil)
}

// Reads all the lines from a file.
func readEntireFile(path string) []string {
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

// Iterate through a file and returns the MD5 digests of the lines grouped by the specified number.
func rollingHashFile(path string, linesNb int) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err == nil {
			defer f.Close()
		} else {
			panic(err)
		}

		lines := make([]string, linesNb)

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)

			if len(lines) > linesNb {
				lines = lines[1:]
			}

			if len(lines) == linesNb {
				_, hash := hashLines(lines)
				ch <- hash
			}
		}

		close(ch)
	}()

	return ch
}

// Checks if the source file content if present in the destination file.
func checkContent(src_path string, dest_path string) bool {
	linesNb, srcHash := hashLines(readEntireFile(src_path))

	for destHash := range rollingHashFile(dest_path, linesNb) {
		if destHash != nil && bytes.Equal(srcHash, destHash) {
			return true
		}
	}

	return false
}

// Add the content of the source file at the end of the destination file if not already present.
func addContent(src_path string, dest_path string) {
	found := checkContent(src_path, dest_path)

	if found {
		fmt.Println("Data found")
	} else {
		fmt.Println("Data not found - appending")
		src := readEntireFile(src_path)
		f, err := os.OpenFile(dest_path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err == nil {
			defer f.Close()
		} else {
			panic(err)
		}

		for _, line := range src {
			_, err = f.WriteString("\n" + line)
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	// fmt.Println(checkContent("test_data/src_found.txt", "test_data/dest.txt"))
	// fmt.Println(checkContent("test_data/src_not_found.txt", "test_data/dest.txt"))

	addContent(os.Args[1], os.Args[2])
}
