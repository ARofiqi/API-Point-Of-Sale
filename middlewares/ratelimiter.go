package middlewares

import (
	"aro-shop/utils"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

var (
	limiters = sync.Map{}
)

func getLimiter(ip string) *rate.Limiter {
	limiter, exists := limiters.Load(ip)
	if !exists {
		newLimiter := rate.NewLimiter(rate.Limit(100), 200)
		limiters.Store(ip, newLimiter)

		go func() {
			time.Sleep(10 * time.Minute)
			limiters.Delete(ip)
		}()
		return newLimiter
	}
	return limiter.(*rate.Limiter)
}

func RateLimiterMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			return utils.Response(
				c,
				http.StatusTooManyRequests,
				"Too many requests, please try again later.",
				nil,
				nil,
				map[string]string{
					"retry_after": time.Now().Add(time.Second * 2).Format(time.RFC3339), // Estimasi waktu tunggu
				},
			)
		}

		return next(c)
	}
}
