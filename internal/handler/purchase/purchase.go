package purchase

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"vr-shope/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ctx context.Context, purchase *models.Purchase) error
	Get(ctx context.Context, id int64) (*models.Purchase, error)
	GetAll(ctx context.Context) ([]*models.Purchase, error)
	Update(ctx context.Context, purchase *models.Purchase) error
	Delete(ctx context.Context, id int64) error
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

func (h *Handler) CreatePurchase() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.PurchaseRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			h.logger.Error("failed to bind request", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		purchase := models.Purchase{
			UserID:    uint64(request.UserID),
			ProductID: uint64(request.ProductID),
			Date:      time.Now(),
		}

		err := h.service.Create(c.Request.Context(), &purchase)
		if err != nil {
			h.logger.Error("failed to create purchase", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create purchase"})
			return
		}

		response := models.PurchaseResponse{
			Message: "purchase created",
		}

		h.logger.Info("purchase created", slog.Any("purchase", response))
		c.JSON(http.StatusCreated, response.Message)
	}
}

func (h *Handler) GetPurchaseByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("invalid id format", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		purchase, err := h.service.Get(c.Request.Context(), int64(id))
		if err != nil {
			h.logger.Error("failed to get purchase", "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
			return
		}

		response := models.PurchaseResponse{
			Message:    "purchase found",
			ID:         purchase.ID,
			UserID:     purchase.UserID,
			ProductID:  purchase.ProductID,
			Date:       purchase.Date,
			WalletUSDT: purchase.WalletUSDT,
			Cost:       purchase.Cost,
		}

		h.logger.Info("purchase found", slog.Any("purchase", response))
		c.JSON(http.StatusOK, response)
	}
}

func (h *Handler) GetAllPurchases() gin.HandlerFunc {
	return func(c *gin.Context) {
		purchases, err := h.service.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("failed to get purchases", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get purchases"})
			return
		}

		responses := make([]models.PurchaseResponse, 0, len(purchases))
		for _, purchase := range purchases {
			responses = append(responses, models.PurchaseResponse{
				Message:    "get purchase",
				ID:         purchase.ID,
				UserID:     purchase.UserID,
				ProductID:  purchase.ProductID,
				Date:       purchase.Date,
				WalletUSDT: purchase.WalletUSDT,
				Cost:       purchase.Cost,
			})
		}

		h.logger.Info("get purchases", slog.Any("purchases", responses))
		c.JSON(http.StatusOK, responses)
	}
}

func (h *Handler) UpdatePurchase() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("invalid id format", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		var request models.PurchaseRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			h.logger.Error("failed to bind request", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		updateData := models.Purchase{
			ID:        uint64(id),
			UserID:    uint64(request.UserID),
			ProductID: uint64(request.ProductID),
			Date:      time.Now(),
		}

		err = h.service.Update(c.Request.Context(), &updateData)
		if err != nil {
			h.logger.Error("failed to update purchase", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update purchase"})
			return
		}

		response := models.PurchaseResponse{
			Message: "purchase updated",
		}
		h.logger.Info("purchase updated", slog.Any("purchase", response))
		c.JSON(http.StatusOK, response.Message)
	}
}

func (h *Handler) DeletePurchase() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			h.logger.Error("invalid id format", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		if err := h.service.Delete(c.Request.Context(), int64(id)); err != nil {
			h.logger.Error("failed to delete purchase", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete purchase"})
			return
		}

		purchaseResp := models.PurchaseResponse{
			Message: "purchase deleted",
		}
		h.logger.Info("purchase deleted", slog.Any("purchase", purchaseResp))
		c.JSON(http.StatusNoContent, purchaseResp.Message)
	}
}
