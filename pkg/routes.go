package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) initRoutes() {
	s.g.Use(gin.Recovery())
	s.g.Use(gin.Logger())
	s.g.Use(cors.Default())
	s.g.GET("/healthz", s.healthz)

	g := s.g.Group("/api/v1")

	// Rings
	g.GET("/r/:ring", s.getRing)
	g.GET("/r/:ring/posts", s.getRingPosts)

	// Posts
	g.GET("/posts/:id", s.getPost)
	g.GET("/posts/:id/comments", s.getComments)

	// Users
	g.GET("/users/:username", s.getUser)
	g.GET("/users/:username/profilePicture", s.getUserProfilePicture)

	// Reddit-compatible API
	s.g.GET("/r/:ring/hot.json", s.getRcRingHot)
	s.g.GET("/r/:ring/about.json", s.getRcRingAbout)
	s.g.GET("/comments/:id", s.getRcComments)
	s.g.GET("/subreddits/search.json", s.getRcRingsSearch)

	// Users
	g.GET("/u/:username", s.getUser)

}
