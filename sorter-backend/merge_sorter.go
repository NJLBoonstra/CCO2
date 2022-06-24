package sorter_backend

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"

	job "cco.bn.edu/shared"
)

func Merge(s1, s2 []string) []string {

	size, i, j := len(s1)+len(s2), 0, 0
	slice := make([]string, size)

	for k := 0; k < size; k++ {
		if i > len(s1)-1 && j <= len(s2)-1 {
			slice[k] = s2[j]
			j++
		} else if j > len(s2)-1 && i <= len(s1)-1 {
			slice[k] = s1[i]
			i++
		} else if strings.Compare(s1[i], s2[j]) == -1 {
			slice[k] = s1[i]
			i++
		} else {
			slice[k] = s2[j]
			j++
		}
	}
	return slice
}

func MergeSort(ctx context.Context, m job.MergePubSub) error {
	chunkBucketName := m.Attributes["chunkBucket"]
	resultBucketName := m.Attributes["resultBucket"]
	origFileName := m.Attributes["jobID"]
	log.Print("jobID:", origFileName, chunkBucketName, resultBucket)
	fileNamesString := string(m.Data)
	log.Print("files:", fileNamesString)
	fileNames := strings.Split(fileNamesString, ",")

	if len(fileNames) < 1 {
		return errors.New("no files to merge?")
	}

	fbClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalf("Could not create a Firestore client: %v", err)
		return err
	}
	defer fbClient.Close()

	myUUID, err := job.AddWorker(origFileName, job.SorterReduce, fbClient, ctx)
	if err != nil {
		log.Printf("could not add worker %v", err)
		return err
	}

	client, err := storage.NewClient(ctx)
	check(err, "Client could not be created")
	bkt := client.Bucket(chunkBucketName)
	newObjectName := origFileName + "-sorted"
	resultObj := client.Bucket(resultBucketName).Object(newObjectName)

	if len(fileNames) == 1 {
		obj := bkt.Object(fileNames[0])

		_, err = resultObj.CopierFrom(obj).Run(ctx)
		if err != nil {
			log.Print("could not copy single chunk to result:", err)
			return err
		}
	}

	chunks := make([][]string, 0)

	for _, fileName := range fileNames {
		obj := bkt.Object(fileName)
		reader, err := obj.NewReader(ctx)
		check(err, "Reader creation failed")
		slurp, err := ioutil.ReadAll(reader)
		check(err, "Reading file failed")
		defer reader.Close()
		str_arr := strings.Split(string(slurp), "\n")
		chunks = append(chunks, str_arr)

		err = obj.Delete(ctx)
		if err != nil {
			log.Print("could not delete chunk:", err)
		}
	}

	for len(chunks) > 1 {
		chunks[len(chunks)-2] = Merge(chunks[len(chunks)-1], chunks[len(chunks)-2])
		chunks = chunks[:len(chunks)-1]
	}

	result := strings.Join(chunks[0], "\n")

	w := resultObj.NewWriter(ctx)
	_, _ = w.Write([]byte(result))
	defer w.Close()

	err = job.UpdateWorker(origFileName, myUUID, job.Completed, fbClient, ctx)
	if err != nil {
		log.Printf("Could not update job: %v", err)
		return err
	}

	err = job.SetState(origFileName, job.Completed, fbClient, ctx)
	if err != nil {
		log.Printf("could not set job state: %v", err)
		return err
	}

	return nil
}
