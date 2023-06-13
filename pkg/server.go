package server

import (
	"backend/pkg/models"
	"backend/pkg/reddit_compat"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strings"
)

type Server struct {
	g      *gin.Engine
	db     *gorm.DB
	logger *logrus.Logger
}

func New(dsn string) (*Server, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	s := Server{
		g:      gin.New(),
		db:     db,
		logger: logrus.New(),
	}

	s.initRoutes()
	err = s.initModels()
	if err != nil {
		return nil, err
	}

	s.fillTestData()
	return &s, nil
}

func (s *Server) SetLogger(logger *logrus.Logger) {
	if logger != nil {
		s.logger = logger
	}
}

func (s *Server) Run(addr string) error {
	return s.g.Run(addr)
}

func (s *Server) initRoutes() {
	s.g.Use(gin.Recovery())
	s.g.Use(gin.Logger())
	s.g.GET("/healthz", s.healthz)

	g := s.g.Group("/api/v1")

	// Rings
	g.GET("/r/:ring", s.getRing)
	g.GET("/r/:ring/posts", s.getRingPosts)

	// Reddit-compatible API
	s.g.GET("/r/:ring/hot.json", s.getRcRingHot)
	s.g.GET("/subreddits/search.json", s.getRcRingsSearch)

	// Users
	g.GET("/u/:user", s.getUser)

}

func (s *Server) healthz(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": "ok",
	})
}

func (s *Server) getRing(context *gin.Context) {
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

func (s *Server) initModels() error {
	// Auto-migrate all the models in `models`
	return s.db.AutoMigrate(
		&models.Post{},
		&models.Ring{},
		&models.User{},
		&models.SocialLink{},
	)
}

func (s *Server) getRingPosts(context *gin.Context) {
	// Gets the posts in ring, sorted by score
	ringName := context.Param("ring")
	if ringName == "" {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "Ring name is required",
		})
		return
	}

	var posts []models.Post
	tx := s.db.Order("score desc").Find(&posts, "ring_name = ?", ringName)
	if tx.Error != nil {
		s.logger.Errorf("Unable to get posts for %s: %v", ringName, tx.Error)
		context.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to get posts",
		})
		return
	}

	context.JSON(200, posts)
}

func (s *Server) getUser(context *gin.Context) {
	username := context.Param("user")
	if username == "" {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "Username is required",
		})
		return
	}

	var user models.User
	tx := s.db.Preload("SocialLinks").First(&user, "username = ?", username)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			context.AbortWithStatusJSON(404, gin.H{
				"error": "User not found",
			})
			return
		}
		s.logger.Errorf("Unable to get user %s: %v", username, tx.Error)
		context.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to get user",
		})
		return
	}

	context.JSON(200, user)
}

func (s *Server) getRcRingHot(context *gin.Context) {
	ringName := context.Param("ring")
	after := context.Query("after")

	if after != "" {
		context.JSON(http.StatusOK, []string{})
		return
	}

	posts, err := s.repoRingPosts(ringName)
	if err != nil {
		context.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to get posts",
		})
		return
	}

	// Convert to Reddit-compatible format
	listing, err := toRedditPosts(posts)
	if err != nil {
		s.logger.Errorf("unable to convert posts: %v", err)
		internalServerError(context)
		return
	}

	context.JSON(200, listing)
}

func (s *Server) getRcRingsSearch(context *gin.Context) {
	q := context.Query("q")
	if q == "" {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "q is required",
		})
		return
	}

	nsfwQuery := context.Query("include_over_18")
	includeNsfw := false
	if nsfwQuery == "1" {
		includeNsfw = true
	}

	rings, err := s.repoRingsSearch(q, includeNsfw)
	if err != nil {
		s.logger.Errorf("unable to search rings: %v", err)
		internalServerError(context)
		return
	}

	// Convert to Reddit-compatible format
	listing, err := toRedditSubreddits(rings)
	if err != nil {
		s.logger.Errorf("unable to convert rings: %v", err)
		internalServerError(context)
		return
	}

	context.JSON(200, listing)
}

func parseNilAsEmpty(post *reddit_compat.Post) *reddit_compat.Post {
	// Given a RedditPosts struct, parse the struct tag for the `json` key and check if it does
	// have the `nilasempty` key. If it does, then set the value to an empty array.
	// This is needed because Reddit expects an empty array instead of null for some fields.

	t := reflect.TypeOf(post).Elem()
	v := reflect.ValueOf(post).Elem()
	num := t.NumField()
	// Iterate over the fields
	for i := 0; i < num; i++ {
		// Get the field
		field := t.Field(i)
		// Get the value of the field
		value := v.Field(i)
		// Get the json tag
		tag := field.Tag.Get("json")
		// Check if the tag has the `nilasempty` key
		if strings.Contains(tag, "nilasempty") {
			value.Set(reflect.MakeSlice(value.Type(), 0, 0))
		}
	}
	return post
}

func internalServerError(context *gin.Context) {
	context.AbortWithStatusJSON(500, gin.H{
		"error": "Internal server error",
	})
}
