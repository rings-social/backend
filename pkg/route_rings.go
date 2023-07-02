package server

import (
	"backend/pkg/models"
	"backend/pkg/request"
	"backend/pkg/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) routeGetRing(context *gin.Context) {
	ringName := context.Param("ring")
	if ringName == "" {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "Ring name is required",
		})
		return
	}

	var ring models.Ring
	tx := s.db.First(&ring, "name = ?", ringName)
	if tx.Error != nil {
		// If tx reports not found:
		if tx.Error == gorm.ErrRecordNotFound {
			context.AbortWithStatusJSON(404, gin.H{
				"error": "Ring not found",
			})
			return
		}
		// Otherwise, it's an internal error:
		s.logger.Errorf("Unable to get ring %s: %v", ringName, tx.Error)
		context.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to get ring",
		})
		return
	}

	context.JSON(200, ring)
}

func (s *Server) routeGetRings(c *gin.Context) {
	// Get paginated request
	pagination, done := s.getPagination(c)
	if done {
		return
	}

	// Check for query parameter
	nsfw := false
	nsfwParam := c.Query("nsfw")
	if strings.ToLower(nsfwParam) == "true" {
		nsfw = true
	}

	total, err := s.repoGetTotalRings()
	if err != nil {
		s.logger.Errorf("Unable to get total rings: %v", err)
		internalServerError(c)
		return
	}

	q := c.Query("q")
	if q != "" {
		rings, err := s.repoRingsSearch(q, nsfw, pagination.After, pagination.Limit)
		if err != nil {
			s.logger.Errorf("Unable to search rings: %v", err)
			internalServerError(c)
			return
		}
		after := ""
		if len(rings) > 0 {
			after = fmt.Sprintf("%d", rings[len(rings)-1].Subscribers)
		}

		returnPaginated(c, after, rings, -1)
		return
	}

	offset := 0
	if pagination.After != "" {
		var err error
		offset, err = strconv.Atoi(pagination.After)
		if err != nil {
			s.logger.Errorf("Unable to parse after %s: %v", pagination.After, err)
			internalServerError(c)
			return
		}
	}
	rings, err := s.repoGetRings(offset, pagination.Limit)
	if err != nil {
		s.logger.Errorf("Unable to get rings: %v", err)
		internalServerError(c)
		return
	}

	after := ""
	afterV := offset + int(pagination.Limit)

	if len(rings) > 0 {
		after = fmt.Sprintf("%d", afterV)
	}

	// Paginated result
	returnPaginated(c, after, rings, total)
}

func returnPaginated[T any](c *gin.Context, after string, items []T, total int64) {
	paginatedResult := response.Paginated[T]{
		After: after,
		Items: items,
		Total: total,
	}
	c.JSON(200, paginatedResult)
}

func (s *Server) routeCreateRing(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		return
	}

	ringName := c.Param("ring")
	isValidRingName := validateRingName(ringName)
	if !isValidRingName {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "Invalid ring name"},
		)
		return
	}

	var ringRequest request.CreateRingRequest
	err := c.BindJSON(&ringRequest)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err = ringRequest.Validate()

	ring := models.Ring{
		Name:          ringName,
		Title:         ringRequest.Title,
		Description:   ringRequest.Description,
		OwnerUsername: username.(string),
		PrimaryColor:  ringRequest.Color,
	}

	tx := s.db.Create(&ring)
	if tx.Error != nil {
		s.logger.Errorf("Unable to create ring: %v", tx.Error)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to create ring",
		})
		return
	}

	c.JSON(200, ring)
}
