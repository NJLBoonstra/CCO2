package job

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JobState int

const (
	Created JobState = iota
	Running
	Completed
	Failed
)

type Job struct {
	ID              string     `json:"id"`
	State           JobState   `json:"state"`
	SortState       []JobState `json:"sortState"`
	PalindromeState []JobState `json:"palindromeState"`
	Error           string     `json:"error"`
}

type PubSubMessage struct {
	Attributes map[string]string `json:"attributes"`
	MessageId  string            `json:"messageId"`
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

func Get(jobID string, fbClient *firestore.Client, ctx context.Context) (Job, error) {
	job := Job{}
	if jobID == "" {
		return job, errors.New("jobID cannot be an empty string")
	}

	data, err := fbClient.Collection("jobs").Doc(jobID).Get(ctx)

	if err != nil && status.Code(err) == codes.NotFound {
		job.Error = "Job with ID '" + jobID + "' not found"
		return job, errors.New("Job with ID '" + jobID + "' not found")
	}

	if err != nil {
		job.Error = err.Error()
		return job, err
	}

	err = data.DataTo(&job)
	if err != nil {
		job.Error = err.Error()
		log.Printf("DataTo error: %v", err)
		log.Printf("DataTo error data: %v", data)
		return job, err
	}

	return job, nil
}

func Update(job Job, chunk int, js JobState, fbClient *firestore.Client, ctx context.Context) error {

	docRef := fbClient.Collection("jobs").Doc(job.ID)
	_, err := docRef.Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return errors.New("cannot update a non-existing Job")
	}

	// UPdate the document yeah
	// _, err = docRef.Update(ctx, []firestore.Update{
	// 	PalindromeState: js,
	// })
	return err
}

func Create(jobID string, numChunks int, fbClient *firestore.Client, ctx context.Context) (Job, error) {
	chunkStatus := make([]JobState, numChunks)
	for i := range chunkStatus {
		chunkStatus[i] = Created
	}

	j := Job{
		ID:              jobID,
		State:           Created,
		SortState:       chunkStatus,
		PalindromeState: chunkStatus,
		Error:           "",
	}
	_, err := fbClient.Collection("jobs").Doc(j.ID).Set(ctx, &j)
	if err != nil {
		log.Printf("could not create/update document %v: %v", j.ID, err)
		return j, err
	}

	return j, nil
}
