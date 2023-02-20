package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
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

// A struct to map a Single Pokemon
type SinglePokemon struct {
	BaseHappiness     int          `json:"base_happiness"`
	CaptureRate       int          `json:"capture_rate"`
	GenderRate        int          `json:"gender_rate"`
	GenderDifferences bool         `json:"has_gender_differences"`
	Description       []FlavorText `json:"flavor_text_entries"`
	Name              string       `json:"name"`
	EvolvesFrom       []Evolution  `json:"evolves_from_species"`
	ShapeOfPkmn       []Shape      `json:"shape"`
	ArchetypeOfPkmn   []Genera     `json:"genera"`
}

// A struct to map the description of the Single Pokemon
type FlavorText struct {
	Text string `json:"flavor_text"`
}

// A struct to map the evolution chain of a Single Pokemon
type Evolution struct {
	Anchestor string `json:"name"`
}

// A struct to map the shape of a Single Pokemon
type Shape struct {
	PkmnShape string `json:"name"`
}

// A struct to map the genera of a Single Pokemon
type Genera struct {
	Genus string `json:"genus"`
}

// Create a global instance of the Template method of the template package.
// You’ll access this template instance from various parts of your program.
var tmplt *template.Template

// I need a struct to display the info on the web page
type DispPokemon struct {
	Headline string
	Body     string
}

// Follow the following tutorial
// https://www.makeuseof.com/go-html-templating/

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")

		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		responseData, err := io.ReadAll(response.Body)
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

		fmt.Fprint(w, html)

	})

	http.HandleFunc("/pokemon", func(w http.ResponseWriter, r *http.Request) {
		pokemonID := r.URL.Query().Get("id")
		response, err := http.Get(pokemonID)

		fmt.Println("109")

		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		fmt.Println("116")

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("123")

		tmplt, _ = template.ParseFiles("description.html")
		var pokemonDescription SinglePokemon
		json.Unmarshal(responseData, &pokemonDescription)
		description := pokemonDescription.Description[0].Text
		description = regexp.MustCompile(`[^a-zA-Z0-9.'Éé ]+`).ReplaceAllString(description, " ")

		fmt.Println("131")

		event := DispPokemon{
			Headline: pokemonDescription.Name,
			Body:     description,
		}

		fmt.Println("138")

		err = tmplt.Execute(w, event)

		if err != nil {
			return
		}

		fmt.Println("142")

		//html := "<html><head><title>" + pokemonDescription.Name + "</title></head><body><h1>" + pokemonDescription.Name + "</h1><p>" + description + "</p></body></html>"
		//fmt.Fprintf(w, html)

	})

	http.ListenAndServe(":8080", nil)

}