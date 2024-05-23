package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(
	userService *service.UserService,
) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterAdmin(
	ctx *fiber.Ctx,
) error {
	var body model.UserRegisterBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[register admin] failed to parse body: %v",
					err,
				),
			},
		)
	}

	userModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[register admin] failed to validate body: %v",
					err,
				),
			},
		)
	}

	token, err := h.userService.RegisterAdmin(
		ctx.Context(),
		userModel,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[register admin] failed to store user: %v",
					err,
				),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"token": token,
		})
}

func (h *UserHandler) LoginAdmin(
	ctx *fiber.Ctx,
) error {
	var body model.UserLoginBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[login admin] failed to parse body: %v",
					err,
				),
			},
		)
	}

	userModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[login admin] failed to validate body: %v",
					err,
				),
			},
		)
	}

	token, err := h.userService.LoginAdmin(
		ctx.Context(),
		userModel,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[login admin] failed to authenticate: %v",
					err,
				),
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (h *UserHandler) RegisterUser(
	ctx *fiber.Ctx,
) error {
	var body model.UserRegisterBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[register user] failed to parse body: %v",
					err,
				),
			},
		)
	}

	userModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[register user] failed to validate body: %v",
					err,
				),
			},
		)
	}

	token, err := h.userService.RegisterUser(
		ctx.Context(),
		userModel,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[register user] failed to store user: %v",
					err,
				),
			},
		)
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(fiber.Map{
			"token": token,
		})
}

func (h *UserHandler) LoginUser(
	ctx *fiber.Ctx,
) error {
	var body model.UserLoginBody
	err := ctx.BodyParser(&body)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[login user] failed to parse body: %v",
					err,
				),
			},
		)
	}

	userModel, err := body.IsValid()
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[login user] failed to validate body: %v",
					err,
				),
			},
		)
	}

	token, err := h.userService.LoginUser(
		ctx.Context(),
		userModel,
	)
	if err != nil {
		return HandleError(
			ctx,
			ErrorResponse{
				error:   err,
				message: err.Error(),
				detail: fmt.Sprintf(
					"[login user] failed to authenticate: %v",
					err,
				),
			},
		)
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}
