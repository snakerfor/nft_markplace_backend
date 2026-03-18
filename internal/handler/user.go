package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nft-marketplace/internal/model"
	"nft-marketplace/internal/service"
	"nft-marketplace/pkg/errors"
	"nft-marketplace/pkg/jwt"
	"nft-marketplace/pkg/response"
)

type UserHandler struct {
	userSvc *service.UserService
	jwtSecret []byte
}

func NewUserHandler(userSvc *service.UserService, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		userSvc:   userSvc,
		jwtSecret: jwtSecret,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, parseValidationErrors(err))
		return
	}

	user, err := h.userSvc.CreateUser(req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response.Success(c, toUserResponse(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, parseValidationErrors(err))
		return
	}

	user, err := h.userSvc.Authenticate(req.Username, req.Password)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	token, err := jwt.GenerateToken(h.jwtSecret, user.ID, user.Username)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  toUserResponse(user),
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.userSvc.GetUserByID(userID.(uint))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response.Success(c, toUserResponse(user))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, parseValidationErrors(err))
		return
	}

	user, err := h.userSvc.UpdateUser(userID.(uint), req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response.Success(c, toUserResponse(user))
}

func toUserResponse(user *model.User) model.UserResponse {
	return model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func parseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	errors["general"] = err.Error()
	return errors
}
