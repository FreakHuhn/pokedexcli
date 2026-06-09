package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		input := cleanInput(strings.ToLower(scanner.Text()))
		if command, ok := commandMap[input[0]]; ok {
			if err := command.callback(); err != nil {
				fmt.Fprintln(os.Stderr, "Error executing command:", err)
			}
		} else {
			fmt.Fprintln(os.Stderr, "Unknown command:", input[0])
		}
	}
}