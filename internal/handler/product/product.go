package product

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"vr-shope/internal/models/dto"
	"vr-shope/internal/models/services"
	"vr-shope/internal/service"
)

type Product interface {
	CreateProduct() gin.HandlerFunc
	GetProductByID() gin.HandlerFunc
	GetAllProducts() gin.HandlerFunc
	UpdateProduct() gin.HandlerFunc
	DeleteProduct() gin.HandlerFunc
	AddLike() gin.HandlerFunc
	RemoveLike() gin.HandlerFunc
	GetProductsWithPagination() gin.HandlerFunc
}

type ProductHandler struct {
	service *service.ProductService
	logger  *slog.Logger
}

func NewHandler(service *service.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ProductHandler) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productReq dto.ProductRequest

		if err := c.ShouldBindJSON(&productReq); err != nil {
			h.logger.Error("Error binding JSON", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		productServ := &services.Product{
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

		productResp := dto.ProductResponse{
			Message: "product created",
		}

		h.logger.Info("Product created", slog.Any("productResp", productResp))
		c.JSON(http.StatusCreated, productResp.Message)
	}
}

func (h *ProductHandler) GetProductByID() gin.HandlerFunc {
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

		productResp := dto.ProductResponse{
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

func (h *ProductHandler) GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := h.service.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error fetching products", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
			return
		}

		var productResponses []dto.ProductResponse
		for _, product := range products {
			productResponses = append(productResponses, dto.ProductResponse{
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

func (h *ProductHandler) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Error parsing product ID", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var productReq dto.ProductRequest

		if err := c.ShouldBindJSON(&productReq); err != nil {
			h.logger.Error("Error binding JSON", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		productServ := &services.Product{
			ID:            uint64(id),
			Name:          productReq.Name,
			Cost:          productReq.Cost,
			QuantityStock: productReq.QuantityStock,
			Guarantees:    productReq.Guarantees,
			Country:       productReq.Country,
			Like:          productReq.Like,
		}

		product, err := h.service.Update(c.Request.Context(), productServ)
		if err != nil {
			h.logger.Error("Error updating product", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update product"})
			return
		}

		productResp := dto.ProductResponse{
			Message:       "product updated",
			ID:            product.ID,
			Name:          product.Name,
			Cost:          product.Cost,
			QuantityStock: product.QuantityStock,
			Guarantees:    product.Guarantees,
			Country:       product.Country,
			Like:          productReq.Like,
		}

		h.logger.Info("Product updated", slog.Any("productResp", productResp))
		c.JSON(http.StatusOK, productResp)
	}
}

func (h *ProductHandler) DeleteProduct() gin.HandlerFunc {
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
		productResp := dto.ProductResponse{
			Message: "product deleted",
		}
		c.JSON(http.StatusOK, productResp.Message)
	}
}

func (h *ProductHandler) AddLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = h.service.AddLike(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error adding like", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding like"})
			return
		}

		message := dto.ProductResponse{
			Message: "Added like",
		}

		h.logger.Info("Product added", slog.Any("productResp", message))
		c.JSON(http.StatusOK, message.Message)
	}
}

func (h *ProductHandler) RemoveLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = h.service.RemoveLike(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error removing like", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing like"})
			return
		}

		message := dto.ProductResponse{
			Message: "Delete like",
		}

		h.logger.Info("Product removed", slog.Any("productResp", message))
		c.JSON(http.StatusOK, message.Message)
	}
}

func (h *ProductHandler) GetProductByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		products, err := h.service.GetProductByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error fetching product by name", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching product by name"})
			return
		}

		var productsResponse []dto.ProductResponse
		for _, product := range products {
			productResponse := dto.ProductResponse{
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

func (h *ProductHandler) GetProductsWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		products, err := h.service.GetProductsWithPagination(c.Request.Context(), limit, offset)
		if err != nil {
			h.logger.Error("Failed to retrieve products with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
			return
		}

		var productsResponse []dto.ProductResponse
		for _, product := range products {
			productResponse := dto.ProductResponse{
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
