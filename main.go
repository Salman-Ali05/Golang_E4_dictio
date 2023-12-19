// main.go
package main

import (
	"fmt"

	"./dictionary"
)

func main() {
	dictionary := dictionary.NewDictionary()

	dictionary.Add("go", "aller")
	dictionary.Add("hello", "bonjour")
	dictionary.Add("world", "monde")
	dictionary.Add("eat", "manger")
	dictionary.Add("drink", "boire")

	def, found := dictionary.Get("go")
	if found {
		fmt.Printf("Translation of 'go': %s\n", def)
	} else {
		fmt.Println("Not found.")
	}

	dictionary.Remove("world")

	wordList := dictionary.List()
	fmt.Println("Dictionary words:")
	for _, entry := range wordList {
		fmt.Println(entry)
	}
}
