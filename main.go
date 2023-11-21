package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/alecthomas/kong"
)

type CLIParser struct {
	GUI         bool   `help:"Run in GUI"`
	PlayerCount int    `arg:"" short:"c" help:"Amount of players. Must be greater than 0"`
	JsonFile    string `arg:"" short:"f" help:"Path to JSON file to be used for words" type:"existingfile"`
}

type Words []map[string]string

func main() {
	var cli CLIParser
	err := kong.Parse(&cli, kong.Description("A game of lies and deceit... but with art!")).Error
	if err != nil {
		fmt.Println("Error occurred in parsing command line input")
		os.Exit(1)
	}

	if cli.PlayerCount <= 0 {
		fmt.Println("Player count must be specified and larger than 0!")
		os.Exit(1)
	}
	if cli.JsonFile == "" {
		fmt.Println("Json file with words must be specified!")
		os.Exit(1)
	}

	f, err := os.ReadFile(cli.JsonFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var w Words
	err = json.Unmarshal(f, &w)
	if err != nil {
		fmt.Println(err)
	}

	FakeDrawerInBrnoCLI(w, cli.PlayerCount)
}

func FakeDrawerInBrnoCLI(w Words, playerCount int) error {
	resultChan := make(chan result)
	go FakeDrawerInBrnoLogic(w, playerCount, resultChan)

	result, ok := <- resultChan
	if !ok {
		return errors.New("logic goroutine closed unexpectedly")
	} else if result.Error != nil {
		return result.Error
	}
	category := result.Message

	for result := range resultChan {
		if result.Error != nil {
			return result.Error
		}
		fmt.Printf("Hey there, press enter to view your role ;)")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Println(result.Message)
		fmt.Printf("\nCategory is %s\n\nPress enter to leave...", category)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Print("\033[H\033[2J")
	}
	return nil
}

type result struct {
    Message string
    Error error
}

func FakeDrawerInBrnoLogic(w Words, playerCount int, out chan result) {
	defer close(out)
	impostor := rand.Intn(playerCount)
	selectedWord := rand.Intn(len(w))
	category, catOk := w[selectedWord]["category"]
	word, wordOk := w[selectedWord]["text"]


	if !catOk || !wordOk {
		out <- result{Error: errors.New("invalid Json file format")}
		return
	}

	out <- result{Message: category}

	for i := 0; i < playerCount; i++ {
		if i == impostor {
			out <- result{Message: "You are the fake :)"}
		} else {
			out <- result{Message: fmt.Sprintf("The word is: %s", word)}
		}
	}
}
