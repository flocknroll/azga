package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func hashLines(lines []string) (int, []byte) {
	nb := 0
	md5 := md5.New()

	for ix, line := range lines {
		io.WriteString(md5, line)
		nb = ix + 1
	}

	return nb, md5.Sum(nil)
}

func readEntireFile(path string) []string {
	f, err := os.OpenFile(path, os.O_RDONLY, 0444)
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

func rollingHashFile(path string, linesNb int) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		f, err := os.OpenFile(path, os.O_RDONLY, 0444)
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

func checkContent(src string, dest string) bool {
	linesNb, srcHash := hashLines(readEntireFile(src))

	for destHash := range rollingHashFile(dest, linesNb) {
		if destHash != nil && bytes.Equal(srcHash, destHash) {
			return true
		}
	}

	return false
}

func main() {
	fmt.Println(checkContent("test_data/src_found.txt", "test_data/dest.txt"))
	fmt.Println(checkContent("test_data/src_not_found.txt", "test_data/dest.txt"))
}
