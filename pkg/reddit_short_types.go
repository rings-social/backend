package server

import "backend/pkg/reddit_compat"

type RedditPosts reddit_compat.KindData[reddit_compat.Listing[reddit_compat.Post]]
type RedditSubreddits reddit_compat.KindData[reddit_compat.Listing[reddit_compat.Subreddit]]
type RedditAbout reddit_compat.KindData[reddit_compat.SubredditDetails]
type RedditComments reddit_compat.KindData[reddit_compat.Listing[reddit_compat.Comment]]
