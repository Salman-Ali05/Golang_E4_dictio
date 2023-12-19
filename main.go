package main

import (
	"fmt"
	"sort"
)

type Dictionary struct {
	entries map[string]string
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		entries: make(map[string]string),
	}
}

func (d *Dictionary) add(word, definition string) {
	d.entries[word] = definition
}

func (d *Dictionary) get(word string) (string, bool) {
	definition, ok := d.entries[word]
	return definition, ok
}

func (d *Dictionary) remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) list() []string {
	var result []string

	for word, definition := range d.entries {
		result = append(result, fmt.Sprintf("%s: %s", word, definition))
	}

	sort.Strings(result)

	return result
}

func main() {
	dictionary := NewDictionary()

	dictionary.add("go", "aller")
	dictionary.add("hello", "bonjour")
	dictionary.add("world", "monde")

	def, found := dictionary.get("go")
	if found {
		fmt.Printf("Translation of 'go': %s\n", def)
	} else {
		fmt.Println("Not found.")
	}

	dictionary.remove("world")

	wordList := dictionary.list()
	fmt.Println("Dictionary words :")
	for _, entry := range wordList {
		fmt.Println(entry)
	}
}
