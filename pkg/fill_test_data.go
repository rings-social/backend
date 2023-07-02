package server

import (
	"backend/pkg/models"
	"gorm.io/gorm/clause"
	"time"
)

func (s *Server) fillTestData() {
	slackLink := "https://join.slack.com/t/ringssocial/shared_invite/zt-1xyl4xys4-fhjfig1CqmALL~cWqiIGcQ"
	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Ring{
		Name:         "news",
		Description:  "News from around the world",
		DisplayName:  "News",
		Title:        "Title",
		Subscribers:  1000,
		CreatedAt:    time.Now(),
		PrimaryColor: "#FFC107",
	})

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Ring{
		Name:         "popular",
		Description:  "Popular Posts",
		CreatedAt:    time.Now(),
		Subscribers:  2000,
		PrimaryColor: "#49545f",
	})

	s.fillTestUsers()

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

	introductionPost := models.Post{
		Model:          models.Model{ID: 5, CreatedAt: time.Now()},
		AuthorUsername: "denvit",
		RingName:       "popular",
		Title:          "Welcome to the Rings.social",
		Link:           createRefString("/about"),
		Domain:         createRefString("rings.social"),
		Score:          10,
		CommentsCount:  1,
		Nsfw:           false,
	}

	notVisitedPost := models.Post{
		Model:          models.Model{ID: 6, CreatedAt: time.Now()},
		AuthorUsername: "denvit",
		RingName:       "popular",
		Title:          "Do not click me",
		Link:           createRefString("https://www.youtube.com/watch?v=dQw4w9WgXcQ"),
		Domain:         createRefString("youtube.com"),
		Score:          -1,
		CommentsCount:  0,
		Nsfw:           false,
	}

	err := s.createOrUpdatePosts([]models.Post{popularPost,
		textPost,
		nsfwPost,
		newsPost,
		introductionPost,
		notVisitedPost,
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

	s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&models.Comment{
		Model:          models.Model{ID: 4, CreatedAt: time.Now()},
		AuthorUsername: "denvit",
		PostId:         introductionPost.ID,
		Ups:            5,
		Downs:          0,
		Score:          5,
		Body: "To learn more about Rings, check out the [about page](/about).  \n\n" +
			"If you want to contribute to the project, check out [GitHub organization](https://github.com/rings-social)" +
			" and [join our Slack channel](" + slackLink + ").  \n\n" +
			"By the way, did you know that our comments support markdown? **Bold**, _italic_, `preformat`\n" +
			"<script>alert('hello')</script>\n" +
			"```js\n" +
			"console.log('hello')\n" +
			"```\n",
	})

	s.db.Exec("ALTER SEQUENCE comments_id_seq RESTART WITH 5;")
}

func createRefString(s string) *string {
	return &s
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

func (s *Server) fillTestUsers() {
	authSubj := "auth0|6498715e20e9a839adbfee8e"
	users := []models.User{
		{
			Username:    "random_dude",
			DisplayName: "Random Dude",
			// ProfilePicture: createRefString("https://images.unsplash.com/photo-1596075780750-81249df16d19?fit=crop&w=200&q=80"),
			ProfilePicture: nil,
			SocialLinks: []models.SocialLink{
				{
					Platform: "twitter",
					Url:      "https://twitter.com/random_dude",
				},
			},
			CreatedAt: time.Now(),
		},
		{
			Username:       "john_doe",
			DisplayName:    "John Doe",
			ProfilePicture: createRefString("https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?fit=crop&w=200&q=80"),
			SocialLinks: []models.SocialLink{
				{
					Platform: "twitter",
					Url:      "https://twitter.com/john_doe",
				},
			},
			CreatedAt: time.Now(),
		},
		{
			Username:       "denvit",
			DisplayName:    "Denys Vitali",
			Admin:          true,
			ProfilePicture: createRefString("https://pbs.twimg.com/profile_images/1441455322455949319/_0xwiskP_400x400.jpg"),
			SocialLinks: []models.SocialLink{
				{
					Platform: "twitter",
					Url:      "https://twitter.com/DenysVitali",
				},
			},
			AuthSubject: &authSubj,
			Badges: []models.Badge{
				{
					Id:              "admin",
					BackgroundColor: "#d70000",
					TextColor:       "#ffffff",
				},
				{
					Id:              "supporter",
					BackgroundColor: "#ffde3f",
					TextColor:       "#895900",
				},
			},
			CreatedAt: time.Now(),
		},
	}

	for _, u := range users {
		tx := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&u)
		if tx.Error != nil {
			s.logger.Fatalf("failed to create test users: %v", tx.Error)
		}
	}
}
