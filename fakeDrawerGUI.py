from audioop import add
from textwrap import fill
import tkinter as tk
from fakeDrawerLogic import fake_artist_goes_to_ny, parse_json, save_words


def play_game() -> None:
    # TODO
    pass


def add_words() -> None:
    # TODO
    pass


def welcome_screen() -> None:
    welcome_screen = tk.Frame()
    tk.Label(welcome_screen, text="Welcome to Fake Artist Goes to Ne- I mean Brno!",
             width=50, height=5).pack(fill=tk.X)

    buttons = tk.Frame(welcome_screen)
    tk.Button(buttons, text="Start Game!", width=20, height=3,
              bg="green").pack(side=tk.LEFT)
    tk.Button(buttons, text="Enter word addition mode",
              width=20, height=3, bg="orange", command= add_words).pack(side=tk.RIGHT)

    buttons.pack()
    welcome_screen.pack()
    window.mainloop()


if __name__ == "__main__":
    window = tk.Tk()
    welcome_screen()
