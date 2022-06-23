package job

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WorkerState int
type WorkerType int

const CollectionJobName string = "jobs"

// const CollectionWorkersName string = "workers"

const (
	Created WorkerState = iota
	Running
	Reducing
	Completed
	Failed
)

const (
	Sorter WorkerType = iota
	Palindrome
	SorterReduce
	PalindromeReduce
)

type WorkerTypeState struct {
	Type  WorkerType  `json:"type"`
	State WorkerState `json:"state"`
}

type Job struct {
	ID      string                     `json:"id"`
	State   WorkerState                `json:"state"`
	Workers map[string]WorkerTypeState `json:"workers"`
	Error   string                     `json:"error"`
}

// WorkerTypeState []WorkerTypeState `json:"workerTypeState" firestore:"WorkerTypeState,omitempty"`

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

	data, err := fbClient.Collection(CollectionJobName).Doc(jobID).Get(ctx)

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

func GetList(fbClient *firestore.Client, ctx context.Context) ([]Job, error) {
	docRefs := fbClient.Collection(CollectionJobName).DocumentRefs(ctx)

	jobs := []Job{}

	for {
		docRef, err := docRefs.Next()

		if err == iterator.Done {
			break
		}

		job := Job{}

		j, err := docRef.Get(ctx)
		if err != nil {
			return nil, err
		}

		err = j.DataTo(&job)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)

	}

	return jobs, nil
}

func AddWorker(jobID string, wt WorkerType, fbClient *firestore.Client, ctx context.Context) (uuid.UUID, error) {
	workerUUID := uuid.New()

	_, err := fbClient.Collection(CollectionJobName).Doc(jobID).Update(ctx, []firestore.Update{
		{Path: "Workers." + workerUUID.String(), Value: WorkerTypeState{Type: wt, State: Created}},
	})

	if err != nil {
		return uuid.Nil, err
	}
	return workerUUID, nil
}

func SetState(jobID string, ws WorkerState, fbClient *firestore.Client, ctx context.Context) error {
	docRef := fbClient.Collection(CollectionJobName).Doc(jobID)
	_, err := docRef.Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return errors.New("cannot update a non-existing Job")
	}

	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "State", Value: ws},
	})

	return err
}

func AllWorkerTypeStates(jobID string, wts WorkerTypeState, fbClient *firestore.Client, ctx context.Context) (bool, error) {
	docRef := fbClient.Collection(CollectionJobName).Doc(jobID)
	j, err := docRef.Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return false, errors.New("cannot check workers of a non-existing Job")
	}

	job := Job{}
	err = j.DataTo(&job)
	if err != nil {
		return false, err
	}

	for _, v := range job.Workers {
		if v != wts {
			return false, nil
		}
	}

	return true, nil
}

func UpdateWorker(jobID string, workerUUID uuid.UUID, ws WorkerState, fbClient *firestore.Client, ctx context.Context) error {
	docRef := fbClient.Collection(CollectionJobName).Doc(jobID)
	j, err := docRef.Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return errors.New("cannot update a non-existing Job")
	}

	// TODO: dit misschien efficienter?
	job := Job{}
	err = j.DataTo(&job)
	if err != nil {
		return err
	}

	workerType := job.Workers[workerUUID.String()].Type

	// UPdate the document yeah
	_, err = fbClient.Collection(CollectionJobName).Doc(jobID).Update(ctx, []firestore.Update{
		{Path: "Workers." + workerUUID.String(), Value: WorkerTypeState{Type: workerType, State: ws}},
	})
	return err
}

func Create(jobID string, numChunks int, fbClient *firestore.Client, ctx context.Context) (Job, error) {
	j := Job{
		ID:      jobID,
		State:   Created,
		Workers: map[string]WorkerTypeState{},
		Error:   "",
	}
	_, err := fbClient.Collection(CollectionJobName).Doc(jobID).Set(ctx, &j)
	if err != nil {
		log.Printf("could not create/update document %v: %v", j.ID, err)
		return j, err
	}

	return j, nil
}
