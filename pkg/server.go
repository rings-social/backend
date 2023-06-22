package server

import (
	"backend/pkg/models"
	"backend/pkg/reddit_compat"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

type Server struct {
	g       *gin.Engine
	db      *gorm.DB
	logger  *logrus.Logger
	baseUrl string
}

func New(dsn string, baseUrl string) (*Server, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	s := Server{
		g:       gin.New(),
		db:      db,
		logger:  logrus.New(),
		baseUrl: baseUrl,
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
		&models.Comment{},
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

	// Get ring
	r, err := s.repoRingAbout(ringName)
	if err != nil {
		s.logger.Errorf("Unable to get ring %s: %v", ringName, err)
		internalServerError(context)
		return
	}

	var posts []models.Post
	tx := s.db.Order("score desc").
		Find(&posts, "ring_name = ?", ringName)
	if tx.Error != nil {
		s.logger.Errorf("Unable to get posts for %s: %v", ringName, tx.Error)
		context.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to get posts",
		})
		return
	}

	context.JSON(200, convertResponsePosts(posts, r))
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
		s.convertToRedditPosts(context, []models.Post{})
		return
	}

	posts, err := s.repoRingPosts(ringName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatusJSON(404, gin.H{
				"error": "Ring not found",
			})
			return
		}
		s.logger.Errorf("Unable to get posts for %s: %v", ringName, err)
		context.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to get posts",
		})
		return
	}

	s.convertToRedditPosts(context, posts)
}

func (s *Server) convertToRedditPosts(context *gin.Context, posts []models.Post) {
	// Convert to Reddit-compatible format
	listing, err := toRedditPosts(posts, s.baseUrl)
	if err != nil {
		s.logger.Errorf("unable to convert posts: %v", err)
		internalServerError(context)
		return
	}

	if listing.Data.Children == nil {
		listing.Data.Children = []reddit_compat.KindData[reddit_compat.Post]{}
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

func parseNilAsEmpty[T any](element T) T {
	// Given a RedditPosts struct, parse the struct tag for the `json` key and check if it does
	// have the `nilasempty` key. If it does, then set the value to an empty array.
	// This is needed because Reddit expects an empty array instead of null for some fields.

	t := reflect.TypeOf(element).Elem()
	v := reflect.ValueOf(element).Elem()
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
	return element
}

func internalServerError(context *gin.Context) {
	context.AbortWithStatusJSON(500, gin.H{
		"error": "Internal server error",
	})
}

func badRequest(context *gin.Context) {
	context.AbortWithStatusJSON(400, gin.H{
		"error": "Bad request",
	})
}
