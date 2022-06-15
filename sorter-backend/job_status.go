package sorter_backend

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
)

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
}

func GetJob(jobID string) Job {
	data, err := fbClient.Collection("jobs").Doc(jobID).Get(context.Background())

	job := Job{}

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
