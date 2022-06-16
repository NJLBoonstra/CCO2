package sorter_backend

import (
	"context"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
)

func partialSort(ctx context.Context, m PubSubMessage) error {
	chunkSizeStr, exists := os.LookupEnv("CHUNK_SIZE")
	if !exists {
		log.Fatal("Make sure CHUNK_SIZE is set!")
	}
	chunkSize, err := strconv.Atoi(chunkSizeStr)
	if err != nil {
		log.Fatalf("Could not convert CHUNK_SIZE to int: %v", err)
	}
	// read pubsub
	chunkIndex, _ := strconv.Atoi(m.Attributes["chunkIdx"])
	margin, _ := strconv.Atoi(m.Attributes["margin"])
	// if err != nil {
	// 	// TODO: Handle error.
	// }

	// read from cloud storage
	partialString := make([]byte, chunkSize+margin)
	extPartialString := make([]byte, margin)
	overRead := 0
	client, _ := storage.NewClient(ctx)
	// if err != nil {
	// 	// TODO: Handle error.
	// }

	bkt := client.Bucket("test")
	obj := bkt.Object(m.Attributes["object"])
	r, _ := obj.NewRangeReader(ctx, int64(chunkSize)*int64(chunkIndex), int64(chunkSize)+int64(margin))
	// if err != nil {
	// 	// TODO: Handle error.
	// }
	_, _ = r.Read(partialString)

	// determine first and last newline of chunk
	str := string(partialString)
	firstNL := strings.Index(str, "\n")
	lastNL := strings.Index(str[chunkSize:], "\n")
	for lastNL == -1 {
		overRead++
		r, _ = obj.NewRangeReader(ctx, int64(chunkSize)*int64(chunkIndex)+int64(margin), int64(chunkSize))
		_, _ = r.Read(extPartialString)
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
	newObjectName := m.Attributes["object"] + "-" + m.Attributes["index"]
	resultObj := bkt.Object(newObjectName)
	w := resultObj.NewWriter(ctx)
	io.WriteString(w, result)

	// determine if this is the last chunk
	// if so, create pub/sub message for merging

	// TODO

	return nil
}
