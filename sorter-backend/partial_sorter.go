package sorter_backend

import (
	"context"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	job "cco.bn.edu/shared"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
)

// HelloPubSub consumes a Pub/Sub message.
func PartialSort(ctx context.Context, m job.PubSubMessage) error {
	bucketName := m.Attributes["bucket"]
	fileName := m.Attributes["jobID"]

	log.Print("bucketName: ", bucketName)
	log.Println("fileName: ", fileName)

	chunkSize, err := strconv.Atoi(m.Attributes["chunkSize"])
	if err != nil {
		log.Fatalf("Could not convert CHUNK_SIZE to int: %v", err)
	}

	// read pubsub
	chunkIndex, _ := strconv.Atoi(m.Attributes["chunkIdx"])
	margin := 128

	log.Println("chunkSize: ", chunkSize)
	log.Println("chunkIndex: ", chunkIndex)
	log.Println("margin: ", margin)

	// read from cloud storage
	partialString := make([]byte, chunkSize+margin)
	extPartialString := make([]byte, margin)
	overRead := 0
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Client could not be created", err)
	}

	bkt := client.Bucket(bucketName)

	obj := bkt.Object(fileName)
	r, err := obj.NewRangeReader(ctx, int64(chunkSize)*int64(chunkIndex), int64(chunkSize)+int64(margin))
	if err != nil {
		log.Fatalf("Reader creation failed for obj: %v in bucket: %v, %v", obj, bucketName, err)
	}
	_, err = r.Read(partialString)
	if err != nil {
		log.Fatal("Reading obj failed", err)
	}

	// determine first and last newline of chunk
	str := string(partialString)
	firstNL := strings.Index(str, "\n")
	lastNL := strings.Index(str[chunkSize:], "\n")
	for lastNL == -1 {
		overRead++
		r, err = obj.NewRangeReader(ctx, int64(chunkSize)*int64(chunkIndex)+int64(margin), int64(margin))
		if err != nil {
			log.Fatalf("Could not create a NewRangeReader: %v", err)
		}
		_, err = r.Read(extPartialString)
		if err != nil {
			log.Fatalf("Reading obj in iteration failed %v", err)
		}
		str += string(extPartialString)
		lastNL = strings.Index(str[chunkSize+margin*overRead:], "\n")
	}
	cut_str := str[firstNL:lastNL]

	// split
	split_str := strings.Fields(cut_str)

	// sort
	sort.Strings(split_str)

	// merge sorts
	result := strings.Join(split_str, " ")

	// store sorting result
	newObjectName := fileName + "-" + strconv.Itoa(chunkIndex)
	resultObj := bkt.Object(newObjectName)
	w := resultObj.NewWriter(ctx)
	_, err = io.WriteString(w, result)
	if err != nil {
		log.Fatal("Writing obj failed", err)
	}

	fbClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Could not create a Firestore client: %v", err)
		return err
	}
	j, _ := job.Get(fileName, fbClient, ctx)
	if err != nil {
		log.Printf("job.Get failed: %v", err)
		return err
	}

	j.SortState[chunkIndex] = job.Completed

	err = job.Update(j, fbClient, ctx)
	if err != nil {
		log.Printf("Could not update the job: %v", err)
		return err
	}

	// determine if this is the last chunk
	// if so, create pub/sub message for merging

	// TODO

	return nil
}
