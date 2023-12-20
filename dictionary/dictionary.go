package dictionary

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Entry représente une entrée dans le dictionnaire.
type Entry struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

// Dictionary représente un dictionnaire.
type Dictionary struct {
	entries    map[string]string
	addChan    chan dictioOps
	removeChan chan dictioOps
}

type dictioOps struct {
	action string
	word   string
	def    string
	res    chan bool
}

// NewDictionary crée une nouvelle instance de Dictionary.
func NewDictionary(filename string) (*Dictionary, error) {
	d := &Dictionary{
		entries:    make(map[string]string),
		addChan:    make(chan dictioOps),
		removeChan: make(chan dictioOps),
	}
	go d.startOperationManager()

	err := d.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// LoadFromFile charge les données du fichier JSON dans le dictionnaire.
func (d *Dictionary) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &d.entries)
	if err != nil {
		return fmt.Errorf("Failed decode JSON: %v", err)
	}

	return nil
}

// SaveToFile sauvegarde les données du dictionnaire dans un fichier JSON.
func (d *Dictionary) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(d.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed encode JSON: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

// Add ajoute une entrée au dictionnaire.
func (d *Dictionary) Add(word, definition string) {
	resChan := make(chan bool)
	d.addChan <- dictioOps{
		action: "add",
		word:   word,
		def:    definition,
		res:    resChan,
	}

	<-resChan
}

// Get récupère une définition par mot.
func (d *Dictionary) Get(word string) (Entry, bool) {
	definition, ok := d.entries[word]
	return Entry{Word: word, Definition: definition}, ok
}

// Remove supprime une entrée du dictionnaire.
func (d *Dictionary) Remove(word string) {
	resChan := make(chan bool)
	d.removeChan <- dictioOps{
		action: "remove",
		word:   word,
		res:    resChan,
	}
	<-resChan
}

// List retourne une liste triée des entrées du dictionnaire.
func (d *Dictionary) List() []Entry {
	var result []Entry

	for word, definition := range d.entries {
		result = append(result, Entry{Word: word, Definition: definition})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Word < result[j].Word
	})

	return result
}

// startOperationManager gère les opérations du dictionnaire de manière asynchrone.
func (d *Dictionary) startOperationManager() {
	for {
		select {
		case operation := <-d.addChan:
			d.entries[operation.word] = operation.def
			operation.res <- true

		case operation := <-d.removeChan:
			delete(d.entries, operation.word)
			operation.res <- true
		}
	}
}
