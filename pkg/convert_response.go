package server

import (
	"backend/pkg/models"
	"backend/pkg/response"
)

func convertResponsePosts(posts []models.Post, r *models.Ring) []response.Post {
	var responsePosts []response.Post
	for _, p := range posts {
		responsePosts = append(responsePosts, response.Post{
			ID:             int(p.ID),
			PostedOn:       p.CreatedAt,
			RingName:       r.Name,
			RingColor:      r.PrimaryColor,
			AuthorUsername: p.AuthorUsername,
			Title:          p.Title,
			Body:           p.Body,
			Link:           p.Link,
			Domain:         p.Domain,
			Score:          p.Score,
			CommentsCount:  p.CommentsCount,
			Ups:            p.Ups,
			Downs:          p.Downs,
			Nsfw:           p.Nsfw,
		})
	}
	return responsePosts
}
