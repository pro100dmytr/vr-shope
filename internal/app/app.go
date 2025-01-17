package app

import (
	"fmt"
	"log/slog"
	"os"
	"vr-shope/internal/config"
	"vr-shope/internal/handler/product"
	"vr-shope/internal/handler/purchase"
	"vr-shope/internal/handler/user"
	"vr-shope/internal/middleware"
	"vr-shope/internal/repository"
	"vr-shope/internal/service"
	"vr-shope/internal/storage/postgresql"

	"github.com/gin-gonic/gin"
)

func Run(configPath string) error {
	serverCfg, dbCfg, loggerCfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	var level slog.Level
	switch loggerCfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)

	db, err := postgresql.OpenConnection(dbCfg)
	if err != nil {
		logger.Error("Error creating database connection", slog.Any("error", err))
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	defer db.Close()

	userStorage, err := repository.NewUserStorage(db)
	if err != nil {
		logger.Error("Error creating user storage", slog.Any("error", err))
		return fmt.Errorf("failed to create user storage: %w", err)
	}

	userService := service.NewUserService(userStorage)
	userHandler := user.NewHandler(userService, logger)

	productStorage, err := repository.NewProductStorage(db)
	if err != nil {
		logger.Error("Error creating track storage", slog.Any("error", err))
		return fmt.Errorf("failed to create track storage: %w", err)
	}

	productService := service.NewProductService(productStorage)
	productHandler := product.NewHandler(productService, logger)

	purchaseStorage, err := repository.NewPurchaseStorage(db)
	if err != nil {
		logger.Error("Error creating playlist storage", slog.Any("error", err))
		return fmt.Errorf("failed to create playlist storage: %w", err)
	}

	purchaseService := service.NewPurchaseService(purchaseStorage)
	purchaseHandler := purchase.NewHandler(purchaseService, logger)

	router := gin.Default()

	router.POST("/users/create", userHandler.CreateUser())
	router.POST("/product/create", productHandler.CreateProduct())
	router.POST("/purchase/create", purchaseHandler.CreatePurchase())
	router.POST("/users/login", userHandler.Login())

	Routes := router.Group("/api/v1")
	Routes.Use(middleware.AuthMiddleware())
	{
		Routes.GET("/users", userHandler.GetAllUsers())
		Routes.GET("/users/:id", userHandler.GetUserByID())
		Routes.GET("/users&email=<user_email>", userHandler.GetUserByEmail())
		Routes.GET("/users?offset=1&limit=10", userHandler.GetUserWithPagination())
		Routes.PUT("/users/:id", userHandler.UpdateUser())
		Routes.DELETE("/users/:id", userHandler.DeleteUser())

		Routes.GET("/product", productHandler.GetAllProducts())
		Routes.GET("/product/:id", productHandler.GetProductByID())
		Routes.GET("/product?name=<product_name>", productHandler.GetProductByName())
		Routes.GET("/product?offset=1&limit=10", productHandler.GetProductsWithPagination())
		Routes.PUT("/product/:id", productHandler.UpdateProduct())
		Routes.DELETE("/product/:id", productHandler.DeleteProduct())

		Routes.GET("/playlists", purchaseHandler.GetAllPurchases())
		Routes.GET("/playlists/:id", purchaseHandler.GetPurchaseByID())
		Routes.PUT("/playlists/:id", purchaseHandler.UpdatePurchase())
		Routes.DELETE("/playlists/:id", purchaseHandler.DeletePurchase())
	}

	if err = router.Run(fmt.Sprintf(":%s", serverCfg.Port)); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	return nil
}
