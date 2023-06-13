package reddit_compat

type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type Award struct {
	AwardSubType                     string  `json:"award_sub_type"`
	AwardType                        string  `json:"award_type"`
	AwardingsRequiredToGrantBenefits any     `json:"awardings_required_to_grant_benefits"`
	CoinPrice                        int     `json:"coin_price"`
	CoinReward                       int     `json:"coin_reward"`
	Count                            int     `json:"count"`
	DaysOfDripExtension              *int    `json:"days_of_drip_extension"`
	DaysOfPremium                    *int    `json:"days_of_premium"`
	Description                      string  `json:"description"`
	EndDate                          any     `json:"end_date"`
	GiverCoinReward                  any     `json:"giver_coin_reward"`
	IconFormat                       *string `json:"icon_format"`
	IconHeight                       int     `json:"icon_height"`
	IconURL                          string  `json:"icon_url"`
	IconWidth                        int     `json:"icon_width"`
	ID                               string  `json:"id"`
	IsEnabled                        bool    `json:"is_enabled"`
	IsNew                            bool    `json:"is_new"`
	Name                             string  `json:"name"`
	PennyDonate                      any     `json:"penny_donate"`
	PennyPrice                       *int    `json:"penny_price"`
	ResizedIcons                     []Image `json:"resized_icons,nilasempty"`
	ResizedStaticIcons               []Image `json:"resized_static_icons,nilasempty"`
	StartDate                        any     `json:"start_date"`
	StaticIconHeight                 int     `json:"static_icon_height"`
	StaticIconURL                    string  `json:"static_icon_url"`
	StaticIconWidth                  int     `json:"static_icon_width"`
	StickyDurationSeconds            any     `json:"sticky_duration_seconds"`
	SubredditCoinReward              int     `json:"subreddit_coin_reward"`
	SubredditID                      any     `json:"subreddit_id"`
	TiersByRequiredAwardings         any     `json:"tiers_by_required_awardings"`
}

type Flair struct {
	A string `json:"a,omitempty"`
	E string `json:"e"`
	T string `json:"t,omitempty"`
	U string `json:"u,omitempty"`
}

type GalleryItem struct {
	ID      int    `json:"id"`
	MediaID string `json:"media_id"`
}

type Gallery struct {
	Items []GalleryItem `json:"items,nilasempty"`
}

type Gilding struct {
	Gid2 int `json:"gid_2,omitempty"`
	Gid3 int `json:"gid_3,omitempty"`
}

type RedditVideo struct {
	BitrateKbps       int    `json:"bitrate_kbps"`
	DashURL           string `json:"dash_url"`
	Duration          int    `json:"duration"`
	FallbackURL       string `json:"fallback_url"`
	HasAudio          bool   `json:"has_audio"`
	Height            int    `json:"height"`
	HlsURL            string `json:"hls_url"`
	IsGif             bool   `json:"is_gif"`
	ScrubberMediaURL  string `json:"scrubber_media_url"`
	TranscodingStatus string `json:"transcoding_status"`
	Width             int    `json:"width"`
}

type Media struct {
	RedditVideo RedditVideo `json:"reddit_video"`
}

type MediaMetadata struct {
	E  string `json:"e"`
	ID string `json:"id"`
	M  string `json:"m"`
	O  []struct {
		U string `json:"u"`
		X int    `json:"x"`
		Y int    `json:"y"`
	} `json:"o"`
	P []struct {
		U string `json:"u"`
		X int    `json:"x"`
		Y int    `json:"y"`
	} `json:"p"`
	S struct {
		U string `json:"u"`
		X int    `json:"x"`
		Y int    `json:"y"`
	} `json:"s"`
	Status string `json:"status"`
}

