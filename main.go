package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/kong"
)

type CLIParser struct {
	// Run struct{} `default:"1" cmd:"" help:"Run the fake drawer app"`
	Play PlayParser `cmd:"" help:"Play A fake drawer in Brno"`
	Add  AddParser  `cmd:"" help:"Add words to JSON file to be used in runs of the game"`
}
type PlayParser struct {
	GUI         bool   `help:"Run in GUI"`
	PlayerCount int    `arg:"" short:"c" help:"Amount of players. Must be greater than 0"`
	JsonFile    string `arg:"" short:"f" help:"Path to JSON file to be used for words" type:"existingfile"`
}
type AddParser struct {
	GUI      bool   `help:"Run in GUI"`
	JsonFile string `arg:"" short:"f" help:"Path to JSON file to have words added" type:"existingfile"`
}

type Words []Word
type Word struct {
	Word     string `json:"text"`
	Category string `json:"category"`
}

func main() {
	var cli CLIParser
	ctx := kong.Parse(&cli, kong.Description("A game of lies and deceit... but with art!"))
	if ctx.Error != nil {
		fmt.Println("Error occurred in parsing command line input")
		os.Exit(1)
	}
	if err := ctx.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

func (p *PlayParser) Run(ctx *kong.Context) error {
	if p.PlayerCount <= 0 {
		return errors.New("player count must be specified and larger than 0")
	}
	if p.JsonFile == "" {
		return errors.New("json file with words must be specified")
	}

	w, err := loadJsonFile(p.JsonFile)
	if err != nil {
		return err
	}

	if p.GUI {
		if err := FakeDrawerInBrnoGUI(w, p.PlayerCount, app.New()); err != nil {
			return err
		}
	} else {
		if err := FakeDrawerInBrnoCLI(w, p.PlayerCount); err != nil {
			return err
		}
	}
	return nil
}

func loadJsonFile(path string) (Words, error) {
	var w Words

	f, err := os.ReadFile(path)
	if err != nil {
		return w, err
	}

	if err := json.Unmarshal(f, &w); err != nil {
		return w, err
	}
	return w, nil
}

func saveToJsonFile(path string, w Words) error {
	input, err := json.Marshal(w)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, input, 0o660); err != nil {
		return err
	}
	return nil
}

func (a *AddParser) Run(ctx *kong.Context) error {
	if a.GUI {
		return addWordsToJSONGUI(a.JsonFile, app.New())
	}
	return addWordsToJSONCLI(a.JsonFile)
}

func loadCategories(w Words) []string {
	cats := make(map[string]bool)
	for i := range w {
		cats[w[i].Category] = true
	}

	result := make([]string, len(cats))
	i := 0
	for k := range cats {
		result[i] = k
		i++
	}
	return result
}

func addWordsToJSONGUI(path string, a fyne.App) error {
	w, err := loadJsonFile(path)
	if err != nil {
		return err
	}

	win := a.NewWindow("A fake artist in Brno word adder")
	categories := widget.NewLabel(fmt.Sprintf("Categories in inputed JSON file: %s", loadCategories(w)))
	addWordQuestion := widget.NewLabel("Do you want to add another word?")
	wordTitle := widget.NewLabel("Enter word:")
	wordEntry := widget.NewEntry()
	catTitle := widget.NewLabel("Enter category of word:")
	catEntry := widget.NewEntry()
	winContainer := container.NewVBox()
	entryContainer := container.NewVBox()

	confirmationContainer := container.NewVBox(
		addWordQuestion,
		container.NewHBox(
			widget.NewButton("Yes", func() {
				winContainer.RemoveAll()
				winContainer.Add(entryContainer)
			}),
			widget.NewButton("No", func() {
				winContainer.RemoveAll()
				if err = saveToJsonFile(path, w); err != nil {
					winContainer.Add(widget.NewLabel("Failed to save JSON file!"))
					winContainer.Add(widget.NewLabel(err.Error()))
				} else {
					winContainer.Add(widget.NewLabel("Successfully saved words to JSON file!"))
				}
				winContainer.Add(widget.NewButton("Close", win.Close))
			})))
	enterButton := widget.NewButton("Enter", func() {
		if strings.TrimSpace(catEntry.Text) == "" || strings.TrimSpace(wordEntry.Text) == "" {
			winContainer.Add(widget.NewLabel("Entry boxes cannot be empty!\nPlease fill out both boxes and try again"))
		} else {
			w = append(w,
				Word{
					Category: strings.TrimSpace(catEntry.Text),
					Word:     strings.TrimSpace(wordEntry.Text)})
			catEntry.SetText("")
			wordEntry.SetText("")
			winContainer.RemoveAll()
			winContainer.Add(confirmationContainer)
		}
	})

	entryContainer = container.NewVBox(categories, wordTitle, wordEntry, catTitle, catEntry, enterButton)

	winContainer.Add(entryContainer)

	win.SetContent(winContainer)
	win.ShowAndRun()

	return err
}

func addWordsToJSONCLI(path string) error {
	w, err := loadJsonFile(path)
	if err != nil {
		return err
	}

	fmt.Printf("Categories in inputed JSON file: %s\n", loadCategories(w))

	run := true
	reader := bufio.NewReader(os.Stdin)

	for run {
		var toAdd Word
		fmt.Println("Enter word to add: ")
		temp, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		toAdd.Word = string(temp)

		fmt.Println("Enter category of word to add: ")
		temp, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		toAdd.Category = string(temp)

		w = append(w, toAdd)

		fmt.Println("Do you want to add another word? (Y/n) ")
		temp, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		input := string(temp)
		if strings.Contains(input, "n") || strings.Contains(input, "N") {
			run = false
		} else if !strings.Contains(input, "y") && !strings.Contains(input, "Y") {
			return errors.New("invalid input given")
		}
	}

	if err := saveToJsonFile(path, w); err != nil {
		return err
	}
	return nil
}

func FakeDrawerInBrnoGUI(w Words, playerCount int, a fyne.App) error {
	resultChan := make(chan result)
	go FakeDrawerInBrnoLogic(w, playerCount, resultChan)

	result, ok := <-resultChan
	if !ok {
		return errors.New("logic goroutine closed unexpectedly")
	} else if result.Error != nil {
		return result.Error
	}
	category := fmt.Sprintf("Category is: %s", result.Message)

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
	category := w[selectedWord].Category
	word := w[selectedWord].Word

	out <- result{Message: category}

	for i := 0; i < playerCount; i++ {
		if i == impostor {
			out <- result{Message: "You are the fake :)"}
		} else {
			out <- result{Message: fmt.Sprintf("The word is: %s", word)}
		}
	}
}
