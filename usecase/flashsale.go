package usecase

import (
	"log"

	"github.com/samaita/double-book/model"
)

type FlashSaleListData struct {
	FlashSale   model.FlashSale `json:"flashsale"`
	ProductList []model.Product `json:"flashsale_product"`
}

func HandleGetFlashSaleByDate(date string) ([]FlashSaleListData, error) {
	var (
		resp          []FlashSaleListData
		flashSaleList []model.FlashSale
		err           error
	)

	if flashSaleList, err = model.GetFlashSaleByDate(date, model.StatusFlashSalePublished); err != nil {
		log.Printf("[handleGetFlashSaleByDate][GetFlashSaleByDate] Input: %s Output: %v", date, err)
		return resp, err
	}

	for _, flashSale := range flashSaleList {
		var (
			metadata FlashSaleListData
		)

		metadata.FlashSale = flashSale

		for _, flashSaleDetail := range flashSale.Detail {
			productID := flashSaleDetail.ProductID
			flashSaleProduct := model.NewProduct(productID)

			if err = flashSaleProduct.Load(); err != nil {
				log.Printf("[handleGetFlashSaleByDate][Product][Load] Input: %d Output: %v", productID, err)
				return resp, err
			}

			if err = flashSaleProduct.GetStock(); err != nil {
				log.Printf("[handleGetFlashSaleByDate][Product][GetStock] Input: %d Output: %v", productID, err)
				return resp, err
			}
			metadata.ProductList = append(metadata.ProductList, flashSaleProduct)
		}

		resp = append(resp, metadata)
	}

	return resp, nil
}
