package server

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Pagination struct {
	After string
	Limit uint64
}

func (s *Server) getPagination(c *gin.Context) (*Pagination, bool) {
	// Pagination is done using the following query parameters:
	// - after: the cursor to start fetching results from
	// - limit: the maximum number of results to return

	after := c.Query("after")
	limit := c.Query("limit")

	var parsedLimit uint64 = 25

	if limit != "" {
		var err error
		parsedLimit, err = strconv.ParseUint(limit, 10, 64)
		if err != nil {
			s.logger.Warnf("unable to parse limit %s: %v", limit, err)
			c.JSON(400, gin.H{"error": "limit must be a number"})
			return nil, true
		}

		if parsedLimit > 100 {
			c.JSON(400, gin.H{"error": "limit must be less than or equal to 100"})
			return nil, true
		}
	}

	return &Pagination{
		After: after,
		Limit: parsedLimit,
	}, false

}
