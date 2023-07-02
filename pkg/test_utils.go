package server

import (
	"os"
	"testing"
)

func testGetNewServer(t *testing.T) *Server {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "host=localhost user=ring password=ring dbname=ring port=5432"
	}
	s, err := New(dbUrl, nil, "http://localhost:8081")
	if err != nil {
		t.Fatalf("unable to create server: %v", err)
	}

	return s
}
