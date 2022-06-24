package sorter_backend

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	job "cco.bn.edu/shared"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func sort_lines(s string) string {
	split_str := strings.Split(s, "\n")
	sort.Strings(split_str)
	return strings.Join(split_str, "\n")
}

func check(e error, message string) {
	if e != nil {
		log.Fatalf("Error: %v. err: %v", message, e)
	}
}

// HelloPubSub consumes a Pub/Sub message.
func PartialSort(ctx context.Context, m job.PubSubMessage) error {
	marginSize, err := strconv.Atoi(m.Attributes["marginSize"])
	check(err, "Could not convert marginSize to int")
	bucketName := m.Attributes["bucket"]
	fileName := m.Attributes["jobID"]
	chunkBucket := m.Attributes["chunkBucket"]
	resultBucket := m.Attributes["resultBucket"]
	chunkSize, err := strconv.Atoi(m.Attributes["chunkSize"])
	check(err, "Could not convert CHUNK_SIZE to int")
	chunkIndex, err := strconv.Atoi(m.Attributes["chunkIdx"])
	check(err, "Could not convert chunk index to int")
	objectSize, err := strconv.Atoi(m.Attributes["objectSize"])
	check(err, "Could not convert objectSize to int")
	// read from cloud storage
	client, err := storage.NewClient(ctx)
	check(err, "Client could not be created")
	defer client.Close()

	fbClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Could not create a Firestore client: %v", err)
		return err
	}
	defer fbClient.Close()

	myUUID, err := job.AddWorker(fileName, job.Sorter, fbClient, ctx)
	if err != nil {
		log.Printf("could not add worker %v", err)
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(fileName)

	EOF := false

	if (chunkIndex+1)*chunkSize+marginSize > objectSize {
		chunkSize = objectSize - chunkIndex*chunkSize
		EOF = true
	}

	// chunk_bytes := make([]byte, chunkSize)
	overRead := 0

	chunk_reader, err := obj.NewRangeReader(ctx, int64(chunkSize*chunkIndex), int64(chunkSize))
	check(err, "Reader creation failed for obj")
	slurp, err := ioutil.ReadAll(chunk_reader)
	check(err, "Reading obj failed")
	chunk_reader.Close()

	chunk_string := string(slurp)

	firstNL := 0
	if chunkIndex != 0 {
		firstNL = strings.Index(chunk_string, "\n")
	}

	// This chunk contains no NL
	if firstNL == -1 {
		return nil
	}

	lastNL := len(chunk_string)

	if !EOF {
		margin_reader, err := obj.NewRangeReader(ctx, int64(chunkSize*(chunkIndex+1)), int64(marginSize))
		check(err, "Reader creation failed for obj")
		margin_bytes, err := ioutil.ReadAll(margin_reader)
		check(err, "Reading obj failed")
		defer margin_reader.Close()

		margin_string := string(margin_bytes)
		chunk_string += margin_string
		lastNL = strings.Index(margin_string, "\n")
		for lastNL == -1 {
			overRead++
			offset := int64((chunkIndex+1)*chunkSize + marginSize*overRead)
			if offset+int64(marginSize) > int64(objectSize) {
				marginSize = int(objectSize) - int(offset)
				EOF = true
				margin_reader, err = obj.NewRangeReader(ctx, offset, int64(marginSize))
				check(err, "Could not create a NewRangeReader")
				margin_bytes, err = ioutil.ReadAll(margin_reader)
				check(err, "Reading obj in iteration failed")
				defer margin_reader.Close()
				margin_string = string(margin_bytes)
				chunk_string += margin_string
				break
			}
			margin_reader, err = obj.NewRangeReader(ctx, int64((chunkIndex+1)*chunkSize+marginSize*overRead), int64(marginSize))
			check(err, "Could not create a NewRangeReader")
			margin_bytes, err = ioutil.ReadAll(margin_reader)
			check(err, "Reading obj in iteration failed")
			defer margin_reader.Close()
			margin_string = string(margin_bytes)
			lastNL = strings.Index(margin_string, "\n")
			chunk_string += margin_string
		}
		lastNL += chunkSize + marginSize*overRead
	}

	if EOF {
		lastNL = len(chunk_string)
	}

	// TODO make lastNL relative to the chunk
	cut_str := chunk_string[firstNL:lastNL]

	result := sort_lines(cut_str)

	// store sorting result
	newObjectName := fileName + "/" + strconv.Itoa(chunkIndex)
	chunkBkt := client.Bucket(chunkBucket)
	resultObj := chunkBkt.Object(newObjectName)
	w := resultObj.NewWriter(ctx)
	_, err = w.Write([]byte(result))
	if err != nil {
		log.Fatal("Writing obj failed", err)
		job.UpdateWorker(fileName, myUUID, job.Failed, fbClient, ctx)
	}
	w.Close()

	err = job.UpdateWorker(fileName, myUUID, job.Completed, fbClient, ctx)
	if err != nil {
		log.Fatalf("Could not update job: %v", err)
		return err
	}

	// determine if this is the last chunk
	// if so, create pub/sub message for merging
	allDone, _ := job.AllWorkerTypeStates(fileName, job.WorkerTypeState{Type: job.Sorter, State: job.Completed}, fbClient, ctx)

	if allDone {
		// Last chunk, do something with merging perhaps
		job.SetState(fileName, job.Reducing, fbClient, ctx)

		q := &storage.Query{
			StartOffset: fileName + "/",
		}

		objects := chunkBkt.Objects(ctx, q)
		chunks := []string{}

		for {
			attrs, err := objects.Next()
			if err == iterator.Done {
				break
			}

			chunks = append(chunks, attrs.Name)

		}

		task := &pubsub.Message{
			Attributes: map[string]string{
				"jobID":        fileName,
				"resultBucket": resultBucket,
				"chunkBucket":  chunkBucket,
			},
			Data: []byte(strings.Join(chunks, ",")),
		}

		r := psClient.Topic("reduceJobs").Publish(ctx, task)
		msgId, err := r.Get(ctx)

		if err != nil {
			log.Printf("could not publish job: %v", err)
		}

		log.Printf("Published message: %v", msgId)

	}

	return nil
}
