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
	marginSize, err := strconv.Atoi(m.Attributes["marginSize"])
	if err != nil {
		log.Fatalf("Could not convert marginSize to int: %v", err)
	}

	// read pubsub
	bucketName := m.Attributes["bucket"]
	fileName := m.Attributes["jobID"]

	log.Print("bucketName: ", bucketName)
	log.Println("fileName: ", fileName)

	chunkSize, err := strconv.Atoi(m.Attributes["chunkSize"])
	if err != nil {
		log.Fatalf("Could not convert CHUNK_SIZE to int: %v", err)
	}

	chunkIndex, _ := strconv.Atoi(m.Attributes["chunkIdx"])

	log.Println("chunkSize: ", chunkSize)
	log.Println("chunkIndex: ", chunkIndex)
	log.Println("margin: ", marginSize)

	// read from cloud storage
	chunk_bytes := make([]byte, chunkSize)
	margin_bytes := make([]byte, marginSize)
	overRead := 0
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Client could not be created", err)
	}
	defer client.Close()

	bkt := client.Bucket(bucketName)

	obj := bkt.Object(fileName)
	chunk_reader, err := obj.NewRangeReader(ctx, int64(chunkSize*chunkIndex), int64(chunkSize))
	if err != nil {
		log.Fatalf("Reader creation failed for obj: %v in bucket: %v, %v", obj, bucketName, err)
	}
	margin_reader, err := obj.NewRangeReader(ctx, int64(chunkSize*(chunkIndex+1)), int64(marginSize))
	if err != nil {
		log.Fatalf("Reader creation failed for obj: %v in bucket: %v, %v", obj, bucketName, err)
	}
	_, err = chunk_reader.Read(chunk_bytes)
	if err != nil {
		log.Fatal("Reading obj failed", err)
	}
	_, err = margin_reader.Read(margin_bytes)
	if err != nil {
		log.Fatal("Reading obj failed", err)
	}
	chunk_reader.Close()
	margin_reader.Close()

	// determine first and last newline of chunk
	chunk_string := string(chunk_bytes)
	margin_string := string(margin_bytes)
	firstNL := strings.Index(chunk_string, "\n")
	if firstNL == -1 {
		return nil
	}
	lastNL := strings.Index(margin_string, "\n")
	chunk_string += margin_string

	for lastNL == -1 {
		overRead++
		margin_reader, err = obj.NewRangeReader(ctx, int64((chunkIndex+1)*chunkSize+marginSize*overRead), int64(marginSize))
		if err != nil {
			log.Fatalf("Could not create a NewRangeReader: %v", err)
		}
		_, err = margin_reader.Read(margin_bytes)
		margin_reader.Close()
		if err != nil {
			log.Fatalf("Reading obj in iteration failed %v", err)
		}
		margin_string = string(margin_bytes)
		lastNL = strings.Index(margin_string, "\n")
		chunk_string += margin_string
	}

	cut_str := chunk_string[firstNL : (chunkIndex+1)*chunkSize+marginSize*overRead+lastNL]

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
	defer w.Close()

	fbClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Could not create a Firestore client: %v", err)
		return err
	}
	defer fbClient.Close()
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
