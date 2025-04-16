package middlewares

import (
	"aro-shop/config"
	"aro-shop/utils"
	"aro-shop/dto"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	cfg       = config.LoadConfig()
	jwtSecret = []byte(cfg.JWTSecret)
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		errorDetails := make(dto.ErrorDetails)

		tokenString := c.Request().Header.Get("Authorization")

		if tokenString == "" {
			errorDetails["authorization"] = "Token tidak ada dalam header"
			return utils.Response(c, http.StatusUnauthorized, "Token tidak ditemukan", nil, nil, errorDetails)
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		} else {
			errorDetails["authorization"] = "Format token harus menggunakan 'Bearer <token>'"
			return utils.Response(c, http.StatusUnauthorized, "Format token tidak valid", nil, nil, errorDetails)
		}

		if tokenString == "" {
			errorDetails["jwt"] = "Token kosong setelah parsing 'Bearer '"
			return utils.Response(c, http.StatusUnauthorized, "Token tidak valid", nil, nil, errorDetails)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			log.Println("Error parsing token:", err)
			errorDetails["jwt"] = "Gagal parsing token"
			return utils.Response(c, http.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa", nil, err, errorDetails)
		}

		if token == nil || !token.Valid {
			errorDetails["jwt"] = "Token tidak valid"
			return utils.Response(c, http.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa", nil, nil, errorDetails)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			errorDetails["jwt"] = "Klaim token tidak valid"
			return utils.Response(c, http.StatusUnauthorized, "Gagal membaca klaim token", nil, nil, errorDetails)
		}

		userID, userIDExists := claims["user_id"]
		if !userIDExists || userID == nil {
			errorDetails["jwt"] = "User ID tidak ditemukan dalam klaim token"
			return utils.Response(c, http.StatusUnauthorized, "User ID tidak ditemukan dalam token", nil, nil, errorDetails)
		}

		// Set data user ke context
		c.Set("user", token)
		c.Set("user_id", userID)

		return next(c)
	}
}
