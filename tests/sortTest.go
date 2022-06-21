package main

import (
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
	chunkIndex := 0
	chunkSize := 152083
	marginSize := 1
	overRead := 0
	f, err := os.Open("alice29.txt")

	fi, _ := f.Stat()
	fileSize := fi.Size()

	EOF := false

	if int64((chunkIndex+1)*chunkSize+marginSize) >= fileSize {
		log.Println("chunk larger than file.")
		chunkSize = int(fileSize) - chunkIndex*chunkSize
		log.Println("new chunksize: ", chunkSize)
		EOF = true
	}

	chunk_bytes := make([]byte, chunkSize)

	_, err = f.ReadAt(chunk_bytes, int64(chunkIndex*chunkSize))
	if err != nil {
		log.Fatal("read failed ", err)
	}
	chunk_string := string(chunk_bytes)

	firstNL := 0
	if chunkIndex != 0 {
		firstNL = strings.Index(chunk_string, "\n")
	}
	if firstNL == -1 {
		return
	}

	lastNL := len(chunk_string)

	if !EOF {
		margin_bytes := make([]byte, marginSize)
		_, err := f.ReadAt(margin_bytes, int64((chunkIndex+1)*chunkSize))
		if err != nil {
			log.Fatal(err)
		}
		margin_string := string(margin_bytes)
		chunk_string += margin_string
		lastNL = strings.Index(margin_string, "\n")
		for lastNL == -1 {
			log.Print("enter loop", lastNL)
			overRead++
			offset := int64((chunkIndex+1)*chunkSize + marginSize*overRead)
			log.Println("margin needs to be extended")
			if offset+int64(marginSize) > fileSize {
				log.Println("EOF reached")
				marginSize = int(fileSize) - int(offset)
				margin_bytes = make([]byte, marginSize)
				EOF = true
				_, err := f.ReadAt(margin_bytes, offset)
				if err != nil {
					log.Fatal(err)
				}
				margin_string = string(margin_bytes)
				lastNL = int(fileSize)
				chunk_string += margin_string
				break
			}
			_, err := f.ReadAt(margin_bytes, offset)
			if err != nil {
				log.Fatal(err)
			}
			margin_string = string(margin_bytes)
			log.Println("added: ", margin_string)
			lastNL = strings.Index(margin_string, "\n")
			log.Println("index:", lastNL)
			chunk_string += margin_string
		}
		lastNL += (chunkIndex+1)*chunkSize + marginSize*(overRead+1)
	}

	if EOF {
		lastNL = int(fileSize)
	}

	cut_str := chunk_string[firstNL:lastNL]
	r1 := []byte(sort_lines(cut_str))
	err = os.WriteFile("result.txt", r1, 0644)
	f.Close()
}
