package server

import (
	"backend/pkg/models"
	"gorm.io/gorm/clause"
	"time"
)

func (s *Server) fillTestData() {
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Ring{
		Name:        "news",
		Description: "News from around the world",
		DisplayName: "News",
		Title:       "Title",
		Subscribers: 3094892,
		CreatedOn:   time.Now(),
	})

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Ring{
		Name:        "popular",
		Description: "Popular Posts",
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

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.User{
		Username:    "john_doe",
		DisplayName: "John Doe",
		SocialLinks: []models.SocialLink{
			{
				Platform: "twitter",
				Url:      "https://twitter.com/john_doe",
			},
		},
		CreatedOn: time.Now(),
	})

	theGuardianLink := "https://www.theguardian.com/us-news/2023/jun/12/harmeet-dhillon-republican-lawyer-rnc-fox-news"
	theGuardianDomain := "theguardian.com"

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Post{
		ID:             1,
		AuthorUsername: "random_dude",
		RingName:       "news",
		Title:          "Republican official appears to have moved $1.3m from nonprofit to own law firm",
		Link:           &theGuardianLink,
		Domain:         &theGuardianDomain,
		PostedOn:       time.Now(),
		Score:          1303,
		Nsfw:           false,
	})
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Post{
		ID:             2,
		AuthorUsername: "random_dude",
		RingName:       "popular",
		Title:          "This is a popular post",
		Link:           &theGuardianLink,
		Domain:         &theGuardianDomain,
		PostedOn:       time.Now(),
		Score:          1234,
		CommentsCount:  15,
		Nsfw:           true,
	})
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Post{
		ID:             3,
		AuthorUsername: "john_doe",
		RingName:       "popular",
		Title:          "This is a text post",
		Body:           "abc",
		PostedOn:       time.Now(),
		Score:          5,
		CommentsCount:  1,
		Nsfw:           false,
	})
}
