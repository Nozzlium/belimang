package model

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/util"
)

type ProductCategory string

const (
	Beverage   ProductCategory = "Beverage"
	Food       ProductCategory = "Food"
	Snack      ProductCategory = "Snack"
	Condiments ProductCategory = "Condiments"
	Additions  ProductCategory = "Additions"
)

type Product struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	MerchantID      uuid.UUID
	Name            string
	Price           float64
	ProductCategory ProductCategory
	ImageURL        string
	CreatedAt       time.Time
}

type ProductRequestBody struct {
	Name            string          `json:"name"`
	ProductCategory ProductCategory `json:"productCategory"`
	Price           float64         `json:"price"`
	ImageUrl        string          `json:"imageUrl"`
}

func (body ProductRequestBody) IsValid() (Product, error) {
	var product Product
	if nameLen := len(body.Name); nameLen < 2 ||
		nameLen > 30 {
		log.Println("di nama")
		return product, constant.ErrBadInput
	}
	product.Name = body.Name

	switch body.ProductCategory {
	case Beverage,
		Food,
		Snack,
		Condiments,
		Additions:
		product.ProductCategory = body.ProductCategory
	default:
		log.Println("di kategori")
		return product, constant.ErrBadInput
	}

	if body.Price < 1 {
		log.Println("di harga")
		return product, constant.ErrBadInput
	}
	product.Price = body.Price

	err := util.ValidateURL(
		body.ImageUrl,
	)
	if err != nil {
		log.Println("di url")
		return product, constant.ErrBadInput
	}
	product.ImageURL = body.ImageUrl

	return product, nil
}
