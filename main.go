package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    startURL := "https://pokeapi.co/api/v2/location-area"
    cfg := &configStruct{
        nextURL: &startURL,
    }
	
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "Error reading input:", err)
            continue
        }
        input := cleanInput(strings.ToLower(scanner.Text()))
        if len(input) == 0 {
            continue
        }
        if command, ok := commandMap[input[0]]; ok {
            if err := command.callback(cfg); err != nil {
                fmt.Fprintln(os.Stderr, "Error executing command:", err)
            }
        } else {
            fmt.Fprintln(os.Stderr, "Unknown command:", input[0])
        }
    }
}