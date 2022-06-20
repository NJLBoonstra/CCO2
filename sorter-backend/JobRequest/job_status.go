package sorter_job_request

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
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
	Error           string     `json:"error,omitempty"`
}

var psClient *pubsub.Client
var fbClient *firestore.Client

func init() {
	var exists bool

	project, exists := os.LookupEnv("GOOGLE_CLOUD_PROJECT")
	if !exists {
		log.Fatalf("Please set GOOGLE_CLOUD_PROJECT")
	}

	var err error

	ctx := context.Background()

	psClient, err = pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Cannot create a Pub/Sub Client: %v", err)
	}

	fbClient, err = firestore.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("Cannot create a Firestore client: %v", err)
	}

	log.Printf("ps %v fb %v", psClient, fbClient)
}

func GetJob(jobID string) Job {
	log.Printf("fb %v", fbClient)
	data, err := fbClient.Collection("jobs").Doc(jobID).Get(context.Background())

	job := Job{}

	if err != nil && status.Code(err) == codes.NotFound {
		job.Error = "Job with ID '" + jobID + "' not found!"
		return job
	}

	if err != nil {
		job.Error = err.Error()

		return job
	}

	err = data.DataTo(job)
	if err != nil {
		job.Error = err.Error()
		return job
	}

	return job
}

func JobRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET supported", http.StatusMethodNotAllowed)
		return
	}

	// Construct the jobID from the URL
	reqURL := strings.Split(r.URL.Path, "/")
	jobID := reqURL[len(reqURL)-1]

	j := GetJob(jobID)

	js, err := json.Marshal(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
