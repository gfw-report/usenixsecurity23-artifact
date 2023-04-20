package readfiles

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

// GetFiles return a slice of opened files.
// usage example: in external function, files := GetFiles(flag.Args())
func GetFiles(filePaths []string) []*os.File {
	files := make([]*os.File, 0)
	if len(filePaths) == 0 {
		files = append(files, os.Stdin)
	}
	for _, path := range filePaths {
		if path == "-" {
			files = append(files, os.Stdin)
		} else {
			matches, err := filepath.Glob(path)
			if err != nil {
				log.Panicln(err)
			}
			for _, p := range matches {
				file, err := os.Open(p)
				if err != nil {
					log.Panicln(err)
				}
				files = append(files, file)
			}
		}
	}

	return files
}

// ReadFiles return a channel and then sequentially pipe all lines in filePaths to the channel in a non-blocking way.
// usage example: in external function, lines := ReadFiles(flag.Args())
func ReadFiles(filePaths []string) chan string {
	const maxNumLines = 100000
	lines := make(chan string, maxNumLines)
	go func() {
		files := GetFiles(filePaths)
		for _, file := range files {
			scanner := bufio.NewScanner(file)
			// optionally, resize scanner's capacity for lines over 64K
			if err := scanner.Err(); err != nil {
				log.Panicln(err)
			}
			for scanner.Scan() {
				lines <- scanner.Text()
			}
			file.Close()
		}
		close(lines)
	}()
	return lines
}
