package main

import (
	"io/ioutil"
	"log"
	"sort"
	"strings"
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

func main() {

	dir := "./textfiles/"
	files, err := ioutil.ReadDir(dir)

	if len(files) == 1 {
		return
	}

	check(err, "error reading files")

	chunks := make([][]string, 0)

	for _, file := range files {
		filePath := dir + file.Name()
		data, _ := ioutil.ReadFile(filePath)
		str_arr := strings.Split(string(data), "\n")
		chunks = append(chunks, str_arr)
	}

	for len(chunks) > 1 {
		chunks[len(chunks)-2] = Merge(chunks[len(chunks)-1], chunks[len(chunks)-2])
		chunks = chunks[:len(chunks)-1]
	}

	print(strings.Join(chunks[0], "\n"))
}
