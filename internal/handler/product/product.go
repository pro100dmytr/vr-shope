package product

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"vr-shope/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ctx context.Context, product *models.Product) error
	Get(ctx context.Context, id int) (*models.Product, error)
	GetAll(ctx context.Context) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	GetProductByName(ctx context.Context, name string) ([]*models.Product, error)
	GetProductsWithPagination(ctx context.Context, limit, offset string) ([]*models.Product, error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productReq models.ProductRequest

		if err := c.ShouldBindJSON(&productReq); err != nil {
			h.logger.Error("Error binding JSON", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		productServ := &models.Product{
			Name:          productReq.Name,
			Cost:          productReq.Cost,
			QuantityStock: productReq.QuantityStock,
			Guarantees:    productReq.Guarantees,
			Country:       productReq.Country,
			Like:          productReq.Like,
		}

		err := h.service.Create(c.Request.Context(), productServ)
		if err != nil {
			h.logger.Error("Error creating product", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		productResp := models.ProductResponse{
			Message: "product created",
		}

		h.logger.Info("Product created", slog.Any("productResp", productResp))
		c.JSON(http.StatusCreated, productResp.Message)
	}
}

func (h *Handler) GetProductByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Error parsing product ID", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		product, err := h.service.Get(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error fetching product", slog.Any("err", err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}

		productResp := models.ProductResponse{
			Message:       "product found",
			ID:            product.ID,
			Name:          product.Name,
			Cost:          product.Cost,
			QuantityStock: product.QuantityStock,
			Guarantees:    product.Guarantees,
			Country:       product.Country,
			Like:          product.Like,
		}

		h.logger.Info("Product found", slog.Any("productResp", productResp))
		c.JSON(http.StatusOK, productResp)
	}
}

func (h *Handler) GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := h.service.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error fetching products", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
			return
		}

		var productResponses []models.ProductResponse
		for _, product := range products {
			productResponses = append(productResponses, models.ProductResponse{
				Message:       "get product",
				ID:            product.ID,
				Name:          product.Name,
				Cost:          product.Cost,
				QuantityStock: product.QuantityStock,
				Guarantees:    product.Guarantees,
				Country:       product.Country,
				Like:          product.Like,
			})
		}

		h.logger.Info("Products found", slog.Any("productResponses", productResponses))
		c.JSON(http.StatusOK, productResponses)
	}
}

func (h *Handler) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Error parsing product ID", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var productReq models.ProductRequest

		if err := c.ShouldBindJSON(&productReq); err != nil {
			h.logger.Error("Error binding JSON", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		productServ := &models.Product{
			ID:            uint64(id),
			Name:          productReq.Name,
			Cost:          productReq.Cost,
			QuantityStock: productReq.QuantityStock,
			Guarantees:    productReq.Guarantees,
			Country:       productReq.Country,
			Like:          productReq.Like,
		}

		err = h.service.Update(c.Request.Context(), productServ)
		if err != nil {
			h.logger.Error("Error updating product", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update product"})
			return
		}

		h.logger.Info("Product updated")
		response := models.ProductResponse{
			Message: "product updated",
		}
		c.JSON(http.StatusOK, response)
	}
}

func (h *Handler) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Error parsing product ID", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		err = h.service.Delete(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error deleting product", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete product"})
			return
		}

		h.logger.Info("Product deleted", slog.Any("productResp", id))
		productResp := models.ProductResponse{
			Message: "product deleted",
		}
		c.JSON(http.StatusOK, productResp.Message)
	}
}

func (h *Handler) GetProductByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		products, err := h.service.GetProductByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error fetching product by name", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching product by name"})
			return
		}

		var productsResponse []models.ProductResponse
		for _, product := range products {
			productResponse := models.ProductResponse{
				Message:       "product by name",
				ID:            product.ID,
				Name:          product.Name,
				Cost:          product.Cost,
				QuantityStock: product.QuantityStock,
				Guarantees:    product.Guarantees,
				Country:       product.Country,
				Like:          product.Like,
			}

			productsResponse = append(productsResponse, productResponse)
		}

		h.logger.Info("Products retrieved", slog.Any("productsResponse", productsResponse))
		c.JSON(http.StatusOK, productsResponse)
	}
}

func (h *Handler) GetProductsWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		products, err := h.service.GetProductsWithPagination(c.Request.Context(), limit, offset)
		if err != nil {
			h.logger.Error("Failed to retrieve products with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
			return
		}

		var productsResponse []models.ProductResponse
		for _, product := range products {
			productResponse := models.ProductResponse{
				Message:       "product by name",
				ID:            product.ID,
				Name:          product.Name,
				Cost:          product.Cost,
				QuantityStock: product.QuantityStock,
				Guarantees:    product.Guarantees,
				Country:       product.Country,
				Like:          product.Like,
			}
			productsResponse = append(productsResponse, productResponse)
		}

		h.logger.Info("Products retrieved", slog.Any("productsResponse", productsResponse))
		c.JSON(http.StatusOK, productsResponse)
	}
}
