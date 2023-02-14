package v2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// A Response struct to map the JSON response
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")

		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

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
