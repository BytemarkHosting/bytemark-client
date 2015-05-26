package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptYesNo provides a y/n prompt. Returns true if the user enters y, false otherwise.
func PromptYesNo(prompt string) bool {
	return Prompt(prompt+" (y/n) ") == "y"
}

// Prompt provides a string prompt, returns the entered string (including CR on Windows probably)
func Prompt(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')

	if err != nil {
		exit(err)
	}
	return strings.TrimSpace(res)
}