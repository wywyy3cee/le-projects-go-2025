package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"
)

type SearchResult struct {
	File string
	Word string
}

// ----------------- VIEW -----------------
func printResults(results <-chan SearchResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for res := range results {
		fmt.Printf("✅ WORD \"%s\" was found in FILE: %s\n", res.Word, res.File)
	}
}

// ----------------- CONTROLLER / LOGIC -----------------
const MaxWorkers = 5

var sem = make(chan struct{}, MaxWorkers) // semaphore for limiting goroutines

func searchFile(filename, word string, results chan<- SearchResult, wg *sync.WaitGroup) {
	defer wg.Done()

	sem <- struct{}{}        // slot capture
	defer func() { <-sem }() // slot release

	content, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	clean := strings.ToLower(strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return r
	}, string(content)))

	if strings.Contains(clean, strings.ToLower(word)) {
		results <- SearchResult{File: filename, Word: word}
	}
}

func searchPath(path, word string, results chan<- SearchResult, wg *sync.WaitGroup) {
	info, err := os.Stat(path)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	if !info.IsDir() {
		wg.Add(1)
		go searchFile(path, word, results, wg)
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Println("Error reading directory:", err)
		return
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		searchPath(entryPath, word, results, wg)
	}
}

// ----------------- MAIN -----------------
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run grepdirrec.go <word> <file_or_directory>")
		return
	}

	word := os.Args[1]
	root := os.Args[2]

	var wg sync.WaitGroup
	results := make(chan SearchResult)

	var printWg sync.WaitGroup
	printWg.Add(1)
	go printResults(results, &printWg)

	searchPath(root, word, results, &wg)

	wg.Wait()
	close(results)

	printWg.Wait()
}
