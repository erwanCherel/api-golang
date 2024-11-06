package authHandler

import (
	"fmt"
	"myapi/database"
	"myapi/models"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func Auth(c *fiber.Ctx) error {
	var authRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&authRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user models.User
	database.DB.Where("email = ?", authRequest.Email).First(&user)
	if user.ID == 0 {
		fmt.Println("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authRequest.Password))
	if err != nil {
		fmt.Println("Invalid Password:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	tokenClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Println("Failed to sign the token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create token",
		})
	}

	tokenRecord := models.Token{
		Code:      tokenString,
		ExpiredAt: time.Now().Add(24 * time.Hour),
		UserID:    user.ID,
	}

	if err := database.DB.Create(&tokenRecord).Error; err != nil {
		fmt.Println("Failed to store token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to store token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":      tokenString,
		"expires_at": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	})
}
