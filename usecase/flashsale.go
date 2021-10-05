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
		log.Printf("[HandleGetFlashSaleByDate][GetFlashSaleByDate] Input: %s Output: %v", date, err)
		return resp, err
	}

	for _, flashSale := range flashSaleList {
		var (
			metadata FlashSaleListData
		)

		if err = flashSale.GetDetail(); err != nil {
			log.Printf("[HandleGetFlashSaleByDate][Flashsale][GetDetail] Input: %d Output: %v", flashSale.FlashSaleID, err)
			return resp, err
		}
		metadata.FlashSale = flashSale

		for _, flashSaleDetail := range flashSale.Detail {
			productID := flashSaleDetail.ProductID
			flashSaleProduct := model.NewProduct(productID)

			if err = flashSaleProduct.Load(); err != nil {
				log.Printf("[HandleGetFlashSaleByDate][Product][Load] Input: %d Output: %v", productID, err)
				return resp, err
			}

			flashSaleShop := model.NewShop(flashSaleProduct.ShopID)
			if err = flashSaleShop.Load(); err != nil {
				log.Printf("[HandleGetFlashSaleByDate][Shop][Load] Input: %d Output: %v", flashSaleProduct.ShopID, err)
				return resp, err
			}
			flashSaleProduct.Shop = flashSaleShop

			if err = flashSaleProduct.GetStock(); err != nil {
				log.Printf("[HandleGetFlashSaleByDate][Product][GetStock] Input: %d Output: %v", productID, err)
				return resp, err
			}
			metadata.ProductList = append(metadata.ProductList, flashSaleProduct)
		}

		resp = append(resp, metadata)
	}

	return resp, nil
}
