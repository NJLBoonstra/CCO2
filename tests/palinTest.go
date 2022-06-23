package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	f, _ := os.Open("/home/niels/CCO2/tests/test1.txt")

	buffer := make([]byte, 1024*1024)

	n, _ := f.Read(buffer)
	palindromes := 0
	longest_pal := 0

	str := string(buffer[:n])
	words := strings.Split(str, " ")

	for _, w := range words {
		w = strings.Trim(w, " \n")

		if len(w) > 0 && CheckPalindrome(w) {
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
