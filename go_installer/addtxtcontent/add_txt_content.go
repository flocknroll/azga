package addtxtcontent

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/flocknroll/azga/go_installer/utils"
)

// Iterate through a file and returns the MD5 digests of the lines grouped by the specified number.
func rollingHashFile(path string, linesNb int) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		f, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err == nil {
			defer f.Close()
		} else {
			log.Fatal(err)
		}

		lines := make([]string, 0, linesNb)

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			lines = append(lines, line)

			if len(lines) > linesNb {
				lines = lines[1:]
			}

			if len(lines) == linesNb {
				_, hash := utils.HashLinesCRC32(lines)
				ch <- hash
			}
		}

		close(ch)
	}()

	return ch
}

// Check if the source file content if present in the destination file.
func CheckContent(srcPath string, destPath string) bool {
	linesNb, srcHash := utils.HashLinesCRC32(utils.ReadEntireFile(srcPath))

	for destHash := range rollingHashFile(destPath, linesNb) {
		if bytes.Equal(srcHash, destHash) {
			return true
		}
	}

	return false
}

// Add the content of the source file at the end of the destination file.
func AddContent(srcPath string, destPath string) {
	src := utils.ReadEntireFile(srcPath)
	f, err := os.OpenFile(destPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err == nil {
		defer f.Close()
	} else {
		log.Fatal(err)
	}

	for _, line := range src {
		_, err = f.WriteString("\n" + line)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Check if a delimited text section is present in the file.
func CheckDelimitedSection(path string, startDelimiter string, endDelimiter string) (startLine int, endLine int, totalLines int, found bool) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err == nil {
		defer f.Close()
	} else {
		log.Fatal(err)
	}

	totalLines = 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		totalLines += 1
		line := scanner.Text()
		if line == startDelimiter {
			startLine = totalLines
		}
		if line == endDelimiter {
			endLine = totalLines
		}
	}

	return startLine, endLine, totalLines, startLine > 0 && endLine > 0
}

// Delete lines included in the range in the target file.
func DeleteLines(path string, start int, end int) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	tmpFile, _ := ioutil.TempFile(os.TempDir(), "tmp-azga*.txt")
	if err == nil {
		defer os.Remove(tmpFile.Name())
	} else {
		log.Fatal(err)
	}

	lineNb := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNb += 1

		line := scanner.Text()

		if lineNb != start-1 {
			line += "\n"
		}
		if lineNb < start || lineNb > end {
			io.WriteString(tmpFile, line)
		}
	}
	f.Close()
	utils.CopyFile(tmpFile.Name(), path)
}
