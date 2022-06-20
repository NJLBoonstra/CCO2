package sorter_handle_upload

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"cco.bn.edu/sorter_job_request"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type PubSubMessage struct {
	Attributes map[string]string `json: "attributes"`
	MessageId  string            `json: "messageId"`
}

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

var psClient *pubsub.Client
var fbClient *firestore.Client
var sClient *storage.Client
var chunkSize int

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

func HandleUpload(ctx context.Context, e GCSEvent) error {
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
	chunkStatus := make([]sorter_job_request.JobState, chunks)
	for i, _ := range chunkStatus {
		chunkStatus[i] = sorter_job_request.Created
	}

	j := &sorter_job_request.Job{
		ID:              jobUUID.String(),
		State:           sorter_job_request.Created,
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
				log.Printf("Could not publish message: %v", err)
				return err
			}

			log.Printf("Chunk %v, published message: %v", i, msgId)
		}

	}
	return nil
}
