package main

import (
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func() error
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
}

var commandHelp = []struct {
    name        string
    description string
}{
    {"exit", "Exit the Pokedex"},
    {"help", "Displays a help message."},
}


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

func commandExit() error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func help() error {
	fmt.Printf("Welcome to the Pokedex!\n")
	for _, command := range commandHelp {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}