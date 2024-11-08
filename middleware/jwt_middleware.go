package middleware

import (
	"fmt"
	"myapi/database"
	"myapi/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header is missing",
		})
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["sub"].(float64)

		expiration := int64(claims["exp"].(float64))
		if time.Now().Unix() > expiration {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token has expired",
			})
		}

		var user models.User
		fmt.Println(userID)
		if err := database.DB.Where("id = ?", int(userID)).First(&user).Error; err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "User not found",
			})
		}

		c.Locals("user", user)

		return c.Next()
	}

	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid token",
	})
}
