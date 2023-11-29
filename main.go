package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

	if cli.GUI {
		if err := FakeDrawerInBrnoGUI(w, cli.PlayerCount); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := FakeDrawerInBrnoCLI(w, cli.PlayerCount); err != nil {
			fmt.Println(err)
		}
	}
}

// UI still very much placeholder, I want to add UI to add words,
// make it nicer and less slideshow-like
func FakeDrawerInBrnoGUI(w Words, playerCount int) error {
	resultChan := make(chan result)
	go FakeDrawerInBrnoLogic(w, playerCount, resultChan)

	result, ok := <-resultChan
	if !ok {
		return errors.New("logic goroutine closed unexpectedly")
	} else if result.Error != nil {
		return result.Error
	}
	category := fmt.Sprintf("Category is: %s", result.Message)

	a := app.New()
	win := a.NewWindow("A fake artist in Brno")

	title := widget.NewLabel("Hey there player, press ok to view your role!")
	titleCat := widget.NewLabel("")
	nextIterShowsRole := true
	win.SetContent(container.NewVBox(
		title,
		titleCat,
		widget.NewButton("Ok!", func() {
			if !nextIterShowsRole {
				if playerCount > 0 {
					title.SetText("Hey there player, press ok to view your role!!")
				} else {
					title.SetText("Enjoy the game!")
				}
				titleCat.SetText("")
				nextIterShowsRole = true
			} else {
				result, ok = <-resultChan
				if !ok {
					win.Close()
				} else if result.Error != nil {
					win.Close()
				} else if nextIterShowsRole {
					title.SetText(result.Message)
					titleCat.SetText(category)

					nextIterShowsRole = false
					playerCount -= 1
				}
			}
		}),
	))
	win.ShowAndRun()

	return nil
}

func FakeDrawerInBrnoCLI(w Words, playerCount int) error {
	resultChan := make(chan result)
	go FakeDrawerInBrnoLogic(w, playerCount, resultChan)

	result, ok := <-resultChan
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
		fmt.Print("\033[H\033[2J")
		fmt.Printf("Hey there, press enter to view your role ;)\n")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		fmt.Println(result.Message)
		fmt.Printf("\nCategory is: %s\n\nPress enter to leave...", category)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	return nil
}

type result struct {
	Message string
	Error   error
}

func FakeDrawerInBrnoLogic(w Words, playerCount int, out chan result) {
	defer close(out)
	if len(w) == 0 {
		out <- result{Error: errors.New("words map cannot be empty")}
		return
	}
	if playerCount <= 0 {
		out <- result{Error: errors.New("player count cannot be less than 1")}
		return
	}
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
