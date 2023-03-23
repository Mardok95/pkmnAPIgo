package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	PkmnName      string
	Headline      string
	Body          string
	ImageURL      string
	ImageAltText  string
	CaptureRate   string
	GenderRate    string
	BaseHappiness string
}

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

		tmplt, _ = template.ParseFiles("templates/pokedex.html")

		err = tmplt.Execute(w, responseObject)
		if err != nil {
			log.Fatal(err)
		}

		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	})

	http.HandleFunc("/pokemon", func(w http.ResponseWriter, r *http.Request) {
		pokemonID := r.URL.Query().Get("id")
		response, err := http.Get(pokemonID)
		parts := strings.Split(pokemonID, "/")
		pokedexNumber := parts[len(parts)-2]

		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		tmplt, _ = template.ParseFiles("templates/description.html")
		var pokemonDescription SinglePokemon
		json.Unmarshal(responseData, &pokemonDescription)
		description := pokemonDescription.Description[0].Text
		description = regexp.MustCompile(`[^a-zA-Z0-9.'Éé ]+`).ReplaceAllString(description, " ")

		imgURL := "static/pokemonSprite/" + pokedexNumber + ".png"

		event := DispPokemon{
			PkmnName:      pokemonDescription.Name,
			Headline:      pokemonDescription.Name,
			Body:          "DESCRIPTION: " + description,
			ImageURL:      imgURL,
			ImageAltText:  "Sprite of " + pokemonDescription.Name,
			CaptureRate:   "Capture rate: " + strconv.Itoa(pokemonDescription.CaptureRate) + "%",
			GenderRate:    "Gender rate: " + strconv.Itoa(pokemonDescription.GenderRate) + "%",
			BaseHappiness: "Base happiness: " + strconv.Itoa(pokemonDescription.BaseHappiness),
		}

		err = tmplt.Execute(w, event)

		if err != nil {
			return
		}

	})

	http.ListenAndServe(":8080", nil)

}
