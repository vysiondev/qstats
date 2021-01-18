package quaverapi_structs

type QuaverLeaderboardResponse struct {
	Status int64             `json:"status"`
	Users  []LeaderboardUser `json:"users"`
}

type LeaderboardUser struct {
	ID             int64            `json:"id"`
	SteamID        string           `json:"steam_id"`
	Username       string           `json:"username"`
	Country        string           `json:"country"`
	Allowed        int64            `json:"allowed"`
	Privileges     int64            `json:"privileges"`
	Usergroups     int64            `json:"usergroups"`
	AvatarURL      string           `json:"avatar_url"`
	TimeRegistered string           `json:"time_registered"`
	LatestActivity string           `json:"latest_activity"`
	Stats          LeaderboardStats `json:"stats"`
}

type LeaderboardStats struct {
	Rank                     int64   `json:"rank"`
	RankedScore              int64   `json:"ranked_score"`
	OverallAccuracy          float64 `json:"overall_accuracy"`
	OverallPerformanceRating float64 `json:"overall_performance_rating"`
	PlayCount                int64   `json:"play_count"`
	MaxCombo                 int64   `json:"max_combo"`
}
