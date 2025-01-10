package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	config2 "vr-shope/internal/config"
	"vr-shope/internal/handler/product"
	"vr-shope/internal/handler/purchase"
	"vr-shope/internal/handler/user"
	"vr-shope/internal/middleware"
	"vr-shope/internal/repository"
	"vr-shope/internal/service"
)

func Run(config string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg, err := config2.LoadConfig(config)
	if err != nil {
		logger.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}

	userStorage, err := repository.NewUserStorage(cfg)
	if err != nil {
		logger.Error("Error creating user storage", slog.Any("error", err))
		os.Exit(1)
	}

	defer userStorage.Close()

	userService := service.NewUserService(userStorage)
	userHandler := user.NewHandler(userService, logger)

	purchaseStorage, err := repository.NewPurchaseStorage(cfg)
	if err != nil {
		logger.Error("Error creating purchase storage", slog.Any("error", err))
		os.Exit(1)
	}

	defer purchaseStorage.Close()

	purchaseService := service.NewPurchaseService(purchaseStorage)
	purchaseHandler := purchase.NewHandler(purchaseService, logger)

	productStorage, err := repository.NewProductStorage(cfg)
	if err != nil {
		logger.Error("Error creating product storage", slog.Any("error", err))
		os.Exit(1)
	}

	defer productStorage.Close()

	productService := service.NewProductService(productStorage)
	productHandler := product.NewHandler(productService, logger)

	router := gin.Default()

	router.POST("/users/login", userHandler.Login())

	userRoutes := router.Group("/api")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/users", userHandler.GetAllUsers())
		userRoutes.GET("/users/:id", userHandler.GetUserByID())
		userRoutes.GET("/users?email=<email>", userHandler.GetUserByEmail())
		userRoutes.GET("/users?offset=1&limit=10", userHandler.GetUserWithPagination())
		userRoutes.POST("/users", userHandler.CreateUser())
		userRoutes.PUT("/users/:id", userHandler.UpdateUser())
		userRoutes.DELETE("/users/:id", userHandler.DeleteUser())
	}

	purchaseRoutes := router.Group("/api")
	purchaseRoutes.Use(middleware.AuthMiddleware())
	{
		purchaseRoutes.GET("/purchase", purchaseHandler.GetAllPurchases())
		purchaseRoutes.GET("/purchase/:id", purchaseHandler.GetPurchaseByID())
		purchaseRoutes.POST("/purchase", purchaseHandler.CreatePurchase())
		purchaseRoutes.PUT("/purchase/:id", purchaseHandler.UpdatePurchase())
		purchaseRoutes.DELETE("/purchase/:id", purchaseHandler.DeletePurchase())
	}

	productRoutes := router.Group("/api")
	productRoutes.Use(middleware.AuthMiddleware())
	{
		productRoutes.GET("/product", productHandler.GetAllProducts())
		productRoutes.GET("/product/:id", productHandler.GetProductByID())
		productRoutes.GET("/users?offset=1&limit=10", productHandler.GetProductsWithPagination())
		productRoutes.POST("/product", productHandler.CreateProduct())
		productRoutes.PUT("/product/:id", productHandler.UpdateProduct())
		productRoutes.DELETE("/product/:id", productHandler.DeleteProduct())
		productRoutes.PATCH("/tracks/:id/like", productHandler.AddLike())
		productRoutes.DELETE("/tracks/:id/like", productHandler.RemoveLike())
	}

	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := router.Run(serverAddr); err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
