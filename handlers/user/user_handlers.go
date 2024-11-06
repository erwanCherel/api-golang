package userHandler

import (
	"myapi/database"
	"myapi/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []models.User

	db.Find(&users)

	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "Not found", "data": nil})
	}

	return c.JSON(fiber.Map{"message": "Ok", "data": users})
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB
	user := new(models.User)

	err := c.BodyParser(user)
	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad Request",
			"code":    fiber.StatusBadRequest,
			"data":    err.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not hash password",
			"code":    fiber.StatusInternalServerError,
			"data":    err.Error(),
		})
	}

	user.Password = string(hashedPassword)

	err = db.Create(&user).Error
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create user",
			"code":    fiber.StatusInternalServerError,
			"data":    err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Created user",
		"data":    user,
	})
}

func GetUser(c *fiber.Ctx) error {
	db := database.DB
	var user models.User

	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {

		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid user ID", "data": nil})
	}

	result := db.First(&user, "id = ?", userID)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}
