package model

import (
	"context"
	"log"
	"time"

	"github.com/samaita/double-book/repository"
)

const (
	invoicePaymentMethodBSI = 1

	statusInvoiceDetailActive = 1

	statusInvoiceUnpaid = 1
	statusInvoicePaid   = 2
)

type Invoice struct {
	InvoiceID     int64           `json:"invoice_id"`
	User          int64           `json:"user_id"`
	Status        int             `json:"status"`
	PaymentMethod int             `json:"payment_method"`
	CreateTime    time.Time       `json:"create_time"`
	Detail        []InvoiceDetail `json:"invoice_detail"`
}

type InvoiceDetail struct {
	InvoiceDetailID int64 `json:"invoice_detail_id"`
	InvoiceID       int64 `json:"invoice_id"`
	ProductID       int64 `json:"product_id"`
	PriceFinal      int   `json:"price_final"`
	Amount          int   `json:"amount"`
	Status          int   `json:"status"`
}

func NewInvoice(id int64) Invoice {
	return Invoice{}
}

func (i *Invoice) Create(userID int64, status, paymentMethod int) error {
	var (
		query string
		err   error
	)

	query = `
		INSERT INTO invoice
		(user_id, status, payment_method, create_time)
		VALUES
		($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err = repository.DB.ExecContext(ctx, query, userID, status, paymentMethod, time.Now()); err != nil {
		log.Printf("[Invoice][Create][Exec] Input: %d Output: %v", userID, err)
		return err
	}

	return nil
}

func (i *Invoice) Add(productID int64, priceFinal, amount, status int) error {
	var (
		query string
		err   error
	)

	query = `
		INSERT INTO invoice_detail
		(product_id, priceFinal, amount, status, create_time)
		VALUES
		($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err = repository.DB.ExecContext(ctx, query, productID, priceFinal, amount, status, time.Now()); err != nil {
		log.Printf("[Invoice][Add][Exec] Input: %d Output: %v", productID, err)
		return err
	}

	return nil
}

func (i *Invoice) UpdateStatus(prevStatus, nextStatus int) error {
	var (
		query string
		err   error
	)

	query = `
		UPDATE invoice
		SET status = $1, update_time = $2
		WHERE invoice_id = $3 AND status = $4
	`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err = repository.DB.ExecContext(ctx, query, nextStatus, time.Now(), i.InvoiceID, prevStatus); err != nil {
		log.Printf("[Invoice][UpdateDetailStatus][Exec] Input: %d Output: %v", i.InvoiceID, err)
		return err
	}

	return nil
}

func (i *Invoice) GetDetail(status int) error {
	var (
		query string
		err   error
	)

	query = `
		SELECT invoice_detail_id, product_id, price_final, amount
		FROM invoice_detail
		WHERE invoice_id = $1 AND status = $2
		ORDER BY create_time DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := repository.DB.QueryContext(ctx, query, i.InvoiceID, status)
	if err != nil {
		log.Printf("[Invoice][GetDetail][Query] Input: %d Output: %v", i.InvoiceID, err)
		return err
	}
	defer result.Close()

	for result.Next() {
		var (
			invoiceDetailID, productID int64
			priceFinal, amount         int
		)

		if err = result.Scan(&invoiceDetailID, &productID, &priceFinal, &amount); err != nil {
			log.Printf("[Invoice][GetDetail][GetListInvoiceByUserID][Scan] Input: %d Output: %v", i.InvoiceID, err)
			return err
		}

		i.Detail = append(i.Detail, InvoiceDetail{
			InvoiceDetailID: invoiceDetailID,
			InvoiceID:       i.InvoiceID,
			ProductID:       productID,
			PriceFinal:      priceFinal,
			Status:          status,
		})
	}

	return nil
}

func GetListInvoiceByUserID(userID int64, status int) ([]Invoice, error) {
	var (
		query       string
		err         error
		listInvoice []Invoice
	)

	query = `
		SELECT invoice_id, payment_method, create_time
		FROM invoice
		WHERE user_id = $1 AND status = $2
		ORDER BY create_time DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := repository.DB.QueryContext(ctx, query, userID, status)
	if err != nil {
		log.Printf("[GetListInvoiceByUserID][Query] Input: %d Output: %v", userID, err)
		return listInvoice, err
	}
	defer result.Close()

	for result.Next() {
		var (
			invoiceID     int64
			paymentMethod int
			createTime    time.Time
		)

		if err = result.Scan(&invoiceID, &paymentMethod, &createTime); err != nil {
			log.Printf("[GetListInvoiceByUserID][Scan] Input: %d Output: %v", userID, err)
			return listInvoice, err
		}

		listInvoice = append(listInvoice, Invoice{
			InvoiceID:     invoiceID,
			PaymentMethod: paymentMethod,
			CreateTime:    createTime,
		})
	}

	return listInvoice, nil
}
