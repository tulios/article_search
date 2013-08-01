package search

import (
	"errors"
	"io/ioutil"
	"strings"
)

type Index struct {
	index map[string][]string
}

const (
	INDEX_START_SIZE = 100
	TOKENS           = ".,`~;-+@#$%*(){}[]\\/!?'\""
	OPERATORS        = "&|"
)

// Public

func NewIndex(files ...string) (*Index, error) {
	storage := new(Index)
	storage.index = make(map[string][]string, INDEX_START_SIZE)

	for _, file := range files {
		content, err := readFileContent(file)
		if err != nil {
			return nil, err
		}

		words := tokenize(content)
		indexWords(storage, words, file)
	}

	return storage, nil
}

func (i *Index) Search(term string) ([]string, error) {
	term = sanitize(term, TOKENS)
	term = strings.ToLower(term)

	words := strings.Split(term, "&")

	if len(words) > 1 {
		return i.and(words...)
	}

	words = strings.Split(term, "|")
	return i.or(words...)
}

// Private

func (i *Index) and(words ...string) ([]string, error) {
	var articles []string
	var join []string

	for _, word := range words {
		files, ok := i.index[strings.TrimSpace(word)]
		if !ok {
			return nil, errors.New("No articles found.")
		}
		join = append(join, files...)
	}

	for word, count := range countItems(join) {
		if count == len(words) {
			articles = append(articles, word)
		}
	}

	return articles, nil
}

func (i *Index) or(words ...string) ([]string, error) {
	var articles []string

	for _, word := range words {
		files, ok := i.index[strings.TrimSpace(word)]
		if !ok {
			return nil, errors.New("No articles found.")
		}
		articles = appendUnique(articles, files...)
	}

	return articles, nil
}

func indexWords(storage *Index, words []string, filepath string) {
	for _, word := range words {
		slice, ok := storage.index[word]

		if !ok {
			slice = []string{filepath}

		} else if !hasItem(slice, filepath) {
			slice = append(slice, filepath)
		}

		storage.index[word] = slice
	}
}

func hasItem(slice []string, item string) bool {
	aux := make(map[string]bool, len(slice))

	for _, element := range slice {
		aux[element] = true
	}

	return aux[item] == true
}

func countItems(items []string) map[string]int {
	result := make(map[string]int)
	for _, item := range items {
		v, ok := result[item]
		if !ok {
			v = 1
		} else {
			v++
		}

		result[item] = v
	}

	return result
}

func appendUnique(slice []string, items ...string) []string {
	for _, item := range items {
		if !hasItem(slice, item) {
			slice = append(slice, item)
		}
	}

	return slice
}

func tokenize(content string) []string {
	content = sanitize(content, TOKENS)
	content = sanitize(content, OPERATORS)
	content = strings.ToLower(content)
	return strings.Fields(content)
}

func sanitize(s string, chars string) string {
	for _, char := range []rune(chars) {
		s = strings.Replace(s, string(char), "", -1)
	}

	return s
}

func readFileContent(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
