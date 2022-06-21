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

	subrequest := "status"
	if len(reqURL) == 0 {
		log.Fatalf("No JobID given!")
	}

	jobID := reqURL[1]
	// log.Printf("jobID = %v", jobID)
	// log.Printf("len(reqURL) = %v", len(reqURL))
	if len(reqURL) == 3 {
		subrequest = reqURL[2]
	}
	// log.Printf("subrequest = %v", subrequest)

	ctx := context.Background()
	fbClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Could not create Firestore client: %v", err)
	}

	var js []byte
	if subrequest == "status" {
		j, _ := job.Get(jobID, fbClient, ctx)
		js, err = json.Marshal(j)
	} else if subrequest == "palindrome" {
		j, _ := job.GetPalindromeResult(jobID, fbClient, ctx)
		js, err = json.Marshal(j)
	}
	// Handle the json.Marhal error here
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
