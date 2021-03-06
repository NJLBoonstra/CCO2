package sorter_backend

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	job "cco.bn.edu/shared"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
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

	jobID := strings.Split(chunkFileName, "/")[0]

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(chunkFileName)
	attrs, _ := obj.Attrs(ctx)
	myUUID, err := uuid.Parse(attrs.Metadata["palindromeWorkerID"])
	check(err, fmt.Sprintf("cannot parse palindromeWorkerID: %v", attrs.Metadata["palindromeWorkerID"]))

	if err != nil {
		log.Printf("Could not read object attributes: %v", err)
		job.UpdateWorker(jobID, myUUID, job.Failed, fbClient, ctx)
	}
	// obj_size := attrs.Size

	reader, err := obj.NewReader(ctx)

	if err != nil {
		log.Printf("Could not open reader: %v", err)
		job.UpdateWorker(jobID, myUUID, job.Failed, fbClient, ctx)
		return err
	}

	buffer, err := ioutil.ReadAll(reader)

	if err != nil {
		log.Printf("Could not read file: %v", err)
		job.UpdateWorker(jobID, myUUID, job.Failed, fbClient, ctx)
		return err
	}

	palindromes := 0
	longest_pal := 0

	str := string(buffer)
	lines := strings.Split(str, "\n")

	for _, l := range lines {
		for _, w := range strings.Split(l, " ") {
			w = strings.Trim(w, "\t \n")

			if len(w) > 0 && CheckPalindrome(w) {
				palindromes++
				if len(w) > longest_pal {
					longest_pal = len(w)
				}
			}
		}
	}

	err = job.UpdateWorker(jobID, myUUID, job.Completed, fbClient, ctx)
	if err != nil {
		log.Fatalf("Could not update job: %v", err)
		return err
	}

	err = job.UpdatePalindromeResult(jobID, myUUID, palindromes, longest_pal, fbClient, ctx)
	if err != nil {
		log.Fatalf("Could not update Palindrome result: %v", err)
		job.UpdateWorker(jobID, myUUID, job.Failed, fbClient, ctx)
		return err
	}

	log.Printf("Palindromes: %d; Longest: %d", palindromes, longest_pal)

	// determine if this is the last chunk
	// if so, create pub/sub message for merging
	allDone, _ := job.AllWorkerTypeStates(jobID, job.WorkerTypeState{Type: job.Palindrome, State: job.Completed}, fbClient, ctx)

	if allDone {
		// Last chunk, do something with merging perhaps

		res, _ := job.GetPalindromeResult(jobID, fbClient, ctx)
		workerResults := res.PalindromeWorkerResult
		longest := 0
		sum := 0
		for _, v := range workerResults {
			if v.LongestPalindrome > longest {
				longest = v.LongestPalindrome
			}
			sum += v.Palindromes
		}
		job.UpdatePalindromeJobResult(jobID, sum, longest, fbClient, ctx)

		job.SetFinish(jobID, time.Now(), "palindrome", fbClient, ctx)

		sortDone, _ := job.AllWorkerTypeStates(jobID, job.WorkerTypeState{Type: job.SorterReduce, State: job.Completed}, fbClient, ctx)

		if sortDone {
			// Delete all chunks, cuz we're always gonna be the latest worker to finish
			q := &storage.Query{
				Prefix: jobID,
			}
			chunks := bkt.Objects(ctx, q)

			for {
				o, err := chunks.Next()

				if err == iterator.Done {
					break
				}

				if err != nil {
					log.Fatalf("cannot iterate: %v", err)
				}

				err = bkt.Object(o.Name).Delete(ctx)
				if err != nil {
					log.Printf("could not delete file: %v", err)
				}
			}
		}
	}

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
