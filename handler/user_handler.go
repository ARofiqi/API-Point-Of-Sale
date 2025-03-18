package handler

import (
	"aro-shop/config"
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	cfg       = config.LoadConfig()
	jwtSecret = []byte(cfg.JWTSecret)
)


func Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation error", nil, err, errorDetails)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Error hashing password", nil, err, nil)
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     models.RoleUser,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to register user", nil, err, nil)
	}

	data := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	return utils.Response(c, http.StatusCreated, "User registered successfully", data, nil, nil)
}

func Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation error", nil, err, errorDetails)
	}

	var user models.User
	email := strings.TrimSpace(strings.ToLower(req.Email))

	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return utils.Response(c, http.StatusUnauthorized, "Invalid email or password", nil, err, nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return utils.Response(c, http.StatusUnauthorized, "Invalid email or password", nil, err, nil)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to generate token", nil, err, nil)
	}

	data := map[string]string{"token": tokenString}
	return utils.Response(c, http.StatusOK, "Login successful", data, nil, nil)
}

func SetUserRole(c echo.Context) error {
	userID := c.Param("id")
	var req struct {
		Role models.Role `json:"role" validate:"required,oneof=user admin"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation error", nil, err, errorDetails)
	}

	if err := db.DB.Model(&models.User{}).Where("id = ?", userID).Update("role", req.Role).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to update role", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "User role updated successfully", nil, nil, nil)
}
