package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	chunkIndex := 0
	chunkSize := 50
	marginSize := 10
	overRead := 0
	f, err := os.Open("text.txt")

	b1 := make([]byte, chunkSize)
	m1 := make([]byte, marginSize)
	_, err = f.ReadAt(b1, int64(chunkIndex*chunkSize))
	_, err = f.ReadAt(m1, int64((chunkIndex+1)*chunkSize))
	if err != nil {
		log.Fatal(err)
	}
	s1 := string(b1)
	margin := string(m1)
	log.Println("loaded string: ", s1)
	log.Println("margin: ", s1[chunkSize:])
	firstNL := strings.Index(s1, "\n")
	if firstNL == -1 {
		return
	}
	lastNL := strings.Index(margin, "\n")
	s1 += margin

	for lastNL == -1 {
		log.Println("lastNL not found")
		overRead++
		_, err := f.ReadAt(m1, int64((chunkIndex+1)*chunkSize+marginSize*overRead))
		if err != nil {
			log.Fatal(err)
		}
		margin = string(m1)
		lastNL = strings.Index(margin, "\n")
		s1 += margin
	}
	cut_str := s1[firstNL : (chunkIndex+1)*chunkSize+marginSize*overRead+lastNL]
	log.Println("slice: ", cut_str)
	split_str := strings.Fields(cut_str)
	sort.Strings(split_str)
	result := strings.Join(split_str, " ")
	fmt.Println("result:", result)
}
