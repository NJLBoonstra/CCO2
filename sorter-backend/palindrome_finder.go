package sorter_backend

import (
	"context"
	"log"
	"strconv"

	job "cco.bn.edu/shared"
	"cloud.google.com/go/storage"
)

func FindPalindromes(ctx context.Context, m job.PubSubMessage) {
	fileName := m.Attributes["jobID"]
	bucketName := m.Attributes["bucketName"]
	chunkIdx, err := strconv.Atoi(m.Attributes["chunkIdx"])
	// chunkSize, err := strconv.Atoi(m.Attributes["chunkSize"])
	if err != nil {
		log.Fatalf("Could not convert the chunkIdx to an int: %v", err)
	}

	// Currently, the palindrome implementation only runs for the whole file
	if chunkIdx > 0 {
		return
	}

	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(fileName)
	attrs, err := obj.Attrs(ctx)

	if err != nil {
		log.Printf("Could not read object attributes: %v", err)
	}
	obj_size := attrs.Size

	reader, err := obj.NewReader(ctx)

	if err != nil {
		log.Printf("Could not open reader: %v", err)
		return
	}

	buffer := make([]byte, obj_size)

	_, err = reader.Read(buffer)

	if err != nil {
		log.Printf("Could not read file: %v", err)
	}

	palindromes := 0
	longest_pal := 0

	word_start := -1
	var word string
	// Iterate over values in buffer to construct words
	// Ik heb gekozen om de []byte niet om te zetten naar string, want dat
	// leek me een onnodige extra stap
	for i, v := range buffer {
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			// Word has started, or inside a word
			if word_start < 0 {
				word_start = i
			}
		} else {
			if word_start > 0 {
				// Word *may* have ended
				word = string(buffer[word_start : i+1])
			}
		}

		if len(word) > 0 {
			pal := CheckPalindrome(word)

			if pal {
				palindromes++

				if len(word) > longest_pal {
					longest_pal = len(word)
				}
			}
			// Reset word
			word = ""
		}
	}

	log.Printf("Palindromes: %d; Longest: %d", palindromes, longest_pal)
}

func CheckPalindrome(word string) bool {
	log.Printf("CheckPalindrome: %s", word)

	for i := 0; i < (len(word)/2)+1; i++ {
		j := len(word) - i - 1
		if word[i] != word[j] {
			return false
		}
	}

	return true
}
