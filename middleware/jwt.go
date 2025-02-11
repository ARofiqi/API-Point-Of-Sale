package middleware

import (
	"aro-shop/config"
	"aro-shop/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	cfg          = config.LoadConfig()
	jwtSecret    = []byte(cfg.JWTSecret)
	errorDetails = make(map[string]string)
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			errorDetails["jwt"] = "Token tidak dapat diparsing atau sudah expired"
			return utils.Response(c, http.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa", nil, err, errorDetails)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			errorDetails["jwt"] = "Klaim token tidak valid"
			return utils.Response(c, http.StatusUnauthorized, "Gagal membaca klaim token", nil, nil, errorDetails)
		}

		c.Set("user_id", claims["user_id"])

		return next(c)
	}
}
