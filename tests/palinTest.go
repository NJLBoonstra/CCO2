package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	buffer, _ := os.ReadFile("./results/result1.txt")

	palindromes := 0
	longest_pal := 0
	longest_par_str := ""

	str := string(buffer)
	lines := strings.Split(str, "\n")

	for _, l := range lines {
		for _, w := range strings.Split(l, " ") {
			w = strings.Trim(w, "\t \n")

			if len(w) > 0 && CheckPalindrome(w) {
				palindromes++
				if len(w) > longest_pal {
					longest_pal = len(w)
					longest_par_str = w
				}
			}
		}
	}

	log.Printf("palins %v, longest %v", palindromes, longest_pal)
	log.Println(longest_par_str)
}

func CheckPalindrome(word string) bool {
	if word == "" {
		return false
	}

	for i := 0; i < (len(word)/2)+1; i++ {
		j := len(word) - i - 1
		if word[i] != word[j] {
			return false
		}
	}

	return true
}
