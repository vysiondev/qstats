package quaverapi_structs

type UserOnline struct {
	Status        int                     `json:"status"`
	IsOnline      bool                    `json:"is_online"`
	UserID        int                     `json:"user_id"`
	CurrentStatus UserOnlineCurrentStatus `json:"current_status"`
}

type UserOnlineCurrentStatus struct {
	Content string `json:"content"`
	Status  int    `json:"status"`
}
