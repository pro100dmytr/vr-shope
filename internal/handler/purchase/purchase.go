package purchase

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"vr-shope/internal/models/dto"
	"vr-shope/internal/models/services"
	"vr-shope/internal/service"
)

type Purchase interface {
	CreatePurchase() gin.HandlerFunc
	GetPurchaseByID() gin.HandlerFunc
	GetAllPurchases() gin.HandlerFunc
	UpdatePurchase() gin.HandlerFunc
	DeletePurchase() gin.HandlerFunc
}

type PurchaseHandler struct {
	service *service.PurchaseService
	logger  *slog.Logger
}

func NewHandler(service *service.PurchaseService, logger *slog.Logger) *PurchaseHandler {
	return &PurchaseHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PurchaseHandler) CreatePurchase() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.PurchaseRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			h.logger.Error("failed to bind request", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		purchase := services.Purchase{
			ID:     0,
			UserID: uint64(request.UserID),
			Cost:   request.Cost,
			Date:   time.Now(),
		}

		err := h.service.Create(c.Request.Context(), &purchase)
		if err != nil {
			h.logger.Error("failed to create purchase", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create purchase"})
			return
		}

		response := dto.PurchaseResponse{
			Message: "purchase created",
		}

		h.logger.Info("purchase created", slog.Any("purchase", response))
		c.JSON(http.StatusCreated, response.Message)
	}
}

func (h *PurchaseHandler) GetPurchaseByID() gin.HandlerFunc {
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

		response := dto.PurchaseResponse{
			Message: "purchase found",
			ID:      purchase.ID,
			UserID:  purchase.UserID,
			Cost:    purchase.Cost,
			Date:    purchase.Date,
		}
		h.logger.Info("purchase found", slog.Any("purchase", response))
		c.JSON(http.StatusOK, response)
	}
}

func (h *PurchaseHandler) GetAllPurchases() gin.HandlerFunc {
	return func(c *gin.Context) {
		purchases, err := h.service.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("failed to get purchases", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get purchases"})
			return
		}

		responses := make([]dto.PurchaseResponse, 0, len(purchases))
		for _, purchase := range purchases {
			responses = append(responses, dto.PurchaseResponse{
				Message: "get purchase",
				ID:      purchase.ID,
				UserID:  purchase.UserID,
				Cost:    purchase.Cost,
				Date:    purchase.Date,
			})
		}
		h.logger.Info("get purchases", slog.Any("purchases", responses))
		c.JSON(http.StatusOK, responses)
	}
}

func (h *PurchaseHandler) UpdatePurchase() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("invalid id format", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		var request dto.PurchaseRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			h.logger.Error("failed to bind request", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		updateData := services.Purchase{
			ID:     uint64(id),
			UserID: uint64(request.UserID),
			Cost:   request.Cost,
			Date:   time.Now(),
		}

		err = h.service.Update(c.Request.Context(), &updateData)
		if err != nil {
			h.logger.Error("failed to update purchase", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update purchase"})
			return
		}

		response := dto.PurchaseResponse{
			Message: "purchase updated",
		}
		h.logger.Info("purchase updated", slog.Any("purchase", response))
		c.JSON(http.StatusOK, response.Message)
	}
}

func (h *PurchaseHandler) DeletePurchase() gin.HandlerFunc {
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

		purchaseResp := dto.PurchaseResponse{
			Message: "purchase deleted",
		}
		h.logger.Info("purchase deleted", slog.Any("purchase", purchaseResp))
		c.JSON(http.StatusNoContent, purchaseResp.Message)
	}
}
