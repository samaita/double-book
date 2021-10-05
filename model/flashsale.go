package model

import (
	"context"
	"log"
	"time"

	"github.com/samaita/double-book/repository"
)

const (
	StatusFlashSalePublished = 2
	StatusFlashSaleSet       = 1
	StatusFlashSaleNonActive = 0
)

type FlashSale struct {
	FlashSaleID  int64             `json:"flashsale_id"`
	Name         string            `json:"name"`
	Status       int               `json:"status"`
	ScheduleTime time.Time         `json:"schedule_time"`
	Detail       []FlashSaleDetail `json:"detail"`
}

type FlashSaleDetail struct {
	FlashSaleDetailID int64 `json:"flashsale_detail_id"`
	FlashSaleID       int64 `json:"flashsale_id"`
	ProductID         int64 `json:"product_id"`
	Status            int64 `json:"status"`
}

func NewFlashSale(id int64) FlashSale {
	return FlashSale{
		FlashSaleID: id,
		Detail:      []FlashSaleDetail{},
	}
}

func (fs *FlashSale) Load() error {
	var (
		query string
		err   error
	)

	query = `
		SELECT name, status, schedule_time
		FROM flashsale
		WHERE flashsale_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = repository.DB.QueryRowContext(ctx, query, fs.FlashSaleID).Scan(&fs.Name, &fs.Status, &fs.ScheduleTime); err != nil {
		log.Printf("[FlashSale][Load] Input: %d Output: %v", fs.FlashSaleID, err)
		return err
	}

	return nil
}

func (fs *FlashSale) GetDetail() error {
	var (
		query string
		err   error
	)

	query = `
		SELECT flashsale_detail_id, product_id, status
		FROM flashsale_detail
		WHERE flashsale_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := repository.DB.QueryContext(ctx, query, fs.FlashSaleID)
	if err != nil {
		log.Printf("[FlashSale][GetDetail][Query] Input: %d Output: %v", fs.FlashSaleID, err)
		return err
	}
	defer result.Close()

	for result.Next() {
		var (
			flashSaleDetail FlashSaleDetail
		)

		if err = result.Scan(&flashSaleDetail.FlashSaleDetailID, &flashSaleDetail.ProductID, &flashSaleDetail.Status); err != nil {
			log.Printf("[FlashSale][GetDetail][Scan] Input: %d Output: %v", fs.FlashSaleID, err)
			return err
		}

		fs.Detail = append(fs.Detail, flashSaleDetail)
	}

	return nil
}

func GetFlashSaleByDate(date string, status int) ([]FlashSale, error) {
	var (
		query         string
		listFlashSale []FlashSale
		timeFlashSale time.Time
		err           error
	)

	if timeFlashSale, err = time.Parse("2006-01-02", date); err != nil {
		log.Printf("[GetFlashSaleByDate][Parse] Input: %s Output: %v", date, err)
		return listFlashSale, err
	}

	timeFlashSaleDayStart := timeFlashSale
	timeFlashSaleDayEnd := timeFlashSale.Add(23 * time.Hour).Add(59 * time.Minute)

	query = `
		SELECT flashsale_id, name, schedule_time
		FROM flashsale
		WHERE schedule_time BETWEEN $1 AND $2 AND status = $3
		ORDER BY schedule_time ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := repository.DB.QueryContext(ctx, query, timeFlashSaleDayStart, timeFlashSaleDayEnd, status)
	if err != nil {
		log.Printf("[GetFlashSaleByDate][Query] Input: %s Output: %v", date, err)
		return listFlashSale, err
	}
	defer result.Close()

	for result.Next() {
		var (
			FlashSaleID  int64
			Name         string
			ScheduleTime time.Time
		)

		if err = result.Scan(&FlashSaleID, &Name, &ScheduleTime); err != nil {
			log.Printf("[GetFlashSaleByDate][Scan] Input: %s Output: %v", date, err)
			return listFlashSale, err
		}

		listFlashSale = append(listFlashSale, FlashSale{
			FlashSaleID:  FlashSaleID,
			Name:         Name,
			Status:       status,
			ScheduleTime: ScheduleTime,
		})
	}

	return listFlashSale, nil
}
