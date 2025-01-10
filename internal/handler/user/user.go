package user

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"vr-shope/internal/models/dto"
	"vr-shope/internal/models/services"
	"vr-shope/internal/service"
)

type User interface {
	CreateUser() gin.HandlerFunc
	GetUserByID() gin.HandlerFunc
	GetAllUsers() gin.HandlerFunc
	GetUserWithPagination() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	GetByEmail() gin.HandlerFunc
}

type UserHandler struct {
	service *service.UserService
	logger  *slog.Logger
}

func NewHandler(service *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userReq dto.UserRequest

		if err := c.ShouldBindJSON(&userReq); err != nil {
			h.logger.Error("Error binding JSON", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		userServ := &services.User{
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

		userResp := dto.UserResponse{
			Message: "User created",
		}

		h.logger.Info("User created", slog.Any("userResp", userResp))
		c.JSON(http.StatusCreated, userResp.Message)
	}
}

func (h *UserHandler) GetUserByID() gin.HandlerFunc {
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

		userResp := dto.UserResponse{
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

func (h *UserHandler) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.service.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error fetching users", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
			return
		}

		var userResponses []dto.UserResponse
		for _, user := range users {
			userResponses = append(userResponses, dto.UserResponse{
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

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Error parsing user id", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing user id"})
			return
		}

		var userReq dto.UserRequest

		if err := c.ShouldBindJSON(&userReq); err != nil {
			h.logger.Error("Error binding JSON", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		userServ := &services.User{
			ID:          id,
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

		userResp := dto.UserResponse{
			Message: "User updated",
		}

		h.logger.Info("User updated", slog.Any("userResp", userResp))
		c.JSON(http.StatusOK, userResp.Message)
	}
}

func (h *UserHandler) DeleteUser() gin.HandlerFunc {
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

		userResp := dto.UserResponse{
			Message: "User deleted",
		}
		c.JSON(http.StatusOK, userResp.Message)
	}
}

func (h *UserHandler) GetUserByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")

		user, err := h.service.GetByEmail(c.Request.Context(), email)
		if err != nil {
			h.logger.Error("Error fetching user", slog.Any("err", err))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		userResp := dto.UserResponse{
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

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user dto.UserRequest

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

		tokenResponse := dto.TokenResponse{
			Token: token,
		}

		h.logger.Info("Token found", slog.Any("tokenResponse", tokenResponse))
		c.JSON(http.StatusOK, tokenResponse)
	}
}

func (h *UserHandler) GetUserWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		users, err := h.service.GetUsersWithPagination(c.Request.Context(), limit, offset)
		if err != nil {
			h.logger.Error("Failed to retrieve users with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		var usersResponse []dto.UserResponse
		for _, user := range users {
			userResponse := dto.UserResponse{
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
