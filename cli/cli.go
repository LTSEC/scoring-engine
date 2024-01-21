package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Cli() {

	var userInput string

	for {
		var currDirectory, err = os.Getwd()
		if err != nil {
			fmt.Println("directory error")
		}
		fmt.Print(currDirectory + "$ ")
		userInput = inputParser()
		// slicing off the new line character for ease in manipulation and such
		userInput = userInput[0 : len(userInput)-1]
		// fmt.Println(userInput + "\n")
		// if exit is typed, we want to exit the program
		if userInput == "exit" {
			break
		}
		userArgs := tokenizer(userInput)

		commandSelector(userArgs)
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

// switch statement for command selection
func commandSelector(tokenizedInput []string) {

	// the switch acts on the first word of the command
	// the idea is that you'd pass the remaining args to the requisit functions
	switch tokenizedInput[0] {
	// test case
	case "hello":
		fmt.Println("it was hello!")
	// default case to pipe into bash
	default:
		bashInjection(tokenizedInput)
	}
}

// function for injecting commands into bash
func bashInjection(command []string) {

	// run command guy with exec
	// the .. thing lets you pass a slice as if it were a hard-coded , separated list
	cmd := exec.Command(command[0], command[1:]...)
	// force the output of cmd to be regular stdout
	cmd.Stdout = os.Stdout

	// check for error and print
	if err := cmd.Run(); err != nil {
		fmt.Println("couldn't run the guy", err)
	}
}
