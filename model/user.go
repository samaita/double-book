package model

const (
	userStatusDeactivate = -1
	userStatusRegistered = 0
	userStatusActive     = 1
	userStatusBanned     = 2

	DefaultThumbnail = "https://s4.anilist.co/file/anilistcdn/character/large/b1336-73LQxWKUWy78.png"
)

type User struct {
	UserID      int64  `json:"user_id"`
	DisplayName string `json:"display_name"`
	ImageURL    string `json:"image_url"`
	Status      int    `json:"status"`
}
