package model

import (
	"context"
	"log"
	"time"

	"github.com/samaita/double-book/repository"
)

const (
	userStatusDeactivate = -1
	userStatusRegistered = 0
	userStatusActive     = 1
	userStatusBanned     = 2

	DefaultUserDisplayName = "User"
	DefaultUserThumbnail   = "https://s4.anilist.co/file/anilistcdn/character/large/b1336-73LQxWKUWy78.png"
)

type User struct {
	UserID      int64  `json:"user_id"`
	DisplayName string `json:"display_name"`
	ImageURL    string `json:"image_url"`
	Status      int    `json:"status"`
}

func NewUser(id int64) User {
	return User{
		UserID:      id,
		DisplayName: DefaultUserDisplayName,
		ImageURL:    DefaultUserThumbnail,
	}
}

func (u *User) Load() error {
	var (
		query string
		err   error
	)

	query = `
		SELECT display_name,  image_url, status
		FROM user
		WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = repository.DB.QueryRowContext(ctx, query, u.UserID).Scan(&u.DisplayName, &u.ImageURL, &u.Status); err != nil {
		log.Printf("[User][Load] Input: %d Output: %v", u.UserID, err)
		return err
	}

	return nil
}
