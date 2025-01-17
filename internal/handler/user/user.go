package user

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"vr-shope/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateUser(ctx context.Context, user *models.User) error
	Get(ctx context.Context, id int) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetUsersWithPagination(ctx context.Context, limit, offset string) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	GetToken(ctx context.Context, login string, password string) (string, error)
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

func (h *Handler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userReq models.UserRequest

		if err := c.ShouldBindJSON(&userReq); err != nil {
			h.logger.Error("Error binding JSON", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		userServ := &models.User{
			Login:       userReq.Login,
			Name:        userReq.Name,
			LastName:    userReq.LastName,
			PhoneNumber: userReq.PhoneNumber,
			Password:    userReq.Password,
			Email:       userReq.Email,
		}

		err := h.service.CreateUser(c.Request.Context(), userServ)
		if err != nil {
			h.logger.Error("Error creating userReq", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		userResp := models.UserResponse{
			Message: "User created",
		}

		h.logger.Info("User created", slog.Any("userResp", userResp))
		c.JSON(http.StatusCreated, userResp.Message)
	}
}

func (h *Handler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Error converting id to int", slog.Any("id", id))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		user, err := h.service.Get(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error fetching user", slog.Any("err", err))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		userResp := models.UserResponse{
			Message:         "get user",
			ID:              user.ID,
			Login:           user.Login,
			Name:            user.Name,
			LastName:        user.LastName,
			PhoneNumber:     user.PhoneNumber,
			Email:           user.Email,
			WalletUSDT:      user.WalletUSDT,
			NumberPurchases: user.NumberPurchases,
		}

		h.logger.Info("User found", slog.Any("userResp", userResp))
		c.JSON(http.StatusOK, userResp)
	}
}

func (h *Handler) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.service.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error fetching users", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
			return
		}

		var userResponses []models.UserResponse
		for _, user := range users {
			userResponses = append(userResponses, models.UserResponse{
				Message:         "get user",
				ID:              user.ID,
				Login:           user.Login,
				Name:            user.Name,
				LastName:        user.LastName,
				PhoneNumber:     user.PhoneNumber,
				Email:           user.Email,
				WalletUSDT:      user.WalletUSDT,
				NumberPurchases: user.NumberPurchases,
			})
		}

		h.logger.Info("Users found", slog.Any("userResponses", userResponses))
		c.JSON(http.StatusOK, userResponses)
	}
}

func (h *Handler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Error parsing user id", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing user id"})
			return
		}

		var userReq models.UserRequest

		if err := c.ShouldBindJSON(&userReq); err != nil {
			h.logger.Error("Error binding JSON", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		userServ := &models.User{
			ID:          uint64(id),
			Login:       userReq.Login,
			Name:        userReq.Name,
			LastName:    userReq.LastName,
			PhoneNumber: userReq.PhoneNumber,
			Password:    userReq.Password,
			Email:       userReq.Email,
			WalletUSDT:  userReq.WalletUSDT,
		}

		err = h.service.Update(c.Request.Context(), userServ)
		if err != nil {
			h.logger.Error("Error updating user", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
			return
		}

		userResp := models.UserResponse{
			Message: "User updated",
		}

		h.logger.Info("User updated", slog.Any("userResp", userResp))
		c.JSON(http.StatusOK, userResp.Message)
	}
}

func (h *Handler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Error parsing user id", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing user id"})
			return
		}

		err = h.service.Delete(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error deleting user", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
			return
		}

		h.logger.Info("User deleted", slog.Any("userResp", id))

		userResp := models.UserResponse{
			Message: "User deleted",
		}
		c.JSON(http.StatusOK, userResp.Message)
	}
}

func (h *Handler) GetUserByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")

		user, err := h.service.GetByEmail(c.Request.Context(), email)
		if err != nil {
			h.logger.Error("Error fetching user", slog.Any("err", err))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		userResp := models.UserResponse{
			Message:         "get user",
			ID:              user.ID,
			Login:           user.Login,
			Name:            user.Name,
			LastName:        user.LastName,
			PhoneNumber:     user.PhoneNumber,
			Email:           user.Email,
			WalletUSDT:      user.WalletUSDT,
			NumberPurchases: user.NumberPurchases,
		}

		h.logger.Info("User found", slog.Any("userResp", userResp))
		c.JSON(http.StatusOK, userResp)
	}
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.UserRequest

		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		token, err := h.service.GetToken(c.Request.Context(), user.Login, user.Password)
		if err != nil {
			h.logger.Error("User not found", slog.Any("error", err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		tokenResponse := models.TokenResponse{
			Token: token,
		}

		h.logger.Info("Token found", slog.Any("tokenResponse", tokenResponse))
		c.JSON(http.StatusOK, tokenResponse)
	}
}

func (h *Handler) GetUserWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		users, err := h.service.GetUsersWithPagination(c.Request.Context(), limit, offset)
		if err != nil {
			h.logger.Error("Failed to retrieve users with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		var usersResponse []models.UserResponse
		for _, user := range users {
			userResponse := models.UserResponse{
				Message:         "get user",
				ID:              user.ID,
				Login:           user.Login,
				Name:            user.Name,
				LastName:        user.LastName,
				PhoneNumber:     user.PhoneNumber,
				Email:           user.Email,
				WalletUSDT:      user.WalletUSDT,
				NumberPurchases: user.NumberPurchases,
			}
			usersResponse = append(usersResponse, userResponse)
		}

		h.logger.Info("Users found", slog.Any("usersResponse", usersResponse))
		c.JSON(http.StatusOK, usersResponse)
	}
}
