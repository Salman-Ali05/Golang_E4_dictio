package main

import (
	"fmt"

	"Golang_E4_dictio/main/dictionary"
)

func main() {
	dictionary, err := dictionary.NewDictionary("dictionary.json")
	if err != nil {
		fmt.Println("Couldn't create file :", err)
		return
	}

	dictionary.Add("go", "aller")
	dictionary.Add("hello", "bonjour")
	dictionary.Add("world", "monde")
	dictionary.Add("eat", "manger")
	dictionary.Add("drink", "boire")
	dictionary.Add("run", "courir")
	dictionary.Add("chai", "thé")

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

	err = dictionary.SaveToFile("dictionary.json")
	if err != nil {
		fmt.Println("Couldn't save data in json file :", err)
	}
}
