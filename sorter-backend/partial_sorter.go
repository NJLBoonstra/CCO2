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
	chunkIndex, err := strconv.Atoi(m.Attributes["chunkIdx"])
	margin, err := strconv.Atoi(m.Attributes["margin"])
	if err != nil {
		// TODO: Handle error.
	}

	// read from cloud storage
	partialString := make([]byte, chunkSize+margin)
	extPartialString := make([]byte, margin)
	overRead := 0
	client, err := storage.NewClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}

	bkt := client.Bucket("test")
	obj := bkt.Object(m.Attributes["object"])
	r, err := obj.NewReader(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	n, err := r.ReadAt(partialString, chunkSize*chunkIndex)

	// determine first and last newline of chunk
	str := string(partialString)
	firstNL := strings.Index(str, "\n")
	lastNL := strings.Index(str[chunkSize:], "\n")
	for lastNL == -1 {
		overRead++
		n, err := r.ReadAt(extPartialString, chunkSize*chunkIndex+margin)
		str += string(extPartialString)
		lastNL = strings.Index(str[chunkSize+margin*overRead:], "\n")
	}
	cut_str := str[firstNL:lastNL]

	// split
	split_str := strings.Fields(cut_str)

	// sort
	sort.Strings(split_str)

	// store sorting result
	newObjectName := m.Attributes["object"] + "-" + m.Attributes["index"]
	resultObj := bkt.Object(newObjectName)
	w := resultObj.NewWriter(ctx)
	io.WriteString(w, split_str)

	// determine if this is the last chunk
	// if so, create pub/sub message for merging

	// TODO

	return nil
}
