package handler

import (
	"net/http"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OTPHandler struct {
	otpUseCase usecase.OTPUseCase
}

func NewOTPHandler(uc usecase.OTPUseCase) *OTPHandler {
	return &OTPHandler{otpUseCase: uc}
}

// SendOTP generates and sends an OTP to the user's email
// @Summary Send OTP
// @Description Generate and send OTP code to email
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body dto.SendOTPRequest true "Send OTP Request"
// @Success 200 {object} dto.OTPResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /otp/send [post]
func (h *OTPHandler) SendOTP(c *gin.Context) {
	var req dto.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "email":
				message = "Invalid email format"
			case "oneof":
				message = e.Field() + " must be one of: email_verification, password_reset, login"
			default:
				message = "Invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	response, err := h.otpUseCase.SendOTP(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// VerifyOTP verifies the OTP code
// @Summary Verify OTP
// @Description Verify OTP code sent to email
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Verify OTP Request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			e := validationErrors[0]
			var message string
			switch e.Tag() {
			case "required":
				message = e.Field() + " is required"
			case "email":
				message = "Invalid email format"
			case "min", "max":
				message = "OTP code must be 6 digits"
			case "oneof":
				message = e.Field() + " must be one of: email_verification, password_reset, login"
			default:
				message = "Invalid " + e.Field()
			}
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Errors: message})
			return
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	response, err := h.otpUseCase.VerifyOTP(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

