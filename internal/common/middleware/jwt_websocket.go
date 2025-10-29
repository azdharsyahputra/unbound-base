package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func WebSocketAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing token")
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)

		// ✅ Gunakan "sub" karena itulah field ID user di token Unbound
		sub, ok := claims["sub"].(string)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid subject claim")
		}

		// Konversi string ke uint
		var userID uint
		_, err = fmt.Sscan(sub, &userID)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid user id in token")
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
