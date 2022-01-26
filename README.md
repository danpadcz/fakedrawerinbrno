# fakedrawerinbrno

This is a terminal implementation of the game "A fake artist goes to New York"
I created this for me and my friends, but also as to try creating my first github project and just to mess around.

## Usage

The script takes 2 arguments, first one is the amount of players and the second one is the path to the json file with the prompts
So for example `python3 fakeDrawerInBrno.py 3 Words/words.json` starts a 3 player game with the json file at Words/words.json

After the script finishes, it creates a copy of the used json file at the same place as the script, but without the word that was chosen for the game. The file will be called `shortened.json` so be careful that the script doesn't rewrite some of your data.

## JSON file format

The file includes one list of json objects with two keys: "category" and "text". The category is what is shown to all players, while text is shown to all but the fake.

This is an example of the correct format of the json file: `[{"category":"Games", "text":"Bang!"}, {"category":"Games", "text":"A Fake Artist Goes To New York"}]`

Have fun! :D
