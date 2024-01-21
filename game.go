package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"math/rand"
	"os"
	"strings"
)

var LETTERS = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'k', 'l',
	'm', 'n', 'o', 'p', 'r', 's', 't', 'u', 'v', 'y'}

const MinFound = 40

func ReadLinesFromFile(filename string) []string {
	file, _ := os.Open(filename)
	var lines []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	_ = file.Close()
	return lines
}

var words = ReadLinesFromFile("dict.txt")

func makeGame(words []string) map[string]interface{} {
	for {
		center := choice(LETTERS)
		letters := sample(LETTERS, center, 6)
		allowed := string(append(letters, center))
		var found []string
		for _, word := range words {
			if !strings.ContainsRune(word, center) {
				continue
			}
			if allRunesInWordAllowed(word, allowed) {
				found = append(found, word)
			}
		}

		if len(found) >= MinFound && anyBingo(found) {
			total := 0
			for _, w := range found {
				if len(w) > 4 {
					total += len(w)
				} else {
					total++
				}
			}
			return map[string]interface{}{
				"letters":  string(letters),
				"center":   string(center),
				"words":    len(found),
				"total":    total,
				"wordlist": found,
			}
		}
	}
}

func choice(options []rune) rune {
	return options[rand.Intn(len(options))]
}

func sample(options []rune, except rune, n int) []rune {
	var res []rune
	count := 0
	rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })
	for _, r := range options {
		if r != except {
			res = append(res, r)
			count++
			if count == n {
				break
			}
		}
	}
	return res
}

func allRunesInWordAllowed(word string, allowed string) bool {
	for _, r := range word {
		if !strings.ContainsRune(allowed, r) {
			return false
		}
	}
	return true
}

func anyBingo(words []string) bool {
	for _, word := range words {
		if len(uniqueRunes(word)) == 7 {
			return true
		}
	}
	return false
}

func uniqueRunes(s string) []rune {
	u := make([]rune, 0, len(s))
	m := make(map[rune]bool)

	for _, val := range s {
		if _, exists := m[val]; !exists {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func main() {
	lambda.Start(HandleRequest)
}

type MyEvent struct{}

func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	jsonData, _ := json.Marshal(makeGame(words))
	message := string(jsonData)
	return &message, nil
}
