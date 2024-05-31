package model

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/util"
)

type (
	ProductCategory string
)

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

type ProductQueries struct {
	ItemID          string          `query:"itemId"`
	Name            string          `query:"name"`
	ProductCategory ProductCategory `query:"productCategory"`
	MerchantId      uuid.UUID
	Limit           int
	Offset          int
	CreatedAt       string
}

func (q *ProductQueries) BuildWhereClauses() ([]string, []interface{}) {
	clauses := make([]string, 0, 4)
	params := make([]interface{}, 0, 4)

	clauses = append(
		clauses,
		"merchant_id = $%d",
	)
	params = append(
		params,
		q.MerchantId,
	)

	itemId, err := uuid.Parse(
		q.ItemID,
	)
	if err == nil {
		clauses = append(
			clauses,
			"id = $%d",
		)
		params = append(
			params,
			itemId,
		)
	}

	if q.Name != "" {
		clauses = append(
			clauses,
			"name ilike '%%' || $%d || '%%'",
		)
		params = append(params, q.Name)
	}

	switch q.ProductCategory {
	case Beverage,
		Food,
		Snack,
		Condiments,
		Additions:
		clauses = append(
			clauses,
			"product_category = $%d",
		)
		params = append(
			params,
			q.ProductCategory,
		)
	}

	return clauses, params
}

func (q *ProductQueries) BuildPagination() (string, []interface{}) {
	var params []interface{}

	limit := 5
	offset := 0
	if q.Limit > 0 {
		limit = q.Limit
	}
	if q.Offset > 0 {
		offset = q.Offset
	}
	params = append(
		params,
		limit,
		offset,
	)

	return "limit $%d offset $%d", params
}

func (q *ProductQueries) BuildOrderByClause() []string {
	var sqlClause []string

	if q.CreatedAt != "" ||
		OrderBy(
			q.CreatedAt,
		).IsValid() {
		sqlClause = append(
			sqlClause,
			fmt.Sprintf(
				"created_at %s",
				q.CreatedAt,
			),
		)
	} else {
		sqlClause = append(
			sqlClause,
			"created_at desc",
		)
	}

	return sqlClause
}

type ProductItemsResponseBody struct {
	Data []ProductData `json:"data"`
	Meta ProductMeta   `json:"meta"`
}

type ProductData struct {
	ItemId          string  `json:"itemId"`
	Name            string  `json:"name"`
	ProductCategory string  `json:"productCategory"`
	Price           float64 `json:"price"`
	ImageURl        string  `json:"imageUrl"`
	CreatedAt       string  `json:"createdAt"`
}

func (product *Product) ToProductData() ProductData {
	return ProductData{
		ItemId: product.ID.String(),
		Name:   product.Name,
		ProductCategory: string(
			product.ProductCategory,
		),
		Price:    product.Price,
		ImageURl: product.ImageURL,
		CreatedAt: util.ToISO8601(
			product.CreatedAt,
		),
	}
}

type ProductMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
