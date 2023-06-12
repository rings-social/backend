package server

import (
	"backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
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

func (s *Server) fillTestData() {
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Ring{
		Name:        "news",
		Description: "News from around the world",
		CreatedOn:   time.Now(),
	})

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.User{
		Username:    "random_dude",
		DisplayName: "Random Dude",
		SocialLinks: []models.SocialLink{
			{
				Platform: "twitter",
				Url:      "https://twitter.com/random_dude",
			},
		},
		CreatedOn: time.Now(),
	})

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Post{
		ID:             1,
		AuthorUsername: "random_dude",
		RingName:       "news",
		Title:          "Republican official appears to have moved $1.3m from nonprofit to own law firm",
		Link:           "https://www.theguardian.com/us-news/2023/jun/12/harmeet-dhillon-republican-lawyer-rnc-fox-news",
		Domain:         "theguardian.com",
		PostedOn:       time.Now(),
		Score:          1303,
	})
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
