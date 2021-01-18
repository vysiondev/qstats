package quaverapi_structs

type QuaverMapResponse struct {
	Status int64     `json:"status"`
	Map    QuaverMap `json:"map"`
}

type QuaverMap struct {
	ID                   int64       `json:"id"`
	MapsetID             int64       `json:"mapset_id"`
	Md5                  string      `json:"md5"`
	AlternativeMd5       interface{} `json:"alternative_md5"`
	CreatorID            int64       `json:"creator_id"`
	CreatorUsername      string      `json:"creator_username"`
	GameMode             int64       `json:"game_mode"`
	RankedStatus         int64       `json:"ranked_status"`
	Artist               string      `json:"artist"`
	Title                string      `json:"title"`
	Source               string      `json:"source"`
	Tags                 string      `json:"tags"`
	Description          string      `json:"description"`
	DifficultyName       string      `json:"difficulty_name"`
	Length               int64       `json:"length"`
	BPM                  float64     `json:"bpm"`
	DifficultyRating     float64     `json:"difficulty_rating"`
	CountHitobjectNormal int64       `json:"count_hitobject_normal"`
	CountHitobjectLong   int64       `json:"count_hitobject_long"`
	PlayCount            int64       `json:"play_count"`
	FailCount            int64       `json:"fail_count"`
	ModsPending          int64       `json:"mods_pending"`
	ModsAccepted         int64       `json:"mods_accepted"`
	ModsDenied           int64       `json:"mods_denied"`
	ModsIgnored          int64       `json:"mods_ignored"`
	DateSubmitted        string      `json:"date_submitted"`
	DateLastUpdated      string      `json:"date_last_updated"`
}
