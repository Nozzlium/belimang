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