type PollData struct {
	IsPrediction bool `json:"is_prediction"`
	Options      []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"options"`
	PredictionStatus    any `json:"prediction_status"`
	ResolvedOptionID    any `json:"resolved_option_id"`
	TotalStakeAmount    any `json:"total_stake_amount"`
	TotalVoteCount      int `json:"total_vote_count"`
	TournamentID        any `json:"tournament_id"`
	UserSelection       any `json:"user_selection"`
	UserWonAmount       any `json:"user_won_amount"`
	VoteUpdatesRemained any `json:"vote_updates_remained"`
	VotingEndTimestamp  int `json:"voting_end_timestamp"`
}

type MultiresImage struct {
	ID          string   `json:"id"`
	Resolutions []Image  `json:"resolutions,nilasempty"`
	Source      Image    `json:"source"`
	Variants    struct{} `json:"variants"`
}

type Preview struct {
	Enabled bool            `json:"enabled"`
	Images  []MultiresImage `json:"images,nilasempty"`
}

type Post struct {
	AllAwardings               []Award                   `json:"all_awardings,nilasempty"`
	AllowLiveComments          bool                      `json:"allow_live_comments"`
	ApprovedAtUtc              any                       `json:"approved_at_utc"`
	ApprovedBy                 any                       `json:"approved_by"`
	Archived                   bool                      `json:"archived"`
	Author                     string                    `json:"author"`
	AuthorFlairBackgroundColor *string                   `json:"author_flair_background_color"`
	AuthorFlairCssClass        *string                   `json:"author_flair_css_class"`
	AuthorFlairRichtext        []Flair                   `json:"author_flair_richtext,nilasempty"`
	AuthorFlairTemplateID      *string                   `json:"author_flair_template_id"`
	AuthorFlairText            *string                   `json:"author_flair_text"`
	AuthorFlairTextColor       *string                   `json:"author_flair_text_color"`
	AuthorFlairType            string                    `json:"author_flair_type"`
	AuthorFullname             string                    `json:"author_fullname"`
	AuthorIsBlocked            bool                      `json:"author_is_blocked"`
	AuthorPatreonFlair         bool                      `json:"author_patreon_flair"`
	AuthorPremium              bool                      `json:"author_premium"`
	Awarders                   []any                     `json:"awarders,nilasempty"`
	BannedAtUtc                any                       `json:"banned_at_utc"`
	BannedBy                   any                       `json:"banned_by"`
	CanGild                    bool                      `json:"can_gild"`
	CanModPost                 bool                      `json:"can_mod_post"`
	Category                   any                       `json:"category"`
	Clicked                    bool                      `json:"clicked"`
	ContentCategories          []string                  `json:"content_categories,nilasempty"`
	ContestMode                bool                      `json:"contest_mode"`
	Created                    int                       `json:"created"`
	CreatedUtc                 int                       `json:"created_utc"`
	DiscussionType             any                       `json:"discussion_type"`
	Distinguished              *string                   `json:"distinguished"`
	Domain                     *string                   `json:"domain"`
	Downs                      int                       `json:"downs"`
	Edited                     any                       `json:"edited"`
	GalleryData                *Gallery                  `json:"gallery_data,omitempty"`
	Gilded                     int                       `json:"gilded"`
	Gildings                   Gilding                   `json:"gildings"`
	Hidden                     bool                      `json:"hidden"`
	HideScore                  bool                      `json:"hide_score"`
	ID                         string                    `json:"id"`
	IsCreatedFromAdsUi         bool                      `json:"is_created_from_ads_ui"`
	IsCrosspostable            bool                      `json:"is_crosspostable"`
	IsGallery                  bool                      `json:"is_gallery,omitempty"`
	IsMeta                     bool                      `json:"is_meta"`
	IsOriginalContent          bool                      `json:"is_original_content"`
	IsRedditMediaDomain        bool                      `json:"is_reddit_media_domain"`
	IsRobotIndexable           bool                      `json:"is_robot_indexable"`
	IsSelf                     bool                      `json:"is_self"`
	IsVideo                    bool                      `json:"is_video"`
	Likes                      any                       `json:"likes"`
	LinkFlairBackgroundColor   string                    `json:"link_flair_background_color"`
	LinkFlairCssClass          *string                   `json:"link_flair_css_class"`
	LinkFlairRichtext          []Flair                   `json:"link_flair_richtext,nilasempty"`
	LinkFlairTemplateID        string                    `json:"link_flair_template_id,omitempty"`
	LinkFlairText              *string                   `json:"link_flair_text"`
	LinkFlairTextColor         string                    `json:"link_flair_text_color"`
	LinkFlairType              string                    `json:"link_flair_type"`
	Locked                     bool                      `json:"locked"`
	Media                      *Media                    `json:"media"`
	MediaEmbed                 struct{}                  `json:"media_embed"`
	MediaMetadata              *map[string]MediaMetadata `json:"media_metadata,omitempty"`
	MediaOnly                  bool                      `json:"media_only"`
	ModNote                    any                       `json:"mod_note"`
	ModReasonBy                any                       `json:"mod_reason_by"`
	ModReasonTitle             any                       `json:"mod_reason_title"`
	ModReports                 []any                     `json:"mod_reports"`
	Name                       string                    `json:"name"`
	NoFollow                   bool                      `json:"no_follow"`
	NumComments                int                       `json:"num_comments"`
	NumCrossposts              int                       `json:"num_crossposts"`
	NumReports                 any                       `json:"num_reports"`
	Over18                     bool                      `json:"over_18"`
	ParentWhitelistStatus      *string                   `json:"parent_whitelist_status"`
	Permalink                  string                    `json:"permalink"`
	Pinned                     bool                      `json:"pinned"`
	PollData                   *PollData                 `json:"poll_data,omitempty"`
	PostHint                   string                    `json:"post_hint,omitempty"`
	Preview                    *Preview                  `json:"preview,omitempty"`
	Pwls                       *int                      `json:"pwls"`
	Quarantine                 bool                      `json:"quarantine"`
	RemovalReason              any                       `json:"removal_reason"`
	RemovedBy                  any                       `json:"removed_by"`
	RemovedByCategory          any                       `json:"removed_by_category"`
	ReportReasons              any                       `json:"report_reasons"`
	Saved                      bool                      `json:"saved"`
	Score                      int                       `json:"score"`
	SecureMedia                *Media                    `json:"secure_media"`
	SecureMediaEmbed           struct {
	} `json:"secure_media_embed"`
	Selftext              string  `json:"selftext"`
	SelftextHtml          *string `json:"selftext_html"`
	SendReplies           bool    `json:"send_replies"`
	Spoiler               bool    `json:"spoiler"`
	Stickied              bool    `json:"stickied"`
	Subreddit             string  `json:"subreddit"`
	SubredditID           string  `json:"subreddit_id"`
	SubredditNamePrefixed string  `json:"subreddit_name_prefixed"`
	SubredditSubscribers  int     `json:"subreddit_subscribers"`
	SubredditType         string  `json:"subreddit_type"`
	SuggestedSort         *string `json:"suggested_sort"`
	Thumbnail             string  `json:"thumbnail"`
	ThumbnailHeight       *int    `json:"thumbnail_height"`
	ThumbnailWidth        *int    `json:"thumbnail_width"`
	Title                 string  `json:"title"`
	TopAwardedType        *string `json:"top_awarded_type"`
	TotalAwardsReceived   int     `json:"total_awards_received"`
	TreatmentTags         []any   `json:"treatment_tags"`
	Ups                   int     `json:"ups"`
	UpvoteRatio           float64 `json:"upvote_ratio"`
	URL                   *string `json:"url"`
	URLOverriddenByDest   string  `json:"url_overridden_by_dest,omitempty"`
	UserReports           []any   `json:"user_reports"`
	ViewCount             any     `json:"view_count"`
	Visited               bool    `json:"visited"`
	WhitelistStatus       *string `json:"whitelist_status"`
	Wls                   *int    `json:"wls"`
}
