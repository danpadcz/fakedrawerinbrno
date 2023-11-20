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

type words []map[string]string

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

	var w words
	err = json.Unmarshal(f, &w)
	if err != nil {
		fmt.Println(err)
	}

	FakeDrawerInBrnoCLI(w, cli.PlayerCount)
}

func FakeDrawerInBrnoCLI(w words, playerCount int) error {
	impostor := rand.Intn(playerCount)
	selectedWord := rand.Intn(len(w))
	category, catOk := w[selectedWord]["category"]
	word, wordOk := w[selectedWord]["text"]
	if !catOk || !wordOk {
		return errors.New("invalid Json file format")
	}

	for i := 0; i < playerCount; i++ {
		fmt.Printf("Hey there, press enter to view your role ;)")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		if i == impostor {
			fmt.Printf("You are the fake :) \n")
		} else {
			fmt.Printf("The word is: %s \n", word)
		}
		fmt.Printf("Category is %s\nPress enter to leave...", category)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Print("\033[H\033[2J")
	}
	return nil
}
