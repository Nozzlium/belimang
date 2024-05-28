package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/service"
)

type MerchantHandler struct {
	merchantService *service.MerchantService
}

func NewMerchantHandler(
	merchantService *service.MerchantService,
) *MerchantHandler {
	return &MerchantHandler{
		merchantService: merchantService,
	}
}

func (h *MerchantHandler) Create(
	ctx *fiber.Ctx,
) error {
	var body model.MerchantRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[create merchant] failed to parse body: %v",
					err,
				),
			},
		)
	}

	merchantModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[create merchant] failed to validate body: %v",
					err,
				),
			},
		)
	}

	merchantId, err := h.merchantService.Create(
		ctx.Context(),
		merchantModel,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[create merchant] failed creating merchant: %v",
					err,
				),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"merchantId": merchantId.String(),
		})
}

func (h *MerchantHandler) FindAll(
	ctx *fiber.Ctx,
) error {
	var queries model.MerchantQueries
	ctx.QueryParser(&queries)
	queries.Limit = ctx.QueryInt(
		"limit",
		5,
	)
	queries.Offset = ctx.QueryInt(
		"offset",
		0,
	)
	queries.CreatedAt = ctx.Query(
		"createdAt",
		"desc",
	)

	merchantData, total, err := h.merchantService.FindAll(
		ctx.Context(),
		queries,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[find merchant] failed to find merchants: %v",
					err,
				),
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"data": merchantData,
		"meta": fiber.Map{
			"limit":  queries.Limit,
			"offset": queries.Offset,
			"total":  total,
		},
	})
}
