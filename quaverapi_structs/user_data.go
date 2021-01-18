package quaverapi_structs

import "time"

type KeymodeStats struct {
	GlobalRank         int `json:"globalRank"`
	CountryRank        int `json:"countryRank"`
	MultiplayerWinRank int `json:"multiplayerWinRank"`
	Stats              struct {
		UserID                   int     `json:"user_id"`
		TotalScore               int     `json:"total_score"`
		RankedScore              int     `json:"ranked_score"`
		OverallAccuracy          float64 `json:"overall_accuracy"`
		OverallPerformanceRating float64 `json:"overall_performance_rating"`
		PlayCount                int     `json:"play_count"`
		FailCount                int     `json:"fail_count"`
		MaxCombo                 int     `json:"max_combo"`
		ReplaysWatched           int     `json:"replays_watched"`
		TotalMarv                int     `json:"total_marv"`
		TotalPerf                int     `json:"total_perf"`
		TotalGreat               int     `json:"total_great"`
		TotalGood                int     `json:"total_good"`
		TotalOkay                int     `json:"total_okay"`
		TotalMiss                int     `json:"total_miss"`
		TotalPauses              int     `json:"total_pauses"`
		MultiplayerWins          int     `json:"multiplayer_wins"`
		MultiplayerLosses        int     `json:"multiplayer_losses"`
		MultiplayerTies          int     `json:"multiplayer_ties"`
	} `json:"stats"`
}

type UserData struct {
	Status int `json:"status"`
	User   struct {
		Info struct {
			ID             int       `json:"id"`
			SteamID        string    `json:"steam_id"`
			Username       string    `json:"username"`
			TimeRegistered time.Time `json:"time_registered"`
			Allowed        int       `json:"allowed"`
			Privileges     int       `json:"privileges"`
			Usergroups     int       `json:"usergroups"`
			MuteEndtime    time.Time `json:"mute_endtime"`
			LatestActivity time.Time `json:"latest_activity"`
			Country        string    `json:"country"`
			AvatarURL      string    `json:"avatar_url"`
			Userpage       string    `json:"userpage"`
			Online         bool      `json:"online"`
		} `json:"info"`
		ProfileBadges []struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"profile_badges"`
		ActivityFeed []struct {
			ID        int       `json:"id"`
			Type      int       `json:"type"`
			Timestamp time.Time `json:"timestamp"`
			Map       struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"map"`
		} `json:"activity_feed"`
		Keys4 KeymodeStats `json:"keys4"`
		Keys7 KeymodeStats `json:"keys7"`
	} `json:"user"`
}
