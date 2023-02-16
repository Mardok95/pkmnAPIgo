package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// We need to create some struct to manage the response of the API

// A Response struct for the entire response
// The name field is mapped into "name" in JSON, the Pokemon filed
// is mapped with the vaule "pokemon_entries" of the JSON response.
// The Pokemon []Pokemon declaration mean that the "Pokemon" field is
// a list of "Pokemon" object
type Response struct {
	Name    string    `json:"name"`
	Pokemon []Pokemon `json:"pokemon_entries"`
}

// A Pokemon struct to map every pokemon
type Pokemon struct {
	EntryNo int            `json:"entry_number"`
	Species PokemonSpecies `json:"pokemon_species"`
}

// A struct to map the Pokemon's Species
type PokemonSpecies struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FlavorText struct {
	Text string `json:"flavor_text"`
}

type PokemonDescription struct {
	Name        string       `json:"name"`
	Description []FlavorText `json:"flavor_text_entries"`
}

func main() {
	// First of all we need to query the API endpoint using http.Get
	// the result is mapped into response and err
	// Also we use and http.HandleFunc to star a server to show the result of the API GET
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		// Perform a conversion of the response body from bytes into something
		// that can be printed out in the console

		// We read all the stream of bytes with ioutil.ReadAll and then
		// convert in string with string(responseData)
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// We need now to Unmarshal the returned JSON string into a new var
		// We declare a new variable of Response type
		var responseObject Response
		json.Unmarshal(responseData, &responseObject)

		html := "<head><head><title>" + responseObject.Name + "</title></head><body><h1>" + responseObject.Name + "</h1><ul>"
		for i := 0; i < len(responseObject.Pokemon); i++ {
			html += "<li><a href='/pokemon?id=" + string(responseObject.Pokemon[i].Species.URL) + "'>" + responseObject.Pokemon[i].Species.Name + "</a></li>"
		}

		html += "</ul></body></html>"

		fmt.Fprintf(w, html)

	})

	http.HandleFunc("/pokemon", func(w http.ResponseWriter, r *http.Request) {
		pokemonID := r.URL.Query().Get("id")
		response, err := http.Get(pokemonID)

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var pokemonDescription PokemonDescription
		json.Unmarshal(responseData, &pokemonDescription)

		html := "<html><head><title>" + pokemonDescription.Name + "</title></head><body><h1>" + pokemonDescription.Name + "</h1><p>" + pokemonDescription.Description[0].Text + "</p></body></html>"
		fmt.Fprintf(w, html)

	})

	http.ListenAndServe(":8080", nil)

}
