package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// checks checks to see if all functions executed without returning an error
func checks(fs ...func() error) {
	for i := len(fs) - 1; i >= 0; i-- {
		if err := fs[i](); err != nil {
			fmt.Println("Received error:", err)
		}
	}
}

func main() {
	var ss []string
	phraseList := []string{"(file attached)", "<Media omitted>"}
	stopWords := getStopWords()

	rStart, err := regexp.Compile("^[0-9]{2}/[0-9]{2}/[0-9]{4}, [0-9]{2}:[0-9]{2} - .*:")
	if err != nil {
		log.Fatalf("Error while trying to compile regex: %s", err)
	}

	file, err := os.Open("resources/chat.txt")
	if err != nil {
		log.Fatalf("Error while trying to open chat.txt: %s", err)
	}
	defer checks(file.Close)

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		ss = append(ss, line)
	}
	ss = deleteLinesWithPhrases(ss, phraseList)
	ss = deleteParts(ss, rStart)
	keys, words := countWords(ss, stopWords)

	for _, k := range keys {
		fmt.Println(k, words[k])
	}
}

func getStopWords() []string {
	var words []string
	file, err := os.Open("resources/stop-words.txt")
	if err != nil {
		log.Fatalf("Error while trying to open stop-words.txt: #{err}")
	}
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		word := fileScanner.Text()
		words = append(words, word)
	}
	return words
}

func deleteLinesWithPhrases(ss []string, phraseList []string) []string {
	var ssNew []string
	for _, s := range ss {
		hasWord := false
		for _, p := range phraseList {
			if strings.Contains(s, p) {
				hasWord = true
			}
		}
		if !hasWord {
			ssNew = append(ssNew, s)
		}
	}
	return ssNew
}

func deleteParts(ss []string, regexp *regexp.Regexp) []string {
	var ssNew []string

	for _, s := range ss {
		s = regexp.ReplaceAllString(s, "")
		ssNew = append(ssNew, s)
	}
	return ssNew
}

func countWords(ss []string, ignore []string) ([]string, map[string]int) {
	m := make(map[string]int)
	for _, s := range ss {
		words := strings.Fields(s)
		for _, word := range words {
			hasWord := false
			for _, i := range ignore {
				if strings.EqualFold(word, i) {
					hasWord = true
				}
			}
			if !hasWord {
				m[word] += 1
			}
		}
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]] > m[keys[j]]
	})

	return keys, m
}
