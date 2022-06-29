from os import path


def find_palindromes(filename: str):
    full_path: str = path.abspath(filename)

    buffer: str = ""
    with open(full_path, "r") as f:
        buffer = f.read()

    buffer = buffer.split(" ")

    palindromes: int = 0
    longest_palindrome: int = 0

    word: str
    for word in buffer:
        word = word.strip().strip(" \n")

        if len(word) > 0 and word == word[::-1]:
            palindromes += 1
            print(word)
            if len(word) > longest_palindrome:
                longest_palindrome = len(word)

    return palindromes, longest_palindrome


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("file", type=str)

    args = parser.parse_args()

    print(f"Filename: {args.file}")

    palindromes, longest_palindrome = find_palindromes(args.file)

    print(f"Palindromes: {palindromes}, longest: {longest_palindrome}")
