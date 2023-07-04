package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) initRoutes() {
	s.g.Use(gin.Recovery())
	s.g.Use(gin.Logger())
	s.g.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "Authorization",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
	}))
	s.g.GET("/healthz", s.healthz)
	s.g.GET("/", s.indexRoute)

	g := s.g.Group("/api/v1")

	g.Use(s.authMiddleware())

	// Rings
	g.GET("/rings", s.routeGetRings)
	g.GET("/r/:ring", s.routeGetRing)
	g.GET("/r/:ring/posts", s.maybeAuthenticatedUser, s.getRingPosts)
	g.POST("/r/:ring", s.authenticatedUser, s.routeCreateRing)

	// Posts
	g.POST("/posts", s.authenticatedUser, s.createPost)
	g.GET("/posts/:id", s.maybeAuthenticatedUser, s.getPost)
	g.GET("/posts/:id/comments", s.getComments)
	g.POST("/posts/:id/comments", s.postComment)
	g.DELETE("/posts/:id/comments/:commentId", s.deleteComment)

	g.PUT("/posts/:id/upvote", s.authenticatedUser, s.routeUpvotePost)
	g.PUT("/posts/:id/downvote", s.authenticatedUser, s.routeDownvotePost)

	// Comments
	g.GET("/comments", s.routeGetRecentComments)
	g.PUT("/posts/:id/comments/:commentId/upvote", s.upvoteComment)
	g.PUT("/posts/:id/comments/:commentId/downvote", s.downvoteComment)

	// Users
	g.GET("/users/me", s.getMe)
	g.GET("/users", s.getUsers)
	g.GET("/users/:username", s.getUser)
	g.GET("/users/:username/profilePicture", s.getUserProfilePicture)

	// SignUp
	g.GET("/signup/usernameAvailability", s.usernameAvailability)
	g.POST("/signup/username", s.signupUsername)

	// Reddit-compatible API
	s.g.GET("/r/:ring/hot.json", s.getRcRingHot)
	s.g.GET("/r/:ring/about.json", s.getRcRingAbout)
	s.g.GET("/comments/:id", s.getRcComments)
	s.g.GET("/subreddits/search.json", s.getRcRingsSearch)

	// Users
	g.GET("/u/:username", s.getUser)

}
