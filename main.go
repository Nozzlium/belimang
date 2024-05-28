package main

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/caarlos0/env/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/nozzlium/belimang/internal/client"
	"github.com/nozzlium/belimang/internal/config"
	"github.com/nozzlium/belimang/internal/handler"
	"github.com/nozzlium/belimang/internal/middleware"
	"github.com/nozzlium/belimang/internal/repository"
	"github.com/nozzlium/belimang/internal/service"
)

func main() {
	fiberApp := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		Prefork:     true,
	})

	err := setupApp(fiberApp)
	if err != nil {
		log.Fatal(err)
	}

	err = fiberApp.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func setupApp(app *fiber.App) error {
	var cfg config.Config
	opts := env.Options{
		TagName: "json",
	}
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		log.Fatalf("%+v\n", err)
		return err
	}

	db, err := client.InitDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
		return err
	}

	userRepository := repository.NewUserRepository(
		db,
	)
	merchantRepository := repository.NewMerchantRepository(
		db,
	)
	productRepository := repository.NewProductRepository(
		db,
	)

	userService := service.NewUserService(
		userRepository,
		cfg.JWTSecret,
		int(cfg.BCryptSalt),
	)
	merchantService := service.NewMerchantService(
		merchantRepository,
	)
	productService := service.NewProductService(
		productRepository,
	)

	userHandler := handler.NewUserHandler(
		userService,
	)
	merchantHandler := handler.NewMerchantHandler(
		merchantService,
	)
	productHandler := handler.NewProductHandler(
		productService,
	)

	admin := app.Group("/admin")
	admin.Post(
		"/register",
		userHandler.RegisterAdmin,
	)
	admin.Post(
		"/login",
		userHandler.LoginAdmin,
	)
	adminProtected := admin.Use(
		middleware.Protected(),
	).Use(middleware.SetClaimsData())
	adminProtected.Post(
		"/merchants",
		merchantHandler.Create,
	)
	adminProtected.Get(
		"/merchants",
		merchantHandler.FindAll,
	)
	adminProtected.Post(
		"/merchants/:merchantId/items",
		productHandler.Create,
	)

	user := app.Group("/user")
	user.Post(
		"/register",
		userHandler.RegisterUser,
	)
	user.Post(
		"/login",
		userHandler.LoginUser,
	)

	return nil
}
