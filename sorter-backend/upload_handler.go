package sorter_backend

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

// var psClient *pubsub.Client
// var fbClient *firestore.Client
var sClient *storage.Client
var chunkSize int

func init() {
	var err error

	ctx := context.Background()

	project, exists := os.LookupEnv("GOOGLE_CLOUD_PROJECT")
	if !exists {
		log.Fatalf("Please set GOOGLE_CLOUD_PROJECT")
	}

	chunkSizeStr, exists := os.LookupEnv("CHUNK_SIZE")
	if !exists {
		log.Fatalf("Please set CHUNk_SIZE")
	}

	chunkSize, err = strconv.Atoi(chunkSizeStr)
	if err != nil {
		log.Fatalf("Could not convert chunkSize: %v", err)
	}

	psClient, err = pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Cannot create Pub/Sub client: %v", err)
	}

	fbClient, err = firestore.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Cannot create Firestore client: %v", err)
	}

	sClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Cannot create a Storage Client: %v", err)
	}
}

func HandleUpload(ctx context.Context, e GCSEvent) {
	fileName := e.Name
	bucketName := e.Bucket

	jobUUID, err := uuid.Parse(fileName)
	if err != nil {
		log.Fatalf("Could not parse provided uuid: %v", err)
	}

	obj := sClient.Bucket(bucketName).Object(fileName)
	objAttr, err := obj.Attrs(ctx)
	if err != nil {
		log.Fatalf("Could not get object attrs: %v", err)
	}

	chunks := int(math.Ceil(float64(objAttr.Size) / float64(chunkSize)))
	chunkStatus := make([]JobState, chunks)
	for i, _ := range chunkStatus {
		chunkStatus[i] = Created
	}

	j := &Job{
		ID:              jobUUID.String(),
		State:           Created,
		SortState:       chunkStatus,
		PalindromeState: chunkStatus,
	}

	js, _ := json.Marshal(j)

	_, err = fbClient.Collection("jobs").Doc(j.ID).Set(ctx, j)
	// Todo: handle error

	// Publish tasks for each chunk
	for i, _ := range chunkStatus {
		task := &pubsub.Message{
			Attributes: map[string]string{
				"jobID":    j.ID,
				"chunkIdx": strconv.Itoa(i),
				"bucket":   bucketName,
			},
			Data: js,
		}

		results := []pubsub.PublishResult{
			*psClient.Topic("sortJobs").Publish(ctx, task),
			*psClient.Topic("palindromeJobs").Publish(ctx, task),
		}

		for _, r := range results {
			msgId, err := r.Get(ctx)
			if err != nil {
				log.Fatalf("Could not publish message: %v", err)
			}

			log.Printf("Chunk %v, published message: %v", i, msgId)
		}

	}
}
