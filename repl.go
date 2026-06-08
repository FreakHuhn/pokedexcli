package main

import (
	"strings"
)

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