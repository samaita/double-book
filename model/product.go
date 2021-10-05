package model

import (
	"context"
	"log"
	"time"

	"github.com/samaita/double-book/repository"
)

const (
	productStockEmpty = 0
)

type Product struct {
	ProductID     int64  `json:"product_id"`
	ShopID        int64  `json:"shop_id"`
	Shop          Shop   `json:"shop"`
	Name          string `json:"name"`
	ImageURL      string `json:"image_url"`
	PriceNormal   int64  `json:"price_normal"`
	PriceDiscount int64  `json:"price_discount"`
	Total         int    `json:"total"`
	Remaining     int    `json:"remaining"`
	Status        int    `json:"status"`
}

func NewProduct(id int64) Product {
	return Product{
		ProductID: id,
	}
}

func (p *Product) Load() error {
	var (
		query string
		err   error
	)

	query = `
		SELECT shop_id, name, image_url, status
		FROM product
		WHERE product_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = repository.DB.QueryRowContext(ctx, query, p.ProductID).Scan(&p.ShopID, &p.Name, &p.ImageURL, &p.Status); err != nil {
		log.Printf("[Product][Load] Input: %d Output: %v", p.ProductID, err)
		return err
	}

	return nil
}

func (p *Product) GetStock() error {
	var (
		query       string
		err         error
		stockStatus int
	)

	query = `
		SELECT price_normal, price_discount, total, remaining, status
		FROM stock
		WHERE product_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = repository.DB.QueryRowContext(ctx, query, p.ProductID).Scan(&p.PriceNormal, &p.PriceDiscount, &p.Total, &p.Remaining, &stockStatus); err != nil {
		log.Printf("[Product][GetStock] Input: %d Output: %v", p.ProductID, err)
		return err
	}

	// force empty
	if stockStatus == productStockEmpty {
		p.Remaining = 0
	}

	return nil
}
