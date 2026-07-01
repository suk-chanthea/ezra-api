package handler

import (
	"context"
	"net/http"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/api/idtoken"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: uc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "min":
				message = e.Field() + " is too short"
			case "max":
				message = e.Field() + " is too long"
			case "email":
				message = e.Field() + " must be a valid email"
			default:
				message = "invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	response, err := h.authUseCase.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "email":
				message = "email must be a valid email"
			default:
				message = "invalid request"
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
		return
	}

	if req.Identifier() == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "username or email is required"})
		return
	}

	response, err := h.authUseCase.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user_id from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.authUseCase.Logout(userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "logged out successfully"})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	// Get user_id from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	if err := h.authUseCase.DeleteUser(userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Message: "user deleted successfully"})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req dto.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	// Verify ID token with Google
	payload, err := idtoken.Validate(context.Background(), req.IDToken, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid Google ID token"})
		return
	}

	// Extract user info from token payload
	googleID, ok := payload.Claims["sub"].(string)
	if !ok || googleID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid token: missing subject"})
		return
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid token: missing email"})
		return
	}

	fullname, _ := payload.Claims["name"].(string)
	if fullname == "" {
		fullname = email // Fallback to email if name not provided
	}

	profilePicture, _ := payload.Claims["picture"].(string)

	// Call use case to login/register user
	response, err := h.authUseCase.GoogleLogin(googleID, email, fullname, profilePicture)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "email":
				message = "Invalid email format"
			case "min":
				if e.Field() == "NewPassword" {
					message = "Password must be at least 6 characters"
				} else {
					message = "OTP code must be 6 digits"
				}
			case "max":
				message = "OTP code must be 6 digits"
			default:
				message = "Invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	response, err := h.authUseCase.ResetPassword(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	user, err := h.authUseCase.GetMe(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authenticated"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "min":
				message = e.Field() + " is too short"
			case "max":
				message = e.Field() + " is too long"
			default:
				message = "invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	user, err := h.authUseCase.UpdateMe(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "profile updated successfully",
		Data:    user,
	})
}