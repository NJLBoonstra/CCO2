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
	"github.com/google/uuid"
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

func PartialSort(ctx context.Context, m job.PubSubMessage) error {
	// read all information from pubsub message
	marginSize, err := strconv.Atoi(m.Attributes["marginSize"])
	check(err, "Could not convert marginSize to int")
	bucketName := m.Attributes["bucket"]
	fileName := m.Attributes["jobID"]
	chunkBucket := m.Attributes["chunkBucket"]
	resultBucket := m.Attributes["resultBucket"]
	myUUID, err := uuid.Parse(m.Attributes["workerID"])
	check(err, "cannot parse uuid")
	chunkSize, err := strconv.Atoi(m.Attributes["chunkSize"])
	check(err, "Could not convert CHUNK_SIZE to int")
	chunkIndex, err := strconv.Atoi(m.Attributes["chunkIdx"])
	check(err, "Could not convert chunk index to int")
	objectSize, err := strconv.Atoi(m.Attributes["objectSize"])
	check(err, "Could not convert objectSize to int")

	// create GCS client
	client, err := storage.NewClient(ctx)
	check(err, "Client could not be created")
	defer client.Close()

	// create firestore client
	fbClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Could not create a Firestore client: %v", err)
		return err
	}
	defer fbClient.Close()

	// create handle for relevant bucket and object
	bkt := client.Bucket(bucketName)
	obj := bkt.Object(fileName)

	// helper-variable to store if the end of the object is reached
	EOF := false

	chunkStart := chunkSize * chunkIndex
	chunkEnd := chunkSize * (chunkIndex + 1)

	// check if this reader will read over the file size -> adjust read size
	if (chunkIndex+1)*chunkSize+marginSize >= objectSize {
		chunkSize = objectSize - chunkIndex*chunkSize
		marginSize = 0
		EOF = true
	}

	// read the chunk bytes and convert to string
	chunk_reader, err := obj.NewRangeReader(ctx, int64(chunkStart), int64(chunkSize+marginSize))
	check(err, "Reader creation failed for obj")
	slurp, err := ioutil.ReadAll(chunk_reader)
	check(err, "Reading obj failed")
	defer chunk_reader.Close()
	chunk_string := string(slurp)
	margin_string := string(slurp[chunkSize : chunkSize+marginSize])

	// if the first chunk is sorted start at first charachter, else start at first found NL
	firstNL := 0
	if chunkIndex != 0 {
		firstNL = strings.Index(chunk_string[:chunkSize], "\n")
	}

	// this chunk contains no NL -> another sorter will process this chunk by extending its margin
	if firstNL == -1 {
		err = job.UpdateWorker(fileName, myUUID, job.Completed, fbClient, ctx)
		check(err, "Could not update job")
	}

	// initialize the last NL variable
	lastNL := len(chunk_string)

	// if we have not reached the end of the file, search for the first NL in the margin
	if !EOF {
		// helper-variable to determine the number of times we had to extend the margin to find a NL
		overRead := 0

		// find first NL in margin
		lastNL = strings.Index(margin_string, "\n")

		// if no NL found in margin, extend the margin
		for lastNL == -1 {
			overRead++
			offset := int64(chunkEnd + marginSize*overRead)

			// if the extended margin surpasses the object size, read only remaining bytes
			if offset+int64(marginSize) > int64(objectSize) {
				EOF = true
				marginSize = int(objectSize) - int(offset)
				margin_reader, err := obj.NewRangeReader(ctx, offset, int64(marginSize))
				check(err, "Could not create a NewRangeReader")
				margin_bytes, err := ioutil.ReadAll(margin_reader)
				check(err, "Reading obj in iteration failed")
				defer margin_reader.Close()
				margin_string = string(margin_bytes)
				chunk_string += margin_string
				break
			}
			margin_reader, err := obj.NewRangeReader(ctx, int64(chunkEnd+marginSize*overRead), int64(marginSize))
			check(err, "Could not create a NewRangeReader")
			margin_bytes, err := ioutil.ReadAll(margin_reader)
			check(err, "Reading obj in iteration failed")
			defer margin_reader.Close()
			margin_string = string(margin_bytes)
			lastNL = strings.Index(margin_string, "\n")
			chunk_string += margin_string
		}
		lastNL += chunkSize + marginSize*overRead
	}

	// if the EOF was reached, set last NL to last character
	if EOF {
		lastNL = len(chunk_string)
	}

	// determine final chunk based on NL
	cut_str := chunk_string[firstNL:lastNL]
	result := sort_lines(cut_str)

	// store sorting result in folder
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

	// Add palindrome worker to firestore
	palindromeWorkerUUID, err := job.AddWorker(fileName, job.Palindrome, fbClient, ctx)
	check(err, "could not add worker")

	resultAttrs, err := resultObj.Attrs(ctx)
	check(err, "cannot fetch resultObj attributes")
	resultObj = resultObj.If(storage.Conditions{MetagenerationMatch: resultAttrs.Metageneration})
	resultAttrsUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"palindromeWorkerID": palindromeWorkerUUID.String(),
		},
	}

	// update status of current job in firestore
	_, err = resultObj.Update(ctx, resultAttrsUpdate)
	check(err, "could not update chunk object attributes!")

	// update status of worker in firestore
	err = job.UpdateWorker(fileName, myUUID, job.Completed, fbClient, ctx)
	check(err, "Could not update job")

	// determine if this is the last chunk
	// if so, create pub/sub message for merging
	allDone, _ := job.AllWorkerTypeStates(fileName, job.WorkerTypeState{Type: job.Sorter, State: job.Completed}, fbClient, ctx)

	if allDone {
		job.SetState(fileName, job.Reducing, fbClient, ctx)

		q := &storage.Query{
			Prefix: fileName + "/",
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
		log.Printf("Published files: %v for sorting", strings.Join(chunks, ","))

		r := psClient.Topic("reduceJobs").Publish(ctx, task)
		msgId, err := r.Get(ctx)

		if err != nil {
			log.Printf("could not publish job: %v", err)
		}

		log.Printf("Published message: %v", msgId)

	}

	return nil
}
