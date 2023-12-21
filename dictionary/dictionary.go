package dictionary

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

var saveQueue = make(chan bool)
var mutex = &sync.Mutex{}

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
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// Le fichier n'existe pas, créez-le avec un dictionnaire vide
		d.entries = make(map[string]string)
		err := d.SaveToFile(filename)
		if err != nil {
			return fmt.Errorf("Error creating new file: %v", err)
		}
		return nil
	} else if err != nil {
		return err
	}

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

func (d *Dictionary) Add(word, definition string) {
	resChan := make(chan bool)
	d.addChan <- dictioOps{
		action: "add",
		word:   word,
		def:    definition,
		res:    resChan,
	}

	<-resChan

	// Ajouter l'opération d'écriture du fichier après chaque ajout
	select {
	case saveQueue <- true:
	default:
	}
}

// Get récupère une définition par mot.
func (d *Dictionary) Get(word string) (Entry, bool) {
	definition, ok := d.entries[word]
	return Entry{Word: word, Definition: definition}, ok
}

// Remove supprime une entrée du dictionnaire.
func (d *Dictionary) Remove(word string) bool {
	resChan := make(chan bool)
	d.removeChan <- dictioOps{
		action: "remove",
		word:   word,
		res:    resChan,
	}
	return <-resChan
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
			// Charger les données existantes depuis le fichier
			err := d.LoadFromFile("dictionary.json")
			if err != nil {
				fmt.Println("Error loading from file:", err)
			}

			// Ajouter la nouvelle entrée
			d.entries[operation.word] = operation.def

			// Sauvegarder toutes les entrées dans le fichier
			err = d.SaveToFile("dictionary.json")
			if err != nil {
				fmt.Println("Error saving to file:", err)
			}

			operation.res <- true

		case operation := <-d.removeChan:
			// Charger les données existantes depuis le fichier
			err := d.LoadFromFile("dictionary.json")
			if err != nil {
				fmt.Println("Error loading from file:", err)
			}

			// Vérifier si l'entrée existe avant de la supprimer
			if _, exists := d.entries[operation.word]; exists {
				// Supprimer l'entrée
				delete(d.entries, operation.word)

				// Sauvegarder toutes les entrées dans le fichier
				err = d.SaveToFile("dictionary.json")
				if err != nil {
					fmt.Println("Error saving to file:", err)
				}

				operation.res <- true
			} else {
				// Le mot n'existe pas, renvoyer false
				operation.res <- false
			}
		}
	}
}
