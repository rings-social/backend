package server

import (
	"backend/pkg/models"
	authenticator "backend/pkg/platform/auth"
	"backend/pkg/reddit_compat"
	"backend/pkg/request"
	"backend/pkg/response"
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"
)

type Server struct {
	g            *gin.Engine
	db           *gorm.DB
	logger       *logrus.Logger
	baseUrl      string
	authProvider *authenticator.Authenticator
}

type Auth0Config struct {
	Domain       string
	ClientId     string
	ClientSecret string
}

func New(dsn string, auth0Config Auth0Config, baseUrl string) (*Server, error) {
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

	authProvider, err := authenticator.New(auth0Config.Domain, auth0Config.ClientId, auth0Config.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("unable to create authentication provider: %v", err)
	}

	s := Server{
		g:            gin.New(),
		db:           db,
		logger:       logrus.New(),
		baseUrl:      baseUrl,
		authProvider: authProvider,
	}

	s.initRoutes()
	err = s.initModels()
	if err != nil {
		return nil, err
	}

	// s.fillTestData()
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

func (s *Server) createRing(c *gin.Context) {
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

	ring := models.Ring{
		Name:          ringName,
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

func validateRingName(name string) bool {
	if len(name) < 3 || len(name) > 20 {
		return false
	}

	// Make sure the ring name is lowercase and can only contain
	// letters, numbers, and underscores
	for _, c := range name {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_' {
			return false
		}
	}

	return true
}

func (s *Server) initModels() error {
	// Auto-migrate all the models in `models`
	// Check if the comment_action enum exists
	var commentActionExists bool
	tx := s.db.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'comment_action');").Scan(&commentActionExists)
	if tx.Error != nil {
		return tx.Error
	}

	if !commentActionExists {
		tx := s.db.Exec("CREATE TYPE comment_action AS ENUM ('upvote', 'downvote');")
		if tx.Error != nil {
			return tx.Error
		}
	}

	return s.db.AutoMigrate(
		&models.Comment{},
		&models.Post{},
		&models.Ring{},
		&models.User{},
		&models.SocialLink{},
		&models.CommentAction{},
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
	tx := s.db.
		Preload("Author").
		Order("score desc").
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

func (s *Server) getMe(c *gin.Context) {
	idToken, done := s.idToken(c)
	if done {
		return
	}

	user, err := s.repoGetUserByAuthSubject(idToken.Subject)
	if err != nil {
		s.handleUserError(c, err)
		return
	}
	c.JSON(200, user)

}

func (s *Server) idToken(c *gin.Context) (*oidc.IDToken, bool) {
	v, exists := c.Get("id_token")
	if !exists {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "You must be authenticated to use this endpoint",
		})
		return nil, true
	}

	idToken, ok := v.(*oidc.IDToken)
	if !ok {
		s.logger.Errorf("Unable to cast id_token to *oidc.IDToken")
		internalServerError(c)
		return nil, true
	}
	return idToken, false
}

func (s *Server) getUser(context *gin.Context) {
	username := context.Param("username")
	if username == "" {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "Username is required",
		})
		return
	}

	user, err := s.repoGetUserByUsername(username)
	if err != nil {
		s.handleUserError(context, err)
		return
	}
	context.JSON(200, user)
}

func (s *Server) handleUserError(context *gin.Context, err error) {
	if err == gorm.ErrRecordNotFound {
		context.AbortWithStatusJSON(404, gin.H{
			"error": "User not found",
		})
		return
	}
	s.logger.Errorf("Unable to get user: %v", err)
	context.AbortWithStatusJSON(500, gin.H{
		"error": "Unable to get user",
	})
	return
}

func (s *Server) repoGetUserByUsername(username string) (models.User, error) {
	var user models.User
	tx := s.db.
		Preload("SocialLinks").
		Preload("Badges").
		First(&user, "username = ?", username)
	return user, tx.Error
}

func (s *Server) getUserProfilePicture(context *gin.Context) {
	username := context.Param("username")
	if username == "" {
		context.AbortWithStatusJSON(400, gin.H{
			"error": "Username is required",
		})
		return
	}

	var user models.User
	tx := s.db.First(&user, "username = ?", username)
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

	if user.ProfilePicture == nil {
		context.Redirect(302, s.baseUrl+"/default-profile-picture.jpg")
		return
	}

	context.Redirect(302, *user.ProfilePicture)
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

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") == "" {
			return
		}

		if strings.HasPrefix(c.Request.Header.Get("Authorization"), "Bearer ") {
			// Parse Bearer token
			token := strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer ")

			idToken, err := s.authProvider.VerifyToken(context.Background(), token)
			if err != nil {
				s.logger.Errorf("Unable to verify ID token: %v", err)
				c.AbortWithStatusJSON(401, gin.H{
					"error": "Unable to verify ID token",
				})
				return
			}

			c.Set("id_token", idToken)
		}
	}
}

func (s *Server) usernameAvailability(c *gin.Context) {
	usernameQuery := c.Query("username")
	valid, err := validateUsername(usernameQuery)
	if !valid {
		c.AbortWithStatusJSON(400, gin.H{
			"error": err,
		})
		return
	}

	// Check if username is available
	tx := s.db.First(&models.User{}, "username = ?", usernameQuery)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			c.JSON(200, gin.H{
				"available": true,
			})
			return
		}
		s.logger.Errorf("Unable to check username availability: %v", tx.Error)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to check username availability",
		})
		return
	}

	c.JSON(200, gin.H{
		"available": false,
	})
}

