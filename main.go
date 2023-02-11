package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
}

func main() {
	// First of all we need to query the API endpoint using http.Get
	// the result is mapped into response and err
	response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Perform a conversion of the response body from bytes into something
	// that can be printed out in the console

	// We read all the stream of bytes with ioutil.ReadAll and then
	// convert in string with string(responseData)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Use the following print to check if the request is OK
	// fmt.Println(string(responseData))

	// We need now to Unmarshal the returned JSON string into a new var
	// We declare a new variable of Response type
	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	// We use the address of the variabile because json.Unmarshal need
	// to modify the data and not a copy of the data.

	// Use the following print to check if the Unmarshal operation is ok

	//fmt.Println(responseObject)
	//fmt.Println(len(responseObject.Pokemon))

	// To list all of our firstGen pokemon we need to create a for loops
	// for every object in our responseObject pokemon array

	for _, pokemon := range responseObject.Pokemon {
		fmt.Println(pokemon.Species.Name)
	}
}
