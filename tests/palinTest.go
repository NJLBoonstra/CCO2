package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	f, _ := os.Open("../../alice29.txt")

	buffer := make([]byte, 1024*1024)

	_, _ = f.Read(buffer)
	palindromes := 0
	longest_pal := 0

	str := string(buffer)
	words := strings.Split(str, " ")
	log.Print(words)

	for _, w := range words {
		w = strings.Trim(w, " \n")

		if len(w) > 2 && CheckPalindrome(w) {
			log.Print("palindrome: ", w)
			palindromes++
			if len(w) > longest_pal {
				longest_pal = len(w)
			}
		}
	}

	log.Printf("palins %v, longest %v", palindromes, longest_pal)
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