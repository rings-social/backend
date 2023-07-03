package server

import (
	"fmt"
	"testing"
)

func TestRepoRings(t *testing.T) {
	s := testGetNewServer(t)
	rings, err := s.repoGetRings(5, 100)
	if err != nil {
		t.Fatalf("unable to get rings: %v", err)
	}

	for _, r := range rings {
		fmt.Printf("%s\n", r.Name)
	}
}
