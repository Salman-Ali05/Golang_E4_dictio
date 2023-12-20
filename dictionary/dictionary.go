package dictionary

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// nil = null

type Dictionary struct {
	entries map[string]string
}

func NewDictionary(filename string) (*Dictionary, error) {
	d := &Dictionary{
		entries: make(map[string]string),
	}

	err := d.LoadFromFile(filename)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Dictionary) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &d.entries) // json decoder
	if err != nil {
		return fmt.Errorf("Failed decode JSON : %v", err)
	}

	return nil
}

func (d *Dictionary) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(d.entries, "", "  ") // set json indent
	if err != nil {
		return fmt.Errorf("Failed encode JSON : %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'écriture dans le fichier : %v", err)
	}

	return nil
}

func (d *Dictionary) Add(word, definition string) {
	d.entries[word] = definition
}

func (d *Dictionary) Get(word string) (string, bool) {
	definition, ok := d.entries[word]
	return definition, ok
}

func (d *Dictionary) Remove(word string) {
	delete(d.entries, word)
}

func (d *Dictionary) List() []string {
	var result []string

	for word, definition := range d.entries {
		result = append(result, fmt.Sprintf("%s: %s", word, definition))
	}

	sort.Strings(result)

	return result
}
