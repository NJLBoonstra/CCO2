package sorter_backend

import (
	"context"
	"io/ioutil"
	"strings"

	"cloud.google.com/go/storage"

	job "cco.bn.edu/shared"
)

func Merge(s1, s2 []string) []string {

	size, i, j := len(s1)+len(s2), 0, 0
	slice := make([]string, size, size)

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
	fileNamesString := string(m.Data)
	fileNames := strings.Split(fileNamesString, ",")

	if len(fileNames) == 1 {
		return nil
	}

	client, err := storage.NewClient(ctx)
	check(err, "Client could not be created")

	chunks := make([][]string, 0)
	bkt := client.Bucket(chunkBucketName)

	for _, fileName := range fileNames {
		obj := bkt.Object(fileName)
		reader, err := obj.NewReader(ctx)
		check(err, "Reader creation failed")
		slurp, err := ioutil.ReadAll(reader)
		check(err, "Reading file failed")
		defer reader.Close()
		str_arr := strings.Split(string(slurp), "\n")
		chunks = append(chunks, str_arr)
	}

	for len(chunks) > 1 {
		chunks[len(chunks)-2] = Merge(chunks[len(chunks)-1], chunks[len(chunks)-2])
		chunks = chunks[:len(chunks)-1]
	}

	result := strings.Join(chunks[0], "\n")

	newObjectName := origFileName + "-sorted"
	resultObj := client.Bucket(resultBucketName).Object(newObjectName)
	w := resultObj.NewWriter(ctx)
	_, err = w.Write([]byte(result))
	defer w.Close()
}
