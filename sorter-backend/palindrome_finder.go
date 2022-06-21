package sorter_backend

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	job "cco.bn.edu/shared"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
)

func FindPalindromes(ctx context.Context, m job.PubSubMessage) error {
	fileName := m.Attributes["jobID"]
	bucketName := m.Attributes["bucket"]
	chunkIdx, err := strconv.Atoi(m.Attributes["chunkIdx"])
	// chunkSize, err := strconv.Atoi(m.Attributes["chunkSize"])
	if err != nil {
		log.Fatalf("Could not convert the chunkIdx to an int: %v", err)
	}

	// Currently, the palindrome implementation only runs for the whole file
	if chunkIdx > 0 {
		log.Printf("Skipping chunk %v, chunking not implemented", chunkIdx)
		return nil
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return err
	}
	defer client.Close()

	fbClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Printf("Could not create Firestore Client %v", err)
		return err
	}
	defer fbClient.Close()

	myUUID, err := job.AddWorker(fileName, job.Palindrome, fbClient, ctx)
	if err != nil {
		log.Printf("could not add worker %v", err)
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(fileName)
	attrs, err := obj.Attrs(ctx)

	if err != nil {
		log.Printf("Could not read object attributes: %v", err)
		job.UpdateWorker(fileName, myUUID, job.Failed, fbClient, ctx)
	}
	obj_size := attrs.Size

	reader, err := obj.NewReader(ctx)

	if err != nil {
		log.Printf("Could not open reader: %v", err)
		job.UpdateWorker(fileName, myUUID, job.Failed, fbClient, ctx)
		return err
	}

	buffer := make([]byte, obj_size)

	_, err = reader.Read(buffer)

	if err != nil {
		log.Printf("Could not read file: %v", err)
		job.UpdateWorker(fileName, myUUID, job.Failed, fbClient, ctx)
		return err
	}

	palindromes := 0
	longest_pal := 0

	str := string(buffer)
	words := strings.Split(str, " ")
	log.Print(words)

	for _, w := range words {
		w = strings.Trim(w, " \n")

		if CheckPalindrome(w) {
			palindromes++
			if len(w) > longest_pal {
				longest_pal = len(w)
			}
		}
	}

	err = job.UpdateWorker(fileName, myUUID, job.Completed, fbClient, ctx)
	if err != nil {
		log.Fatalf("Could not update job: %v", err)
		return err
	}

	err = job.AddPalindromeResult(fileName, palindromes, longest_pal, fbClient, ctx)
	if err != nil {
		log.Fatalf("Could not update Palindrome result: %v", err)
		job.UpdateWorker(fileName, myUUID, job.Failed, fbClient, ctx)
		return err
	}

	log.Printf("Palindromes: %d; Longest: %d", palindromes, longest_pal)
	return nil
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
