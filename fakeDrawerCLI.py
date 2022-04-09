import sys
from fakeDrawerLogic import parse_json, fake_artist_goes_to_ny

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
