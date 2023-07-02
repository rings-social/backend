package server

import (
	"backend/pkg/models"
	"backend/pkg/reddit_compat"
	"fmt"
	"strings"
)

// toRedditPosts converts a slice of models.Post to a RedditPosts strucrandom_dudet
func toRedditPosts(posts []models.Post, baseUrl string) (RedditPosts, error) {
	var listing RedditPosts
	listing.Kind = "Listing"
	for _, post := range posts {
		p := convertToRedditPost(&post, baseUrl)
		listing.Data.Children = append(listing.Data.Children, p)
	}

	if len(listing.Data.Children) > 0 {
		last := "t3_" + listing.Data.Children[len(listing.Data.Children)-1].Data.ID
		listing.Data.After = &last
	}

	return listing, nil
}

func convertToRedditPost(post *models.Post, baseUrl string) reddit_compat.KindData[reddit_compat.Post] {
	postHint := "text"
	if post.Link != nil {
		postHint = "link"
	}
	p := reddit_compat.KindData[reddit_compat.Post]{
		Kind: "t3",
		Data: reddit_compat.Post{
			ID:                    fmt.Sprintf("%d", post.ID),
			Name:                  fmt.Sprintf("t3_%d", post.ID),
			Title:                 post.Title,
			Selftext:              post.Body,
			SelftextHtml:          &post.Body,
			Subreddit:             post.RingName,
			SubredditNamePrefixed: prefixSubreddit(post.RingName),
			Author:                post.AuthorUsername,
			Permalink:             fmt.Sprintf("/r/%s/comments/%d/%s", post.RingName, post.ID, seoTitle(post.Title)),
			Ups:                   post.Ups,
			Downs:                 post.Downs,
			Score:                 post.Score,
			NumComments:           post.CommentsCount,
			URL:                   post.Link,
			Domain:                post.Domain,
			Created:               int(post.CreatedAt.Unix()),
			CreatedUtc:            int(post.CreatedAt.UTC().Unix()),
			Over18:                post.Nsfw,
			PostHint:              postHint,
		},
	}

	if p.Data.Domain == nil {
		myUrl := baseUrl + p.Data.Permalink
		selfRingName := "self." + post.RingName
		p.Data.Domain = &selfRingName
		p.Data.Thumbnail = "self"
		p.Data.URL = &myUrl
		p.Data.AuthorFlairType = "text"
		p.Data.LinkFlairType = "text"
	}

	p.Data = *parseNilAsEmpty(&p.Data)
	return p
}

func prefixSubreddit(name string) string {
	return "r/" + name
}

func toRedditPost(post *models.Post, baseUrl string) (reddit_compat.KindData[reddit_compat.Post], error) {
	return convertToRedditPost(post, baseUrl), nil
}

func toRedditSubreddits(rings []models.Ring) (RedditSubreddits, error) {
	var listing RedditSubreddits
	listing.Kind = "Listing"
	for _, ring := range rings {
		red := "#FF0000"
		iconImg := "https://a.thumbs.redditmedia.com/E0Bkwgwe5TkVLflBA7WMe9fMSC7DV2UOeff-UpNJeb0.png"
		subscribers := int(ring.Subscribers)
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
				Subscribers:           &subscribers,
			},
		}
		listing.Data.Children = append(listing.Data.Children, s)
	}

	return listing, nil
}

func toRingAbout(ring *models.Ring) RedditAbout {
	subscribers := int(ring.Subscribers)
	return RedditAbout{
		Kind: "t5",
		Data: reddit_compat.SubredditDetails{
			ID:                  ring.Name,
			Title:               ring.Title,
			Name:                ring.Name,
			DisplayName:         ring.Name,
			Description:         ring.Description,
			Over18:              ring.Nsfw,
			URL:                 "/r/" + ring.Name,
			DisplayNamePrefixed: "r/" + ring.Name,
			Subscribers:         subscribers,
			DescriptionHtml:     ring.Description,
			Created:             int(ring.CreatedAt.Unix()),
			CreatedUtc:          int(ring.CreatedAt.UTC().Unix()),
			PrimaryColor:        ring.PrimaryColor,
			ActiveUserCount:     19,
		},
	}
}

func toRedditComments(post *models.Post, comments []models.Comment, baseUrl string) ([]any, error) {
	if post == nil {
		return nil, fmt.Errorf("post is nil")
	}
	redditPost, err := toRedditPost(post, baseUrl)
	if err != nil {
		return nil, err
	}

	listing := toRedditCommentsInner(post, comments, 0)

	if listing.Data.Children == nil {
		listing.Data.Children = []reddit_compat.KindData[reddit_compat.Comment]{}
	}
	if len(listing.Data.Children) > 0 {
		listing.Data.After = &listing.Data.Children[len(listing.Data.Children)-1].Data.ID
	}

	return []any{
		wrapListing(
			[]reddit_compat.KindData[reddit_compat.Post]{redditPost},
		),
		listing,
	}, nil
}

func toRedditCommentsInner(post *models.Post, comments []models.Comment, depth int) RedditComments {
	var listing RedditComments
	listing.Kind = "Listing"
	for _, comment := range comments {
		c := reddit_compat.KindData[reddit_compat.Comment]{
			Kind: "t1",
			Data: reddit_compat.Comment{
				ID:                    fmt.Sprintf("%d", comment.ID),
				Body:                  comment.Body,
				BodyHtml:              comment.Body,
				Author:                comment.AuthorUsername,
				Subreddit:             post.RingName,
				SubredditNamePrefixed: "r/" + post.RingName,
				Permalink:             getCommentPermalink(post.RingName, comment.PostId, comment.ID),
				LinkID:                fmt.Sprintf("t3_%d", post.ID),
				Score:                 int(comment.Ups - comment.Downs),
				Created:               int(comment.CreatedAt.Unix()),
				CreatedUtc:            int(comment.CreatedAt.UTC().Unix()),
				Replies:               toRedditCommentsInner(post, comment.Replies, depth+1),
				Depth:                 depth,
			},
		}
		c.Data = *parseNilAsEmpty(&c.Data)
		listing.Data.Children = append(listing.Data.Children, c)
	}
	if listing.Data.Children == nil {
		listing.Data.Children = []reddit_compat.KindData[reddit_compat.Comment]{}
	}
	return listing
}

func getCommentPermalink(name string, postId uint, commentId uint) string {
	return fmt.Sprintf("/r/%s/comments/%d/%d", name, postId, commentId)
}

func wrapListing[T any](posts []reddit_compat.KindData[T]) reddit_compat.KindData[reddit_compat.Listing[T]] {
	return reddit_compat.KindData[reddit_compat.Listing[T]]{
		Kind: "Listing",
		Data: reddit_compat.Listing[T]{
			Children: posts,
		},
	}
}

// seoTitle converts the Title into a URL friendly string for the permalink
func seoTitle(title string) string {
	return strings.Replace(strings.ToLower(title), " ", "-", -1)
}
