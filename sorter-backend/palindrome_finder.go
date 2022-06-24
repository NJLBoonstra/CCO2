package sorter_backend

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"

	job "cco.bn.edu/shared"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
)

func FindPalindromes(ctx context.Context, e job.GCSEvent) error {
	chunkFileName := e.Name
	bucketName := e.Bucket

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

	fileName := strings.Split(chunkFileName, "/")[0]

	myUUID, err := job.AddWorker(fileName, job.Palindrome, fbClient, ctx)
	if err != nil {
		log.Printf("could not add worker %v", err)
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(fileName)
	// attrs, err := obj.Attrs(ctx)

	if err != nil {
		log.Printf("Could not read object attributes: %v", err)
		job.UpdateWorker(fileName, myUUID, job.Failed, fbClient, ctx)
	}
	// obj_size := attrs.Size

	reader, err := obj.NewReader(ctx)

	if err != nil {
		log.Printf("Could not open reader: %v", err)
		job.UpdateWorker(fileName, myUUID, job.Failed, fbClient, ctx)
		return err
	}

	buffer, err := ioutil.ReadAll(reader)

	if err != nil {
		log.Printf("Could not read file: %v", err)
		job.UpdateWorker(fileName, myUUID, job.Failed, fbClient, ctx)
		return err
	}

	palindromes := 0
	longest_pal := 0

	str := string(buffer)
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