// signupUsername creates a user with the given username
// and associates it with the ID token
// It expects the username to be passed as a JSON body
func (s *Server) signupUsername(c *gin.Context) {
	idToken, done := s.idToken(c)
	if done {
		return
	}

	var request struct {
		Username string `json:"username"`
	}
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Username is required",
		})
		return
	}

	valid, errMsg := validateUsername(request.Username)
	if !valid {
		c.AbortWithStatusJSON(400, gin.H{
			"error": errMsg,
		})
		return
	}

	// Create a user with the username
	user := models.User{
		Username:    request.Username,
		AuthSubject: &idToken.Subject,
	}
	tx := s.db.Create(&user)
	if tx.Error != nil {
		s.logger.Errorf("Unable to create user: %v", tx.Error)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to create user",
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *Server) usernameForIdToken(token *oidc.IDToken) (string, error) {
	if token == nil {
		return "", fmt.Errorf("token is nil")
	}
	var user models.User
	tx := s.db.First(&user, "auth_subject = ?", token.Subject)
	if tx.Error != nil {
		return "", tx.Error
	}
	return user.Username, nil
}

func (s *Server) addComment(postId uint, username string, request request.Comment) (models.Comment, error) {
	// TODO: Check if user can comment here
	comment := models.Comment{
		PostId:         postId,
		Body:           request.Content,
		AuthorUsername: username,
		ParentId:       request.ParentId,
	}

	if request.ParentId != nil {
		// Get original parent comment
		var parentComment models.Comment
		tx := s.db.First(&parentComment, "id = ?", request.ParentId)
		if tx.Error != nil {
			return comment, tx.Error
		}

		comment.Depth = parentComment.Depth + 1
	}

	tx := s.db.Create(&comment)
	if tx.Error == nil {
		// Add +1 to comment count in post
		tx = s.db.Model(&models.Post{}).
			Where("id = ?", postId).
			Update("comments_count",
				gorm.Expr("comments_count + ?", 1),
			)
		if tx.Error != nil {
			return comment, tx.Error
		}

		// Fetch comment with Preload("Author")
		tx = s.db.Preload("Author").First(&comment, comment.ID)
	}
	return comment, tx.Error
}

func (s *Server) isAdmin(username string) bool {
	var user models.User
	tx := s.db.First(&user, "username = ?", username)
	if tx.Error != nil {
		// Cannot be found / other error
		return false
	}

	return user.Admin
}

func (s *Server) repoCommentVoteAction(commentId uint, username string, action models.VoteAction) (models.CommentAction, int, error) {
	commentAction := models.CommentAction{
		Username:  username,
		CommentId: commentId,
		Action:    action,
	}
	addScore := 0

	// Check if user has already voted
	var existingAction models.CommentAction
	tx := s.db.First(&existingAction, "username = ? AND comment_id = ?", username, commentId)
	if tx.Error == nil {
		if action == existingAction.Action {
			// User has already voted
			return commentAction, 0, fmt.Errorf("user has already voted")
		}

		// User has changed their vote
		if existingAction.Action == models.ActionUpvote && action == models.ActionDownvote {
			addScore = -2
		} else if existingAction.Action == models.ActionDownvote && action == models.ActionUpvote {
			addScore = 2
		}
	} else {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			if action == models.ActionUpvote {
				addScore = 1
			} else if action == models.ActionDownvote {
				addScore = -1
			}
		}
	}

	tx = s.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}, {Name: "comment_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"action"}),
		},
	).Create(&commentAction)
	return commentAction, addScore, tx.Error
}

func (s *Server) repoDownvoteComment(commentId uint, username string) (models.CommentAction, error) {
	commentAction := models.CommentAction{
		Username:  username,
		CommentId: commentId,
		Action:    models.ActionDownvote,
	}
	tx := s.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "username"}, {Name: "comment_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"action"}),
		},
	).Create(&commentAction)
	return commentAction, tx.Error
}

func (s *Server) repoIncreaseCommentScore(commentId uint, amount int) error {
	tx := s.db.Model(&models.Comment{}).
		Where("id = ?", commentId).
		UpdateColumn("score", gorm.Expr("score + ?", amount))
	return tx.Error
}

func (s *Server) hasIdToken(c *gin.Context) bool {
	_, exists := c.Get("id_token")
	return exists
}

func (s *Server) repoCommentActions(username string, postId int64) []models.CommentAction {
	var commentActions []models.CommentAction
	tx := s.db.Model(&models.CommentAction{}).
		Where("username = ?", username).
		Joins("JOIN comments ON comments.id = comment_actions.comment_id").
		Where("comments.post_id = ?", postId).
		Find(&commentActions)
	if tx.Error != nil {
		return []models.CommentAction{}
	}
	return commentActions
}

// getUsers returns a paginated list of users
func (s *Server) getUsers(c *gin.Context) {
	var users []models.User
	var count int64
	pagination, done := s.getPagination(c)
	if done {
		return
	}

	tx := s.db.
		Model(&models.User{}).
		Where("deleted_at IS NULL").
		Count(&count)
	if tx.Error != nil {
		s.logger.Errorf("Failed to get users: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get users",
		})
		return
	}

	tx = s.db.
		Preload("SocialLinks").
		Order("username ASC").
		Limit(int(pagination.Limit)).
		Where("username > ? AND deleted_at IS NULL", pagination.After).
		Find(&users)
	if tx.Error != nil {
		s.logger.Errorf("Failed to get users: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get users",
		})
		return
	}

	responseUsers := cleanupUsers(users)
	outputPagination := response.Paginated[models.User]{
		Items: responseUsers,
		Total: count,
	}
	if len(responseUsers) > 0 {
		outputPagination.After = responseUsers[len(responseUsers)-1].Username
	}

	c.JSON(http.StatusOK, outputPagination)
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
