package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/fikrihkll/chat-app/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		tokenString := authHeader[len("Bearer "):]

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Unexpected signing method")
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		if float64(time.Now().Unix()) > (*claims)["exp"].(float64) {
			return echo.NewHTTPError(http.StatusUnauthorized, common.UnauthorizedError.Error())
		}
		
		userID, ok := (*claims)["user_id"].(string)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid User ID in Token")
		}
		
		c.Set("id", userID)
		c.Set("email", (*claims)["email"])

		return next(c)
	}
}
