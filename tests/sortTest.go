package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func sort_lines(s string) string {
	split_str := strings.Split(s, "\n")
	sort.Strings(split_str)
	return strings.Join(split_str, "\n")
}

func main() {
	for i := 0; i < 5; i++ {
		partial_sort(i)
	}
}

func partial_sort(index int) {
	chunkIndex := index
	chunkSize := 4194304
	marginSize := 64
	overRead := 0
	f, err := os.Open("large_text.txt")

	fi, _ := f.Stat()
	fileSize := fi.Size()

	EOF := false

	chunkStart := chunkSize * chunkIndex
	chunkEnd := chunkSize * (chunkIndex + 1)

	if int64((chunkIndex+1)*chunkSize+marginSize) >= fileSize {
		chunkSize = int(fileSize) - chunkIndex*chunkSize
		marginSize = 0
		EOF = true
	}

	chunk_bytes := make([]byte, chunkSize+marginSize)

	println("chunkbytes: ", len(chunk_bytes))

	_, err = f.ReadAt(chunk_bytes, int64(chunkStart))
	if err != nil {
		log.Fatal("read failed ", err)
	}
	chunk_string := string(chunk_bytes)
	margin_string := string(chunk_bytes[chunkSize : chunkSize+marginSize])

	firstNL := 0
	if chunkIndex != 0 {
		firstNL = strings.Index(chunk_string, "\n")
	}
	if firstNL == -1 {
		return
	}

	lastNL := len(chunk_string)

	if !EOF {
		lastNL = strings.Index(margin_string, "\n")
		for lastNL == -1 {
			log.Println("margin needs to be extended")
			overRead++
			offset := int64(chunkEnd + marginSize*overRead)

			if offset+int64(marginSize) > fileSize {
				EOF = true
				log.Println("EOF reached")
				marginSize = int(fileSize) - int(offset)
				margin_bytes := make([]byte, marginSize)
				_, err := f.ReadAt(margin_bytes, offset)
				if err != nil {
					log.Fatal(err)
				}
				margin_string = string(margin_bytes)
				chunk_string += margin_string
				break
			}
			margin_bytes := make([]byte, marginSize)
			_, err := f.ReadAt(margin_bytes, offset)
			if err != nil {
				log.Fatal(err)
			}
			margin_string = string(margin_bytes)
			lastNL = strings.Index(margin_string, "\n")
			chunk_string += margin_string
		}
		lastNL += chunkSize + marginSize*overRead
	}

	if EOF {
		lastNL = len(chunk_string)
	}

	println("lastNl: ", lastNL)
	cut_str := chunk_string[firstNL:lastNL]
	r1 := []byte(sort_lines(cut_str))
	err = os.WriteFile("results/result"+fmt.Sprint(chunkIndex)+".txt", r1, 0644)
	f.Close()
}
