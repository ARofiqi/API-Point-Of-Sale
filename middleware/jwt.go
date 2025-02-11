package middleware

import (
	"aro-shop/config"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var cfg = config.LoadConfig()
var jwtSecret = []byte(cfg.JWTSecret)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Token tidak ditemukan",
			})
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		} else {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Format token tidak valid",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": "Gagal membaca klaim token",
			})
		}

		c.Set("user_id", claims["user_id"])

		return next(c)
	}
}
