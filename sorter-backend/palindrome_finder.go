package sorter_backend

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

func FindPalindromes(ctx context.Context, e GCSEvent) {
	filename := e.Name
	bucket := e.Bucket

	client, err := storage.NewClient(ctx)

	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	bkt := client.Bucket(bucket)
	obj := bkt.Object(filename)
	attrs, err := obj.Attrs(ctx)

	if err != nil {
		log.Printf("Could not read object attributes: %v", err)
	}
	obj_size := attrs.Size

	reader, err := obj.NewReader(ctx)

	if err != nil {
		log.Printf("Could not open reader: %v", err)
		return
	}

	buffer := make([]byte, obj_size)

	bytes_read, err := reader.Read(buffer)

	if err != nil {
		log.Printf("Could not read file: %v", err)
	}

	palindromes := 0
	longest_pal := 0

	word_start := -1
	var word string
	// Iterate over values in buffer to construct words
	// Ik heb gekozen om de []byte niet om te zetten naar string, want dat
	// leek me een onnodige extra stap
	for i, v := range buffer {
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			// Word has started, or inside a word
			if word_start < 0 {
				word_start = i
			}
		} else {
			if word_start > 0 {
				// Word *may* have ended
				word = string(buffer[word_start : i+1])
			}
		}

		if len(word) > 0 {
			pal := CheckPalindrome(word)

			if pal {
				palindromes++

				if len(word) > longest_pal {
					longest_pal = len(word)
				}
			}
			// Reset word
			word = ""
		}
	}

	log.Printf("Palindromes: %d; Longest: %d", palindromes, longest_pal)
}

func CheckPalindrome(word string) bool {
	log.Printf("CheckPalindrome: %s", word)

	for i := 0; i < (len(word)/2)+1; i++ {
		j := len(word) - i - 1
		if word[i] != word[j] {
			return false
		}
	}

	return true
}
