package model

import (
	"context"
	"log"
	"time"

	"github.com/samaita/double-book/repository"
)

const (
	shopStatusDeleted = -1
	shopStatusActive  = 1
	shopStatusBanned  = 2

	DefaultShopName      = "Toko"
	DefaultShopThumbnail = "https://i1.sndcdn.com/avatars-000370077197-grlz13-t240x240.jpgg"
)

type Shop struct {
	ShopID   int64  `json:"shop_id"`
	UserID   int64  `json:"user_id"`
	Name     string `json:"display_name"`
	ImageURL string `json:"image_url"`
	Domain   string `json:"domain"`
	Status   int    `json:"status"`
}

func NewShop(id int64) Shop {
	return Shop{
		ShopID:   id,
		Name:     DefaultShopName,
		ImageURL: DefaultShopThumbnail,
	}
}

func (s *Shop) Load() error {
	var (
		query string
		err   error
	)

	query = `
		SELECT user_id, name, image_url, domain, status
		FROM shop
		WHERE shop_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = repository.DB.QueryRowContext(ctx, query, s.ShopID).Scan(&s.UserID, &s.Name, &s.ImageURL, &s.Domain, &s.Status); err != nil {
		log.Printf("[Shop][Load] Input: %d Output: %v", s.ShopID, err)
		return err
	}

	return nil
}
