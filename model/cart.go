package model

import (
	"context"
	"log"
	"time"

	"github.com/samaita/double-book/repository"
)

const (
	StatusCartActive = 1

	StatusCartDetailActive  = 1
	StatusCartDetailRemoved = 0

	StatusSuccessATC = 1
	StatusFailedATC  = -1
)

type Cart struct {
	CartID int64                `json:"cart_id"`
	UserID int64                `json:"user_id"`
	Status int                  `json:"status"`
	Detail map[int64]CartDetail `json:"detail"`
}

type CartDetail struct {
	CartDetailID int64 `json:"cart_detail_id"`
	CartID       int64 `json:"cart_id"`
	ProductID    int64 `json:"product_id"`
	Amount       int   `json:"amount"`
	Status       int   `json:"status"`
}

func NewCart(id int64) Cart {
	return Cart{
		CartID: id,
	}
}

func (c *Cart) Create(userID int64, status int) error {
	var (
		query string
		err   error
	)

	query = `
		INSERT INTO cart
		(user_id, status, create_time)
		VALUES
		($1, $2, $3)
		RETURNING cart_id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = repository.DB.QueryRowContext(ctx, query, userID, status, time.Now()).Scan(&c.CartID); err != nil {
		log.Printf("[Cart][Create][Exec] Input: %d Output: %v", c.UserID, err)
		return err
	}
	c.UserID = userID

	return nil
}

func (c *Cart) LoadByUser(userID int64) error {
	var (
		query string
		err   error
	)

	query = `
		SELECT cart_id, status 
		FROM cart
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = repository.DB.QueryRowContext(ctx, query, userID).Scan(&c.CartID, &c.Status); err != nil {
		log.Printf("[Cart][LoadByUser] Input: %d Output: %v", c.CartID, err)
		return err
	}

	return nil
}

func (c *Cart) Add(productID int64, amount int) error {
	var (
		query        string
		err          error
		rowsAffected int64
	)

	tx, err := repository.DB.Begin()
	if err != nil {
		log.Printf("[Cart][Add][Begin] Input: %d Output: %v", productID, err)
		return err
	}

	query = `
		UPDATE stock 
		SET remaining = remaining - 1
		WHERE product_id = $1 AND remaining > 0`

	result, err := tx.Exec(query, productID)
	if err != nil {
		log.Printf("[Cart][Add][Exec] Input: %d Output: %v", productID, err)
		tx.Rollback()
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil || rowsAffected == 0 {
		log.Printf("[Cart][Add][Exec] Input: %d Output: %v", productID, err)
		tx.Rollback()
		return err
	}

	query = `
		INSERT INTO cart_detail
		(cart_id, product_id, amount, status, create_time)
		VALUES
		($1, $2, $3, $4, $5)`

	if _, err = tx.Exec(query, c.CartID, productID, amount, StatusCartDetailActive, time.Now()); err != nil {
		log.Printf("[Cart][Add][Exec] Input: %d Output: %v", c.UserID, err)
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[Cart][Add][Commit] Input: %d Output: %v", productID, err)
		return err
	}

	return nil
}

func (c *Cart) Remove(productID int64, amount int) error {
	var (
		query          string
		err            error
		rowsAffected   int64
		existingAmount int
	)

	tx, err := repository.DB.Begin()
	if err != nil {
		log.Printf("[Cart][Remove][Begin] Input: %d Output: %v", productID, err)
		return err
	}

	query = `
		UPDATE stock 
		SET remaining = remaining + $1
		WHERE product_id = $2 AND remaining + $3 <= total`

	result, err := tx.Exec(query, amount, productID, amount)
	if err != nil {
		log.Printf("[Cart][Remove][Exec][Stock] Input: %d Output: %v", productID, err)
		tx.Rollback()
		return err
	}

	if rowsAffected, err = result.RowsAffected(); err != nil || rowsAffected == 0 {
		log.Printf("[Cart][Remove][RowsAffected][Stock] Input: %d Output: %v", productID, err)
		tx.Rollback()
		return err
	}

	query = `
		UPDATE cart_detail
		SET amount = amount - $1
		WHERE cart_id = $2 AND product_id = $3 AND amount > 0
		RETURNING amount`

	err = tx.QueryRow(query, amount, c.CartID, productID, amount, StatusCartDetailRemoved, time.Now()).Scan(&existingAmount)
	if err != nil {
		log.Printf("[Cart][Remove][QueryRow][Cart Detail] Input: %d Output: %v", c.UserID, err)
		tx.Rollback()
		return err
	}

	if existingAmount == 0 {
		query = `
		UPDATE cart_detail 
		SET status = $1
		WHERE product_id = $2 AND remaining + $3 <= total`

		if _, err := tx.Exec(query, StatusCartDetailRemoved, amount, productID, amount); err != nil {
			log.Printf("[Cart][Remove][Exect][Cart Detail] Input: %d Output: %v", productID, err)
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[Cart][Remove][Commit] Input: %d Output: %v", productID, err)
		return err
	}

	return nil
}

func (c *Cart) UpdateDetailStatus(prevStatus, nextStatus int) error {
	var (
		query string
		err   error
	)

	query = `
		UPDATE cart_detail
		SET status = $1, update_time = $2
		WHERE cart_id = $3 AND status = $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err = repository.DB.ExecContext(ctx, query, nextStatus, time.Now(), c.CartID, prevStatus); err != nil {
		log.Printf("[Cart][UpdateDetailStatus][Exec] Input: %d Output: %v", c.CartID, err)
		return err
	}

	return nil
}

func (c *Cart) GetDetail(status int) error {
	var (
		query string
		err   error
	)

	query = `
		SELECT cart_detail_id, product_id, amount 
		FROM cart_detail
		WHERE cart_id = $ AND status = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := repository.DB.QueryContext(ctx, query, c.CartID, status)
	if err != nil {
		log.Printf("[Cart][GetDetail][Query] Input: %d Output: %v", c.CartID, err)
		return err
	}
	defer result.Close()

	for result.Next() {
		var (
			cartDetail CartDetail
		)
		if err = result.Scan(&cartDetail.CartDetailID, &cartDetail.ProductID, &cartDetail.Amount); err != nil {
			log.Printf("[Cart][GetDetail][Scan] Input: %d Output: %v", c.CartID, err)
			return err
		}

		cartDetail.CartID = c.CartID
		cartDetail.Status = status

		c.Detail[cartDetail.ProductID] = cartDetail
	}

	return nil
}

func (c *Cart) IsExist(productID int64) (bool, error) {
	var (
		query string
		err   error
		exist int
	)

	query = `
		SELECT 1 
		FROM cart_detail
		WHERE cart_id = $1 AND user_id = $2  AND prduct_id = $3 AND status = $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = repository.DB.QueryRowContext(ctx, query, c.CartID, c.UserID, productID, StatusCartDetailActive).Scan(&exist); err != nil {
		log.Printf("[Cart][GetDetail][Query] Input: %d Output: %v", c.CartID, err)
		return false, err
	}

	return exist == 1, nil
}
