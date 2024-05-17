package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/LTSEC/scoring-engine/config"
	"github.com/LTSEC/scoring-engine/score_holder"
	"github.com/LTSEC/scoring-engine/scoring"
)

var yamlConfig *config.Yaml
var ScoringStarted = false

// The CLI takes in user input from stdin to execute predetermined commands.
// This is intended to be the primary method of control for the scoring engine.
//
// Any input is tokenized into a slice, of which the first word is meant to act as the command.
// The subsequent inputs are meant to be passed to a later function that is called by the command if applicable.
//
// If input does not match any commands for the engine, then the entire command is passed into bash for handling.
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
		userInput = strings.TrimSuffix(userInput, "\r\n")
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

	HelpOutput := `Available commands:
	hello (Testing output)
	help (Outputs some helpful information)
	config (Recieves a path and parses the yaml config given)
	checkconfig (Outputs the currently parsed yaml config)
	score (Starts the scoring engine)
	togglescore (Toggles the scoring engine on/off if its already started)
	exit (exits the CLI)
	`

	// the switch acts on the first word of the command
	// the idea is that you'd pass the remaining args to the requisit functions
	switch tokenizedInput[0] {
	// test case
	case "hello":
		fmt.Println("it was hello!")
	case "help":
		fmt.Println(HelpOutput)
	case "config":
		if len(tokenizedInput) != 2 {
			fmt.Println("config requires a path")
		} else {
			yamlConfig = config.Parse(tokenizedInput[1])
		}
	case "checkconfig":
		fmt.Printf("%+v\n", yamlConfig)
	case "score":
		if ScoringStarted == false {
			ScoringStarted = true
			if yamlConfig != nil {
				go scoring.ScoringStartup(yamlConfig)
			} else {
				fmt.Println("Provide a config first.")
			}
		} else {
			fmt.Println("The scoring engine has already been initalized.")
		}
	case "togglescore":
		if ScoringStarted == true {
			if scoring.ScoringOn == true {
				scoring.ScoringToggle(false, yamlConfig)
			} else if scoring.ScoringOn == false {
				go scoring.ScoringToggle(true, yamlConfig)
			}
		} else {
			fmt.Println("Initalize the scoring engine first.")
		}
	case "listscore":
		if ScoringStarted == true {
			for TeamName, TeamData := range score_holder.GetMap() {
				for _, Data := range TeamData {
					for dtype, Score := range Data {
						if reflect.TypeOf(Score) == reflect.TypeOf(1) {
							fmt.Println(fmt.Sprint(TeamName, ": ", dtype, " : ", Score))
						}
					}
				}
			}
		} else {
			fmt.Println("Initalize the scoring engine first.")
		}
	case "exit":
		os.Exit(0)

	default:
		bashInjection(tokenizedInput)
	}
}

// function for injecting commands into bash
func bashInjection(command []string) {

	// run command guy with exec
	// the .. thing lets you pass a slice as if it were a hard-coded , separated list
	if command[0] != "cd" {
		cmd := exec.Command(command[0], command[1:]...)
		// force the output of cmd to be regular stdout
		cmd.Stdout = os.Stdout

		// check for error and print
		if err := cmd.Run(); err != nil {
			fmt.Println("couldn't run the guy", err)
		}
	} else {
		if len(command) == 2 {
			os.Chdir(command[1])
		} else {
			fmt.Println("please include dir")
		}
	}
}
