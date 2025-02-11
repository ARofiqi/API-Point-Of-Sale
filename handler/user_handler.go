package handler

import (
	"aro-shop/config"
	"aro-shop/db"
	"aro-shop/models"
	"aro-shop/utils"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	cfg       = config.LoadConfig()
	jwtSecret = []byte(cfg.JWTSecret)
)

func parseValidationErrors(err error) map[string]string {
	errorDetails := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			switch fieldErr.Tag() {
			case "required":
				errorDetails[fieldErr.Field()] = "Tidak boleh kosong"
			case "email":
				errorDetails[fieldErr.Field()] = "Format email tidak valid"
			case "min":
				errorDetails[fieldErr.Field()] = "Minimal " + fieldErr.Param() + " karakter"
			}
		}
	}
	return errorDetails
}

func Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := parseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation error", nil, err, errorDetails)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Error hashing password", nil, err, nil)
	}

	result, err := db.DB.Exec("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)", req.Name, req.Email, string(hashedPassword), "user")
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to register user", nil, err, nil)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to retrieve user ID", nil, err, nil)
	}

	data := map[string]interface{}{
		"id":    userID,
		"name":  req.Name,
		"email": req.Email,
	}

	return utils.Response(c, http.StatusCreated, "User registered successfully", data, nil, nil)
}

func Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := parseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation error", nil, err, errorDetails)
	}

	var user models.User
	email := strings.TrimSpace(strings.ToLower(req.Email))

	row := db.DB.QueryRow("SELECT id, name, email, password, role FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
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
		Role string `json:"role" validate:"required,oneof=user admin"`
	}

	if err := c.Bind(&req); err != nil {
		return utils.Response(c, http.StatusBadRequest, "Invalid request", nil, err, nil)
	}

	if err := validate.Struct(req); err != nil {
		errorDetails := parseValidationErrors(err)
		return utils.Response(c, http.StatusBadRequest, "Validation error", nil, err, errorDetails)
	}

	_, err := db.DB.Exec("UPDATE users SET role = ? WHERE id = ?", req.Role, userID)
	if err != nil {
		return utils.Response(c, http.StatusInternalServerError, "Failed to update role", nil, err, nil)
	}

	return utils.Response(c, http.StatusOK, "User role updated successfully", nil, nil, nil)
}
