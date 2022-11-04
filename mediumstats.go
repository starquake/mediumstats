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

func main() {
	var ss []string
	phraseList := []string{"(file attached)", "<Media omitted>"}
	ignoreWords := []string{
		"van", "de", "een", "ik", "jij", "het", "-", "je", "is", "dat", "en", "ook", "in", "op", "wel", "met", "voor",
		"of", "er", "als", "dan", "niet", "die", "maar", "te", "nog", "heb", "kan", "zijn", "wat", "aan", "om", "bij",
		"al", "nu", "was", "ben", "zo", "mijn", "meer", "ze", "we", "mij", "naar", "dit", "dus", "heeft", "even", "geen",
	}

	rStart, err := regexp.Compile("^[0-9]{2}/[0-9]{2}/[0-9]{4}, [0-9]{2}:[0-9]{2} - .*:")
	if err != nil {
		log.Fatalf("Error while trying to compile regex: %s", err)
	}

	file, err := os.Open("resources/chat.txt")
	if err != nil {
		log.Fatalf("Error while trying to open chat.txt: %s", err)
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		ss = append(ss, line)
	}
	ss = deleteLinesWithPhrases(ss, phraseList)
	ss = deleteParts(ss, rStart)
	keys, words := countWords(ss, ignoreWords)

	//fmt.Println(strings.Join(ss, "\n"))
	for _, k := range keys {
		fmt.Println(k, words[k])
	}
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

func countWords(ss []string, ignore []string) ([]string, map[string]int) {
	m := make(map[string]int)
	for _, s := range ss {
		words := strings.Fields(s)
		for _, word := range words {
			hasWord := false
			for _, i := range ignore {
				if strings.ToLower(word) == i {
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

func deleteParts(ss []string, regexp *regexp.Regexp) []string {
	var ssNew []string

	for _, s := range ss {
		s = regexp.ReplaceAllString(s, "")
		ssNew = append(ssNew, s)
	}
	return ssNew
}
