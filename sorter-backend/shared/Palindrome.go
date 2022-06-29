package job

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const CollectionPalindrome string = "palindromeResults"

type PalindromeResult struct {
	Palindromes       int `json:"palindromes"`
	LongestPalindrome int `json:"longestPalindrome"`
}

type PalindromeJob struct {
	ID                     string                      `json:"id"`
	PalindromeWorkerResult map[string]PalindromeResult `json:"palindromeWorkerResult"`
	Palindromes            int                         `json:"palindromes"`
	LongestPalindrome      int                         `json:"longestPalindrome"`
}

func CreatePalindromeTable(jobID string, fbClient *firestore.Client, ctx context.Context) error {

	res := PalindromeJob{
		ID:                     jobID,
		PalindromeWorkerResult: map[string]PalindromeResult{},
		Palindromes:            0,
		LongestPalindrome:      0,
	}
	_, err := fbClient.Collection(CollectionPalindrome).Doc(jobID).Set(ctx, &res)

	return err
}

func UpdatePalindromeResult(jobID string, workerUUID uuid.UUID, palindromes int, longest int, fbClient *firestore.Client, ctx context.Context) error {
	docRef := fbClient.Collection(CollectionPalindrome).Doc(jobID)
	j, err := docRef.Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return errors.New("cannot update a non-existing PalindromeJob")
	}

	// TODO: dit misschien efficienter?
	job := PalindromeJob{}
	err = j.DataTo(&job)
	if err != nil {
		return err
	}

	// UPdate the document yeah
	_, err = fbClient.Collection(CollectionPalindrome).Doc(jobID).Update(ctx, []firestore.Update{
		{Path: "PalindromeWorkerResult." + workerUUID.String(), Value: PalindromeResult{Palindromes: palindromes, LongestPalindrome: longest}},
	})
	return err
}

func UpdatePalindromeJobResult(jobID string, palindromes int, longest int, fbClient *firestore.Client, ctx context.Context) error {
	docRef := fbClient.Collection(CollectionPalindrome).Doc(jobID)
	j, err := docRef.Get(ctx)
	if err != nil && status.Code(err) == codes.NotFound {
		return errors.New("cannot update a non-existing PalindromeJob")
	}

	// TODO: dit misschien efficienter?
	job := PalindromeJob{}
	err = j.DataTo(&job)
	if err != nil {
		return err
	}

	// UPdate the document yeah
	_, err = fbClient.Collection(CollectionPalindrome).Doc(jobID).Update(ctx, []firestore.Update{
		{Path: "Palindromes", Value: palindromes},
	})
	_, err = fbClient.Collection(CollectionPalindrome).Doc(jobID).Update(ctx, []firestore.Update{
		{Path: "LongestPalindrome", Value: longest},
	})
	return err
}

func GetPalindromeResult(jobID string, fbClient *firestore.Client, ctx context.Context) (PalindromeJob, error) {
	res := PalindromeJob{}
	data, err := fbClient.Collection(CollectionPalindrome).Doc(jobID).Get(ctx)
	if err != nil {
		return res, err
	}

	err = data.DataTo(&res)

	return res, err
}
