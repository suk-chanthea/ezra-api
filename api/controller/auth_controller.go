package controller

import (
    "net/http"
    "github.com/suk-chanthea/ezra/service"

    "github.com/gin-gonic/gin"
)

type AuthController struct {
    authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
    return &AuthController{authService}
}

func (c *AuthController) Register(ctx *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Fullname string `json:"fullname"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    token, err := c.authService.Register(req.Username, req.Fullname, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user registered successfully",
		"token":   token,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    token, err := c.authService.Login(req.Username, req.Password)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token})
}
