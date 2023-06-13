package server

import (
	"backend/pkg/models"
	"backend/pkg/reddit_compat"
	"fmt"
	"strings"
)

// toRedditPosts converts a slice of models.Post to a RedditPosts strucrandom_dudet
func toRedditPosts(posts []models.Post) (RedditPosts, error) {
	var listing RedditPosts
	for _, post := range posts {
		p := reddit_compat.KindData[reddit_compat.Post]{
			Kind: "t3",
			Data: reddit_compat.Post{
				ID:           fmt.Sprintf("%d", post.ID),
				Title:        post.Title,
				Selftext:     post.Body,
				SelftextHtml: &post.Body,
				Subreddit:    post.RingName,
				Author:       post.AuthorUsername,
				Permalink:    fmt.Sprintf("/r/%s/comments/%d/%s", post.RingName, post.ID, seoTitle(post.Title)),
				Ups:          post.Ups,
				Downs:        post.Downs,
				Score:        post.Score,
				NumComments:  post.CommentsCount,
				URL:          post.Link,
				Created:      int(post.PostedOn.Unix()),
				CreatedUtc:   int(post.PostedOn.UTC().Unix()),
				Over18:       post.Nsfw,
			},
		}

		if p.Data.Domain == nil {
			myUrl := "http://192.168.1.134:8081" + p.Data.Permalink
			selfRingName := "self." + post.RingName
			p.Data.Domain = &selfRingName
			p.Data.Thumbnail = "self"
			p.Data.URL = &myUrl
			p.Data.AuthorFlairType = "text"
			p.Data.LinkFlairType = "text"
		}

		p.Data = *parseNilAsEmpty(&p.Data)
		listing.Data.Children = append(listing.Data.Children, p)
	}

	listing.Data.After = "t3_" + listing.Data.Children[len(listing.Data.Children)-1].Data.ID

	return listing, nil
}

func toRedditSubreddits(rings []models.Ring) (RedditSubreddits, error) {
	var listing RedditSubreddits
	listing.Kind = "Listing"
	for _, ring := range rings {
		red := "#FF0000"
		iconImg := "https://a.thumbs.redditmedia.com/E0Bkwgwe5TkVLflBA7WMe9fMSC7DV2UOeff-UpNJeb0.png"
		s := reddit_compat.KindData[reddit_compat.Subreddit]{
			Kind: "t5",
			Data: reddit_compat.Subreddit{
				ID:                    ring.Name,
				Title:                 ring.Title,
				Name:                  ring.Name,
				DisplayName:           ring.Name,
				Description:           &ring.Description,
				Over18:                &ring.Nsfw,
				URL:                   "/r/" + ring.Name,
				DisplayNamePrefixed:   "r/" + ring.Name,
				BannerBackgroundColor: &red,
				IconImg:               &iconImg,
				Subscribers:           &ring.Subscribers,
			},
		}
		listing.Data.Children = append(listing.Data.Children, s)
	}

	return listing, nil
}

// seoTitle converts the Title into a URL friendly string for the permalink
func seoTitle(title string) string {
	return strings.Replace(strings.ToLower(title), " ", "-", -1)
}
