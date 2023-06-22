package server

import (
	"backend/pkg/models"
	"gorm.io/gorm/clause"
	"time"
)

func (s *Server) fillTestData() {
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Ring{
		Name:         "news",
		Description:  "News from around the world",
		DisplayName:  "News",
		Title:        "Title",
		Subscribers:  3094892,
		CreatedAt:    time.Now(),
		PrimaryColor: "#FFC107",
	})

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Ring{
		Name:         "popular",
		Description:  "Popular Posts",
		CreatedAt:    time.Now(),
		Subscribers:  139843,
		PrimaryColor: "#49545f",
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
		CreatedAt: time.Now(),
	})

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.User{
		Username:    "john_doe",
		DisplayName: "John Doe",
		SocialLinks: []models.SocialLink{
			{
				Platform: "twitter",
				Url:      "https://twitter.com/john_doe",
			},
		},
		CreatedAt: time.Now(),
	})

	theGuardianLink := "https://www.theguardian.com/us-news/2023/jun/12/harmeet-dhillon-republican-lawyer-rnc-fox-news"
	theGuardianDomain := "theguardian.com"

	popularPost := models.Post{
		Model:          models.Model{ID: 1, CreatedAt: time.Now()},
		AuthorUsername: "random_dude",
		RingName:       "popular",
		Title:          "This is a popular post",
		Link:           &theGuardianLink,
		Domain:         &theGuardianDomain,
		Score:          1233,
		CommentsCount:  14,
		Nsfw:           false,
	}

	textPost := models.Post{
		Model:          models.Model{ID: 2, CreatedAt: time.Now()},
		AuthorUsername: "john_doe",
		RingName:       "popular",
		Title:          "This is a text post",
		Body:           "abc",
		Score:          4,
		CommentsCount:  0,
		Nsfw:           false,
	}

	nsfwPost := models.Post{
		Model:          models.Model{ID: 3, CreatedAt: time.Now()},
		AuthorUsername: "john_doe",
		RingName:       "news",
		Title:          "This is a NSFW post",
		Body:           "NSFW content goes here",
		Score:          1,
		Ups:            1,
		Downs:          0,
		CommentsCount:  0,
		Nsfw:           true,
	}

	newsPost := models.Post{
		Model:          models.Model{ID: 4, CreatedAt: time.Now()},
		AuthorUsername: "random_dude",
		RingName:       "news",
		Title:          "Republican official appears to have moved $0.3m from nonprofit to own law firm",
		Link:           &theGuardianLink,
		Domain:         &theGuardianDomain,
		Score:          1302,
		Nsfw:           false,
	}

	err := s.createOrUpdatePosts([]models.Post{popularPost,
		textPost,
		nsfwPost,
		newsPost,
	})

	if err != nil {
		s.logger.Fatalf("failed to create test posts: %v", err)
	}

	comment1 := models.Comment{
		Model:          models.Model{ID: 1, CreatedAt: time.Now()},
		AuthorUsername: "john_doe",
		PostId:         popularPost.ID,
		Body:           "Thanks for sharing!",
		Score:          99,
		Ups:            100,
		Downs:          1,
	}
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&comment1)
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Comment{
		Model:          models.Model{ID: 2, CreatedAt: time.Now()},
		AuthorUsername: "random_dude",
		PostId:         popularPost.ID,
		ParentId:       &comment1.ID,
		Body:           "You're welcome :)",
		Score:          32,
		Ups:            42,
		Downs:          10,
	})
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Comment{
		Model:          models.Model{ID: 3, CreatedAt: time.Now()},
		AuthorUsername: "john_doe",
		PostId:         popularPost.ID,
		Body:           "This comment doesn't have any replies",
		Score:          -1,
		Ups:            0,
		Downs:          1,
	})
}

func (s *Server) createOrUpdatePosts(posts []models.Post) error {
	for _, p := range posts {
		tx := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&p)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}
