package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"bufio"

	"github.com/alecthomas/kong"
)

var CLI struct {
	CLI         bool
	PlayerCount int
	JsonFile    string
}

type words []map[string]string

func main() {
	err := kong.Parse(&CLI).Error
	if err != nil {
		fmt.Println("Error occurred in parsing command line input")
		os.Exit(1)
	}

	f, err := os.ReadFile(CLI.JsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var w words
	err = json.Unmarshal(f, &w)
	if err != nil {
		fmt.Println(err)
	}

	FakeDrawerInBrnoCLI(w, CLI.PlayerCount)
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
