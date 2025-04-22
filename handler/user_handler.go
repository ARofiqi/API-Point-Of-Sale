package handler

import (
	"aro-shop/config"
	"aro-shop/db"
	"aro-shop/dto"
	"aro-shop/models"
	"aro-shop/utils"
	"errors"
	"net/http"
	"strings"
	
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	cfg       = config.LoadConfig()
	jwtSecret = []byte(cfg.JWTSecret)
)

func Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation error", nil, err, errorDetails)
	}

	if req.Password != req.ConfirmPassword {
		return utils.Response(c, http.StatusBadRequest, "Password and Confirm Password do not match", nil, nil, nil)
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return utils.Response(c, http.StatusBadRequest, "Email is already registered", nil, nil, nil)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.Response(c, http.StatusInternalServerError, "Failed to check existing email", nil, err, nil)
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

	data := dto.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return utils.Response(c, http.StatusCreated, "User registered successfully", data, nil, nil)
}

func RegisterAdmin(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := utils.ParseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation error", nil, err, errorDetails)
	}

	if req.Password != req.ConfirmPassword {
		return utils.Response(c, http.StatusBadRequest, "Password and Confirm Password do not match", nil, nil, nil)
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return utils.Response(c, http.StatusBadRequest, "Email is already registered", nil, nil, nil)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.Response(c, http.StatusInternalServerError, "Failed to check existing email", nil, err, nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Error hashing password", nil, err, nil)
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     models.RoleAdmin,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to register user", nil, err, nil)
	}

	data := dto.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return utils.Response(c, http.StatusCreated, "User registered successfully", data, nil, nil)
}

func Login(c echo.Context) error {
	var req dto.LoginRequest
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
		// "exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to generate token", nil, err, nil)
	}

	data := map[string]string{"token": tokenString}
	return utils.Response(c, http.StatusOK, "Login successful", data, nil, nil)
}
