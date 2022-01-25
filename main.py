from random import choice, randrange
from typing import Dict, List
import sys
from json import load, dump
import os


def clear_console() -> None:
    command = 'clear'
    if os.name in ('nt', 'dos'):  # If Machine is running on Windows, use cls
        command = 'cls'
    os.system(command)


def fake_artist_goes_to_ny(players: int, words: Dict[str, List[str]],
                           gui: bool = False) -> None:
    category = choice(list(words))
    word = choice(words[category])
    impostor = randrange(players)

    for i in range(players):
        clear_console()
        input("Press enter to see your role...")
        clear_console()
        if i == impostor:
            print("You are the FAKE")
        else:
            print("You are an ARTIST",
                  f"The word is {word.upper()}.")
        input("Press enter to confirm...")
        clear_console()
    print(f"Category is... {category.upper()}!")


def parse_json(filename: str) -> Dict[str, List[str]]:
    with open(filename, "r") as f:
        words = load(f)
    return words


def test_words(filename: str) -> None:
    words = {
        "food": ["Pizza", "Croissant", "Spaghetti"]
    }
    with open(filename, "w") as f:
        dump(words, f)


def main() -> None:
    if len(sys.argv) < 2:
        print("Player count not specified!")
        return
    if len(sys.argv) < 3:
        print("json file with words not specified!")
        return
    words = parse_json(sys.argv[2])
    fake_artist_goes_to_ny(int(sys.argv[1]), words)


if __name__ == '__main__':
    main()
