package quaverapi_structs

type UserSearch struct {
	Status int `json:"status"`
	Users  []struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		SteamID   string `json:"steam_id"`
		AvatarURL string `json:"avatar_url"`
	} `json:"users"`
}
