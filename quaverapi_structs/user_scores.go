package quaverapi_structs

type UserScores struct {
	Status int64   `json:"status"`
	Scores []Score `json:"scores"`
}

type Score struct {
	ID                int64   `json:"id"`
	Time              string  `json:"time"`
	Mode              int64   `json:"mode"`
	Mods              int64   `json:"mods"`
	ModsString        string  `json:"mods_string"`
	PerformanceRating float64 `json:"performance_rating"`
	PersonalBest      bool    `json:"personal_best"`
	TotalScore        int64   `json:"total_score"`
	Accuracy          float64 `json:"accuracy"`
	Grade             string  `json:"grade"`
	MaxCombo          int64   `json:"max_combo"`
	CountMarv         int64   `json:"count_marv"`
	CountPerf         int64   `json:"count_perf"`
	CountGreat        int64   `json:"count_great"`
	CountGood         int64   `json:"count_good"`
	CountOkay         int64   `json:"count_okay"`
	CountMiss         int64   `json:"count_miss"`
	ScrollSpeed       int64   `json:"scroll_speed"`
	Ratio             float64 `json:"ratio"`
	Map               Map     `json:"map"`
}

type Map struct {
	ID              int64  `json:"id"`
	MapsetID        int64  `json:"mapset_id"`
	Md5             string `json:"md5"`
	Artist          string `json:"artist"`
	Title           string `json:"title"`
	DifficultyName  string `json:"difficulty_name"`
	CreatorID       int64  `json:"creator_id"`
	CreatorUsername string `json:"creator_username"`
	RankedStatus    int64  `json:"ranked_status"`
}

type UserScoresResponse struct {
	Scores []Score
	User   UserData
}
