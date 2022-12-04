package gochess

import (
	"log"
	"os"
	"strings"
)

func AppendToFile(path string, content string) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func DeleteFile(path string) {
	if e := os.Remove(path); e != nil {
		log.Print(e)
	}
}

func ReadLines(path string) ([]string, error) {
	data, err := Asset(path)
	content := string(data)
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		if line != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines, err
}
