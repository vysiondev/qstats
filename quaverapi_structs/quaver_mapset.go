package quaverapi_structs

type QuaverMapsetResponse struct {
	Status int64        `json:"status"`
	Mapset QuaverMapset `json:"mapset"`
}

type QuaverMapset struct {
	ID                      int64       `json:"id"`
	CreatorID               int64       `json:"creator_id"`
	CreatorUsername         string      `json:"creator_username"`
	CreatorAvatarURL        string      `json:"creator_avatar_url"`
	Artist                  string      `json:"artist"`
	Title                   string      `json:"title"`
	Source                  string      `json:"source"`
	Tags                    string      `json:"tags"`
	Description             interface{} `json:"description"`
	DateSubmitted           string      `json:"date_submitted"`
	DateLastUpdated         string      `json:"date_last_updated"`
	RankingQueueStatus      interface{} `json:"ranking_queue_status"`
	RankingQueueLastUpdated string      `json:"ranking_queue_last_updated"`
	RankingQueueVoteCount   interface{} `json:"ranking_queue_vote_count"`
	Maps                    []QuaverMap `json:"maps"`
}
