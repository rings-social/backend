package reddit_compat_test

import (
	"backend/pkg/reddit_compat"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestAllHot(t *testing.T) {
	f, err := os.Open("../../resources/reddit/r/popular/hot.json")
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()

	var res reddit_compat.KindData[reddit_compat.Listing[reddit_compat.Post]]

	err = dec.Decode(&res)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range res.Data.Children {
		fmt.Printf("%s [r/%s]\n\t%s\n", v.Data.Title, v.Data.Subreddit, "https://reddit.com"+v.Data.Permalink)
	}
}
