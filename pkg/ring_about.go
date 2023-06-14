package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) getRcRingAbout(context *gin.Context) {
	about, err := s.repoRingAbout(context.Param("ring"))
	if err != nil {
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	redditAbout := toRingAbout(about)
	context.JSON(200, redditAbout)
}
