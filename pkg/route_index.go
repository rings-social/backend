package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) indexRoute(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the rings-social backend.\n\n"+
		"This backend is Reddit API compatible.\n\n"+
		"Learn more at https://github.com/rings-social/backend.\n"+
		"Connect your app (Sync, Apollo, etc.) to this endpoint and enjoy the rings-social experience.\n"+
		"Alternatively, visit https://rings.social to use the web client.")
}
