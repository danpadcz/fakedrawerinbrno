package main

import (
	"fmt"
	"testing"

	"github.com/smarty/assertions"
)

func TestLogicNoWords(t *testing.T) {
	var w Words
	resultChan := make(chan result)

	go FakeDrawerInBrnoLogic(w, 3, resultChan)

	r := <-resultChan
	if assertions.ShouldBeError(r.Error, "words map cannot be empty")  != "" {
		t.Fatalf(fmt.Sprintf("Expected 'words map cannot be empty' error, got '%s'", r.Error))
	}
}

func TestLogicInvalidPlayerCount(t *testing.T) {
	var w Words
	w = append(w, Word{Word: "cat", Category: "animal"})
	resultChan := make(chan result)
	go FakeDrawerInBrnoLogic(w, -1, resultChan)

	r := <-resultChan
	if assertions.ShouldBeError(r.Error, "player count cannot be less than 1")  != "" {
		t.Fatalf(fmt.Sprintf("Expected 'player count cannot be less than 1' error, got '%s'", r.Error))
	}

	resultChan = make(chan result)
	go FakeDrawerInBrnoLogic(w, 0, resultChan)

	r = <-resultChan
	if assertions.ShouldBeError(r.Error, "player count cannot be less than 1") != "" {
		t.Fatalf(fmt.Sprintf("Expected 'player count cannot be less than 1' error, got '%s'", r.Error))
	}
}

func TestLogicValid(t *testing.T) {
	var w Words
	w = append(w, Word{Word: "cat", Category: "animal"})
	resultChan := make(chan result)
	go FakeDrawerInBrnoLogic(w, 3, resultChan)

	result, ok := <-resultChan
	if !ok {
		t.Fatalf("goroutine closed unexpectedly")
	} else if result.Error != nil {
		t.Fatalf(result.Error.Error())
	}
	if result.Message != "animal" {
		t.Fatalf(fmt.Sprintf("Expected category 'animal' got '%s'", result.Message))
	}

	impostorAppeared := false
	for result := range resultChan {
		if result.Error != nil {
			t.Fatalf(result.Error.Error())
		}
		if result.Message == "You are the fake :)" && !impostorAppeared{
			impostorAppeared = true
		} else if result.Message == "You are the fake :)" {
			t.Fatalf("impostor appeared twice")
		} else if result.Message != "The word is: cat" {
			t.Fatalf("Given incorrect string, expected 'The word is: cat' gotten '%s'", result.Message)
		}

	}
}
