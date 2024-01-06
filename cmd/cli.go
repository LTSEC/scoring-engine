package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Cli() {

	var userInput string

	for {
		fmt.Print("Scoring/Engine:> ")
		userInput = inputParser()
		// slicing off the new line character for ease in manipulation and such
		userInput = userInput[0 : len(userInput)-1]
		fmt.Println(userInput + "\n")
		// if exit is typed, we want to exit the program
		if userInput == "exit" {
			break
		}
		userArgs := tokenizer(userInput)
		fmt.Println(userArgs)
	}

}

func inputParser() string {

	inputReader := bufio.NewReader(os.Stdin)
	userInput, err := inputReader.ReadString('\n')

	if err != nil {
		return "Something went wrong"
	} else {
		return userInput
	}

}

func tokenizer(userInput string) []string {

	return strings.Split(userInput, " ")

}
