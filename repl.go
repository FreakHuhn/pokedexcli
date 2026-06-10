package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	pokecache "github.com/FreakHuhn/pokedexcli/internal"
)

type cliCommand struct {
	name 		string
	description string
	callback 	func(*configStruct, []string) error
}

type configStruct struct {
	nextURL 	*string 
	previousURL *string
	cache 		*pokecache.Cache
}

type locationAreasResponse struct {
	Count 		int `json:"count"`
	Next 		*string `json:"next"`
	Previous 	*string `json:"previous"`
	Results 	[]struct {
		Name 	string `json:"name"`
		URL 	string `json:"url"`
	} `json:"results"`
}

type locationAreaDetailResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

var commandMap = map[string]cliCommand{
	"exit": {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	},
	"help": {
		name : "help",
		description : "Displays a help message.",
		callback : help,
	},
	"map": {
		name: "map",
		description: "Displays the next or previous page of Pokemon.",
		callback: mapFunction,
	},
	"mapb": {
		name: "mapb",
		description: "Displays the previous page of Pokemon.",
		callback: mapBackFunction,
	},
	"explore": {
		name: "explore",
		description: "Explores the specified zone and lists the Pokemon that can be found there.",
		callback: explore,
	},
}

var commandHelp = []struct {
    name        string
    description string
}{
    {"exit", "Exit the Pokedex"},
    {"help", "Displays a help message."},
	{"map", "Displays the next page of Zones."},
	{"mapb", "Displays the previous page of Zones."},
	{"explore <zone>", "Explores the specified zone and lists the Pokemon that can be found there."},
}

// Bereinigt die Benutzereingabe, indem sie in Kleinbuchstaben 
// umgewandelt und in Wörter aufgeteilt wird
func cleanInput(text string) []string {
	list := strings.Split(text," ")
	for i := 0; i < len(list); i++ {
		if list[i] == "" || list[i] == " " {
			list = append(list[:i], list[i+1:]...)
			i--
		}
	}
	return list
}


// Beendet die Anwendung
func commandExit(cfg *configStruct, args []string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

// Gibt die liste der Befehle aus
func help(cfg *configStruct, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	for _, command := range commandHelp {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}


// Macht Get request, unmarshalt die Antwort und gibt Struct zurück
func makeRequest(cfg *configStruct, url string) ([]byte, error) {
	if cachedData, ok := cfg.cache.Get(url); ok {
		return cachedData, nil
    }
    response, err := http.Get(url)
    if err != nil {
		return nil, err
    }
    defer response.Body.Close()
    body, err := io.ReadAll(response.Body)
    if err != nil {
		return nil, err
    }
    cfg.cache.Add(url, body)
	return body, nil
}

// Gibt die Namen der Zonen der nächsten Seite aus (mapb zurück).
func mapFunction(cfg *configStruct, args []string) error {
	var url *string
	if cfg.nextURL != nil {
		url = cfg.nextURL
	} else {
		url = cfg.previousURL
	}
	body, err := makeRequest(cfg, *url)
	if err != nil {
		return err
	}
	var request locationAreasResponse
	err = json.Unmarshal(body, &request)
	if err != nil {
		return err
	}
	cfg.nextURL = request.Next
	cfg.previousURL = request.Previous
	for _, result := range request.Results {
		fmt.Printf("%s\n", result.Name)
	}
	return nil
}

// Das gleiche wie mapFunction, aber umgekehrt, d.h. es wird die vorherige Seite der Zonen ausgegeben (map zurück).
func mapBackFunction(cfg *configStruct, args []string) error {
	var url *string
	if cfg.previousURL != nil {
		url = cfg.previousURL
	} else {
		url = cfg.nextURL
	}
	body, err := makeRequest(cfg, *url)
	if err != nil {
		return err
	}
	var request locationAreasResponse
	err = json.Unmarshal(body, &request)
	if err != nil {
		return err
	}
	cfg.nextURL = request.Next
	cfg.previousURL = request.Previous
	for _, result := range request.Results {
		fmt.Printf("%s\n", result.Name)
	}
	return nil
}

// Erkundet das angegebene Gebiet nach Eingabe von "explore <zone>" 
// und gibt die Pokemon aus, die in diesem Gebiet gefunden werden können.
func explore(cfg *configStruct, args []string) error {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("usage: explore <zone>")
	}

	zone := args[0]
	fmt.Printf("Exploring %s...\n", zone)
	body, err := makeRequest(cfg, fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", zone))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error exploring zone:", err)
		return err
	}

	var request locationAreaDetailResponse
	err = json.Unmarshal(body, &request)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range request.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}