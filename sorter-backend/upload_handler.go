package sorter_backend

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"os"
	"strconv"

	job "cco.bn.edu/shared"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

var psClient *pubsub.Client
var fbClient *firestore.Client
var sClient *storage.Client
var chunkBucket string
var resultBucket string
var chunkSize int
var marginSize int

func init() {
	log.Print("Upload_handler: init()")
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
	chunkBucket, exists = os.LookupEnv("CHUNK_BUCKET")
	if !exists {
		log.Fatalf("Please set CHUNK_BUCKET")
	}

	resultBucket, exists = os.LookupEnv("RESULT_BUCKET")
	if !exists {
		log.Fatalf("Please set RESULT_BUCKET1")
	}

	chunkSize, err = strconv.Atoi(chunkSizeStr)
	if err != nil {
		log.Fatalf("Could not convert chunkSize: %v", err)
	}

	marginSizeStr, exists := os.LookupEnv("MARGIN_SIZE")
	if !exists {
		log.Fatalf("Please set MARGIN_SIZE")
	}

	marginSize, err = strconv.Atoi(marginSizeStr)
	if err != nil {
		log.Fatalf("Could not convert marginSize: %v", err)
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

func HandleUpload(ctx context.Context, e job.GCSEvent) error {
	fileName := e.Name
	bucketName := e.Bucket

	jobUUID, err := uuid.Parse(fileName)
	if err != nil {
		log.Printf("Could not parse provided uuid: %v", err)
		return err
	}

	obj := sClient.Bucket(bucketName).Object(fileName)
	objAttr, err := obj.Attrs(ctx)
	if err != nil {
		log.Printf("Could not get object attrs: %v", err)
		return err
	}

	chunks := int(math.Ceil(float64(objAttr.Size) / float64(chunkSize)))
	j, err := job.Create(jobUUID.String(), objAttr.Metadata["original-filename"], chunks, fbClient, ctx)
	if err != nil {
		log.Printf("Could not get the job: %v", err)
		return err
	}

	err = job.CreatePalindromeTable(jobUUID.String(), fbClient, ctx)
	if err != nil {
		log.Printf("Could not get the job: %v", err)
		return err
	}

	js, _ := json.Marshal(j)

	// Todo: handle error

	// Publish tasks for each chunk
	for i := 0; i < chunks; i++ {
		workerID, err := job.AddWorker(fileName, job.Sorter, fbClient, ctx)
		if err != nil {
			log.Printf("could not add worker %v", err)
			return err
		}
		task := &pubsub.Message{
			Attributes: map[string]string{
				"jobID":        j.ID,
				"workerID":     workerID.String(),
				"chunkIdx":     strconv.Itoa(i),
				"bucket":       bucketName,
				"resultBucket": resultBucket,
				"chunkBucket":  chunkBucket,
				"chunkSize":    strconv.Itoa(chunkSize),
				"marginSize":   strconv.Itoa(marginSize),
				"objectSize":   strconv.FormatInt(objAttr.Size, 10),
			},
			Data: js,
		}

		r := *psClient.Topic("sortJobs").Publish(ctx, task)

		msgId, err := r.Get(ctx)
		if err != nil {
			log.Printf("Could not publish message: %v", err)
			return err
		}

		log.Printf("Chunk %v, published message: %v", i, msgId)

	}
	return nil
}
