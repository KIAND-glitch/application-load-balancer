package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type Dictionary map[string]string

type response struct {
	Word       string
	Definition string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./server <port>")
		return
	}

	port := os.Args[1]

	dictionary, err := loadDictionary("dictionary.json")
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}

	templates := template.Must(template.ParseFiles("lookup.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "lookup.html", nil)
	})

	http.HandleFunc("/lookup", func(w http.ResponseWriter, r *http.Request) {
		word := r.URL.Query().Get("word")
		definition, found := dictionary[word]

		if !found {
			definition = "Word not found"
		}

		res := response{
			Word:       word,
			Definition: definition,
		}

		templates.ExecuteTemplate(w, "lookup.html", res)
	})

	fmt.Println("Dictionary server is listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

func loadDictionary(filename string) (Dictionary, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dictionary Dictionary
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dictionary); err != nil {
		return nil, err
	}

	return dictionary, nil
}
