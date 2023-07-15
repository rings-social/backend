package server

import (
	"backend/pkg/models"
	"backend/pkg/response"
)

func convertResponsePosts(posts []models.Post, r *models.Ring) []response.Post {
	var responsePosts []response.Post
	for _, p := range posts {
		responsePost := response.Post{
			ID:             int(p.ID),
			CreatedAt:      p.CreatedAt,
			RingName:       p.RingName,
			AuthorUsername: p.AuthorUsername,
			Author:         p.Author,
			Title:          p.Title,
			Body:           p.Body,
			Link:           p.Link,
			Domain:         p.Domain,
			Score:          p.Score,
			CommentsCount:  p.CommentsCount,
			Ups:            p.Ups,
			Downs:          p.Downs,
			Nsfw:           p.Nsfw,
			VotedUp:        p.VotedUp,
			VotedDown:      p.VotedDown,
		}
		if p.Ring != nil {
			responsePost.RingColor = p.Ring.PrimaryColor
		} else {
			responsePost.RingColor = r.PrimaryColor
		}
		responsePosts = append(responsePosts, responsePost)
	}
	return responsePosts
}
