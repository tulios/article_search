package main

import (
  "article_search/search"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	store, err := search.NewIndex(os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	fmt.Print("> ")
	for scanner.Scan() {
		articles, err := store.Search(scanner.Text())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s.\n", strings.Join(articles, ", "))
		}
		fmt.Print("> ")
	}
	fmt.Println()
}
