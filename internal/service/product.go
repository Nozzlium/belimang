package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/repository"
	"github.com/nozzlium/belimang/internal/util"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(
	productRepository *repository.ProductRepository,
) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) Create(
	ctx context.Context,
	product model.Product,
) (uuid.UUID, error) {
	userIDString := ctx.Value("userID").(string)
	userID, err := uuid.Parse(
		userIDString,
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	productId, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
	}

	currentDate := util.Now()
	product.ID = productId
	product.UserID = userID
	product.CreatedAt = currentDate

	err = s.productRepository.Insert(
		ctx,
		product,
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	return productId, nil
}

func (s *ProductService) FindAll(
	ctx context.Context,
	queries model.ProductQueries,
) (model.ProductItemsResponseBody, error) {
	products, total, err := s.productRepository.FindAll(
		ctx,
		queries,
	)
	if err != nil {
		return model.ProductItemsResponseBody{}, err
	}

	productData := make(
		[]model.ProductData,
		0,
		len(products),
	)
	for _, product := range products {
		productData = append(
			productData,
			product.ToProductData(),
		)
	}

	productResponse := model.ProductItemsResponseBody{
		Data: productData,
		Meta: model.ProductMeta{
			Limit:  queries.Limit,
			Offset: queries.Offset,
			Total:  total,
		},
	}

	return productResponse, nil
}
