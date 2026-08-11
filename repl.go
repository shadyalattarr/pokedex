package main

import "strings"

// to split on white space and lower all
func cleanInput(text string) []string {
	// var input []string
	input := strings.Fields(strings.ToLower(text))
	return input
}
