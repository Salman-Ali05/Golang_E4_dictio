package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"Golang_E4_dictio/main/dictionary"
)

var dict *dictionary.Dictionary

func main() {
	var err error
	dict, err = dictionary.NewDictionary("dictionary.json")
	if err != nil {
		fmt.Println("Couldn't create dictionary :", err)
		return
	}

	// routes
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/get/", GetHandler)
	http.HandleFunc("/remove/", RemoveHandler)
	http.HandleFunc("/list", ListHandler)

	fmt.Println("Server is running on :8081")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var entry dictionary.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dict.Add(entry.Word, entry.Definition)
	w.WriteHeader(http.StatusOK)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	word := r.URL.Path[len("/get/"):]
	definition, found := dict.Get(word)
	if !found {
		http.Error(w, "Word not found", http.StatusNotFound)
		return
	}

	response := dictionary.Entry{Word: word, Definition: definition.Definition}
	json.NewEncoder(w).Encode(response)
}

func RemoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	word := r.URL.Path[len("/remove/"):]
	dict.Remove(word)
	w.WriteHeader(http.StatusOK)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	wordList := dict.List()
	for _, entry := range wordList {
		fmt.Fprintln(w, entry)
	}
}

/*

SMALL DOC ON API
To list everything : localhost:8081/list
To get one word : localhost:8081/get/word (example : localhost:8081/get/chai)
To add a new word : localhost:8081/add in POST method (you can try in POSTMAN)
json :{
  "word": "word",
  "definition": "definition"
}
To remove a word : localhost:8081/remove in DELETE method (you can try in POSTMAN)
json :{
  "word": "word"
}
*/
