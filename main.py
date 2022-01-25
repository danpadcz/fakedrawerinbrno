from random import randint
from typing import Dict, List
import sys
import json


def fake_artist_goes_to_ny(players: int, words: Dict[str, List[str]],
                           gui: bool = False) -> None:
    pass


def parse_json(filename: str) -> Dict[str, List[str]]:
    with open(filename, "r") as f:
        loaded = f.read()

    words: Dict[str, List[str]] = {}



    return words


def main() -> None:
    if len(sys.argv) < 2:
        print("Player count not specified!")
        return
    if len(sys.argv) < 3:
        print("json file with words not specified!")
        return
    words = parse_json(sys.argv[2])
    fake_artist_goes_to_ny(int(sys.argv[1]))


if __name__ == '__main__':
    main()
