package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/alecthomas/kong"
)

type CLIParser struct {
	App AppParser `cmd:"" default:"1" help:"Run GUI app"`
	CLI struct {
		Play PlayParser `cmd:"" help:"Play A fake drawer in Brno"`
		Add  AddParser  `cmd:"" help:"Add words to JSON file to be used in runs of the game"`
	} `cmd:"" help:"Run individual modules from command line"`
}
type AppParser struct{}
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

func (a *AppParser) Run(ctx *kong.Context) error {
	app := app.New()
	win := app.NewWindow("A fake drawer in Brno")
	win.Resize(fyne.NewSize(800, 600))
	var content *fyne.Container
	var result error
	content = container.NewVBox(widget.NewCard("Welcome to a fake drawer in Brno!", "Please select an option:", nil),
		container.NewHBox(
			widget.NewButtonWithIcon("Play", theme.MediaPlayIcon(), func() {
				content.RemoveAll()
				var path string
				playerCount := widget.NewEntry()
				var fileButton *widget.Button
				fileDialog := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
					if err != nil || uc == nil {
						return
					}
					path = uc.URI().Path()
					fileButton.SetText("File selected!")
				}, win)
				fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
				fileButton = widget.NewButton("Select file", func() {
					fileDialog.Show()
				})
				form := widget.NewForm(&widget.FormItem{Text: "Path to word JSON file", Widget: fileButton},
					&widget.FormItem{Text: "Player count", Widget: playerCount})
				form.OnCancel = win.Close
				form.OnSubmit = func() {
					w, err := loadJsonFile(strings.TrimSpace(path))
					if err != nil {
						e := dialog.NewError(err, win)
						e.Show()
						return
					}
					playerCountInt, err := strconv.Atoi(playerCount.Text)
					if err != nil {
						e := dialog.NewError(errors.New("please enter a valid integer"), win)
						e.Show()
						return
					}
					if playerCountInt < 3 {
						e := dialog.NewError(errors.New("player count cannot be less than 3, please try again"), win)
						e.Show()
						return
					}
					content.RemoveAll()
					result = FakeDrawerInBrnoGUI(w, playerCountInt, win)
				}
				content.Add(form)
			}),

			widget.NewButtonWithIcon("Add words to JSON file", theme.ContentAddIcon(), func() {
				content.RemoveAll()
				var path string
				var fileButton *widget.Button
				fileButton = widget.NewButton("Select file", func() {
					fileDialog := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
						if err != nil || uc == nil {
							return
						}
						fileButton.SetText("File selected!")
						path = uc.URI().Path()
					}, win)
					fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
					fileDialog.Show()
				})
				form := widget.NewForm(&widget.FormItem{
					Text:   "Please select JSON file to edit",
					Widget: fileButton})
				form.OnCancel = win.Close
				form.OnSubmit = func() {
					result = addWordsToJSONGUI(path, win)
				}
				content.Add(form)
			}),

			widget.NewButtonWithIcon("Game help", theme.HelpIcon(), func() {
				info := dialog.NewInformation("How to use",
					"This is a game moderator for the board game A fake artist goes to New York.\nThe concept of the game is that out of all players there is one"+
						" fake artist and the rest of the players are normal artists.\nAll players get a word category and only the normal artists get"+
						" the actual word they have to draw.\nThen all players take turns drawing maximum one line on a shared piece of paper twice.\nThat is"+
						" there will be two rounds of drawing one line. After these two rounds finish\neverybody votes on who they think is the fake. Then the fake"+
						" reveals themselves and guesses what the word was.\nThe fake wins if they aren't guessed by the other players or if they guess the correct word."+
						"\n\nThis app uses JSON files to load the words for the game. \nTo play the game or add words to a json file, you will be prompted\n"+
						"to select a fitting JSON file.",
					win)
				info.Show()
			}),

			widget.NewButtonWithIcon("Quit", theme.CancelIcon(), app.Quit)))
	win.SetContent(content)
	win.ShowAndRun()
	return result
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
		app := app.New()
		win := app.NewWindow("A fake drawer in Brno")
		if err := FakeDrawerInBrnoGUI(w, p.PlayerCount, win); err != nil {
			return err
		}
		win.ShowAndRun()
	} else {
		if err := FakeDrawerInBrnoCLI(w, p.PlayerCount); err != nil {
			return err
		}
	}
	return nil
}

func (a *AddParser) Run(ctx *kong.Context) error {
	if !a.GUI {
		return addWordsToJSONCLI(a.JsonFile)
	}

	app := app.New()
	win := app.NewWindow("A fake drawer in Brno")
	if err := addWordsToJSONGUI(a.JsonFile, win); err != nil {
		return err
	}
	win.ShowAndRun()
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

func addWordsToJSONGUI(path string, win fyne.Window) error {
	w, err := loadJsonFile(path)
	if err != nil {
		return err
	}

	categories := widget.NewLabel(fmt.Sprintf("Categories in inputed JSON file: %s", loadCategories(w)))
	addWordQuestion := widget.NewLabel("Do you want to add another word?")
	wordEntry := widget.NewEntry()
	catEntry := widget.NewEntry()
	form := widget.NewForm(
		&widget.FormItem{Text: "Word", Widget: wordEntry},
		&widget.FormItem{Text: "Category", Widget: catEntry})
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

	form.OnSubmit = func() {
		if strings.TrimSpace(catEntry.Text) == "" || strings.TrimSpace(wordEntry.Text) == "" {
			e := dialog.NewError(errors.New("entry boxes cannot be empty"), win)
			e.Show()
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
	}
	form.OnCancel = win.Close

	entryContainer = container.NewVBox(categories, form)

	winContainer.Add(entryContainer)

	win.SetContent(winContainer)

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

func FakeDrawerInBrnoGUI(w Words, playerCount int, win fyne.Window) error {
	resultChan := make(chan result)
	go FakeDrawerInBrnoLogic(w, playerCount, resultChan)

	result, ok := <-resultChan
	if !ok {
		return errors.New("logic goroutine closed unexpectedly")
	} else if result.Error != nil {
		return result.Error
	}
	category := fmt.Sprintf("Category is: %s", result.Message)

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
