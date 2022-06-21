package job

import (
	"context"

	"cloud.google.com/go/firestore"
)

const CollectionPalindrome string = "palindromeResults"

type PalindromeResult struct {
	JobID             string `json:"jobId"`
	Palindromes       int    `json:"palindromes"`
	LongestPalindrome int    `json:"longestPalindrome"`
}

func AddPalindromeResult(jobID string, palindromes int, longest int, fbClient *firestore.Client, ctx context.Context) error {
	res := PalindromeResult{
		JobID:             jobID,
		Palindromes:       palindromes,
		LongestPalindrome: longest,
	}

	_, err := fbClient.Collection(CollectionPalindrome).Doc(jobID).Set(ctx, &res)

	return err
}

func GetPalindromeResult(jobID string, fbClient *firestore.Client, ctx context.Context) (PalindromeResult, error) {
	res := PalindromeResult{}
	data, err := fbClient.Collection(CollectionPalindrome).Doc(jobID).Get(ctx)
	if err != nil {
		return res, err
	}

	err = data.DataTo(&res)

	return res, err
}
