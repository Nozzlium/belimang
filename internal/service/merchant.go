package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/repository"
	"github.com/nozzlium/belimang/internal/util"
)

type MerchantService struct {
	merchantRepository *repository.MerchantRepository
}

func NewMerchantService(
	merchantRepository *repository.MerchantRepository,
) *MerchantService {
	return &MerchantService{
		merchantRepository: merchantRepository,
	}
}

func (s *MerchantService) Create(
	ctx context.Context,
	merchant model.Merchant,
) (uuid.UUID, error) {
	userIDString := ctx.Value("userID").(string)
	userID, err := uuid.Parse(
		userIDString,
	)
	if err != nil {
		return uuid.UUID{}, err
	}

	merchantId, err := uuid.NewV7()
	if err != nil {
		return merchantId, err
	}

	currentDate := util.Now()
	merchant.ID = merchantId
	merchant.CreatedAt = currentDate
	merchant.UserID = userID

	_, err = s.merchantRepository.Insert(
		ctx,
		merchant,
	)
	if err != nil {
		return merchantId, err
	}

	return merchantId, nil
}

func (s *MerchantService) FindAll(
	ctx context.Context,
	merchantQueries model.MerchantQueries,
) ([]model.MerchantResponaeBody, int, error) {
	merchantData := make(
		[]model.MerchantResponaeBody,
		0,
		merchantQueries.Limit,
	)
	merchants, total, err := s.merchantRepository.FindAll(
		ctx,
		merchantQueries,
	)
	if err != nil {
		return nil, 0, err
	}

	for _, merchant := range merchants {
		merchantData = append(
			merchantData,
			merchant.ToResponseBody(),
		)
	}

	return merchantData, total, nil
}
