package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func(*configStruct) error
}

type configStruct struct {
	nextURL *string 
	previousURL *string
}

type requestStruct struct {
	Count int `json:"count"`
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		URL string `json:"url"`
	} `json:"results"`
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
}

var commandHelp = []struct {
    name        string
    description string
}{
    {"exit", "Exit the Pokedex"},
    {"help", "Displays a help message."},
	{"map", "Displays the next page of Zones."},
	{"mapb", "Displays the previous page of Zones."},
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
func commandExit(cfg *configStruct) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

// Gibt die liste der Befehle aus
func help(cfg *configStruct) error {
	fmt.Printf("Welcome to the Pokedex!\n")
	for _, command := range commandHelp {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}


// Macht Get request, unmarshalt die Antwort und gibt Struct zurück
func makeRequest(url string) (requestStruct, error) {
	var request requestStruct
	response, err := http.Get(url)
	if err != nil {
		return request, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&request)
	if err != nil {
		return request, err
	}
	return request, nil
}

// Gibt die Namen der Zonen der nächsten Seite aus (mapb zurück).
func mapFunction(cfg *configStruct) error {
	var url *string
	if cfg.nextURL != nil {
		url = cfg.nextURL
	} else {
		url = cfg.previousURL
	}
	request, err := makeRequest(*url)
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

func mapBackFunction(cfg *configStruct) error {
	var url *string
	if cfg.previousURL != nil {
		url = cfg.previousURL
	} else {
		url = cfg.nextURL
	}
	request, err := makeRequest(*url)
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
