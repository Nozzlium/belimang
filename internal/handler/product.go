package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(
	productService *service.ProductService,
) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) Create(
	ctx *fiber.Ctx,
) error {
	merchantIdString := ctx.Params(
		"merchantId",
		"",
	)
	var body model.ProductRequestBody
	err := ctx.BodyParser(&body)
	if err != nil {
		err = constant.ErrBadInput
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[create product] failed to parse body: %v",
					err,
				),
			},
		)
	}

	merchantId, err := uuid.Parse(
		merchantIdString,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[create product] failed to parse merchantId: %v",
					err,
				),
			},
		)
	}

	productModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[create product] failed to validate body: %v",
					err,
				),
			},
		)
	}
	productModel.MerchantID = merchantId

	productId, err := h.productService.Create(
		ctx.Context(),
		productModel,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[create product] failed to create product: %v",
					err,
				),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"itemId": productId.String(),
		})
}

func (h *ProductHandler) FindAll(
	ctx *fiber.Ctx,
) error {
	var queries model.ProductQueries
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

	merchantIdString := ctx.Params(
		"merchantId",
	)
	merchantId, err := uuid.Parse(
		merchantIdString,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[find products] failed to find parse merchant id: %v",
					err,
				),
			},
		)
	}
	queries.MerchantId = merchantId

	productResp, err := h.productService.FindAll(
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
					"[find products] failed to find products: %v",
					err,
				),
			},
		)
	}

	return ctx.JSON(productResp)
}
