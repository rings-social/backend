package server

import (
	"backend/pkg/models"
	"fmt"
	"testing"
)

func TestRepoCommentsTree(t *testing.T) {
	s := testGetNewServer(t)
	comments, err := s.repoGetTopComments(1)
	if err != nil {
		t.Fatalf("unable to get comments: %v", err)
	}

	for _, c := range comments {
		fmt.Printf("%s: %s\n", c.AuthorUsername, c.Body)
	}
}
func TestRetrieveComments(t *testing.T) {
	s := testGetNewServer(t)
	comments, done := s.retrieveComments(nil, 1, nil, map[uint]models.CommentAction{})
	if done {
		t.Fatalf("done should be false")
	}

	for _, c := range comments {
		fmt.Printf("%s: %s\n", c.AuthorUsername, c.Body)
	}
}
