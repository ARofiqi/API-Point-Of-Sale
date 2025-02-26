package middlewares

import (
	"net/http"

	"aro-shop/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RoleMiddleware(allowedRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userToken := c.Get("user")
			if userToken == nil {
				return utils.Response(c, http.StatusUnauthorized, "Unauthorized", nil, nil, nil)
			}

			token, ok := userToken.(*jwt.Token)
			if !ok || !token.Valid {
				return utils.Response(c, http.StatusUnauthorized, "Invalid token", nil, nil, nil)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return utils.Response(c, http.StatusUnauthorized, "Invalid claims", nil, nil, nil)
			}

			role, ok := claims["role"].(string)
			if !ok {
				return utils.Response(c, http.StatusUnauthorized, "Role not found in token", nil, nil, nil)
			}

			if role != allowedRole {
				return utils.Response(c, http.StatusForbidden, "Kamu bukan "+allowedRole, nil, nil, nil)
			}

			return next(c)
		}
	}
}
