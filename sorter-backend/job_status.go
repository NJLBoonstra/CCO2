package sorter_backend

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	job "cco.bn.edu/shared"
	"cloud.google.com/go/firestore"
)

func JobRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET supported", http.StatusMethodNotAllowed)
		return
	}

	// Construct the jobID from the URL
	reqURL := strings.Split(r.URL.Path, "/")
	jobID := reqURL[len(reqURL)-1]

	ctx := context.Background()
	fbClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Could not create Firestore client: %v", err)
	}

	j, _ := job.Get(jobID, fbClient, ctx)

	js, err := json.Marshal(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
