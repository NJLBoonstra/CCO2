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

func GetJob(jobID string) (j *Job, err error) {
	data, err := fbClient.Collection("jobs").Doc(jobID).Get(context.Background())

	if err != nil {
		return nil, err
	}

	err = data.DataTo(&j)
	if err != nil {
		return nil, err
	}

	return j, err
}

func JobRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET supported", http.StatusMethodNotAllowed)
		return
	}

	// Construct the jobID from the URL
	reqURL := strings.Split(r.URL.Path, "/")
	jobID := reqURL[len(reqURL)-1]

	if jobID == "" {
		// Generate new UUID?
		http.Error(w, "Please supply a job ID!", http.StatusBadRequest)
		return
	}

	j, err := GetJob(jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
